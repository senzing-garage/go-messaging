package parser

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testCasesForMessage = []struct {
	name                  string
	message               string
	expectedDetails       map[string]string
	expectedDetailsNumber int
	expectedDuration      int64
	expectedErrors        []string
	expectedErrorsNumber  int
	expectedId            string
	expectedLevel         string
	expectedLocation      string
	expectedStatus        string
	expectedText          string
	expectedTime          time.Time
}{
	{
		name:                  "parser-0001",
		message:               `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`,
		expectedDetails:       map[string]string{"1": "Bob", "2": "Jane"},
		expectedDetailsNumber: 2,
		expectedDuration:      int64(1234),
		expectedErrors:        []string{"error1", "error2"},
		expectedErrorsNumber:  2,
		expectedId:            "senzing-99990001",
		expectedLevel:         "TRACE",
		expectedLocation:      "In func1() at messenger_test.go:173",
		expectedStatus:        "OK",
		expectedText:          "Bob works with Jane",
		expectedTime:          time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	},
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error = nil
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Internal functions - names begin with lowercase letter
// ----------------------------------------------------------------------------

func testError(test *testing.T, testObject ParserInterface, err error) {
	if err != nil {
		assert.Fail(test, err.Error())
	}
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

// -- Test New() method ---------------------------------------------------------

func TestParserImpl_Parse(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		test.Run(testCase.name, func(test *testing.T) {
			testObject, err := Parse(testCase.message)
			testError(test, testObject, err)
			if testCase.expectedDetails != nil {
				assert.Equal(test, testCase.expectedDetails, testObject.GetDetails(), testCase.name)
				assert.Equal(test, testCase.expectedDetailsNumber, len(testObject.GetDetails()), testCase.name)
			}
			if testCase.expectedDuration > 0 {
				assert.Equal(test, testCase.expectedDuration, testObject.GetDuration(), testCase.name)
			}
			if testCase.expectedErrors != nil {
				assert.Equal(test, testCase.expectedErrors, testObject.GetErrors(), testCase.name)
				assert.Equal(test, testCase.expectedErrorsNumber, len(testObject.GetErrors()), testCase.name)
			}
			if testCase.expectedId != "" {
				assert.Equal(test, testCase.expectedId, testObject.GetId(), testCase.name)
			}
			if testCase.expectedLevel != "" {
				assert.Equal(test, testCase.expectedLevel, testObject.GetLevel(), testCase.name)
			}
			if testCase.expectedLocation != "" {
				assert.Equal(test, testCase.expectedLocation, testObject.GetLocation(), testCase.name)
			}
			if testCase.expectedStatus != "" {
				assert.Equal(test, testCase.expectedStatus, testObject.GetStatus(), testCase.name)
			}
			if testCase.expectedText != "" {
				assert.Equal(test, testCase.expectedText, testObject.GetText(), testCase.name)
			}
			if !testCase.expectedTime.IsZero() {
				assert.Equal(test, testCase.expectedTime, testObject.GetTime(), testCase.name)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleParserImpl_GetDetails() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetDetails())
	//Output: map[1:Bob 2:Jane]
}

func ExampleParserImpl_GetDuration() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetDuration())
	//Output: 1234
}

func ExampleParserImpl_GetErrors() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetErrors())
	//Output: [error1 error2]
}

func ExampleParserImpl_GetId() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetId())
	//Output: senzing-99990001
}

func ExampleParserImpl_GetLevel() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetLevel())
	//Output: TRACE
}

func ExampleParserImpl_GetLocation() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetLocation())
	//Output: In func1() at messenger_test.go:173
}

func ExampleParserImpl_GetStatus() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetStatus())
	//Output: OK
}

func ExampleParserImpl_GetText() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetText())
	//Output: Bob works with Jane
}

func ExampleParserImpl_GetTime() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetTime().Format(time.RFC3339))
	//Output: 2000-01-01T00:00:00Z
}
