/**
** author: duan hongxing
** date: 17 Jul 2024
** desc:
**	logger
**/

package logger

import (
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once

var log zerolog.Logger

func Get() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339

		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zerolog.InfoLevel) // default to INFO
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if os.Getenv("APP_ENV") != "development" {
			fileLogger := &lumberjack.Logger{
				Filename:   "logs/server.log",
				MaxSize:    5, //
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		/*
			var gitRevision string

			buildInfo, ok := debug.ReadBuildInfo()
			if ok {
				for _, v := range buildInfo.Settings {
					if v.Key == "vcs.revision" {
						gitRevision = v.Value
						break
					}
				}
			}
		*/

		/*
			output.FormatMessage = func(i interface{}) string {
				return fmt.Sprintf("***%s****", i)
			}
			output.FormatFieldName = func(i interface{}) string {
				return fmt.Sprintf("%s:", i)
			}
			output.FormatFieldValue = func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("%s", i))
			}
		*/

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			//Str("git_revision", gitRevision).
			//Str("go_version", buildInfo.GoVersion).
			Logger()
	})

	return log
}
