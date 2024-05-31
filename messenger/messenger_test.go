package messenger

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

const (
	badLevel = 99999
	jsonTest = `{"Key": "Value"}`
)

var jsonRawMessage json.RawMessage

var idMessages = map[int]string{
	0001: "TRACE: %s works with %s",
	1001: "DEBUG: %s works with %s",
	2001: "INFO: %s works with %s",
	3001: "WARN: %s works with %s",
	4001: "ERROR: %s works with %s",
	5001: "FATAL: %s works with %s",
	6001: "PANIC: %s works with %s",
	7001: "PANIC: %s works with %s",
}

var idStatuses = map[int]string{
	0001: "status-0001",
	1001: "status-1001",
	2001: "status-2001",
	3001: "status-3001",
	4001: "status-4001",
	5001: "status-5001",
	6001: "status-6001",
	7001: "status-7001",
}

var testCasesForMessage = []struct {
	name                string
	messageNumber       int
	options             []interface{}
	details             []interface{}
	expectedMessageJSON string
	expectedMessageSlog []interface{}
	expectedText        string
	expectedSlogLevel   slog.Level
}{
	{
		name:                "messenger-0001",
		messageNumber:       1,
		options:             []interface{}{getOptionIDMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"TRACE","id":"SZSDK99990001","text":"TRACE: Bob works with Jane","location":"In func1() at messenger_test.go:268","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "SZSDK99990001", "location", "In func1() at messenger_test.go:281", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "TRACE: Bob works with Jane",
		expectedSlogLevel:   LevelTraceSlog,
	},
	{
		name:                "messenger-1001",
		messageNumber:       1001,
		options:             []interface{}{getOptionIDMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"DEBUG","id":"SZSDK99991001","text":"DEBUG: Bob works with Jane","location":"In func1() at messenger_test.go:268","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "SZSDK99991001", "location", "In func1() at messenger_test.go:281", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "DEBUG: Bob works with Jane",
		expectedSlogLevel:   LevelDebugSlog,
	},
	{
		name:                "messenger-2001",
		messageNumber:       2001,
		options:             []interface{}{getOptionIDMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"INFO","id":"SZSDK99992001","text":"INFO: Bob works with Jane","location":"In func1() at messenger_test.go:268","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "SZSDK99992001", "location", "In func1() at messenger_test.go:281", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "INFO: Bob works with Jane",
		expectedSlogLevel:   LevelInfoSlog,
	},
	{
		name:                "messenger-3001",
		messageNumber:       3001,
		options:             []interface{}{getOptionIDMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"WARN","id":"SZSDK99993001","text":"WARN: Bob works with Jane","location":"In func1() at messenger_test.go:268","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "SZSDK99993001", "location", "In func1() at messenger_test.go:281", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "WARN: Bob works with Jane",
		expectedSlogLevel:   LevelWarnSlog,
	},
	{
		name:                "messenger-4001",
		messageNumber:       4001,
		options:             []interface{}{getOptionIDMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"ERROR","id":"SZSDK99994001","text":"ERROR: Bob works with Jane","location":"In func1() at messenger_test.go:268","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "SZSDK99994001", "location", "In func1() at messenger_test.go:281", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "ERROR: Bob works with Jane",
		expectedSlogLevel:   LevelErrorSlog,
	},
	{
		name:                "messenger-5001",
		messageNumber:       5001,
		options:             []interface{}{getOptionIDMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"FATAL","id":"SZSDK99995001","text":"FATAL: Bob works with Jane","location":"In func1() at messenger_test.go:268","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "SZSDK99995001", "location", "In func1() at messenger_test.go:281", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "FATAL: Bob works with Jane",
		expectedSlogLevel:   LevelFatalSlog,
	},
	{
		name:                "messenger-6001",
		messageNumber:       6001,
		options:             []interface{}{getOptionIDMessages(), getOptionCallerSkip()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"PANIC","id":"SZSDK99996001","text":"PANIC: Bob works with Jane","location":"In func1() at messenger_test.go:268","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "SZSDK99996001", "location", "In func1() at messenger_test.go:281", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "PANIC: Bob works with Jane",
		expectedSlogLevel:   LevelPanicSlog,
	},
	{
		name:                "messenger-7001",
		messageNumber:       7001,
		options:             []interface{}{getOptionIDMessages(), getOptionCallerSkip(), getOptionIDStatuses(), getOptionSenzingComponentID(), getOptionMessageIDTemplate()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"PANIC","id":"Template: 7001","text":"PANIC: Bob works with Jane","status":"status-7001","location":"In func1() at messenger_test.go:268","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"id", "Template: 7001", "status", "status-7001", "location", "In func1() at messenger_test.go:281", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "PANIC: Bob works with Jane",
		expectedSlogLevel:   LevelPanicSlog,
	},
}

var testCasesForMessageDetails = []struct {
	name     string
	input    any
	expected []Detail
}{
	{
		name:     "nil",
		input:    nil,
		expected: []Detail{{Position: 1, Type: "nil"}},
	},
	{
		name:     "int",
		input:    1,
		expected: []Detail{{Position: 1, Type: "integer", Value: "1", ValueRaw: 1}},
	},
	{
		name:     "float64",
		input:    0.6,
		expected: []Detail{{Position: 1, Type: "float", Value: "0.6", ValueRaw: 0.6}},
	},
	{
		name:     "string",
		input:    "a string",
		expected: []Detail{{Position: 1, Type: "string", Value: "a string"}},
	},
	{
		name:     "bool",
		input:    true,
		expected: []Detail{{Position: 1, Type: "boolean", Value: "true", ValueRaw: true}},
	},
	{
		name:     "error",
		input:    fmt.Errorf("test error"),
		expected: []Detail{{Position: 1, Type: "error", Value: "test error", ValueRaw: nil}},
	},
	{
		name:     "map[string]string",
		input:    map[string]string{"string1": "string2"},
		expected: []Detail{{Position: 1, Type: "map[string]string", Key: "string1", Value: "string2", ValueRaw: nil}},
	},
	{
		name:     "int64",
		input:    int64(1),
		expected: []Detail{{Position: 1, Type: "int64", Value: "1", ValueRaw: int64(1)}},
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
	var err error
	return err
}

func teardown() error {
	var err error
	return err
}

// ----------------------------------------------------------------------------
// Internal functions - names begin with lowercase letter
// ----------------------------------------------------------------------------

func getOptionCallerSkip() *OptionCallerSkip {
	return &OptionCallerSkip{
		Value: 2,
	}
}

func getOptionIDMessages() *OptionIDMessages {
	return &OptionIDMessages{
		Value: idMessages,
	}
}

func getOptionIDStatuses() *OptionIDStatuses {
	return &OptionIDStatuses{
		Value: idStatuses,
	}
}

func getOptionMessageIDTemplate() *OptionMessageIDTemplate {
	return &OptionMessageIDTemplate{
		Value: "Template: %04d",
	}
}

func getOptionSenzingComponentID() *OptionSenzingComponentID {
	return &OptionSenzingComponentID{
		Value: 9999,
	}
}

func getTimestamp() *MessageTime {
	return &MessageTime{
		Value: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}

// ----------------------------------------------------------------------------
// Test interface methods
// ----------------------------------------------------------------------------

// -- Test New() method ---------------------------------------------------------

func TestImplementation_NewJSON(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessageJSON) > 0 {
			test.Run(testCase.name+"-NewJson", func(test *testing.T) {
				testObject, err := New(testCase.options...)
				require.NoError(test, err)
				actual := testObject.NewJSON(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedMessageJSON, actual, testCase.name)
			})
		}
	}
}

func TestImplementation_NewSlogLevel(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessageSlog) > 0 {
			test.Run(testCase.name+"-NewSlog", func(test *testing.T) {
				testObject, err := New(testCase.options...)
				require.NoError(test, err)
				message, slogLevel, actual := testObject.NewSlogLevel(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedText, message, testCase.name)
				assert.Equal(test, testCase.expectedSlogLevel, slogLevel, testCase.name)
				assert.Equal(test, testCase.expectedMessageSlog, actual, testCase.name)
			})
		}
	}
}

func Test_New_badComponentID(test *testing.T) {
	_, err := New(&OptionSenzingComponentID{Value: 99999})
	require.ErrorIs(test, err, ErrBadComponentID)
}

func Test_New_badIdMessages(test *testing.T) {
	_, err := New(&OptionIDMessages{})
	require.ErrorIs(test, err, ErrEmptyMessages)
}

func Test_New_badIdStatuses(test *testing.T) {
	_, err := New(&OptionIDStatuses{})
	require.ErrorIs(test, err, ErrEmptyStatuses)
}

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

func Test_cleanErrorString(test *testing.T) {
	err := fmt.Errorf("\n\tError\t\n") //revive:disable-line:error-strings
	cleanErrorString := cleanErrorString(err)
	assert.Equal(test, "Error", cleanErrorString)
}

func Test_interfaceAsString(test *testing.T) {
	assert.Equal(test, "<nil>", interfaceAsString(nil))
	assert.Equal(test, `{"json": "string"}`, interfaceAsString(`{"json": "string"}`))
	assert.Equal(test, "a string", interfaceAsString("a string"))
	assert.Equal(test, "5", interfaceAsString(5))
	assert.Equal(test, "0.6", interfaceAsString(0.6))
	assert.Equal(test, "true", interfaceAsString(true))
	assert.Equal(test, "An error", interfaceAsString(fmt.Errorf("An error")))
	assert.Equal(test, "5", interfaceAsString(int64(5)))
}

func Test_jsonAsInterface(test *testing.T) {
	jsonString := `{"json": "string"}`
	jsonAsInterface := jsonAsInterface(jsonString)
	assert.NotNil(test, jsonAsInterface)
}

func Test_jsonAsInterface_badJSON(test *testing.T) {
	jsonString := `}{`
	assert.Panics(test, func() { jsonAsInterface(jsonString) })
}

func Test_messageDetails(test *testing.T) {
	for _, testCase := range testCasesForMessageDetails {
		test.Run(testCase.name, func(test *testing.T) {
			actual := messageDetails(testCase.input)
			assert.Equal(test, testCase.expected, actual)
		})
	}
}

func Test_messageDetails_empty(test *testing.T) {
	expected := []Detail{}
	actual := messageDetails()
	assert.Equal(test, expected, actual)
}

func Test_messageDetails_errJSON(test *testing.T) {
	err := json.Unmarshal([]byte(jsonTest), &jsonRawMessage)
	require.NoError(test, err)
	expected := []Detail{{Position: 1, Type: "error", Value: jsonTest, ValueRaw: jsonRawMessage}}
	testErr := fmt.Errorf(jsonTest)
	actual := messageDetails(testErr)
	assert.Equal(test, expected, actual)
}

func Test_messageDetails_mapstringstringJSON(test *testing.T) {
	input := map[string]string{"string1": jsonTest}
	err := json.Unmarshal([]byte(jsonTest), &jsonRawMessage)
	require.NoError(test, err)
	expected := []Detail{{Position: 1, Type: "map[string]string", Key: "string1", Value: jsonTest, ValueRaw: jsonRawMessage}}
	actual := messageDetails(input)
	assert.Equal(test, expected, actual)
}

func Test_messageDetails_stringJSON(test *testing.T) {
	err := json.Unmarshal([]byte(jsonTest), &jsonRawMessage)
	require.NoError(test, err)
	expected := []Detail{{Position: 1, Type: "string", Value: jsonTest, ValueRaw: jsonRawMessage}}
	actual := messageDetails(jsonTest)
	assert.Equal(test, expected, actual)
}
