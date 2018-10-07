package s3pool

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/papa-rugi/go-filesiphon/file"
)

type s3Pool struct {
	client *s3.S3
	config *aws.Config
	params map[string]string
}

//Siphonable type is one that can Get or Put files.
type Siphonable interface {
	Get(path string) (io.Reader, error)
	Put(path string, file io.Reader) error
}

//FilePool type can do implements fundamental file operations.
type FilePool interface {
	Info() string
	ParsePath(p string) (string, string)
	Ls(path string) ([]os.FileInfo, error)
	Get(path string) (io.Reader, error)
	Put(path string, file io.Reader) error
	Mkdir(path string) error
	Rm(path string) error
	Cp(src string, dest string) error
	Mv(src string, dest string) error
}

//AwsCredentials provides a struct to write credentials in from a
type AwsCredentials struct {
	SecretAccessKey string
	AccessKey       string
	Region          string
}

//News3Pool creates a new s3Pool object
func News3Pool(params map[string]string) FilePool {
	if params["region"] == "" {
		params["region"] = "us-east-2"
	}

	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(params["access_key_id"], params["secret_access_key"], ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(params["region"]),
	}

	s3Pool := &s3Pool{
		config: config,
		client: s3.New(session.New(config)),
		params: params,
	}

	return s3Pool
}

//Returns info about this pool.
func (s s3Pool) Info() string {
	return "s3"
}

//ParsePath returns a bucket name and file path from a give string path.
func (s s3Pool) ParsePath(p string) (string, string) {
	sp := strings.Split(p, "/")
	bucket := ""
	if len(sp) > 1 {
		bucket = sp[1]
	}
	path := ""
	if len(sp) > 2 {
		path = strings.Join(sp[2:], "/")
	}

	return bucket, path
}

//Ls lists files from the client object, using a filepath as input.
func (s s3Pool) Ls(path string) ([]os.FileInfo, error) {
	bucketPath, filePath := s.ParsePath(path)
	files := make([]os.FileInfo, 0)

	if bucketPath == "" {
		b, err := s.client.ListBuckets(&s3.ListBucketsInput{})
		if err != nil {
			return nil, err
		}
		for _, bucket := range b.Buckets {
			files = append(files, &File{
				FName:   *bucket.Name,
				FType:   "directory",
				FTime:   bucket.CreationDate.UnixNano() / 1000,
				CanMove: false,
			})
		}
		return files, nil
	}

	objs, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(bucketPath),
		Prefix:    aws.String(filePath),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		return nil, err
	}

	for _, object := range objs.Contents {
		files = append(files, &File{
			FName: filepath.Base(*object.Key),
			FType: "file",
			FTime: object.LastModified.UnixNano() / 1000,
			FSize: *object.Size,
		})
	}

	for _, object := range objs.CommonPrefixes {
		files = append(files, &File{
			FName: filepath.Base(*object.Prefix),
			FType: "directory",
		})
	}

	return files, nil
}

//Get returns an io.Reader with the object data from a given string path.
func (s s3Pool) Get(path string) (io.Reader, error) {
	bucketPath, filePath := s.ParsePath(path)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketPath),
		Key:    aws.String(filePath),
	}

	//Get encryption working
	// if s.params["encryption_key"] != "" {
	// 	input.SSECustomerAlgorithm = aws.String("AES256")
	// 	input.SSECustomerKey = aws.String(s.params["encryption_key"])
	// }

	obj, err := s.client.GetObject(input)
	if err != nil {
		return nil, err
	}

	return obj.Body, nil
}

//Put writes to an io.Reader with the object data from a given string path.
func (s s3Pool) Put(path string, file io.Reader) error {
	bucketPath, filePath := s.ParsePath(path)

	if bucketPath == "" {
		return errors.New("Can't do that on S3")
	}
	uploader := s3manager.NewUploader(session.New(s.config))
	input := s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(bucketPath),
		Key:    aws.String(filePath),
	}

	_, err := uploader.Upload(&input)
	return err
}

//Mkdir creates a directory at string path indication.
func (s s3Pool) Mkdir(path string) error {
	bucketPath, filePath := s.ParsePath(path)

	if filePath == "" {
		_, err := s.client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(path),
		})
		return err
	}
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketPath),
		Key:    aws.String(filePath),
	})

	return err
}

//Rm removes an object at a given string path.
func (s s3Pool) Rm(path string) error {
	bucketPath, filePath := s.ParsePath(path)

	if bucketPath == "" {
		return errors.New("Doesn't exist")
	}

	objs, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(bucketPath),
		Prefix:    aws.String(filePath),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return err
	}
	for _, obj := range objs.Contents {
		_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketPath),
			Key:    obj.Key,
		})
		if err != nil {
			return err
		}
	}
	for _, pref := range objs.CommonPrefixes {
		s.Rm("/" + bucketPath + "/" + *pref.Prefix)
		_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketPath),
			Key:    pref.Prefix,
		})
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	if path == "" {
		_, err := s.client.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucketPath),
		})
		return err
	}
	_, err = s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketPath),
		Key:    aws.String(path),
	})

	return err
}

//Cp copies an object from a given src path, to a given dest path.
func (s s3Pool) Cp(src string, dest string) error {
	sBucket, sPath := s.ParsePath(src)
	dBucket, dPath := s.ParsePath(dest)

	if src == "" {
		return errors.New("Can't move this")
	}

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(dBucket),
		CopySource: aws.String(sBucket + "/" + sPath),
		Key:        aws.String(dPath),
	}

	// if s.params["encryption_key"] != "" {
	// 	input.CopySourceSSECustomerAlgorithm = aws.String("AES256")
	// 	input.CopySourceSSECustomerKey = aws.String(s.params["encryption_key"])
	// 	input.SSECustomerAlgorithm = aws.String("AES256")
	// 	input.SSECustomerKey = aws.String(s.params["encryption_key"])
	// }

	_, err := s.client.CopyObject(input)
	if err != nil {
		return err
	}

	return err
}

//Mv copies an object from a given src path, to a given dest path.
func (s s3Pool) Mv(src string, dest string) error {
	sBucket, sPath := s.ParsePath(src)
	dBucket, dPath := s.ParsePath(dest)

	if src == "" {
		return errors.New("Can't move this")
	}

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(dBucket),
		CopySource: aws.String(sBucket + "/" + sPath),
		Key:        aws.String(dPath),
	}
	// if s.params["encryption_key"] != "" {
	// 	input.CopySourceSSECustomerAlgorithm = aws.String("AES256")
	// 	input.CopySourceSSECustomerKey = aws.String(s.params["encryption_key"])
	// 	input.SSECustomerAlgorithm = aws.String("AES256")
	// 	input.SSECustomerKey = aws.String(s.params["encryption_key"])
	// }

	_, err := s.client.CopyObject(input)
	if err != nil {
		return err
	}

	return s.Rm(src)
}
