package plugin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
)

type Config struct {
	SKIP_JWT              bool
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	//AWS_DEFAULT_REGION    string
	AWS_S3_REGION string
	AWS_REGION    string
	//AWS_BUCKET            string
	S3_MOCK             bool
	S3_BUCKET           string
	S3_ENDPOINT         string
	S3_DISABLE_SSL      bool
	S3_FORCE_PATH_STYLE bool
	//REDIS_HOST            string
	//REDIS_PORT            string
	//REDIS_PASSWORD        string
	//SQS_ENDPOINT          string
}

func (c Config) EnvironmentVariables() []string {
	ret := make([]string, 12)
	ret[0] = fmt.Sprintf("SKIP_JWT=%v", c.SKIP_JWT)
	ret[1] = "AWS_ACCESS_KEY_ID=" + c.AWS_ACCESS_KEY_ID
	ret[2] = "AWS_SECRET_ACCESS_KEY=" + c.AWS_SECRET_ACCESS_KEY
	ret[3] = "AWS_DEFAULT_REGION=" + c.AWS_DEFAULT_REGION
	ret[4] = "AWS_S3_REGION=" + c.AWS_S3_REGION
	ret[5] = "AWS_REGION=" + c.AWS_REGION
	ret[6] = "AWS_BUCKET=" + c.AWS_BUCKET
	ret[7] = fmt.Sprintf("S3_MOCK=%v", c.S3_MOCK)
	ret[8] = "S3_BUCKET=" + c.S3_BUCKET
	ret[9] = "S3_ENDPOINT=" + c.S3_ENDPOINT
	ret[10] = fmt.Sprintf("S3_DISABLE_SSL=%v", c.S3_DISABLE_SSL)
	ret[11] = fmt.Sprintf("S3_FORCE_PATH_STYLE=%v", c.S3_FORCE_PATH_STYLE)
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
func (c Config) BatchEnvironmentVariables() []*batch.KeyValuePair {
	ret := make([]*batch.KeyValuePair, 9)
	ret[0] = toBatchKeyValuePair("AWS_ACCESS_KEY_ID", c.AWS_ACCESS_KEY_ID)
	ret[1] = toBatchKeyValuePair("AWS_SECRET_ACCESS_KEY", c.AWS_SECRET_ACCESS_KEY)
	ret[2] = toBatchKeyValuePair("AWS_DEFAULT_REGION", c.AWS_DEFAULT_REGION)
	ret[3] = toBatchKeyValuePair("AWS_S3_REGION", c.AWS_S3_REGION)
	ret[4] = toBatchKeyValuePair("S3_MOCK", fmt.Sprintf("%v", c.S3_MOCK))
	ret[5] = toBatchKeyValuePair("S3_BUCKET", c.S3_BUCKET)
	ret[6] = toBatchKeyValuePair("S3_ENDPOINT", c.S3_ENDPOINT)
	ret[7] = toBatchKeyValuePair("S3_DISABLE_SSL", fmt.Sprintf("%v", c.S3_DISABLE_SSL))
	ret[8] = toBatchKeyValuePair("S3_FORCE_PATH_STYLE", fmt.Sprintf("%v", c.S3_FORCE_PATH_STYLE))
	//ret[9] = toBatchKeyValuePair("REDIS_HOST", c.REDIS_HOST)
	//ret[10] = toBatchKeyValuePair("REDIS_PORT", c.REDIS_PORT)
	//ret[11] = toBatchKeyValuePair("REDIS_PASSWORD", c.REDIS_PASSWORD)
	//ret[12] = toBatchKeyValuePair("SQS_ENDPOINT", c.SQS_ENDPOINT)
	return ret
}
