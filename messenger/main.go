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

// The Messenger interface has methods for creating different
// representations of a message.
type Messenger interface {
	NewError(messageNumber int, details ...interface{}) error
	NewJSON(messageNumber int, details ...interface{}) string
	NewSlog(messageNumber int, details ...interface{}) (string, []interface{})
	NewSlogLevel(messageNumber int, details ...interface{}) (string, slog.Level, []interface{})
}

// ----------------------------------------------------------------------------
// Types - struct
// ----------------------------------------------------------------------------

// Fields in the formatted message.
// Order is important.
// It should be time, level, id, text, code, reason, status, duration, location, errors, details.
type MessageFormat struct {
	Time     string      `json:"time,omitempty"`     // Time of message in UTC.
	Level    string      `json:"level,omitempty"`    // Level:  TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC.
	ID       string      `json:"id,omitempty"`       // Message identifier.
	Text     string      `json:"text,omitempty"`     // Message text.
	Code     string      `json:"code,omitempty"`     // Underlying reason code.
	Reason   string      `json:"reason,omitempty"`   // Underlying reason.
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

// Value of the "code" field.
type MessageCode struct {
	Value string // Underlying message code.
}

// Value of the "details" field.
type MessageDetails struct {
	Value interface{} // All instances passed into the message.
}

// Value of the "duration" field.
type MessageDuration struct {
	Value int64 // Duration in nanoseconds
}

// Value of the "id" field.
type MessageID struct {
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

// Value of the "reason" field.
type MessageReason struct {
	Value string // Underlying message reason.
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

// The component identifier.
// See https://github.com/senzing-garage/knowledge-base/blob/main/lists/senzing-product-ids.md
type OptionComponentID struct {
	Value int // Component issuing message.
}

// Map of message number to message templates.
type OptionIDMessages struct {
	Value map[int]string // Message number to message template map.
}

// Map of message number to status values.
type OptionIDStatuses struct {
	Value map[int]string // Message number to status map
}

// List of fields included in final message.
type OptionMessageField struct {
	Value string // One of AllMessageFields values.
}

// List of fields included in final message.
type OptionMessageFields struct {
	Value []string // One or more of AllMessageFields values.
}

// Format of the unique id.
type OptionMessageIDTemplate struct {
	Value string // Format string.
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
var IDLevelRangesAsString = map[int]string{
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

var (
	ErrBadComponentID = errors.New("componentIdentifier must be in range 1..9999. See https://github.com/senzing-garage/knowledge-base/blob/main/lists/senzing-product-ids.md")
	ErrEmptyMessages  = errors.New("messages must be a map[int]string")
	ErrEmptyStatuses  = errors.New("statuses must be a map[int]string")
)

// Order is important in AllMessageFields. Should match order in MessageFormat.
var AllMessageFields = []string{"time", "level", "id", "text", "code", "reason", "status", "duration", "location", "errors", "details"}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

/*
The New function creates a new instance of MessengerInterface.
Adding options can be used to modify subcomponents.
*/
func New(options ...interface{}) (Messenger, error) {

	var err error
	var result Messenger

	// Default values.

	var (
		callerSkip          int
		idMessages          = map[int]string{}
		idStatuses          = map[int]string{}
		componentIdentifier = 9999
		messageIDTemplate   = fmt.Sprintf("SZSDK%04d", componentIdentifier) + "%04d"
		messageFields       []string
	)

	// Process options.

	for _, value := range options {
		switch typedValue := value.(type) {
		case *OptionCallerSkip:
			callerSkip = typedValue.Value
		case *OptionComponentID:
			componentIdentifier = typedValue.Value
			messageIDTemplate = fmt.Sprintf("SZSDK%04d", componentIdentifier) + "%04d"
		case *OptionIDMessages:
			idMessages = typedValue.Value
		case *OptionIDStatuses:
			idStatuses = typedValue.Value
		case *OptionMessageFields:
			messageFields = typedValue.Value
		case *OptionMessageIDTemplate:
			messageIDTemplate = typedValue.Value
		}
	}

	// Detect incorrect option values.

	if componentIdentifier <= 0 || componentIdentifier >= 10000 {
		return result, ErrBadComponentID
	}

	if idMessages == nil {
		return result, ErrEmptyMessages
	}

	if idStatuses == nil {
		return result, ErrEmptyStatuses
	}

	// Create MessengerInterface.

	result = &BasicMessenger{
		callerSkip:        callerSkip,
		idMessages:        idMessages,
		idStatuses:        idStatuses,
		messageFields:     messageFields,
		messageIDTemplate: messageIDTemplate,
	}
	return result, err
}
