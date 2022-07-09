package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/USACE/filestore"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Level uint8

const (
	INFO Level = iota + 1
	WARN
	ERROR
	DEBUG
	FATAL
	PANIC
	DISABLED
)

type GlobalLogger struct {
	logger zerolog.Logger
	Level  //i believe this will be global to the container each container having its own possible level (and wat having its own level too.)
}
type GlobalConfig struct {
	HasInitialized bool
	Config
	stores map[string]filestore.FileStore
}

var PluginConfig = GlobalConfig{
	HasInitialized: false,
}
var Logger = GlobalLogger{
	Level: INFO,
}

type Log struct {
	Message string `json:"message"`
	Level   Level  `json:"loglevel"`
	Sender  string `json:"sender"`
}

//zeroLog is a struct to parse the returned log from zerolog for the purpose of styling log outputs if SetStyle is used.
type zeroLog struct {
	Message string `json:"message"`
	Level   string `json:"level"`
	Sender  string `json:"sender"` //custom string feild
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
				Logger.Log(Log{
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
			Logger.Log(log)
		}
		PluginConfig.stores[bucketName] = fs
	}

	return fs, nil
}
func ReportProgress(report ProgressReport, linkedManifestId string) {
	//can be placeholder.
	log.Info().Msg(fmt.Sprintf("Manifest: %v\n\tProgress: %v, %v", linkedManifestId, report.Progress, report.Message))
}
func ReportStatus(report StatusReport, linkedManifestId string) {
	//can be placeholder.
	log.Info().Msg(fmt.Sprintf("Manifest: %v\n\tStatus: %v, %v", linkedManifestId, report.Status.String(), report.Message))
}

type logWriter struct {
}

func (w logWriter) Write(b []byte) (n int, err error) {
	log := zeroLog{}
	errjson := json.Unmarshal(b, &log)
	if errjson != nil {
		return 1, errjson
	}
	fmt.Printf("%v issues %v\n\t%v\n", log.Sender, log.Level, log.Message)
	return 0, nil
}
func (logger *GlobalLogger) SetStyle() {
	w := logWriter{}
	logger.logger = zerolog.New(w)
}
func (logger *GlobalLogger) SetLogLevel(logLevel Level) {
	switch logLevel {
	case DEBUG:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case INFO:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case WARN:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case ERROR:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case FATAL:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case PANIC:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case DISABLED:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	logger.Level = logLevel
}
func (logger GlobalLogger) Log(LogMessage Log) {
	//using zerolog is a placeholder, could use SQS or Redis or whatever we want.
	if logger.Level <= LogMessage.Level {
		switch LogMessage.Level {
		case DEBUG:
			logger.logger.Debug().Str("sender", LogMessage.Sender).Msg(LogMessage.Message)
		case INFO:
			logger.logger.Info().Str("sender", LogMessage.Sender).Msg(LogMessage.Message)
		case WARN:
			logger.logger.Warn().Str("sender", LogMessage.Sender).Msg(LogMessage.Message)
		case ERROR:
			logger.logger.Error().Str("sender", LogMessage.Sender).Msg(LogMessage.Message)
		case FATAL:
			logger.logger.Fatal().Str("sender", LogMessage.Sender).Msg(LogMessage.Message)
		case PANIC:
			logger.logger.Panic().Str("sender", LogMessage.Sender).Msg(LogMessage.Message)
		case DISABLED:
			//log.Info().Msg(message)
		default:
			logger.logger.Info().Str("sender", LogMessage.Sender).Msg(LogMessage.Message)
		}
	}
}
func LoadPayload(filepath string) (ModelPayload, error) {
	Logger.Log(Log{
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
		Logger.Log(Log{
			Message: fmt.Sprintf("error reading:%v", filepath),
			Level:   ERROR,
			Sender:  "Plugin Services",
		})
		return payload, err
	}

	return payload, nil
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
