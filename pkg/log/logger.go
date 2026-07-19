package log

import (
	"io"
	"log/slog"
	"os"
	"github.com/lmittmann/tint"
	"time"
)

func InitLogger(isProd bool) io.Closer {
	var logOutput io.Writer = os.Stdout
	var fileToClose io.Closer

	if isProd {
		logFile, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("failed to open log file", "error", err)
			os.Exit(1)
		}
		fileToClose = logFile
		logOutput = io.MultiWriter(os.Stdout, logFile)
	}

	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	var handler slog.Handler
	if isProd {
		handler = slog.NewJSONHandler(logOutput, opts)
	} else {
        handler = tint.NewHandler(logOutput, &tint.Options{
            Level:      slog.LevelDebug,
            TimeFormat: time.Kitchen,
            AddSource:  true,
        })
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return fileToClose
}