package parser

import (
	"encoding/json"
	"fmt"

	"github.com/senzing-garage/go-messaging/go/typedef"
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
	if err != nil {
		err = fmt.Errorf("parser.Parse error: %w", err)
	}

	return result, err
}
