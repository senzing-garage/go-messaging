package appmessage

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slog"
)

// ----------------------------------------------------------------------------
// Types - interface
// ----------------------------------------------------------------------------

// The AppMessageInterface interface is...
type AppMessageInterface interface {
	NewJson(messageNumber int, details ...interface{}) string
	NewSlog(messageNumber int, details ...interface{}) (string, []interface{})
	NewSlogLevel(messageNumber int, details ...interface{}) (string, slog.Level, []interface{})
}

// ----------------------------------------------------------------------------
// Types - struct
// ----------------------------------------------------------------------------

type AppMessageOptionCallerSkip struct {
	Value int
}

type AppMessageDetails struct {
	Value interface{}
}

type AppMessageDuration struct {
	Value int64
}

// Fields in the formatted message.
// Order is important.
// It should be date, time, level, id, status, text, duration, location, errors, details.
type AppMessageFormat struct {
	Time     string      `json:"time,omitempty"`     // Time of message in UTC.
	Level    string      `json:"level,omitempty"`    // Level:  TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC.
	Id       string      `json:"id,omitempty"`       // Message identifier.
	Text     interface{} `json:"text,omitempty"`     // Message text.
	Status   string      `json:"status,omitempty"`   // Status information.
	Duration int64       `json:"duration,omitempty"` // Duration in nanoseconds
	Location string      `json:"location,omitempty"` // Location in the code issuing message.
	Errors   interface{} `json:"errors,omitempty"`   // List of errors.
	Details  interface{} `json:"details,omitempty"`  // All instances passed into the message.
}

type AppMessageId struct {
	Value string
}

type AppMessageLevel struct {
	Value string
}

type AppMessageLocation struct {
	Value string
}

type AppMessageStatus struct {
	Value string
}

type AppMessageText struct {
	Value interface{}
}

type AppMessageTimestamp struct {
	Value time.Time
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	LevelTraceInt int = -8
	LevelDebugInt int = -4
	LevelInfoInt  int = 0
	LevelWarnInt  int = 4
	LevelErrorInt int = 8
	LevelFatalInt int = 12
	LevelPanicInt int = 16
)

const (
	LevelDebugSlog = slog.LevelDebug
	LevelErrorSlog = slog.LevelError
	LevelFatalSlog = slog.Level(LevelFatalInt)
	LevelInfoSlog  = slog.LevelInfo
	LevelPanicSlog = slog.Level(LevelPanicInt)
	LevelTraceSlog = slog.Level(LevelTraceInt)
	LevelWarnSlog  = slog.LevelWarn
)

// Strings representing the supported logging levels.
const (
	LevelDebugName = "DEBUG"
	LevelErrorName = "ERROR"
	LevelFatalName = "FATAL"
	LevelInfoName  = "INFO"
	LevelPanicName = "PANIC"
	LevelTraceName = "TRACE"
	LevelWarnName  = "WARN"
)

// Map from string representation to Log level as typed integer.
var TextToLevelMap = map[string]slog.Level{
	LevelDebugName: LevelDebugSlog,
	LevelErrorName: LevelErrorSlog,
	LevelFatalName: LevelFatalSlog,
	LevelInfoName:  LevelInfoSlog,
	LevelPanicName: LevelPanicSlog,
	LevelTraceName: LevelTraceSlog,
	LevelWarnName:  LevelWarnSlog,
}

var LevelToTextMap = map[slog.Level]string{
	LevelDebugSlog: LevelDebugName,
	LevelErrorSlog: LevelErrorName,
	LevelFatalSlog: LevelFatalName,
	LevelInfoSlog:  LevelInfoName,
	LevelPanicSlog: LevelPanicName,
	LevelTraceSlog: LevelTraceName,
	LevelWarnSlog:  LevelWarnName,
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// An example variable.
var ExampleVariable = map[int]string{
	1: "Just a string",
}

var IdLevelRangesAsString = map[int]string{
	0000: "TRACE",
	1000: "DEBUG",
	2000: "INFO",
	3000: "WARN",
	4000: "ERROR",
	5000: "FATAL",
	6000: "PANIC",
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

/*
The NewSenzingLogger function creates a new instance of MessageLoggerInterface
that is tailored to Senzing applications.
Like New(), adding parameters can be used to modify subcomponents.
*/
func New(productIdentifier int, idMessages map[int]string, idStatuses map[int]string, options ...interface{}) (AppMessageInterface, error) {
	var err error = nil
	var result AppMessageInterface = nil

	// Detect incorrect parameter values.

	if productIdentifier <= 0 || productIdentifier >= 10000 {
		err := errors.New("productIdentifier must be in range 1..9999. See https://github.com/Senzing/knowledge-base/blob/main/lists/senzing-product-ids.md")
		return result, err
	}

	if idMessages == nil {
		err := errors.New("messages must be a map[int]string")
		return result, err
	}

	// Process options

	var callerSkip int = 0
	for _, value := range options {
		switch typedValue := value.(type) {
		case *AppMessageOptionCallerSkip:
			callerSkip = typedValue.Value
		}
	}

	// Create AppMessage

	result = &AppMessageImpl{
		idMessages:        idMessages,
		idStatuses:        idStatuses,
		messageIdTemplate: fmt.Sprintf("senzing-%04d", productIdentifier) + "%04d",
		callerSkip:        callerSkip,
	}

	return result, err
}

/*
The HandlerOptions function returns a slog handler that includes TRACE, FATAL, and PANIC.
See: https://go.googlesource.com/exp/+/refs/heads/master/slog/example_custom_levels_test.go
*/
func HandlerOptions(leveler slog.Leveler) *slog.HandlerOptions {
	if leveler == nil {
		leveler = LevelInfoSlog
	}
	handlerOptions := &slog.HandlerOptions{
		Level: leveler,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "text"
			}
			if a.Key == slog.LevelKey {
				level := ""
				switch typedValue := a.Value.Any().(type) {
				case string:
					level = typedValue
				case slog.Level:
					level = typedValue.String()
				}
				switch {
				case level == "DEBUG-4":
					a.Value = slog.StringValue("TRACE")
				case level == "ERROR+4":
					a.Value = slog.StringValue("FATAL")
				case level == "ERROR+8":
					a.Value = slog.StringValue("PANIC")
				}
			}
			return a
		},
	}
	return handlerOptions
}
