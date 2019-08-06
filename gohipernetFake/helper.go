package gohipernetFake

import "go.uber.org/zap/zapcore"

var NTELIB_LOG_DEBUG func(msg string, fields ...zapcore.Field)
var NTELIB_LOG_INFO func(msg string, fields ...zapcore.Field)
var NTELIB_LOG_ERROR func(msg string, fields ...zapcore.Field)

func wrapLoggerFunc() {
	NTELIB_LOG_DEBUG = Logger.Debug
	NTELIB_LOG_INFO = Logger.Info
	NTELIB_LOG_ERROR = Logger.Error
}
