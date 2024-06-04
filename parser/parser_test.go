package parser

import (
	"fmt"
	"testing"
	"time"

	"github.com/senzing-garage/go-messaging/go/typedef"
	"github.com/stretchr/testify/assert"
)

var testCasesForMessage = []struct {
	name                  string
	message               string
	expectedDetails       typedef.Details
	expectedDetailsNumber int
	expectedDuration      int64
	expectedError         string
	expectedErrors        []string
	expectedErrorsNumber  int
	expectedID            string
	expectedLevel         string
	expectedLocation      string
	expectedStatus        string
	expectedText          string
	expectedTime          time.Time
}{
	{
		name:          "parser-0001",
		message:       "",
		expectedError: "unexpected end of JSON input",
	},
	{
		name:    "parser-0002",
		message: "{}",
	},
	{
		name:          "parser-0003",
		message:       "{Not really JSON}",
		expectedError: "invalid character 'N' looking for beginning of object key string",
	},
	{
		name:          "parser-0004",
		message:       `{"text":"Bob works with Jane", But not really JSON}`,
		expectedError: "invalid character 'B' looking for beginning of object key string",
	},
	{
		name:                  "parser-0010",
		message:               `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"SZSDK99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedDetails:       typedef.Details{{Position: 1, Value: "Bob"}, {Position: 2, Value: "Jane"}},
		expectedDetailsNumber: 2,
		expectedDuration:      int64(1234),
		expectedErrors:        []string{"error1", "error2"},
		expectedErrorsNumber:  2,
		expectedID:            "SZSDK99990001",
		expectedLevel:         "TRACE",
		expectedLocation:      "In func1() at messenger_test.go:173",
		expectedStatus:        "OK",
		expectedText:          "Bob works with Jane",
		expectedTime:          time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	},
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestParse(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		test.Run(testCase.name, func(test *testing.T) {
			parsedMessage, err := Parse(testCase.message)
			if err != nil {
				assert.Equal(test, testCase.expectedError, err.Error(), testCase.name+"-ExpectedError")
			}
			assert.Equal(test, testCase.expectedDetails, parsedMessage.Details, testCase.name+"-Details")
			assert.Len(test, parsedMessage.Details, testCase.expectedDetailsNumber, testCase.name+"-DetailsNum")
			assert.Equal(test, testCase.expectedDuration, parsedMessage.Duration, testCase.name+"-Duration")
			assert.Equal(test, testCase.expectedErrors, parsedMessage.Errors, testCase.name+"-Errors")
			assert.Len(test, parsedMessage.Errors, testCase.expectedErrorsNumber, testCase.name+"-ErrorsNum")
			assert.Equal(test, testCase.expectedID, parsedMessage.ID, testCase.name+"-ID")
			assert.Equal(test, testCase.expectedLevel, parsedMessage.Level, testCase.name+"-Level")
			assert.Equal(test, testCase.expectedLocation, parsedMessage.Location, testCase.name+"-Location")
			assert.Equal(test, testCase.expectedStatus, parsedMessage.Status, testCase.name+"-Status")
			assert.Equal(test, testCase.expectedText, parsedMessage.Text, testCase.name+"-Text")
			assert.Equal(test, testCase.expectedTime, parsedMessage.Time, testCase.name+"-Time")
		})
	}
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleParse_details() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/parser/parser_test.go
	message := `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"SZSDK99990001","text":"Bob works with Jane","status":"OK","duration":1234,"errors":["error1","error2"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`
	parsedMessage, err := Parse(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedMessage.Details)
	//Output: [{ 1  Bob <nil>} { 2  Jane <nil>}]
}

func ExampleParse_duration() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/parser/parser_test.go
	message := `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"SZSDK99990001","text":"Bob works with Jane","status":"OK","duration":1234,"errors":["error1","error2"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`
	parsedMessage, err := Parse(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedMessage.Duration)
	//Output: 1234
}

func ExampleParse_errors() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/parser/parser_test.go
	message := `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"SZSDK99990001","text":"Bob works with Jane","status":"OK","duration":1234,"errors":["error1","error2"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`
	parsedMessage, err := Parse(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedMessage.Errors)
	//Output: [error1 error2]
}

func ExampleParse_id() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/parser/parser_test.go
	message := `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"SZSDK99990001","text":"Bob works with Jane","status":"OK","duration":1234,"errors":["error1","error2"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`
	parsedMessage, err := Parse(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedMessage.ID)
	//Output: SZSDK99990001
}

func ExampleParse_level() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/parser/parser_test.go
	message := `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"SZSDK99990001","text":"Bob works with Jane","status":"OK","duration":1234,"errors":["error1","error2"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`
	parsedMessage, err := Parse(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedMessage.Level)
	//Output: TRACE
}

func ExampleParse_status() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/parser/parser_test.go
	message := `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"SZSDK99990001","text":"Bob works with Jane","status":"OK","duration":1234,"errors":["error1","error2"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`
	parsedMessage, err := Parse(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedMessage.Status)
	//Output: OK
}

func ExampleParse_text() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/parser/parser_test.go
	message := `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"SZSDK99990001","text":"Bob works with Jane","status":"OK","duration":1234,"errors":["error1","error2"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`
	parsedMessage, err := Parse(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedMessage.Text)
	//Output: Bob works with Jane
}

func ExampleParse_time() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/parser/parser_test.go
	message := `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"SZSDK99990001","text":"Bob works with Jane","status":"OK","duration":1234,"errors":["error1","error2"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`
	parsedMessage, err := Parse(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedMessage.Time.Format(time.RFC3339Nano))
	//Output: 2000-01-01T00:00:00Z
}
