package appmessage

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var idMessages = map[int]string{
	2001: "%s knows %s",
	3001: "%s knows %s",
	4001: "%s knows %s",
	2:    "%s does not know %s",
}

var testCasesForMessage = []struct {
	name              string
	productIdentifier int
	idMessages        map[int]string
	idStatuses        map[int]string
	messageNumber     int
	details           []interface{}
	expectedMessage   string
}{
	{
		name:              "appmessage-1",
		productIdentifier: 9999,
		idMessages:        idMessages,
		messageNumber:     1,
		details:           []interface{}{"A", 1, getTimestamp()},
		expectedMessage:   `{"date":"2000-01-01","time":"00:00:00.000000000","level":"TRACE","id":"senzing-99990001","location":"In AFunction() at somewhere.go:1234","details":{"1":"A","2":1}}`,
	},
	{
		name:              "appmessage-2",
		productIdentifier: 9999,
		idMessages:        idMessages,
		messageNumber:     2,
		details:           []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessage:   `{"date":"2000-01-01","time":"00:00:00.000000000","id":"senzing-99990002","text":"Bob does not know Jane","details":["Bob","Jane"]}`,
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

func testError(test *testing.T, testObject AppMessageInterface, err error) {
	if err != nil {
		assert.Fail(test, err.Error())
	}
}

func getTimestamp() *AppMessageTimestamp {
	return &AppMessageTimestamp{
		Value: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

// -- Test New() method ---------------------------------------------------------

func TestAppMessageImpl_NewJson(test *testing.T) {
	callerSkip := &AppMessageCallerSkip{Value: 1}
	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessage) > 0 {
			test.Run(testCase.name+"-NewJson", func(test *testing.T) {
				testObject, err := New(testCase.productIdentifier, testCase.idMessages, testCase.idStatuses, callerSkip)
				testError(test, testObject, err)
				actual := testObject.NewJson(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedMessage, actual, testCase.name)
			})
		}
	}
}

func TestAppMessageImpl_NewSlog(test *testing.T) {
	callerSkip := &AppMessageCallerSkip{Value: 1}
	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessage) > 0 {
			test.Run(testCase.name+"-NewSlog", func(test *testing.T) {
				testObject, err := New(testCase.productIdentifier, testCase.idMessages, testCase.idStatuses, callerSkip)
				testError(test, testObject, err)
				message, actual := testObject.NewSlog(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedMessage, message, testCase.name)
				assert.Equal(test, testCase.expectedMessage, actual, testCase.name)
			})
		}
	}
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleAppMessageImpl_NewJson() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/appmessages/appmessage_test.go
	example := &AppMessageImpl{
		idMessages: idMessages,
	}
	fmt.Print(example.NewJson(2001, "Bob", "Jane"))
	//Output:
	//examplePackage: I'm here
}
