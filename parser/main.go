package parser

import (
	"time"
)

// ----------------------------------------------------------------------------
// Types - interface
// ----------------------------------------------------------------------------

// The MessengerInterface interface has methods for creating different
// representations of a message.
type ParserInterface interface {
	GetDetails() interface{}
	GetDuration() int64
	GetErrors() interface{}
	GetId() string
	GetLevel() string
	GetLocation() string
	GetStatus() string
	GetText() interface{}
	GetTime() time.Time
}

// ----------------------------------------------------------------------------
// Types - struct
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

/*
The New function creates a new instance of MessengerInterface.
Adding options can be used to modify subcomponents.
*/
func Parse(message string) (ParserInterface, error) {
	var result = &ParserImpl{
		message: message,
	}
	err := result.initialize()
	return result, err
}
