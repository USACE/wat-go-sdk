package wat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/usace/wat-go-sdk/plugin"
)

type AwsConfig struct {
	Name                  string `json:"aws_config_name,omitempty"`
	IsPrimary             bool   `json:"is_primary_config"` //where payloads would get stored?
	AWS_ACCESS_KEY_ID     string `json:"aws_access_key_id"`
	AWS_SECRET_ACCESS_KEY string `json:"aws_secret_access_key_id"`
	AWS_REGION            string `json:"aws_region"`
	AWS_BUCKET            string `json:"aws_bucket"`
	S3_MOCK               bool   `json:"aws_mock,omitempty"`             //for testing with minio
	S3_ENDPOINT           string `json:"aws_endpoint,omitempty"`         //for testing with minio
	S3_DISABLE_SSL        bool   `json:"aws_disable_ssl,omitempty"`      //for testing with minio
	S3_FORCE_PATH_STYLE   bool   `json:"aws_force_path_style,omitempty"` //for testing with minio
}

type Config struct {
	CloudProvider Provider    `json:"cloud_provider_type"`
	AwsConfigs    []AwsConfig `json:"aws_configs"`
}

func (c Config) PrimaryConfig() (AwsConfig, error) {

	for _, ac := range c.AwsConfigs {
		if ac.IsPrimary {
			return ac, nil
		}
	}
	nilconfig := AwsConfig{}
	return nilconfig, errors.New("no config was marked as primary")
}
func InitConfig(path string) (Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return cfg, err
	}
	pc, err := cfg.PrimaryConfig()
	if err != nil {
		return cfg, err
	}
	pac := plugin.AwsConfig{
		Name:                  "wat-config",
		IsPrimary:             true,
		AWS_ACCESS_KEY_ID:     pc.AWS_ACCESS_KEY_ID,
		AWS_SECRET_ACCESS_KEY: pc.AWS_SECRET_ACCESS_KEY,
		AWS_REGION:            pc.AWS_REGION,
		AWS_BUCKET:            pc.AWS_BUCKET,
		S3_MOCK:               pc.S3_MOCK,
		S3_ENDPOINT:           pc.S3_ENDPOINT,
		S3_DISABLE_SSL:        pc.S3_DISABLE_SSL,
		S3_FORCE_PATH_STYLE:   pc.S3_FORCE_PATH_STYLE,
	}
	pacs := make([]plugin.AwsConfig, 1)
	pacs[0] = pac
	pcfg := plugin.Config{
		AwsConfigs: pacs,
	}
	err = plugin.InitConfig(pcfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
func (c AwsConfig) EnvironmentVariables() []string {
	ret := make([]string, 8)
	ret[0] = "AWS_ACCESS_KEY_ID=" + c.AWS_ACCESS_KEY_ID
	ret[1] = "AWS_SECRET_ACCESS_KEY=" + c.AWS_SECRET_ACCESS_KEY
	ret[2] = "AWS_BUCKET=" + c.AWS_BUCKET
	ret[3] = "AWS_REGION=" + c.AWS_REGION
	ret[4] = fmt.Sprintf("S3_MOCK=%v", c.S3_MOCK)
	ret[5] = "S3_ENDPOINT=" + c.S3_ENDPOINT
	ret[6] = fmt.Sprintf("S3_DISABLE_SSL=%v", c.S3_DISABLE_SSL)
	ret[7] = fmt.Sprintf("S3_FORCE_PATH_STYLE=%v", c.S3_FORCE_PATH_STYLE)
	//ret[13] = "REDIS_HOST=" + c.REDIS_HOST
	//ret[14] = "REDIS_PORT=" + c.REDIS_PORT
	//ret[15] = "REDIS_PASSWORD=" + c.REDIS_PASSWORD
	//ret[16] = "SQS_ENDPOINT=" + c.SQS_ENDPOINT
	return ret
}
