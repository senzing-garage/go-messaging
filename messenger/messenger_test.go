package messenger_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/senzing-garage/go-messaging/messenger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

// const (
// 	badLevel = 99999
// 	jsonTest = `{"Key": "Value"}`
// )

var (
	durationTest time.Duration = 100000000
	// err                        = fmt.Errorf("") .
	errTest1 = errors.New("error 1")
	errTest2 = errors.New("error 2")
	// jsonRawMessage json.RawMessage .
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
	comment             string
	options             []interface{}
	details             []interface{}
	expectedMessageJSON string
	expectedMessageSlog []interface{}
	expectedText        string
	expectedSlogLevel   slog.Level
}{
	{
		name:          "messenger-0001",
		messageNumber: 1,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"TRACE","id":"SZSDK99990001","text":"TRACE: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"TRACE",
			"id",
			"SZSDK99990001",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText: "TRACE: Bob works with Jane", expectedSlogLevel: messenger.LevelTraceSlog,
	},
	{
		name:          "messenger-1001",
		messageNumber: 1001,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9998),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"DEBUG","id":"SZSDK99981001","text":"DEBUG: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"DEBUG",
			"id",
			"SZSDK99981001",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "DEBUG: Bob works with Jane",
		expectedSlogLevel: messenger.LevelDebugSlog,
	},
	{
		name:          "messenger-2001",
		messageNumber: 2001,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"INFO","id":"SZSDK99992001","text":"INFO: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"INFO",
			"id",
			"SZSDK99992001",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "INFO: Bob works with Jane",
		expectedSlogLevel: messenger.LevelInfoSlog,
	},
	{
		name:          "messenger-3001",
		messageNumber: 3001,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993001","text":"WARN: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"WARN",
			"id",
			"SZSDK99993001",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "WARN: Bob works with Jane",
		expectedSlogLevel: messenger.LevelWarnSlog,
	},
	{
		name:          "messenger-3002",
		messageNumber: 3002,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane", errTest1, errTest2},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993002","text":"WARN: Bob works with Jane","errors":["error 1","error 2"],"details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"},{"position":3,"type":"error","value":"error 1"},{"position":4,"type":"error","value":"error 2"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"WARN",
			"id",
			"SZSDK99993002",
			"errors",
			[]interface{}{"error 1", "error 2"},
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
				{Key: "", Position: 3, Type: "error", Value: "error 1", ValueRaw: interface{}(nil)},
				{Key: "", Position: 4, Type: "error", Value: "error 2", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "WARN: Bob works with Jane",
		expectedSlogLevel: messenger.LevelWarnSlog,
	},
	{
		name:          "messenger-3003",
		messageNumber: 3003,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane", int64(12345)},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993003","text":"WARN: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"},{"position":3,"type":"int64","value":"12345","valueRaw":12345}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"WARN",
			"id",
			"SZSDK99993003",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
				{Key: "", Position: 3, Type: "int64", Value: "12345", ValueRaw: int64(12345)},
			},
		},
		expectedText:      "WARN: Bob works with Jane",
		expectedSlogLevel: messenger.LevelWarnSlog,
	},
	{
		name:          "messenger-3004",
		messageNumber: 3004,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993004","text":"{\"bob\": \"Bob\", \"jane\": \"Jane\"}","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"},{"key":"text","position":3,"type":"map[string]string","value":"{\"bob\": \"Bob\", \"jane\": \"Jane\"}","valueRaw":{"bob":"Bob","jane":"Jane"}}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"WARN",
			"id",
			"SZSDK99993004",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
				{
					Key:      "text",
					Position: 4,
					Type:     "map[string]string",
					Value:    "{\"bob\": \"Bob\", \"jane\": \"Jane\"}",
					ValueRaw: json.RawMessage{
						0x7b,
						0x22,
						0x62,
						0x6f,
						0x62,
						0x22,
						0x3a,
						0x20,
						0x22,
						0x42,
						0x6f,
						0x62,
						0x22,
						0x2c,
						0x20,
						0x22,
						0x6a,
						0x61,
						0x6e,
						0x65,
						0x22,
						0x3a,
						0x20,
						0x22,
						0x4a,
						0x61,
						0x6e,
						0x65,
						0x22,
						0x7d,
					},
				},
			},
		},
		expectedText:      `{"bob": "Bob", "jane": "Jane"}`,
		expectedSlogLevel: messenger.LevelWarnSlog,
	},
	{
		name:          "messenger-3005",
		messageNumber: 3005,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details: []interface{}{
			"Bob",
			"Jane",
			getMessageDuration(),
			getMessageID(),
			getMessageLevel(),
			getMessageLocation(),
			getMessageStatus(),
			getMessageText(),
			getOptionCallerSkip(),
		},
		expectedMessageJSON: `{"level":"TEST_LEVEL","id":"Test-ID-1","text":"Test text","status":"TestStatus","duration":100000000,"details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"TEST_LEVEL",
			"id",
			"Test-ID-1",
			"status",
			"TestStatus",
			"duration",
			int64(100000000),
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "Test text",
		expectedSlogLevel: messenger.LevelPanicSlog,
	},
	{
		name:          "messenger-3006",
		messageNumber: 3006,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane", durationTest},
		expectedMessageJSON: `{"level":"WARN","id":"SZSDK99993006","duration":100000000,"details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"WARN",
			"id",
			"SZSDK99993006",
			"duration",
			int64(100000000),
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "",
		expectedSlogLevel: messenger.LevelWarnSlog,
	},
	{
		name:          "messenger-3007",
		messageNumber: 3007,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFieldsAll(),
			getOptionIDMessages(),
		},
		details: []interface{}{
			"Bob",
			"Jane",
			getMessageReason(),
			getMessageCode(),
			getMessageLocation(),
			getTimestamp(),
		},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"WARN","id":"SZSDK99993007","code":"MessageCode1","reason":"TestMessageReason1","location":"Test location","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"time",
			"2000-01-01T00:00:00Z",
			"level",
			"WARN",
			"id",
			"SZSDK99993007",
			"code",
			"MessageCode1",
			"reason",
			"TestMessageReason1",
			"location",
			"Test location",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "",
		expectedSlogLevel: messenger.LevelWarnSlog,
	},
	{
		name:          "messenger-4001",
		messageNumber: 4001,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"ERROR","id":"SZSDK99994001","text":"ERROR: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"ERROR",
			"id",
			"SZSDK99994001",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "ERROR: Bob works with Jane",
		expectedSlogLevel: messenger.LevelErrorSlog,
	},
	{
		name:          "messenger-5001",
		messageNumber: 5001,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"FATAL","id":"SZSDK99995001","text":"FATAL: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"FATAL",
			"id",
			"SZSDK99995001",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "FATAL: Bob works with Jane",
		expectedSlogLevel: messenger.LevelFatalSlog,
	},
	{
		name:          "messenger-6001",
		messageNumber: 6001,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFieldsWithTime(),
			getOptionIDMessages(),
		},
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJSON: `{"time":"2000-01-01T00:00:00Z","level":"PANIC","id":"SZSDK99996001","text":"PANIC: Bob works with Jane","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"time",
			"2000-01-01T00:00:00Z",
			"level",
			"PANIC",
			"id",
			"SZSDK99996001",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "PANIC: Bob works with Jane",
		expectedSlogLevel: messenger.LevelPanicSlog,
	},
	{
		name:          "messenger-7001",
		messageNumber: 7001,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
			getOptionCallerSkip(),
			getOptionIDStatuses(),
		},
		details:             []interface{}{"Bob", "Jane"},
		expectedMessageJSON: `{"level":"PANIC","id":"SZSDK99997001","text":"PANIC: Bob works with Jane","status":"status-7001","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"PANIC",
			"id",
			"SZSDK99997001",
			"status",
			"status-7001",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "string", Value: "Bob", ValueRaw: interface{}(nil)},
				{Key: "", Position: 2, Type: "string", Value: "Jane", ValueRaw: interface{}(nil)},
			},
		},
		expectedText:      "PANIC: Bob works with Jane",
		expectedSlogLevel: messenger.LevelPanicSlog,
	},
	{
		name:          "messenger-9001",
		messageNumber: 9001,
		comment:       "",
		options: []interface{}{
			getOptionMessageIDTemplate(9999),
			getOptionMessageFields(),
			getOptionIDMessages(),
			getOptionCallerSkip(),
			getOptionIDStatuses(),
		},
		details:             []interface{}{123, true, 1.23},
		expectedMessageJSON: `{"level":"PANIC","id":"SZSDK99999001","details":[{"position":1,"type":"integer","value":"123","valueRaw":123},{"position":2,"type":"boolean","value":"true","valueRaw":true},{"position":3,"type":"float","value":"1.23","valueRaw":1.23}]}`,
		expectedMessageSlog: []interface{}{
			"level",
			"PANIC",
			"id",
			"SZSDK99999001",
			"details",
			[]messenger.Detail{
				{Key: "", Position: 1, Type: "integer", Value: "123", ValueRaw: 123},
				{Key: "", Position: 2, Type: "boolean", Value: "true", ValueRaw: true},
				{Key: "", Position: 3, Type: "float", Value: "1.23", ValueRaw: 1.23},
			},
		},
		expectedText:      "",
		expectedSlogLevel: messenger.LevelPanicSlog,
	},
}

// var testCasesForMessageDetails = []struct {
// 	name     string
// 	input    any
// 	expected []messenger.Detail
// }{
// 	{
// 		name:     "nil",
// 		input:    nil,
// 		expected: []messenger.Detail{{Position: 1, Type: "nil"}},
// 	},
// 	{
// 		name:     "int",
// 		input:    1,
// 		expected: []messenger.Detail{{Position: 1, Type: "integer", Value: "1", ValueRaw: 1}},
// 	},
// 	{
// 		name:     "float64",
// 		input:    0.6,
// 		expected: []messenger.Detail{{Position: 1, Type: "float", Value: "0.6", ValueRaw: 0.6}},
// 	},
// 	{
// 		name:     "string",
// 		input:    "a string",
// 		expected: []messenger.Detail{{Position: 1, Type: "string", Value: "a string"}},
// 	},
// 	{
// 		name:     "bool",
// 		input:    true,
// 		expected: []messenger.Detail{{Position: 1, Type: "boolean", Value: "true", ValueRaw: true}},
// 	},
// 	{
// 		name:     "error",
// 		input:    fmt.Errorf("test error: %w", err),
// 		expected: []messenger.Detail{{Position: 1, Type: "error", Value: "test error: ", ValueRaw: nil}},
// 	},
// 	{
// 		name:  "map[string]string",
// 		input: map[string]string{"string1": "string2"},
// 		expected: []messenger.Detail{
// 			{Position: 1, Type: "map[string]string", Key: "string1", Value: "string2", ValueRaw: nil},
// 		},
// 	},
// 	{
// 		name:     "int64",
// 		input:    int64(1),
// 		expected: []messenger.Detail{{Position: 1, Type: "int64", Value: "1", ValueRaw: int64(1)}},
// 	},
// }

// ----------------------------------------------------------------------------
// Test interface methods
// ----------------------------------------------------------------------------

func Test_XXX(test *testing.T) {
	test.Parallel()

	testObject, err := messenger.New(getOptionMessageIDTemplate(9999), getOptionMessageFields(), getOptionIDMessages())
	require.NoError(test, err)

	err1 := testObject.NewError(1, getMessageStatusValue("a1"))
	err2 := testObject.NewError(2, getMessageStatusValue("b2"))
	_ = testObject.NewError(2, err1, err2)
}

// -- Test New() method ---------------------------------------------------------

func Test_NewError(test *testing.T) {
	test.Parallel()

	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessageJSON) > 0 {
			test.Run(testCase.name+"-NewError", func(test *testing.T) {
				test.Parallel()

				testObject, err := messenger.New(testCase.options...)
				require.NoError(test, err)

				actual := testObject.NewError(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedMessageJSON, actual.Error(), testCase.name)
			})
		}
	}
}

func Test_NewJSON(test *testing.T) {
	test.Parallel()

	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessageJSON) > 0 {
			test.Run(testCase.name+"-NewJson", func(test *testing.T) {
				test.Parallel()

				testObject, err := messenger.New(testCase.options...)
				require.NoError(test, err)

				actual := testObject.NewJSON(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedMessageJSON, actual, testCase.name)
			})
		}
	}
}

func Test_NewJSON_envvar(test *testing.T) {
	test.Setenv("SENZING_MESSAGE_FIELDS", "id, level, text")

	options := []interface{}{getOptionMessageIDTemplate(9999), getOptionMessageFields(), getOptionIDMessages()}
	details := []interface{}{"Bob", "Jane"}
	expected := `{"level":"WARN","id":"SZSDK99993003","text":"WARN: Bob works with Jane"}`
	testObject, err := messenger.New(options...)
	require.NoError(test, err)

	actual := testObject.NewJSON(3003, details...)
	assert.Equal(test, expected, actual)
}

func Test_NewJSON_envvar_all(test *testing.T) {
	test.Setenv("SENZING_MESSAGE_FIELDS", "all")

	options := []interface{}{getOptionMessageIDTemplate(9999), getOptionMessageFields(), getOptionIDMessages()}
	details := []interface{}{"Bob", "Jane"}
	testObject, err := messenger.New(options...)
	require.NoError(test, err)

	_ = testObject.NewJSON(3003, details...)
}

func Test_NewSlogLevel(test *testing.T) {
	test.Parallel()

	for _, testCase := range testCasesForMessage {
		if len(testCase.expectedMessageSlog) > 0 {
			test.Run(testCase.name+"-NewSlog", func(test *testing.T) {
				test.Parallel()

				testObject, err := messenger.New(testCase.options...)
				require.NoError(test, err)

				message, slogLevel, actual := testObject.NewSlogLevel(testCase.messageNumber, testCase.details...)
				assert.Equal(test, testCase.expectedText, message, testCase.name)
				assert.Equal(test, testCase.expectedSlogLevel, slogLevel, testCase.name)
				assert.Equal(test, testCase.expectedMessageSlog, actual, testCase.name)
			})
		}
	}
}

// -- Test error conditions -----------------------------------------------------

func Test_YYY(test *testing.T) {
	test.Parallel()
	testObject := &messenger.BasicMessenger{}
	_, _, _ = testObject.NewSlogLevel(1)
}

func Test_New_badIdMessages(test *testing.T) {
	test.Parallel()

	_, err := messenger.New(messenger.OptionIDMessages{})
	require.ErrorIs(test, err, messenger.ErrEmptyMessages)
}

func Test_New_badIdStatuses(test *testing.T) {
	test.Parallel()

	_, err := messenger.New(messenger.OptionIDStatuses{})
	require.ErrorIs(test, err, messenger.ErrEmptyStatuses)
}

// ----------------------------------------------------------------------------
// Test private methods
// ---------------------------------------------------------------------------

// func TestBasicMessenger_getLevel(test *testing.T) {
// 	messenger := &messenger.BasicMessenger{}
// 	actual := messenger.getLevel(-1)
// 	assert.Equal(test, "UNKNOWN", actual)
// }

// func TestBasicMessenger_populateFields(test *testing.T) {
// 	_ = test
// 	messenger := &messenger.BasicMessenger{}
// 	messenger.populateMessageFields("")
// }

// func TestBasicMessenger_populateMessageFields(test *testing.T) {
// 	test.Setenv("SENZING_MESSAGE_FIELDS", "id, duration, text")
// 	messenger := &messenger.BasicMessenger{}
// 	messenger.populateMessageFields("id, duration, text")
// }

// func TestBasicMessenger_populateMessageFields_all(test *testing.T) {
// 	test.Setenv("SENZING_MESSAGE_FIELDS", "all")
// 	messenger := &messenger.BasicMessenger{}
// 	messenger.populateMessageFields("all")
// }

// func TestBasicMessenger_populateStructure(test *testing.T) {
// 	_ = test
// 	messenger := &messenger.BasicMessenger{}
// 	messenger.populateStructure(1234, getMessageStatus())
// }

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

// func Test_cleanErrorString(test *testing.T) {
// 	err := fmt.Errorf("\n\tError: %w\t\n", err) //revive:disable-line:error-strings
// 	cleanErrorString := cleanErrorString(err)
// 	assert.Equal(test, "Error: ", cleanErrorString)
// }

// func Test_interfaceAsString(test *testing.T) {
// 	assert.Equal(test, "<nil>", interfaceAsString(nil))
// 	assert.JSONEq(test, `{"json": "string"}`, interfaceAsString(`{"json": "string"}`))
// 	assert.Equal(test, "a string", interfaceAsString("a string"))
// 	assert.Equal(test, "5", interfaceAsString(5))
// 	assert.Equal(test, "0.6", interfaceAsString(0.6))
// 	assert.Equal(test, "true", interfaceAsString(true))
// 	assert.Equal(test, "An error: ", interfaceAsString(fmt.Errorf("An error: %w", err)))
// 	assert.Equal(test, "5", interfaceAsString(int64(5)))
// }

// func Test_jsonAsInterface(test *testing.T) {
// 	jsonString := `{"json": "string"}`
// 	jsonAsInterface := jsonAsInterface(jsonString)
// 	assert.NotNil(test, jsonAsInterface)
// }

// func Test_jsonAsInterface_badJSON(test *testing.T) {
// 	jsonString := `}{`
// 	assert.Panics(test, func() { jsonAsInterface(jsonString) })
// }

// func Test_messageDetails(test *testing.T) {
// 	for _, testCase := range testCasesForMessageDetails {
// 		test.Run(testCase.name, func(test *testing.T) {
// 			actual := messageDetails(testCase.input)
// 			assert.Equal(test, testCase.expected, actual)
// 		})
// 	}
// }

// func Test_messageDetails_empty(test *testing.T) {
// 	expected := []Detail{}
// 	actual := messageDetails()
// 	assert.Equal(test, expected, actual)
// }

// func Test_messageDetails_errJSON(test *testing.T) {
// 	err := json.Unmarshal([]byte(jsonTest), &jsonRawMessage)
// 	require.NoError(test, err)
// 	expected := []Detail{{Position: 1, Type: "error", Value: jsonTest, ValueRaw: jsonRawMessage}}
// 	testErr := errors.New(jsonTest) //nolint
// 	actual := messageDetails(testErr)
// 	assert.Equal(test, expected, actual)
// }

// func Test_messageDetails_mapstringstringJSON(test *testing.T) {
// 	input := map[string]string{"string1": jsonTest}
// 	err := json.Unmarshal([]byte(jsonTest), &jsonRawMessage)
// 	require.NoError(test, err)
// 	expected := []Detail{
// 		{Position: 1, Type: "map[string]string", Key: "string1", Value: jsonTest, ValueRaw: jsonRawMessage},
// 	}
// 	actual := messageDetails(input)
// 	assert.Equal(test, expected, actual)
// }

// func Test_messageDetails_stringJSON(test *testing.T) {
// 	err := json.Unmarshal([]byte(jsonTest), &jsonRawMessage)
// 	require.NoError(test, err)
// 	expected := []Detail{{Position: 1, Type: "string", Value: jsonTest, ValueRaw: jsonRawMessage}}
// 	actual := messageDetails(jsonTest)
// 	assert.Equal(test, expected, actual)
// }

// ----------------------------------------------------------------------------
// Internal functions - names begin with lowercase letter
// ----------------------------------------------------------------------------

func getMessageCode() messenger.MessageCode {
	return messenger.MessageCode{
		Value: "MessageCode1",
	}
}

func getMessageDuration() messenger.MessageDuration {
	var duration time.Duration = 100000000

	return messenger.MessageDuration{
		Value: duration.Nanoseconds(),
	}
}

func getMessageID() messenger.MessageID {
	return messenger.MessageID{
		Value: "Test-ID-1",
	}
}

func getMessageLevel() messenger.MessageLevel {
	return messenger.MessageLevel{
		Value: "TEST_LEVEL",
	}
}

func getMessageLocation() messenger.MessageLocation {
	return messenger.MessageLocation{
		Value: "Test location",
	}
}

func getMessageReason() messenger.MessageReason {
	return messenger.MessageReason{
		Value: "TestMessageReason1",
	}
}

func getMessageStatus() messenger.MessageStatus {
	return messenger.MessageStatus{
		Value: "TestStatus",
	}
}

func getMessageStatusValue(value string) messenger.MessageStatus {
	return messenger.MessageStatus{
		Value: value,
	}
}

func getMessageText() messenger.MessageText {
	return messenger.MessageText{
		Value: "Test text",
	}
}

func getOptionCallerSkip() messenger.OptionCallerSkip {
	return messenger.OptionCallerSkip{
		Value: 2,
	}
}

func getOptionIDMessages() messenger.OptionIDMessages {
	return messenger.OptionIDMessages{
		Value: idMessages,
	}
}

func getOptionIDStatuses() messenger.OptionIDStatuses {
	return messenger.OptionIDStatuses{
		Value: idStatuses,
	}
}

func getOptionMessageFields() messenger.OptionMessageFields {
	messageFields := make([]string, len(messenger.AllMessageFields))
	_ = copy(messageFields, messenger.AllMessageFields)

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

	return messenger.OptionMessageFields{
		Value: messageFields,
	}
}

func getOptionMessageFieldsAll() messenger.OptionMessageFields {
	return messenger.OptionMessageFields{
		Value: messenger.AllMessageFields,
	}
}

func getOptionMessageFieldsWithTime() messenger.OptionMessageFields {
	// Remove "location".
	messageFields := make([]string, len(messenger.AllMessageFields))
	_ = copy(messageFields, messenger.AllMessageFields)

	indexOfLocation := slices.Index(messageFields, "location")
	if indexOfLocation > 0 {
		messageFields = slices.Delete(messageFields, indexOfLocation, indexOfLocation+1)
	}

	return messenger.OptionMessageFields{
		Value: messageFields,
	}
}

func getOptionMessageIDTemplate(componentID int) messenger.OptionMessageIDTemplate {
	return messenger.OptionMessageIDTemplate{
		Value: fmt.Sprintf("SZSDK%04d", componentID) + "%04d",
	}
}

func getTimestamp() messenger.MessageTime {
	return messenger.MessageTime{
		Value: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}
