package wat

import "errors"

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
type Provider string

const (
	BATCH Provider = "AWS Batch"
)

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
	return nilconfig, errors.New("No config was marked as primary.")
}
