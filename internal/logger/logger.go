package logger

import (
  "github.com/sirupsen/logrus"
  "os"
)

func New(levelName string) (*logrus.Logger, error) {
  level, err := logrus.ParseLevel(levelName)
  if err != nil {
    return nil, err
  }

  logger := logrus.New()
  logger.Out = os.Stdout
  logger.SetFormatter(&logrus.JSONFormatter{})
  logger.SetLevel(level)

  return logger, nil
}
