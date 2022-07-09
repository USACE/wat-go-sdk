package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/USACE/filestore"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Level uint8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
	PANIC
	DISABLED
)

func (l Level) String() string {
	switch l {
	case INFO:
		return "some Information"
	case WARN:
		return "a Warning"
	case ERROR:
		return "an Error"
	case DEBUG:
		return "a Debug statement"
	case FATAL:
		return "a Fatal message"
	case PANIC:
		return "a Panic'ed state"
	case DISABLED:
		return ""
	default:
		return "Unknown Level"
	}
}

type GlobalLogger struct {
	Level //i believe this will be global to the container each container having its own possible level (and wat having its own level too.)
}
type GlobalConfig struct {
	HasInitialized bool
	Config
	stores map[string]filestore.FileStore
}

var PluginConfig = GlobalConfig{
	HasInitialized: false,
}
var logger = GlobalLogger{
	Level: INFO,
}

type Log struct {
	Message string `json:"message"`
	Level   Level  `json:"loglevel"`
	Sender  string `json:"sender"`
}

type Status uint8

const (
	COMPUTING Status = iota
	FAILED
	SUCCEEDED
)

func (s Status) String() string {
	switch s {
	case COMPUTING:
		return "Computing"
	case FAILED:
		return "Failed"
	case SUCCEEDED:
		return "Succeeded"
	default:
		return "Unknown Status"
	}
}

type StatusReport struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
}
type ProgressReport struct {
	Progress int8   `json:"progress"` //whole integers from 0 to 100...
	Message  string `json:"message"`
}

func initConfig() error {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return err
	}
	PluginConfig.Config = cfg
	PluginConfig.stores = make(map[string]filestore.FileStore)
	return nil
}

func EnvironmentVariables() []string {
	return PluginConfig.EnvironmentVariables()
}
func GetConfig() Config {
	return PluginConfig.Config
}
func getStore(bucketName string) (filestore.FileStore, error) {
	fs, ok := PluginConfig.stores[bucketName]
	if !ok {
		//check if config exists.
		if !PluginConfig.HasInitialized {
			err := initConfig()
			if err != nil {
				SubmitLog(Log{
					Message: "Could not Initialize Plugin Configurations, do you have an .env file",
					Level:   FATAL,
					Sender:  "Plugin Utilities",
				})
			}
		}
		//initalize S3 Store
		mock := PluginConfig.S3_MOCK
		s3Conf := filestore.S3FSConfig{
			S3Id:     PluginConfig.AWS_ACCESS_KEY_ID,
			S3Key:    PluginConfig.AWS_SECRET_ACCESS_KEY,
			S3Region: PluginConfig.AWS_REGION,
			S3Bucket: bucketName, //why would more than one bucket have the same keys?
		}
		if mock {
			s3Conf.Mock = mock
			s3Conf.S3DisableSSL = PluginConfig.S3_DISABLE_SSL
			s3Conf.S3ForcePathStyle = PluginConfig.S3_FORCE_PATH_STYLE
			s3Conf.S3Endpoint = PluginConfig.S3_ENDPOINT
		}
		nfs, err := filestore.NewFileStore(s3Conf)
		fs = nfs
		if err != nil {
			log := Log{
				Message: err.Error(),
				Level:   FATAL,
				Sender:  "Plugin Services",
			}
			SubmitLog(log)
		}
		PluginConfig.stores[bucketName] = fs
	}

	return fs, nil
}
func ReportProgress(report ProgressReport, linkedManifestId string) {
	//can be placeholder.
	log := Log{
		Message: fmt.Sprintf("Manifest: %v\n\tProgress: %v, %v", linkedManifestId, report.Progress, report.Message),
		Level:   INFO,
		Sender:  "Progress Reporter",
	}
	logger.write(log)
}
func ReportStatus(report StatusReport, linkedManifestId string) {
	//can be placeholder.
	log := Log{
		Message: fmt.Sprintf("Manifest: %v\n\tStatus: %v, %v", linkedManifestId, report.Status.String(), report.Message),
		Level:   INFO,
		Sender:  "Status Reporter",
	}
	logger.write(log)
}

func (l GlobalLogger) write(log Log) (n int, err error) {
	sender := ""
	if log.Sender == "" {
		sender = "Unknown Sender"
	} else {
		sender = log.Sender
	}
	fmt.Printf("%v issues %v\n\t%v\n", sender, log.Level.String(), log.Message)
	return 0, nil
}

func SetLogLevel(logLevel Level) {
	logger.Level = logLevel
}
func SubmitLog(LogMessage Log) {
	//using zerolog is a placeholder, could use SQS or Redis or whatever we want.
	if logger.Level <= LogMessage.Level {
		logger.write(LogMessage)
	}
}
func LoadPayload(filepath string) (ModelPayload, error) {
	SubmitLog(Log{
		Message: fmt.Sprintf("reading:%v", filepath),
		Level:   INFO,
		Sender:  "Plugin Services",
	})
	payload := ModelPayload{}
	fs, err := getStore(PluginConfig.AWS_BUCKET)
	if err != nil {
		return payload, err
	}
	data, err := fs.GetObject(filepath)
	if err != nil {
		return payload, err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return payload, err
	}

	err = yaml.Unmarshal(body, &payload)
	if err != nil {
		SubmitLog(Log{
			Message: fmt.Sprintf("error reading:%v", filepath),
			Level:   ERROR,
			Sender:  "Plugin Services",
		})
		return payload, err
	}

	return payload, nil
}
func FetchPayloadInputs(payload ModelPayload, localRoot string) error {
	for _, fileData := range payload.Inputs {
		bytes, err := DownloadObject(fileData.ResourceInfo)
		if err != nil {
			return err
		}
		//write bytes.
		writeBytes(bytes, localRoot, fileData.ResourceInfo.Path)
		//check for other files?
		if len(fileData.InternalPaths) > 0 {
			for _, internalPath := range fileData.InternalPaths {
				bytes, err := DownloadObject(internalPath.ResourceInfo)
				if err != nil {
					return err
				}
				writeBytes(bytes, localRoot, internalPath.ResourceInfo.Path)
			}
		}
	}
	return nil
}
func writeBytes(b []byte, destinationRoot string, destinationPath string) error {
	if _, err := os.Stat(destinationRoot); os.IsNotExist(err) {
		os.MkdirAll(destinationRoot, 0644) //do i need to trim filename?
	}
	err := os.WriteFile(destinationPath, b, 0644)
	if err != nil {
		SubmitLog(Log{
			Message: fmt.Sprintf("failure to write local file: %v\n\terror:%v", destinationPath, err),
			Level:   ERROR,
			Sender:  "Plugin Utilities",
		})
		return err
	}
	return nil
}
func DownloadObject(resource ResourceInfo) ([]byte, error) {
	SubmitLog(Log{
		Message: fmt.Sprintf("reading:%v", resource.Path),
		Level:   INFO,
		Sender:  "Plugin Services",
	})
	fs, err := getStore(resource.Root)
	if err != nil {
		return nil, err
	}
	data, err := fs.GetObject(resource.Path)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// UpLoadFile
func UpLoadFile(resource ResourceInfo, fileBytes []byte) error {
	if resource.Store != "S3" {
		//check if local?
		return errors.New("the resource is not defined as s3")
	}
	if strings.Contains(resource.Path, "../") {
		return errors.New("it is against policy to have relative paths for an s3 store")
	}
	fs, err := getStore(resource.Root) //how can we be sure we have the right secrets?
	if err != nil {
		return err
	}
	_, err = fs.PutObject(resource.Path, fileBytes)
	if err != nil {
		return err
	}

	return err
}
