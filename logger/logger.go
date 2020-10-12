package logger

import (
    "github.com/sirupsen/logrus"
)

type LogrusLog struct {
}

var (
    Logger *LogrusLog
)

func (*LogrusLog) Info(args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Info(args)
}

func (*LogrusLog) Infof(format string, args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Infof(format, args)
}

func (*LogrusLog) Error(args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Error(args)
}

func (*LogrusLog) Errorf(format string, args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Errorf(format, args)
}

func (*LogrusLog) Debug(args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Debug(args)
}

func (*LogrusLog) Debugf(format string, args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Debugf(format, args)
}

func (*LogrusLog) Panic(args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Panic(args)
}

func (*LogrusLog) Panicf(format string, args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Panicf(format, args)
}

func (*LogrusLog) Fatal(args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Fatal(args)
}

func (*LogrusLog) Fatalf(format string, args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Fatalf(format, args)
}

func (*LogrusLog) Warn(args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Warn(args)
}

func (*LogrusLog) Warnf(format string, args ...interface{}) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.Warnf(format, args)
}
