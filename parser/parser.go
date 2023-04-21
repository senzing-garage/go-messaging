package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing/go-messaging/messenger"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ParserImpl is an type-struct for an implementation of the ParserInterface.
type ParserImpl struct {
	message       string
	parsedMessage messenger.MessageFormat
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Determine if string is syntactically JSON.
func isJson(unknownString string) bool {
	unknownStringUnescaped, err := strconv.Unquote(unknownString)
	if err != nil {
		unknownStringUnescaped = unknownString
	}
	var jsonString json.RawMessage
	return json.Unmarshal([]byte(unknownStringUnescaped), &jsonString) == nil
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (parser *ParserImpl) initialize() error {
	if !isJson(parser.message) {
		return fmt.Errorf("string is not JSON")
	}
	err := json.Unmarshal([]byte(parser.message), &parser.parsedMessage)
	return err
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetDetails() interface{} {
	return nil
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetDuration() int64 {
	return 0
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetErrors() interface{} {
	return nil
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetId() string {
	return parser.GetId()
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetLevel() string {
	return ""
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetLocation() string {
	return ""
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetStatus() string {
	return ""
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetText() interface{} {
	return nil
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetTime() time.Time {
	return time.Now()
}
