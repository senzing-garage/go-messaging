package typedef

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/senzing/go-messaging/messenger"
	"github.com/stretchr/testify/assert"
)

var (
	logger messenger.MessengerInterface

	callerSkip          = 0
	componentIdentifier = 9999
	isDebug             = false
	messageIdTemplate   = "test-%04d"

	idMessages = map[int]string{
		0001: "TRACE: %s works with %s",
		0002: "The person is %s",
		1001: "DEBUG: %s works with %s",
		2001: "INFO: %s works with %s",
		2002: "%s",
		2003: "JSON: %s",
		3001: "WARN: %s works with %s",
		4001: "ERROR: %s works with %s",
		5001: "FATAL: %s works with %s",
		6001: "PANIC: %s works with %s",
	}

	idStatuses = map[int]string{
		0001: "Status for 0001",
		1000: "Status for 1000",
	}

	jsonStrings = map[int]string{
		1: `{"a":{"b":{"c":{"d":{"e":"f"},"g":{"h":"i"},"j":{}}},"k":{"m":{"n":"o"}},"p":{"q":"r"},"s":{"t":{"u":"v"}}},"w":"x"}`,
	}
)

var testCasesForTypedef = []struct {
	name                string
	messageNumber       int
	details             []interface{}
	expectedMessageJson string
	expectedText        string
}{
	{
		name:                "typedef-0001",
		messageNumber:       1,
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"TRACE","id":"senzing-99990001","text":"TRACE: Bob works with Jane","status":"Status for 0001","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedText:        "TRACE: Bob works with Jane",
	},
	{
		name:                "typedef-1001",
		messageNumber:       1001,
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"DEBUG","id":"senzing-99991001","text":"DEBUG: Bob works with Jane","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedText:        "DEBUG: Bob works with Jane",
	},
	{
		name:                "typedef-2001",
		messageNumber:       2001,
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"INFO","id":"senzing-99992001","text":"INFO: Bob works with Jane","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedText:        "INFO: Bob works with Jane",
	},
	{
		name:                "typedef-2002",
		messageNumber:       2002,
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"INFO","id":"senzing-99992002","text":"Bob","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedText:        "Bob",
	},
	{
		name:                "typedef-2003",
		messageNumber:       2002,
		details:             []interface{}{jsonStrings[1], getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"INFO","id":"senzing-99992002","text":"{\"a\":{\"b\":{\"c\":{\"d\":{\"e\":\"f\"},\"g\":{\"h\":\"i\"},\"j\":{}}},\"k\":{\"m\":{\"n\":\"o\"}},\"p\":{\"q\":\"r\"},\"s\":{\"t\":{\"u\":\"v\"}}},\"w\":\"x\"}","details":[{"position":1,"value":"{\"a\":{\"b\":{\"c\":{\"d\":{\"e\":\"f\"},\"g\":{\"h\":\"i\"},\"j\":{}}},\"k\":{\"m\":{\"n\":\"o\"}},\"p\":{\"q\":\"r\"},\"s\":{\"t\":{\"u\":\"v\"}}},\"w\":\"x\"}","valueRaw":{"a":{"b":{"c":{"d":{"e":"f"},"g":{"h":"i"},"j":{}}},"k":{"m":{"n":"o"}},"p":{"q":"r"},"s":{"t":{"u":"v"}}},"w":"x"}},{"key":"text","position":2,"value":"{\"a\":{\"b\":{\"c\":{\"d\":{\"e\":\"f\"},\"g\":{\"h\":\"i\"},\"j\":{}}},\"k\":{\"m\":{\"n\":\"o\"}},\"p\":{\"q\":\"r\"},\"s\":{\"t\":{\"u\":\"v\"}}},\"w\":\"x\"}","valueRaw":{"a":{"b":{"c":{"d":{"e":"f"},"g":{"h":"i"},"j":{}}},"k":{"m":{"n":"o"}},"p":{"q":"r"},"s":{"t":{"u":"v"}}},"w":"x"}}]}`,
		expectedText:        `{"a":{"b":{"c":{"d":{"e":"f"},"g":{"h":"i"},"j":{}}},"k":{"m":{"n":"o"}},"p":{"q":"r"},"s":{"t":{"u":"v"}}},"w":"x"}`,
	},
	{
		name:                "typedef-2004",
		messageNumber:       2003,
		details:             []interface{}{jsonStrings[1], getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"INFO","id":"senzing-99992003","text":"JSON: {\"a\":{\"b\":{\"c\":{\"d\":{\"e\":\"f\"},\"g\":{\"h\":\"i\"},\"j\":{}}},\"k\":{\"m\":{\"n\":\"o\"}},\"p\":{\"q\":\"r\"},\"s\":{\"t\":{\"u\":\"v\"}}},\"w\":\"x\"}","details":[{"position":1,"value":"{\"a\":{\"b\":{\"c\":{\"d\":{\"e\":\"f\"},\"g\":{\"h\":\"i\"},\"j\":{}}},\"k\":{\"m\":{\"n\":\"o\"}},\"p\":{\"q\":\"r\"},\"s\":{\"t\":{\"u\":\"v\"}}},\"w\":\"x\"}","valueRaw":{"a":{"b":{"c":{"d":{"e":"f"},"g":{"h":"i"},"j":{}}},"k":{"m":{"n":"o"}},"p":{"q":"r"},"s":{"t":{"u":"v"}}},"w":"x"}}]}`,
		expectedText:        `JSON: {"a":{"b":{"c":{"d":{"e":"f"},"g":{"h":"i"},"j":{}}},"k":{"m":{"n":"o"}},"p":{"q":"r"},"s":{"t":{"u":"v"}}},"w":"x"}`,
	},
	{
		name:                "typedef-3001",
		messageNumber:       3001,
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"WARN","id":"senzing-99993001","text":"WARN: Bob works with Jane","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedText:        "WARN: Bob works with Jane",
	},
	{
		name:                "typedef-4001",
		messageNumber:       4001,
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"ERROR","id":"senzing-99994001","text":"ERROR: Bob works with Jane","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedText:        "ERROR: Bob works with Jane",
	},
	{
		name:                "typedef-5001",
		messageNumber:       5001,
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"FATAL","id":"senzing-99995001","text":"FATAL: Bob works with Jane","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedText:        "FATAL: Bob works with Jane",
	},
	{
		name:                "typedef-6001",
		messageNumber:       6001,
		details:             []interface{}{"Bob", "Jane", getTimestamp()},
		expectedMessageJson: `{"time":"2000-01-01 00:00:00 +0000 UTC","level":"PANIC","id":"senzing-99996001","text":"PANIC: Bob works with Jane","details":[{"position":1,"value":"Bob"},{"position":2,"value":"Jane"}]}`,
		expectedText:        "PANIC: Bob works with Jane",
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
	messengerOptions := []interface{}{
		&messenger.OptionCallerSkip{Value: callerSkip},
		&messenger.OptionIdMessages{Value: idMessages},
		&messenger.OptionIdStatuses{Value: idStatuses},
		&messenger.OptionMessageIdTemplate{Value: messageIdTemplate},
		&messenger.OptionSenzingComponentId{Value: componentIdentifier},
	}
	logger, err = messenger.New(messengerOptions...)
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Internal functions - names begin with lowercase letter
// ----------------------------------------------------------------------------

func testError(test *testing.T, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}
func printRequest(test *testing.T, request string) {
	if false {
		test.Log(request)
	}
}

func printResult(test *testing.T, result *SenzingMessage) {
	if isDebug {
		test.Logf("-------- \n")
		test.Logf("      ID: %s\n", result.ID)
		test.Logf("   Level: %s\n", result.Level)
		test.Logf("    Time: %s\n", result.Time)
		test.Logf("    Text: %s\n", result.Text)
		test.Logf("  Status: %s\n", result.Status)
		test.Logf("Duration: %d\n", result.Duration)
		test.Logf("Location: %s\n", result.Location)

		if len(result.Details) > 0 {
			test.Logf(" Details:\n")
			for _, detail := range result.Details {
				test.Logf("         Position: %d\n", detail.Position)
				test.Logf("              Key: %s\n", detail.Key)
				test.Logf("   ValueAsString: %s\n", detail.ValueAsString)
				test.Logf("           Value: %s\n", detail.Value)
			}
		}

		if len(result.Errors) > 0 {
			test.Logf(" Errors:\n")
			for _, detail := range result.Errors {
				test.Logf("           Error: %s\n", detail)
			}
		}

		test.Logf("\n")
	}
}

func getTimestamp() *messenger.MessageTime {
	return &messenger.MessageTime{
		Value: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
}

// ----------------------------------------------------------------------------
// --- Test cases
// ----------------------------------------------------------------------------

func TestMessengerImpl_NewJson(test *testing.T) {
	for _, testCase := range testCasesForTypedef {
		if len(testCase.expectedMessageJson) > 0 {
			test.Run(testCase.name+"-NewJson", func(test *testing.T) {
				jsonString := logger.NewJson(testCase.messageNumber, testCase.details...)
				result := &SenzingMessage{}
				err := json.Unmarshal([]byte(jsonString), result)
				testError(test, err)
				printResult(test, result)
				assert.Equal(test, testCase.expectedMessageJson, jsonString, testCase.name)
				assert.Equal(test, testCase.expectedText, result.Text, testCase.name)
			})
		}
	}
}

func TestSenzingMessageSimple(test *testing.T) {
	jsonString := logger.NewJson(2, "Bob", "Mary")
	result := &SenzingMessage{}
	err := json.Unmarshal([]byte(jsonString), result)
	testError(test, err)
	printResult(test, result)
	assert.Equal(test, "The person is Bob", result.Text)
	assert.Equal(test, "Bob", result.Details[0].Value)
}

func TestSenzingMessageErrors(test *testing.T) {
	err1 := errors.New("example error #1")
	err2 := errors.New("example error #2")
	jsonString := logger.NewJson(2, "Bob", "Mary", err1, err2)
	result := &SenzingMessage{}
	err := json.Unmarshal([]byte(jsonString), result)
	testError(test, err)
	printResult(test, result)
	assert.Equal(test, "The person is Bob", result.Text)
	assert.Equal(test, "Bob", result.Details[0].Value)
}

func TestSenzingMessageMap(test *testing.T) {
	aMap := map[string]string{
		"BobKey":  "BobValue",
		"MaryKey": "MaryValue",
	}
	jsonString := logger.NewJson(2, "Bob", "Mary", aMap)
	result := &SenzingMessage{}
	err := json.Unmarshal([]byte(jsonString), result)
	testError(test, err)
	printResult(test, result)
	assert.Equal(test, "The person is Bob", result.Text)
	assert.Equal(test, "Bob", result.Details[0].Value)
}

// func TestSenzingMessageJson3(test *testing.T) {
// 	jsonDetail := `{"a":{"b":{"c":{"d":{"e":"f"},"g":{"h":"i"},"j":{}}},"k":{"m":{"n":"o"}},"p":{"q":"r"},"s":{"t":{"u":"v"}}},"w":"x"}`
// 	jsonString := logger.NewJson(2, "Bob", "Mary", jsonDetail)
// 	printRequest(test, jsonString)
// 	result := &SenzingMessage{}
// 	err := json.Unmarshal([]byte(jsonString), result)
// 	testError(test, err)
// 	printResult(test, result)
// 	assert.Equal(test, "The person is Bob", result.Text)
// 	assert.Equal(test, "Bob", result.Details[0].Value)

// 	aMap, ok := result.Details[2].Value.(map[string]interface{})
// 	if ok {
// 		assert.Equal(test, "x", aMap["w"])
// 	} else {
// 		assert.FailNow(test, "Map not OK")
// 	}
// }
