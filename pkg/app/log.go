package app

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/deckhouse/deckhouse/pkg/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Use info level with timestamps and a text output by default
var (
	L                = log.NewNop()
	LogLevel         = "info"
	LogNoTime        = false
	LogType          = "text"
	LogProxyHookJSON = false
)

// ForcedDurationForDebugLevel - force expiration for debug level.
const (
	ForcedDurationForDebugLevel = 30 * time.Minute
	ProxyJsonLogKey             = "proxyJsonLog"
)

func levelFromString(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// DefineLoggingFlags defines flags for logger settings.
func DefineLoggingFlags(cmd *kingpin.CmdClause) {
	cmd.Flag("log-level", "Logging level: debug, info, error. Default is info. Can be set with $LOG_LEVEL.").
		Envar("LOG_LEVEL").
		Default(LogLevel).
		StringVar(&LogLevel)
	cmd.Flag("log-type", "Logging formatter type: json, text or color. Default is text. Can be set with $LOG_TYPE.").
		Envar("LOG_TYPE").
		Default(LogType).
		StringVar(&LogType)
	cmd.Flag("log-no-time", "Disable timestamp logging if flag is present. Useful when output is redirected to logging system that already adds timestamps. Can be set with $LOG_NO_TIME.").
		Envar("LOG_NO_TIME").
		BoolVar(&LogNoTime)
	cmd.Flag("log-proxy-hook-json", "Delegate hook stdout/ stderr JSON logging to the hooks and act as a proxy that adds some extra fields before just printing the output").
		Envar("LOG_PROXY_HOOK_JSON").
		BoolVar(&LogProxyHookJSON)
}

// SetupLogging sets logger formatter and level.
func SetupLogging() {
	opts := []log.Option{
		log.WithOutput(os.Stderr),
		log.WithLevel(levelFromString(LogLevel)),
	}

	switch LogType {
	case "json":
		opts = append(opts, log.WithHandlerType(log.JSONHandlerType))
	case "color", "text":
		opts = append(opts, log.WithHandlerType(log.TextHandlerType))
	default:
		opts = append(opts, log.WithHandlerType(log.TextHandlerType))
	}

	L = log.NewLogger(opts...)
}
