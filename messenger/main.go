package messenger

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slog"
)

// ----------------------------------------------------------------------------
// Types - interface
// ----------------------------------------------------------------------------

// The MessengerInterface interface has methods for creating different
// representations of a message.
type MessengerInterface interface {
	NewJson(messageNumber int, details ...interface{}) string
	NewSlog(messageNumber int, details ...interface{}) (string, []interface{})
	NewSlogLevel(messageNumber int, details ...interface{}) (string, slog.Level, []interface{})
}

// ----------------------------------------------------------------------------
// Types - struct
// ----------------------------------------------------------------------------

// Fields in the formatted message.
// Order is important.
// It should be time, level, id, text, status, duration, location, errors, details.
type MessageFormat struct {
	Time     string      `json:"time,omitempty"`     // Time of message in UTC.
	Level    string      `json:"level,omitempty"`    // Level:  TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC.
	Id       string      `json:"id,omitempty"`       // Message identifier.
	Text     string      `json:"text,omitempty"`     // Message text.
	Status   string      `json:"status,omitempty"`   // Status information.
	Duration int64       `json:"duration,omitempty"` // Duration in nanoseconds
	Location string      `json:"location,omitempty"` // Location in the code issuing message.
	Errors   interface{} `json:"errors,omitempty"`   // List of errors.
	Details  []Detail    `json:"details,omitempty"`  // All instances passed into the message.
}

type Detail struct {
	Key      string      `json:"key,omitempty"`
	Position int32       `json:"position,omitempty"`
	Type     string      `json:"type,omitempty"`
	Value    string      `json:"value,omitempty"`
	ValueRaw interface{} `json:"valueRaw,omitempty"`
}

// --- Override values when creating messages ---------------------------------

// Value of the "details" field.
type MessageDetails struct {
	Value interface{} // All instances passed into the message.
}

// Value of the "duration" field.
type MessageDuration struct {
	Value int64 // Duration in nanoseconds
}

// Value of the "id" field.
type MessageId struct {
	Value string // Message identifier.
}

// Value of the "level" field.
type MessageLevel struct {
	Value string // Level:  TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC.
}

// Value of the "location" field.
type MessageLocation struct {
	Value string // Location in the code issuing message.
}

// Value of the "status" field.
type MessageStatus struct {
	Value string // Status information.
}

// Value of the "text" field.
type MessageText struct {
	Value string // Message text.
}

// Value of the "time" field.
type MessageTime struct {
	Value time.Time // Time of message in UTC.
}

// --- Options for New() ------------------------------------------------------

// Number of callers to skip when determining location.
type OptionCallerSkip struct {
	Value int // Number of callers to skip in the stack trace when determining the location.
}

// Map of message number to message templates.
type OptionIdMessages struct {
	Value map[int]string // Message number to message template map.
}

// Map of message number to status values.
type OptionIdStatuses struct {
	Value map[int]string // Message number to status map
}

// Format of the unique id.
type OptionMessageIdTemplate struct {
	Value string // Format string.
}

// The component identifier.
// See https://github.com/Senzing/knowledge-base/blob/main/lists/senzing-product-ids.md
type OptionSenzingComponentId struct {
	Value int // Component issuing message.
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Log levels as integers.
// Compatible with golang.org/x/exp/slog.
const (
	LevelTraceInt int = -8
	LevelDebugInt int = -4
	LevelInfoInt  int = 0
	LevelWarnInt  int = 4
	LevelErrorInt int = 8
	LevelFatalInt int = 12
	LevelPanicInt int = 16
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

// Existing and new log levels used with slog.Level.
const (
	LevelDebugSlog = slog.LevelDebug
	LevelErrorSlog = slog.LevelError
	LevelFatalSlog = slog.Level(LevelFatalInt)
	LevelInfoSlog  = slog.LevelInfo
	LevelPanicSlog = slog.Level(LevelPanicInt)
	LevelTraceSlog = slog.Level(LevelTraceInt)
	LevelWarnSlog  = slog.LevelWarn
)

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message ID Low-bound for message levels
// i.e. a message in range 0 - 999 is a TRACE message.
var IdLevelRangesAsString = map[int]string{
	0000: LevelTraceName,
	1000: LevelDebugName,
	2000: LevelInfoName,
	3000: LevelWarnName,
	4000: LevelErrorName,
	5000: LevelFatalName,
	6000: LevelPanicName,
}

// Map from slog.Level to string representation.
var LevelToTextMap = map[slog.Level]string{
	LevelDebugSlog: LevelDebugName,
	LevelErrorSlog: LevelErrorName,
	LevelFatalSlog: LevelFatalName,
	LevelInfoSlog:  LevelInfoName,
	LevelPanicSlog: LevelPanicName,
	LevelTraceSlog: LevelTraceName,
	LevelWarnSlog:  LevelWarnName,
}

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

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

/*
The New function creates a new instance of MessengerInterface.
Adding options can be used to modify subcomponents.
*/
func New(options ...interface{}) (MessengerInterface, error) {
	var err error = nil
	var result MessengerInterface = nil

	// Default values.

	var (
		callerSkip          int            = 0
		idMessages          map[int]string = map[int]string{}
		idStatuses          map[int]string = map[int]string{}
		componentIdentifier int            = 9999
		messageIdTemplate   string         = fmt.Sprintf("senzing-%04d", componentIdentifier) + "%04d"
	)

	// Process options.

	for _, value := range options {
		switch typedValue := value.(type) {
		case *OptionCallerSkip:
			callerSkip = typedValue.Value
		case *OptionIdMessages:
			idMessages = typedValue.Value
		case *OptionIdStatuses:
			idStatuses = typedValue.Value
		case *OptionSenzingComponentId:
			componentIdentifier = typedValue.Value
			messageIdTemplate = fmt.Sprintf("senzing-%04d", componentIdentifier) + "%04d"
		case *OptionMessageIdTemplate:
			messageIdTemplate = typedValue.Value
		}
	}

	// Detect incorrect option values.

	if componentIdentifier <= 0 || componentIdentifier >= 10000 {
		err := errors.New("componentIdentifier must be in range 1..9999. See https://github.com/Senzing/knowledge-base/blob/main/lists/senzing-product-ids.md")
		return result, err
	}

	if idMessages == nil {
		err := errors.New("messages must be a map[int]string")
		return result, err
	}

	if idStatuses == nil {
		err := errors.New("statuses must be a map[int]string")
		return result, err
	}

	// Create MessengerInterface.

	result = &MessengerImpl{
		callerSkip:        callerSkip,
		idMessages:        idMessages,
		idStatuses:        idStatuses,
		messageIdTemplate: messageIdTemplate,
	}
	return result, err
}
