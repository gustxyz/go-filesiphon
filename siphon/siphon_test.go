package siphon

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/papa-rugi/go-filesiphon/pools/s3pool"
)

func TestSiphonFile(t *testing.T) {
	file, _ := os.Open("../awsCreds1.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	awsCredentials := AwsCredentials{}
	err := decoder.Decode(&awsCredentials)
	if err != nil {
		t.Log("error:", err)
	}

	s3Source := News3Pool(map[string]string{
		"region":            awsCredentials.Region,
		"access_key_id":     awsCredentials.AccessKey,
		"secret_access_key": awsCredentials.SecretAccessKey,
	})

	err1 := s3Source.Put("/upload-testing-to-test/siphontest/srcSiphon.txt", strings.NewReader("This file was uploaded to the source filepool at path /upload-testing-to-test/siphontest/srcSiphon.txt and then siphoned to the dest pool"))

	if err1 != nil {
		t.Log("Error", err1)
	}

	t.Log("File placed at source.")

	s3Dest := News3Pool(map[string]string{
		"region":            awsCredentials.Region,
		"access_key_id":     awsCredentials.AccessKey,
		"secret_access_key": awsCredentials.SecretAccessKey,
	})

	err2 := siphonFile(s3Source, "/upload-testing-to-test/siphontest/srcSiphon.txt",
		s3Dest, "/upload-testing-to-test/siphontest/destput/srcSiphon.txt")

	if err2 != nil {
		t.Log("Error", err2)
	}

	objectData, err3 := s3Dest.Get("/upload-testing-to-test/siphontest/destput/srcSiphon.txt")

	if err3 != nil {
		t.Log("Error", err3)
	}

	data, err3 := ioutil.ReadAll(objectData)

	t.Log(string(data))
}
