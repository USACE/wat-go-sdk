package plugin

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
type Config struct {
	AwsConfigs []AwsConfig `json:"aws_configs"`
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

/*
func (c AwsConfig) EnvironmentVariables() []string {
	ret := make([]string, 7)
	ret[0] = "AWS_ACCESS_KEY_ID=" + c.AWS_ACCESS_KEY_ID
	ret[1] = "AWS_SECRET_ACCESS_KEY=" + c.AWS_SECRET_ACCESS_KEY
	ret[2] = "AWS_BUCKET=" + c.AWS_BUCKET
	ret[3] = fmt.Sprintf("S3_MOCK=%v", c.S3_MOCK)
	ret[4] = "S3_ENDPOINT=" + c.S3_ENDPOINT
	ret[5] = fmt.Sprintf("S3_DISABLE_SSL=%v", c.S3_DISABLE_SSL)
	ret[6] = fmt.Sprintf("S3_FORCE_PATH_STYLE=%v", c.S3_FORCE_PATH_STYLE)
	//ret[13] = "REDIS_HOST=" + c.REDIS_HOST
	//ret[14] = "REDIS_PORT=" + c.REDIS_PORT
	//ret[15] = "REDIS_PASSWORD=" + c.REDIS_PASSWORD
	//ret[16] = "SQS_ENDPOINT=" + c.SQS_ENDPOINT
	return ret
}
func toBatchKeyValuePair(key string, value string) *batch.KeyValuePair {
	keyvalue := batch.KeyValuePair{
		Name:  aws.String(key),
		Value: aws.String(value),
	}
	return &keyvalue
}

//this is realy useful in WAT (but maybe not in plugin.utilities)
func (c AwsConfig) BatchEnvironmentVariables() []*batch.KeyValuePair {
	ret := make([]*batch.KeyValuePair, 7)
	ret[0] = toBatchKeyValuePair("AWS_ACCESS_KEY_ID", c.AWS_ACCESS_KEY_ID)
	ret[1] = toBatchKeyValuePair("AWS_SECRET_ACCESS_KEY", c.AWS_SECRET_ACCESS_KEY)
	ret[2] = toBatchKeyValuePair("AWS_BUCKET", c.AWS_BUCKET)
	ret[3] = toBatchKeyValuePair("S3_MOCK", fmt.Sprintf("%v", c.S3_MOCK))
	ret[4] = toBatchKeyValuePair("S3_ENDPOINT", c.S3_ENDPOINT)
	ret[5] = toBatchKeyValuePair("S3_DISABLE_SSL", fmt.Sprintf("%v", c.S3_DISABLE_SSL))
	ret[6] = toBatchKeyValuePair("S3_FORCE_PATH_STYLE", fmt.Sprintf("%v", c.S3_FORCE_PATH_STYLE))
	//ret[9] = toBatchKeyValuePair("REDIS_HOST", c.REDIS_HOST)
	//ret[10] = toBatchKeyValuePair("REDIS_PORT", c.REDIS_PORT)
	//ret[11] = toBatchKeyValuePair("REDIS_PASSWORD", c.REDIS_PASSWORD)
	//ret[12] = toBatchKeyValuePair("SQS_ENDPOINT", c.SQS_ENDPOINT)
	return ret
}
*/
