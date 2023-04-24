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
	expectedIsJson        bool
	expectedLevel         string
	expectedLocation      string
	expectedMessage       string
	expectedMessageText   string
	expectedStatus        string
	expectedText          string
	expectedTime          time.Time
}{
	{
		name:            "parser-0001",
		message:         "",
		expectedDetails: map[string]string{},
		expectedErrors:  []string{},
		expectedTime:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		name:                "parser-0002",
		message:             "{}",
		expectedDetails:     map[string]string{},
		expectedErrors:      []string{},
		expectedIsJson:      true,
		expectedMessage:     "{}",
		expectedMessageText: "{}",
		expectedTime:        time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		name:                "parser-0003",
		message:             "{Not really JSON}",
		expectedDetails:     map[string]string{},
		expectedErrors:      []string{},
		expectedMessage:     "{Not really JSON}",
		expectedMessageText: "{Not really JSON}",
		expectedTime:        time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		name:                "parser-0003",
		message:             `{"text":"Bob works with Jane", But not really JSON}`,
		expectedDetails:     map[string]string{},
		expectedErrors:      []string{},
		expectedMessage:     `{"text":"Bob works with Jane", But not really JSON}`,
		expectedMessageText: `{"text":"Bob works with Jane", But not really JSON}`,
		expectedTime:        time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		name:                  "parser-0010",
		message:               `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`,
		expectedDetails:       map[string]string{"1": "Bob", "2": "Jane"},
		expectedDetailsNumber: 2,
		expectedDuration:      int64(1234),
		expectedErrors:        []string{"error1", "error2"},
		expectedErrorsNumber:  2,
		expectedId:            "senzing-99990001",
		expectedIsJson:        true,
		expectedLevel:         "TRACE",
		expectedLocation:      "In func1() at messenger_test.go:173",
		expectedMessage:       `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`,
		expectedMessageText:   "Bob works with Jane",
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

// -- Test Parse() method -----------------------------------------------------

func TestParserImpl_Parse(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		test.Run(testCase.name, func(test *testing.T) {
			testObject := Parse(testCase.message)
			assert.Equal(test, testCase.expectedDetails, testObject.GetDetails(), testCase.name+"-GetMessage()")
			assert.Equal(test, testCase.expectedDetailsNumber, len(testObject.GetDetails()), testCase.name+"-GetDetails()")
			assert.Equal(test, testCase.expectedDuration, testObject.GetDuration(), testCase.name+"-GetDuration()")
			assert.Equal(test, testCase.expectedErrors, testObject.GetErrors(), testCase.name+"-GetErrors()")
			assert.Equal(test, testCase.expectedErrorsNumber, len(testObject.GetErrors()), testCase.name+"-GetErrors()")
			assert.Equal(test, testCase.expectedId, testObject.GetId(), testCase.name+"-GetId()")
			assert.Equal(test, testCase.expectedIsJson, testObject.IsJson(), testCase.name+"-IsJson()")
			assert.Equal(test, testCase.expectedLevel, testObject.GetLevel(), testCase.name+"-GetLevel()")
			assert.Equal(test, testCase.expectedLocation, testObject.GetLocation(), testCase.name+"-GetLocation()")
			assert.Equal(test, testCase.expectedMessage, testObject.GetMessage(), testCase.name+"-GetMessage()")
			assert.Equal(test, testCase.expectedMessageText, testObject.GetMessageText(), testCase.name+"-GetMessageText()")
			assert.Equal(test, testCase.expectedStatus, testObject.GetStatus(), testCase.name+"-GetStatus()")
			assert.Equal(test, testCase.expectedText, testObject.GetText(), testCase.name+"-GetText()")
			assert.Equal(test, testCase.expectedTime, testObject.GetTime(), testCase.name+"-GetTime()")
		})
	}
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleParserImpl_GetDetails() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	fmt.Print(Parse(exampleMesssage).GetDetails())
	//Output: map[1:Bob 2:Jane]
}

func ExampleParserImpl_GetDuration() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	fmt.Print(Parse(exampleMesssage).GetDuration())
	//Output: 1234
}

func ExampleParserImpl_GetErrors() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	fmt.Print(Parse(exampleMesssage).GetErrors())
	//Output: [error1 error2]
}

func ExampleParserImpl_GetId() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	fmt.Print(Parse(exampleMesssage).GetId())
	//Output: senzing-99990001
}

func ExampleParserImpl_GetLevel() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	fmt.Print(Parse(exampleMesssage).GetLevel())
	//Output: TRACE
}

func ExampleParserImpl_GetLocation() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	fmt.Print(Parse(exampleMesssage).GetLocation())
	//Output: In func1() at messenger_test.go:173
}

func ExampleParserImpl_GetStatus() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	fmt.Print(Parse(exampleMesssage).GetStatus())
	//Output: OK
}

func ExampleParserImpl_GetText() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	fmt.Print(Parse(exampleMesssage).GetText())
	//Output: Bob works with Jane
}

func ExampleParserImpl_GetTime() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"OK","duration":1234,"errors":["error1", "error2"],"details":{"1":"Bob","2":"Jane"}}`
	fmt.Print(Parse(exampleMesssage).GetTime().Format(time.RFC3339))
	//Output: 2000-01-01T00:00:00Z
}
