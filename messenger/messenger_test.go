package messenger

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
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
	name                string
	messageNumber       int
	options             []interface{}
	details             []interface{}
	expectedMessageJson string
	expectedMessageSlog []interface{}
	expectedText        string
	expectedSlogLevel   slog.Level
}{
	{
		name:                "messenger-0001",
		messageNumber:       1,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"senzing-99990001","text":"TRACE: Bob works with Jane","location":"In func1() at messenger_test.go:173","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99990001", "location", "In func1() at messenger_test.go:186", "details", []Detail{{Key: "", Position: 1, Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "TRACE: Bob works with Jane",
		expectedSlogLevel:   LevelTraceSlog,
	},
	{
		name:                "messenger-1001",
		messageNumber:       1001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01T00:00:00Z","level":"DEBUG","id":"senzing-99991001","text":"DEBUG: Bob works with Jane","location":"In func1() at messenger_test.go:173","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99991001", "location", "In func1() at messenger_test.go:186", "details", []Detail{{Key: "", Position: 1, Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "DEBUG: Bob works with Jane",
		expectedSlogLevel:   LevelDebugSlog,
	},
	{
		name:                "messenger-2001",
		messageNumber:       2001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01T00:00:00Z","level":"INFO","id":"senzing-99992001","text":"INFO: Bob works with Jane","location":"In func1() at messenger_test.go:173","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99992001", "location", "In func1() at messenger_test.go:186", "details", []Detail{{Key: "", Position: 1, Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "INFO: Bob works with Jane",
		expectedSlogLevel:   LevelInfoSlog,
	},
	{
		name:                "messenger-3001",
		messageNumber:       3001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01T00:00:00Z","level":"WARN","id":"senzing-99993001","text":"WARN: Bob works with Jane","location":"In func1() at messenger_test.go:173","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99993001", "location", "In func1() at messenger_test.go:186", "details", []Detail{{Key: "", Position: 1, Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "WARN: Bob works with Jane",
		expectedSlogLevel:   LevelWarnSlog,
	},
	{
		name:                "messenger-4001",
		messageNumber:       4001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01T00:00:00Z","level":"ERROR","id":"senzing-99994001","text":"ERROR: Bob works with Jane","location":"In func1() at messenger_test.go:173","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99994001", "location", "In func1() at messenger_test.go:186", "details", []Detail{{Key: "", Position: 1, Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "ERROR: Bob works with Jane",
		expectedSlogLevel:   LevelErrorSlog,
	},
	{
		name:                "messenger-5001",
		messageNumber:       5001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01T00:00:00Z","level":"FATAL","id":"senzing-99995001","text":"FATAL: Bob works with Jane","location":"In func1() at messenger_test.go:173","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99995001", "location", "In func1() at messenger_test.go:186", "details", []Detail{{Key: "", Position: 1, Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "FATAL: Bob works with Jane",
		expectedSlogLevel:   LevelFatalSlog,
	},
	{
		name:                "messenger-6001",
		messageNumber:       6001,
		options:             []interface{}{getOptionIdMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01T00:00:00Z","level":"PANIC","id":"senzing-99996001","text":"PANIC: Bob works with Jane","location":"In func1() at messenger_test.go:173","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "senzing-99996001", "location", "In func1() at messenger_test.go:186", "details", []Detail{{Key: "", Position: 1, Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "PANIC: Bob works with Jane",
		expectedSlogLevel:   LevelPanicSlog,
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

func TestMessengerImpl_NewSlogLevel(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessageSlog) > 0 {
			test.Run(testCase.name+"-NewSlog", func(test *testing.T) {
				testObject, err := New(testCase.options...)
				testError(test, testObject, err)
				message, slogLevel, actual := testObject.NewSlogLevel(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedText, message, testCase.name)
				assert.Equal(test, testCase.expectedSlogLevel, slogLevel, testCase.name)
				assert.Equal(test, testCase.expectedMessageSlog, actual, testCase.name)
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
	//Output: {"time":"2000-01-01T00:00:00Z","level":"INFO","id":"senzing-99992001","location":"In ExampleMessengerImpl_NewJson() at messenger_test.go:205","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}
}

func ExampleMessengerImpl_NewSlog() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/messenger/messenger_test.go
	example, err := New()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(example.NewSlog(2001, "Bob", "Jane", getTimestamp(), getOptionCallerSkip()))
	//Output: [id senzing-99992001 location In NewSlog() at messenger.go:370 details [{ 1 Bob <nil>} { 2 Jane <nil>}]]
}

func ExampleMessengerImpl_NewSlogLevel() {
	// For more information, visit https://github.com/Senzing/go-messaging/blob/main/messenger/messenger_test.go
	example, err := New()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(example.NewSlogLevel(2001, "Bob", "Jane", getTimestamp(), getOptionCallerSkip()))
	//Output: INFO [id senzing-99992001 location In ExampleMessengerImpl_NewSlogLevel() at messenger_test.go:225 details [{ 1 Bob <nil>} { 2 Jane <nil>}]]
}
