package parser

import (
	"encoding/json"

	"github.com/senzing/go-messaging/go/typedef"
)

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

/*
The Parse function creates a new instance of ParserInterface and an error.
*/
func Parse(message string) (*typedef.SenzingMessage, error) {
	result := &typedef.SenzingMessage{}
	err := json.Unmarshal([]byte(message), result)
	return result, err
}
