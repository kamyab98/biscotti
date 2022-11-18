package handler

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type Storage interface {
	Store([]byte) error
	Flush() error
}

type LoggerStorage struct {
	logger            *zap.Logger
	bucketSize        int32
	numberOfUnflushed int32
}

type LoggerStorageConfig struct {
	Location string
}

var loggerStorage *LoggerStorage

func GetLoggerStorage() *LoggerStorage {
	if loggerStorage == nil {
		logsPath := filepath.Join(".", "logs")
		err := os.MkdirAll(logsPath, os.ModePerm)
		if err != nil {
			panic(err)
		}

		loggerStorageConfig := LoggerStorageConfig{Location: "logs/logs.log"}
		loggerStorage, _ = loggerStorageConfig.Build()
	}
	return loggerStorage
}

func (l LoggerStorageConfig) Build() (*LoggerStorage, error) {
	var cfg zap.Config
	loggerConfig := []byte(`{
    "level": "debug",
    "encoding": "json",
    "outputPaths": ["` + l.Location + `"],
    "errorOutputPaths": ["stderr"],
    "encoderConfig": {
      "messageKey": "message"
    }
  }`)
	if err := json.Unmarshal(loggerConfig, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	 return &LoggerStorage{logger: logger, bucketSize: 1000, numberOfUnflushed: 0}, nil
}

func (l LoggerStorage) Store(data []byte) error {
	l.logger.Info(string(data))
	l.numberOfUnflushed++
	if l.numberOfUnflushed > l.bucketSize {
		err := l.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l LoggerStorage) Flush() error {
	err := l.logger.Sync()
	if err != nil {
		return err
	}
	return nil
}
