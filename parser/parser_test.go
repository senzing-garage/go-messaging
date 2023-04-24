package parser

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testCasesForMessage = []struct {
	name             string
	message          string
	expectedDetails  map[string]interface{}
	expectedDuration int64
	expectedErrors   string
	expectedId       string
	expectedLevel    string
	expectedLocation string
	expectedStatus   string
	expectedText     string
	expectedTime     time.Time
}{
	{
		name:             "parser-0001",
		message:          `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"TRACE: Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"none","duration":1234,"details":{"1":"Bob","2":"Jane"}}`,
		expectedDetails:  map[string]interface{}{"1": "Bob", "2": "Jane"},
		expectedDuration: int64(1234),
		expectedId:       "senzing-99990001",
		expectedLevel:    "TRACE",
		expectedLocation: "In func1() at messenger_test.go:173",
		expectedStatus:   "none",
		expectedText:     "TRACE: Bob works with Jane",
		expectedTime:     time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
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

func TestMessengerImpl_NewJson(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		test.Run(testCase.name, func(test *testing.T) {
			testObject, err := Parse(testCase.message)
			testError(test, testObject, err)
			if testCase.expectedDetails != nil {
				assert.Equal(test, testCase.expectedDetails, testObject.GetDetails(), testCase.name)
			}
			if testCase.expectedDuration > 0 {
				assert.Equal(test, testCase.expectedDuration, testObject.GetDuration(), testCase.name)
			}
			if testCase.expectedErrors != "" {
				assert.Equal(test, testCase.expectedErrors, testObject.GetErrors(), testCase.name)
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

func ExampleParserImpl_GetId() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/parser/parser_test.go
	exampleMesssage := `{"time":"2000-01-01T00:00:00.00000000Z","level":"TRACE","id":"senzing-99990001","text":"TRACE: Bob works with Jane","location":"In func1() at messenger_test.go:173","status":"none","duration":1234,"details":{"1":"Bob","2":"Jane"}}`
	parsedMessage, err := Parse(exampleMesssage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(parsedMessage.GetId())
	//Output: senzing-99990001
}
