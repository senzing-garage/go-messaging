package messenger

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var idMessages = map[int]string{
	0001: "TRACE: %s knows %s",
	1001: "DEBUG: %s knows %s",
	2001: "INFO: %s knows %s",
	3001: "WARN: %s knows %s",
	4001: "ERROR: %s knows %s",
	5001: "FATAL: %s knows %s",
	6001: "PANIC: %s knows %s",
}

var testCasesForMessage = []struct {
	name                string
	messageNumber       int
	options             []interface{}
	details             []interface{}
	expectedMessageJson string
	expectedMessageSlog []interface{}
	expectedText        string
}{
	{
		name:                "messenger-0001",
		messageNumber:       1,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"TRACE","id":"senzing-99990001","text":"TRACE: Bob knows Jane","location":"In func1() at messenger_test.go:164","details":{"1":"Bob","2":"Jane"}}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99990001", "location", "In func1() at messenger_test.go:177", "details", map[string]interface{}{"1": "Bob", "2": "Jane"}},
		expectedText:        "TRACE: Bob knows Jane",
	},
	{
		name:                "messenger-1001",
		messageNumber:       1001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"DEBUG","id":"senzing-99991001","text":"DEBUG: Bob knows Jane","location":"In func1() at messenger_test.go:164","details":{"1":"Bob","2":"Jane"}}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99991001", "location", "In func1() at messenger_test.go:177", "details", map[string]interface{}{"1": "Bob", "2": "Jane"}},
		expectedText:        "DEBUG: Bob knows Jane",
	},
	{
		name:                "messenger-2001",
		messageNumber:       2001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"INFO","id":"senzing-99992001","text":"INFO: Bob knows Jane","location":"In func1() at messenger_test.go:164","details":{"1":"Bob","2":"Jane"}}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99992001", "location", "In func1() at messenger_test.go:177", "details", map[string]interface{}{"1": "Bob", "2": "Jane"}},
		expectedText:        "INFO: Bob knows Jane",
	},
	{
		name:                "messenger-3001",
		messageNumber:       3001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"WARN","id":"senzing-99993001","text":"WARN: Bob knows Jane","location":"In func1() at messenger_test.go:164","details":{"1":"Bob","2":"Jane"}}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99993001", "location", "In func1() at messenger_test.go:177", "details", map[string]interface{}{"1": "Bob", "2": "Jane"}},
		expectedText:        "WARN: Bob knows Jane",
	},
	{
		name:                "messenger-4001",
		messageNumber:       4001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"ERROR","id":"senzing-99994001","text":"ERROR: Bob knows Jane","location":"In func1() at messenger_test.go:164","details":{"1":"Bob","2":"Jane"}}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99994001", "location", "In func1() at messenger_test.go:177", "details", map[string]interface{}{"1": "Bob", "2": "Jane"}},
		expectedText:        "ERROR: Bob knows Jane",
	},
	{
		name:                "messenger-5001",
		messageNumber:       5001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"FATAL","id":"senzing-99995001","text":"FATAL: Bob knows Jane","location":"In func1() at messenger_test.go:164","details":{"1":"Bob","2":"Jane"}}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99995001", "location", "In func1() at messenger_test.go:177", "details", map[string]interface{}{"1": "Bob", "2": "Jane"}},
		expectedText:        "FATAL: Bob knows Jane",
	},
	{
		name:                "messenger-6001",
		messageNumber:       6001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"PANIC","id":"senzing-99996001","text":"PANIC: Bob knows Jane","location":"In func1() at messenger_test.go:164","details":{"1":"Bob","2":"Jane"}}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99996001", "location", "In func1() at messenger_test.go:177", "details", map[string]interface{}{"1": "Bob", "2": "Jane"}},
		expectedText:        "PANIC: Bob knows Jane",
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

func testError(test *testing.T, testObject MessengerInterface, err error) {
	if err != nil {
		assert.Fail(test, err.Error())
	}
}

func getOptionIdMessages() *OptionIdMessages {
	return &OptionIdMessages{
		Value: idMessages,
	}
}

func getOptionCallerSkip() *OptionCallerSkip {
	return &OptionCallerSkip{
		Value: 2,
	}
}

func getTimestamp() *MessageTime {
	return &MessageTime{
		Value: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

// -- Test New() method ---------------------------------------------------------

func TestMessengerImpl_NewJson(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessageJson) > 0 {
			test.Run(testCase.name+"-NewJson", func(test *testing.T) {
				testObject, err := New(testCase.options...)
				testError(test, testObject, err)
				actual := testObject.NewJson(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedMessageJson, actual, testCase.name)
			})
		}
	}
}

func TestMessengerImpl_NewSlog(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessageSlog) > 0 {
			test.Run(testCase.name+"-NewSlog", func(test *testing.T) {
				testObject, err := New(testCase.options...)
				testError(test, testObject, err)
				message, _, actual := testObject.NewSlogLevel(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedMessageSlog, actual, testCase.name)
				assert.Equal(test, testCase.expectedText, message, testCase.name)
			})
		}
	}
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleMessengerImpl_NewJson() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/messenger/messenger_test.go
	example, err := New()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(example.NewJson(2001, "Bob", "Jane", getTimestamp(), getOptionCallerSkip()))
	//Output: {"time":"2000-01-01 00:00:00 +0000 UTC","level":"INFO","id":"senzing-99992001","location":"In ExampleMessengerImpl_NewJson() at messenger_test.go:195","details":{"1":"Bob","2":"Jane"}}
}
