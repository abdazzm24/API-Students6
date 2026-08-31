package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() *slog.Logger {

	_ = os.MkdirAll(
		"logs",
		0755,
	)

	file := &lumberjack.Logger{
		Filename: filepath.Join(
			"logs",
			"app.log",
		),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     14,
		Compress:   true,
	}

	writer := io.MultiWriter(
		os.Stdout,
		file,
	)

	handler := slog.NewJSONHandler(
		writer,
		nil,
	)

	return slog.New(handler)
}