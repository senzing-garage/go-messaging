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
func (parser *ParserImpl) GetDetails() map[string]string {
	result := map[string]string{}
	if parser.parsedMessage.Details != nil {
		parsedDetails, ok := parser.parsedMessage.Details.(map[string]interface{})
		if ok {
			for key, value := range parsedDetails {
				result[key] = fmt.Sprint(value)
			}
		}
	}
	return result
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetDuration() int64 {
	return parser.parsedMessage.Duration
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetErrors() []string {
	result := []string{}
	if parser.parsedMessage.Errors != nil {
		parsedDetails, ok := parser.parsedMessage.Errors.([]interface{})
		if ok {
			for _, value := range parsedDetails {
				result = append(result, fmt.Sprint(value))
			}
		}
	}
	return result
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetId() string {
	return parser.parsedMessage.Id
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetLevel() string {
	return parser.parsedMessage.Level
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetLocation() string {
	return parser.parsedMessage.Location
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetStatus() string {
	return parser.parsedMessage.Status
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetText() string {
	return fmt.Sprint(parser.parsedMessage.Text)
}

/*
The GetXxxxxx method returns...

# Input

Output
*/
func (parser *ParserImpl) GetTime() time.Time {
	result, err := time.Parse(time.RFC3339Nano, parser.parsedMessage.Time)
	if err != nil {
		fmt.Println(err.Error())
		result = time.Time{}
	}
	return result
}
