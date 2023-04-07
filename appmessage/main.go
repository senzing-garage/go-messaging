package appmessage

import (
	"errors"
	"fmt"
	"time"
)

// ----------------------------------------------------------------------------
// Types - interface
// ----------------------------------------------------------------------------

// The AppMessageInterface interface is...
type AppMessageInterface interface {
	NewJson(messageNumber int, details ...interface{}) string
	NewSlog(messageNumber int, details ...interface{}) (string, []interface{})
}

// ----------------------------------------------------------------------------
// Types - struct
// ----------------------------------------------------------------------------

type AppMessageCallerSkip struct {
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
	Date     string      `json:"date,omitempty"`     // Date of message in UTC.
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

// An example constant.
const ExampleConstant = 1

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

	var callerSkip int = 0

	for _, value := range options {
		switch typedValue := value.(type) {
		case *AppMessageCallerSkip:
			callerSkip = typedValue.Value
		}
	}

	result = &AppMessageImpl{
		idMessages:        idMessages,
		idStatuses:        idStatuses,
		messageIdTemplate: fmt.Sprintf("senzing-%04d", productIdentifier) + "%04d",
		callerSkip:        callerSkip,
	}

	return result, err
}
