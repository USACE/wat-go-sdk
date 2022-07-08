package plugin

import (
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
	Level
}

var Logger = GlobalLogger{
	Level: INFO,
}

type Log struct {
	Message string `json:"message"`
	Level   Level  `json:"level"`
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
type Services struct {
	config   Config
	stores   map[string]filestore.FileStore //should this be an array of file store? indexed by bucket name?
	loglevel Level

	//sqs
	//redis
	//paul-bunyan

}

func InitServices(prefix string) (Services, error) {
	var cfg Config
	//zerolog.SetGlobalLevel(zerolog.InfoLevel) //set from config?
	s := Services{}
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return s, err
	}
	s.config = cfg
	s.stores = make(map[string]filestore.FileStore)
	_, err := s.getStore(cfg.S3_BUCKET)
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
func (s *Services) getStore(bucketName string) (filestore.FileStore, error) {
	fs, ok := s.stores[bucketName]
	if !ok {
		//initalize S3 Store
		mock := s.config.S3_MOCK
		s3Conf := filestore.S3FSConfig{
			S3Id:     s.config.AWS_ACCESS_KEY_ID,
			S3Key:    s.config.AWS_SECRET_ACCESS_KEY,
			S3Region: s.config.AWS_REGION,
			S3Bucket: bucketName,
		}
		if mock {
			s3Conf.Mock = mock
			s3Conf.S3DisableSSL = s.config.S3_DISABLE_SSL
			s3Conf.S3ForcePathStyle = s.config.S3_FORCE_PATH_STYLE
			s3Conf.S3Endpoint = s.config.S3_ENDPOINT
		}
		//fmt.Println(s3Conf)

		nfs, err := filestore.NewFileStore(s3Conf)
		fs = nfs
		if err != nil {
			log := Log{
				Message: err.Error(),
				Level:   FATAL,
			}
			Logger.Log(log)
		}
		s.stores[bucketName] = fs
	}

	return fs, nil
}
func (s Services) ReportProgress(report ProgressReport) {
	//can be placeholder.
	log.Info().Msg(fmt.Sprintf("Progress: %v, %v", report.Progress, report.Message))
}
func (s Services) ReportStatus(report StatusReport) {
	//can be placeholder.
	log.Info().Msg(fmt.Sprintf("Status: %v, %v", report.Status.String(), report.Message))
}
func (logger *GlobalLogger) SetLogLevel(logLevel Level) {
	switch logLevel {
	case DEBUG:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case INFO:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case WARN:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case ERROR:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case FATAL:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case PANIC:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case DISABLED:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	logger.Level = logLevel
}
func (logger GlobalLogger) Log(logmessage Log) {
	//using zerolog is a placeholder, could use SQS or Redis or whatever we want.
	if logger.Level >= logmessage.Level {
		switch logmessage.Level {
		case DEBUG:
			log.Debug().Msg(logmessage.Message)
		case INFO:
			log.Info().Msg(logmessage.Message)
		case WARN:
			log.Warn().Msg(logmessage.Message)
		case ERROR:
			log.Error().Msg(logmessage.Message)
		case FATAL:
			log.Fatal().Msg(logmessage.Message)
		case PANIC:
			log.Panic().Msg(logmessage.Message)
		case DISABLED:
			//log.Info().Msg(message)
		default:
			log.Info().Msg(logmessage.Message)
		}
	}
}
func (s *Services) LoadPayload(filepath string) (ModelPayload, error) {
	Logger.Log(Log{
		Message: fmt.Sprintf("reading:%v", filepath),
		Level:   INFO,
	})
	payload := ModelPayload{}
	fs, err := s.getStore(s.config.S3_BUCKET)
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
		})
		return payload, err
	}

	return payload, nil
}

// UpLoadFile
func (s *Services) UpLoadFile(resource ResourceInfo, fileBytes []byte) error {
	if resource.Store != "S3" {
		//check if local?
		return errors.New("the resource is not defined as s3")
	}
	if strings.Contains(resource.Path, "../") {
		return errors.New("it is against policy to have relative paths for an s3 store")
	}
	fs, err := s.getStore(resource.Root)
	if err != nil {
		return err
	}
	_, err = fs.PutObject(resource.Path, fileBytes)
	if err != nil {
		return err
	}

	return err
}
