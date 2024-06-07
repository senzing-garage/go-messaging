package messenger

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
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

var (
	durationTest   time.Duration = 100000000
	errTest1                     = fmt.Errorf("error 1")
	errTest2                     = fmt.Errorf("error 2")
	jsonRawMessage json.RawMessage
)

var idMessages = map[int]string{
	0001: "TRACE: %s works with %s",
	1001: "DEBUG: %s works with %s",
	2001: "INFO: %s works with %s",
	3001: "WARN: %s works with %s",
	3002: "WARN: %s works with %s",
	3003: "WARN: %s works with %s",
	3004: `{"bob": "%s", "jane": "%s"}`,
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
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"TRACE","id":"SZSDK99990001","text":"TRACE: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"level", "TRACE", "id", "SZSDK99990001", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "TRACE: Bob works with Jane", expectedSlogLevel: LevelTraceSlog,
	},
	{
		name:                "messenger-1001",
		messageNumber:       1001,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"DEBUG","id":"SZSDK99991001","text":"DEBUG: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"level", "DEBUG", "id", "SZSDK99991001", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "DEBUG: Bob works with Jane",
		expectedSlogLevel:   LevelDebugSlog,
	},
	{
		name:                "messenger-2001",
		messageNumber:       2001,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"INFO","id":"SZSDK99992001","text":"INFO: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"level", "INFO", "id", "SZSDK99992001", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "INFO: Bob works with Jane",
		expectedSlogLevel:   LevelInfoSlog,
	},
	{
		name:                "messenger-3001",
		messageNumber:       3001,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993001","text":"WARN: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"level", "WARN", "id", "SZSDK99993001", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "WARN: Bob works with Jane",
		expectedSlogLevel:   LevelWarnSlog,
	},
	{
		name:                "messenger-3002",
		messageNumber:       3002,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane", errTest1, errTest2},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993002","text":"WARN: Bob works with Jane","errors":["error 1","error 2"],"details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"},{"position":3,"type":"error","value":"error 1"},{"position":4,"type":"error","value":"error 2"}]}`,
		expectedMessageSlog: []interface{}{"level", "WARN", "id", "SZSDK99993002", "errors", []interface{}{"error 1", "error 2"}, "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}, {Key: "", Position: 3, Type: "error", Value: "error 1", ValueRaw: interface{}(nil)}, {Key: "", Position: 4, Type: "error", Value: "error 2", ValueRaw: interface{}(nil)}}},
		expectedText:        "WARN: Bob works with Jane",
		expectedSlogLevel:   LevelWarnSlog,
	},
	{
		name:                "messenger-3003",
		messageNumber:       3003,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane", int64(12345)},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993003","text":"WARN: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"},{"position":3,"type":"int64","value":"12345","valueRaw":12345}]}`,
		expectedMessageSlog: []interface{}{"level", "WARN", "id", "SZSDK99993003", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}, {Key: "", Position: 3, Type: "int64", Value: "12345", ValueRaw: int64(12345)}}},
		expectedText:        "WARN: Bob works with Jane",
		expectedSlogLevel:   LevelWarnSlog,
	},
	{
		name:                "messenger-3004",
		messageNumber:       3004,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993004","text":"{\"bob\": \"Bob\", \"jane\": \"Jane\"}","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"},{"key":"text","position":3,"type":"map[string]string","value":"{\"bob\": \"Bob\", \"jane\": \"Jane\"}","valueRaw":{"bob":"Bob","jane":"Jane"}}]}`,
		expectedMessageSlog: []interface{}{"level", "WARN", "id", "SZSDK99993004", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}, {Key: "text", Position: 4, Type: "map[string]string", Value: "{\"bob\": \"Bob\", \"jane\": \"Jane\"}", ValueRaw: json.RawMessage{0x7b, 0x22, 0x62, 0x6f, 0x62, 0x22, 0x3a, 0x20, 0x22, 0x42, 0x6f, 0x62, 0x22, 0x2c, 0x20, 0x22, 0x6a, 0x61, 0x6e, 0x65, 0x22, 0x3a, 0x20, 0x22, 0x4a, 0x61, 0x6e, 0x65, 0x22, 0x7d}}}},
		expectedText:        `{"bob": "Bob", "jane": "Jane"}`,
		expectedSlogLevel:   LevelWarnSlog,
	},
	{
		name:                "messenger-3005",
		messageNumber:       3005,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane", getMessageDuration(), getMessageID(), getMessageLevel(), getMessageLocation(), getMessageStatus(), getMessageText(), getOptionCallerSkip()},
		expectedMessageJSON: `{"level":"TEST_LEVEL","id":"Test-ID-1","text":"Test text","status":"TestStatus","duration":100000000,"details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"level", "TEST_LEVEL", "id", "Test-ID-1", "status", "TestStatus", "duration", int64(100000000), "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "Test text",
		expectedSlogLevel:   LevelPanicSlog,
	},
	{
		name:                "messenger-3006",
		messageNumber:       3006,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane", durationTest},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993006","duration":100000000,"details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"level", "WARN", "id", "SZSDK99993006", "duration", int64(100000000), "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "",
		expectedSlogLevel:   LevelWarnSlog,
	},
	{
		name:                "messenger-3007",
		messageNumber:       3007,
		options:             []interface{}{getOptionMessageFieldsAll(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane", getMessageReason(), getMessageCode(), getMessageLocation(), getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"WARN","id":"SZSDK99993007","code":"MessageCode1","reason":"TestMessageReason1","location":"Test location","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"time", "2000-01-01T00:00:00Z", "level", "WARN", "id", "SZSDK99993007", "code", "MessageCode1", "reason", "TestMessageReason1", "location", "Test location", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "",
		expectedSlogLevel:   LevelWarnSlog,
	},
	{
		name:                "messenger-4001",
		messageNumber:       4001,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"ERROR","id":"SZSDK99994001","text":"ERROR: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"level", "ERROR", "id", "SZSDK99994001", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "ERROR: Bob works with Jane",
		expectedSlogLevel:   LevelErrorSlog,
	},
	{
		name:                "messenger-5001",
		messageNumber:       5001,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"FATAL","id":"SZSDK99995001","text":"FATAL: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"level", "FATAL", "id", "SZSDK99995001", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "FATAL: Bob works with Jane",
		expectedSlogLevel:   LevelFatalSlog,
	},
	{
		name:                "messenger-6001",
		messageNumber:       6001,
		options:             []interface{}{getOptionMessageFieldsWithTime(), getOptionIDMessages()},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"PANIC","id":"SZSDK99996001","text":"PANIC: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"time", "2000-01-01T00:00:00Z", "level", "PANIC", "id", "SZSDK99996001", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
		expectedText:        "PANIC: Bob works with Jane",
		expectedSlogLevel:   LevelPanicSlog,
	},
	{
		name:                "messenger-7001",
		messageNumber:       7001,
		options:             []interface{}{getOptionMessageFields(), getOptionIDMessages(), getOptionCallerSkip(), getOptionIDStatuses(), getOptionComponentID(), getOptionMessageIDTemplate()},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"PANIC","id":"Template: 7001","text":"PANIC: Bob works with Jane","status":"status-7001","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{"level", "PANIC", "id", "Template: 7001", "status", "status-7001", "details", []Detail{{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)}, {Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)}}},
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
// Test interface methods
// ----------------------------------------------------------------------------

// -- Test New() method ---------------------------------------------------------

func Test_NewError(test *testing.T) {
	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessageJSON) > 0 {
			test.Run(testCase.name+"-NewError", func(test *testing.T) {
				testObject, err := New(testCase.options...)
				require.NoError(test, err)
				actual := testObject.NewError(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedMessageJSON, actual.Error(), testCase.name)
			})
		}
	}
}

func Test_NewJSON(test *testing.T) {
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

func Test_NewSlogLevel(test *testing.T) {
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
	_, err := New(&OptionComponentID{Value: 99999})
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
// Test private methods
// ---------------------------------------------------------------------------

func TestBasicMessenger_getLevel(test *testing.T) {
	messenger := &BasicMessenger{}
	actual := messenger.getLevel(-1)
	assert.Equal(test, "UNKNOWN", actual)
}

func TestBasicMessenger_populateFields(test *testing.T) {
	_ = test
	messenger := &BasicMessenger{}
	messenger.populateMessageFields()
}

func TestBasicMessenger_populateMessageFields(test *testing.T) {
	test.Setenv("SENZING_MESSAGE_FIELDS", "id, duration, text")
	messenger := &BasicMessenger{}
	messenger.populateMessageFields()
}

func TestBasicMessenger_populateMessageFields_all(test *testing.T) {
	test.Setenv("SENZING_MESSAGE_FIELDS", "all")
	messenger := &BasicMessenger{}
	messenger.populateMessageFields()
}

func TestBasicMessenger_populateStructure(test *testing.T) {
	_ = test
	messenger := &BasicMessenger{}
	messenger.populateStructure(1234, getMessageStatus())
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

// ----------------------------------------------------------------------------
// Internal functions - names begin with lowercase letter
// ----------------------------------------------------------------------------

func getMessageCode() *MessageCode {
	return &MessageCode{
		Value: "MessageCode1",
	}
}

func getMessageDuration() *MessageDuration {
	var duration time.Duration = 100000000
	return &MessageDuration{
		Value: duration.Nanoseconds(),
	}
}

func getMessageID() *MessageID {
	return &MessageID{
		Value: "Test-ID-1",
	}
}

func getMessageLevel() *MessageLevel {
	return &MessageLevel{
		Value: "TEST_LEVEL",
	}
}

func getMessageLocation() *MessageLocation {
	return &MessageLocation{
		Value: "Test location",
	}
}

func getMessageReason() *MessageReason {
	return &MessageReason{
		Value: "TestMessageReason1",
	}
}

func getMessageStatus() *MessageStatus {
	return &MessageStatus{
		Value: "TestStatus",
	}
}

func getMessageText() *MessageText {
	return &MessageText{
		Value: "Test text",
	}
}

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

func getOptionMessageFields() *OptionMessageFields {

	messageFields := make([]string, len(AllMessageFields))
	_ = copy(messageFields, AllMessageFields)

	// Remove "location".

	indexOfLocation := slices.Index(messageFields, "location")
	if indexOfLocation > 0 {
		messageFields = slices.Delete(messageFields, indexOfLocation, indexOfLocation+1)
	}

	// Remove "time".

	indexOfLocation = slices.Index(messageFields, "time")
	if indexOfLocation >= 0 {
		messageFields = slices.Delete(messageFields, indexOfLocation, indexOfLocation+1)
	}
	return &OptionMessageFields{
		Value: messageFields,
	}
}

func getOptionMessageFieldsAll() *OptionMessageFields {
	return &OptionMessageFields{
		Value: AllMessageFields,
	}
}

func getOptionMessageFieldsWithTime() *OptionMessageFields {

	// Remove "location".

	messageFields := make([]string, len(AllMessageFields))
	_ = copy(messageFields, AllMessageFields)

	indexOfLocation := slices.Index(messageFields, "location")
	if indexOfLocation > 0 {
		messageFields = slices.Delete(messageFields, indexOfLocation, indexOfLocation+1)
	}
	return &OptionMessageFields{
		Value: messageFields,
	}
}

func getOptionMessageIDTemplate() *OptionMessageIDTemplate {
	return &OptionMessageIDTemplate{
		Value: "Template: %04d",
	}
}

func getOptionComponentID() *OptionComponentID {
	return &OptionComponentID{
		Value: 9999,
	}
}

func getTimestamp() *MessageTime {
	return &MessageTime{
		Value: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}
