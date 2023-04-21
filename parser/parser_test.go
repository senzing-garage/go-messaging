package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var idMessages = map[int]string{
	0001: "TRACE: %s works with %s",
	1001: "DEBUG: %s works with %s",
	2001: "INFO: %s works with %s",
	3001: "WARN: %s works with %s",
	4001: "ERROR: %s works with %s",
	5001: "FATAL: %s works with %s",
	6001: "PANIC: %s works with %s",
}

var testCasesForMessage = []struct {
	name       string
	message    string
	expectedId string
}{
	{
		name:       "parser-0001",
		message:    `{"id", "senzing-99990001", "location", "In func1() at messenger_test.go:186", "details", map[string]interface{}{"1": "Bob", "2": "Jane"}}`,
		expectedId: "senzing-99990001",
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
			assert.Equal(test, testCase.expectedId, testObject.GetId(), testCase.name)
		})
	}
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

// func ExampleMessengerImpl_NewJson() {
// 	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/messenger/messenger_test.go
// 	example, err := New()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Print(example.NewJson(2001, "Bob", "Jane", getTimestamp(), getOptionCallerSkip()))
// 	//Output: {"time":"2000-01-01 00:00:00 +0000 UTC","level":"INFO","id":"senzing-99992001","location":"In ExampleMessengerImpl_NewJson() at messenger_test.go:205","details":{"1":"Bob","2":"Jane"}}
// }
