package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/USACE/filestore"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Services struct {
	config Config
	fs     filestore.FileStore
	//sqs
	//redis
	//paul-bunyan
}

func InitServices(prefix string) (Services, error) {
	var cfg Config
	s := Services{}
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return s, err
	}
	s.config = cfg
	err := s.initStore()
	if err != nil {
		return s, err
	}
	return s, nil
}

func (s Services) EnvironmentVariables() []string {
	return s.config.EnvironmentVariables()
}
func (s Services) Config() Config {
	return s.config
}
func (s *Services) initStore() error {
	//initalize S3 Store
	mock := s.config.S3_MOCK
	s3Conf := filestore.S3FSConfig{
		S3Id:     s.config.AWS_ACCESS_KEY_ID,
		S3Key:    s.config.AWS_SECRET_ACCESS_KEY,
		S3Region: s.config.AWS_REGION,
		S3Bucket: s.config.S3_BUCKET,
	}
	if mock {
		s3Conf.Mock = mock
		s3Conf.S3DisableSSL = s.config.S3_DISABLE_SSL
		s3Conf.S3ForcePathStyle = s.config.S3_FORCE_PATH_STYLE
		s3Conf.S3Endpoint = s.config.S3_ENDPOINT
	}
	fmt.Println(s3Conf)

	fs, err := filestore.NewFileStore(s3Conf)

	if err != nil {
		log.Fatal(err)
	}
	s.fs = fs
	return nil
}

func (s Services) LoadJsonFile(filepath string, spec interface{}) error {
	fmt.Println("reading:", filepath)
	data, err := s.fs.GetObject(filepath)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	errjson := json.Unmarshal(body, &spec)
	if errjson != nil {
		fmt.Println("error:", errjson)
		return errjson
	}

	return nil

}
func (s Services) LoadYamlFile(filepath string, spec interface{}) error {
	fmt.Println("reading:", filepath)
	data, err := s.fs.GetObject(filepath)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	errjson := yaml.Unmarshal(body, &spec)
	if errjson != nil {
		fmt.Println("error:", errjson)
		return errjson
	}

	return nil

}

// UpLoadToS3
func (s Services) UpLoadFile(newS3Path string, fileBytes []byte) (filestore.FileOperationOutput, error) {
	var repsonse *filestore.FileOperationOutput
	var err error
	repsonse, err = s.fs.PutObject(newS3Path, fileBytes)
	if err != nil {
		return *repsonse, err
	}

	return *repsonse, err
}
