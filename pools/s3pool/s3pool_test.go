package s3pool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLs(t *testing.T) {
	file, _ := os.Open("awsCreds1.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	awsCredentials := AwsCredentials{}
	err := decoder.Decode(&awsCredentials)
	if err != nil {
		fmt.Println("error:", err)
	}

	s3 := News3Pool(map[string]string{
		"region":            awsCredentials.Region,
		"access_key_id":     awsCredentials.AccessKey,
		"secret_access_key": awsCredentials.SecretAccessKey,
	})

	//List at path
	objectList, err := s3.Ls("/upload-testing-to-test")

	if err != nil {
		t.Log("Error", err)
	}

	for _, o := range objectList {
		t.Logf("%s\n", o.Name())
	}

	t.Log("List recieved")
}

func TestGet(t *testing.T) {
	file, _ := os.Open("awsCreds1.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	awsCredentials := AwsCredentials{}
	err := decoder.Decode(&awsCredentials)
	if err != nil {
		t.Log("error:", err)
	}

	s3 := News3Pool(map[string]string{
		"region":            awsCredentials.Region,
		"access_key_id":     awsCredentials.AccessKey,
		"secret_access_key": awsCredentials.SecretAccessKey,
	})

	objectData, err := s3.Get("/upload-testing-to-test/gotest.js")

	if err != nil {
		t.Log("Error", err)
	}

	data, err := ioutil.ReadAll(objectData)

	t.Log(string(data))
	t.Log("File retrieved.")
}

func TestPut(t *testing.T) {
	file, _ := os.Open("../../awsCreds1.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	awsCredentials := AwsCredentials{}
	decoder.Decode(&awsCredentials)

	s3 := News3Pool(map[string]string{
		"region":            awsCredentials.Region,
		"access_key_id":     awsCredentials.AccessKey,
		"secret_access_key": awsCredentials.SecretAccessKey,
	})

	err := s3.Put("/upload-testing-to-test/go-test2", strings.NewReader("Go is great, File-Siphon is good"))

	if err != nil {
		t.Log("Error", err)
	}

	t.Log("File put.")
}

func TestMkdir(t *testing.T) {
	file, _ := os.Open("awsCreds1.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	awsCredentials := AwsCredentials{}
	decoder.Decode(&awsCredentials)

	s3 := News3Pool(map[string]string{
		"region":            awsCredentials.Region,
		"access_key_id":     awsCredentials.AccessKey,
		"secret_access_key": awsCredentials.SecretAccessKey,
	})

	err := s3.Mkdir("/upload-testing-to-test/gotest")

	if err != nil {
		fmt.Println("Error", err)
	}

	t.Log("Dir created.")
}

func TestRm(t *testing.T) {
	file, _ := os.Open("awsCreds1.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	awsCredentials := AwsCredentials{}
	decoder.Decode(&awsCredentials)

	s3 := News3Pool(map[string]string{
		"region":            awsCredentials.Region,
		"access_key_id":     awsCredentials.AccessKey,
		"secret_access_key": awsCredentials.SecretAccessKey,
	})

	err := s3.Rm("/upload-testing-to-test/myfile")

	if err != nil {
		t.Log("Error", err)
	}

	t.Log("File removed")
}

func TestCp(t *testing.T) {
	file, _ := os.Open("awsCreds1.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	awsCredentials := AwsCredentials{}
	decoder.Decode(&awsCredentials)

	s3 := News3Pool(map[string]string{
		"region":            awsCredentials.Region,
		"access_key_id":     awsCredentials.AccessKey,
		"secret_access_key": awsCredentials.SecretAccessKey,
	})

	err := s3.Mv("/upload-testing-to-test/go-test2", "/upload-testing-to-test/gotest.js")

	if err != nil {
		t.Log("Error", err)
	}

	t.Log("File copied.")
}

func TestMv(t *testing.T) {
	file, _ := os.Open("awsCreds1.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	awsCredentials := AwsCredentials{}
	decoder.Decode(&awsCredentials)

	s3 := News3Pool(map[string]string{
		"region":            awsCredentials.Region,
		"access_key_id":     awsCredentials.AccessKey,
		"secret_access_key": awsCredentials.SecretAccessKey,
	})

	err := s3.Mv("/upload-testing-to-test/gotest.js", "/upload-testing-to-test/gotestmv.js")

	if err != nil {
		t.Log("Error", err)
	}

	t.Log("File moved.")
}
