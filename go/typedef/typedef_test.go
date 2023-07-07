package typedef

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-messaging/messenger"
	"github.com/stretchr/testify/assert"
)

var (
	logger messenger.MessengerInterface

	callerSkip          = 1
	componentIdentifier = 9999
	isDebug             = true
	messageIdTemplate   = "test-%04d"

	idMessages = map[int]string{
		1: "The person is %s",
	}

	idStatuses = map[int]string{
		0001: "Status for 0001",
		1000: "Status for 1000",
	}
)

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

// ----------------------------------------------------------------------------
// --- Test cases
// ----------------------------------------------------------------------------

func TestSenzingMessageSimple(test *testing.T) {
	jsonString := logger.NewJson(0001, "Bob", "Mary")
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
	jsonString := logger.NewJson(0001, "Bob", "Mary", err1, err2)
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
	jsonString := logger.NewJson(0001, "Bob", "Mary", aMap)
	result := &SenzingMessage{}
	err := json.Unmarshal([]byte(jsonString), result)
	testError(test, err)
	printResult(test, result)
	assert.Equal(test, "The person is Bob", result.Text)
	assert.Equal(test, "Bob", result.Details[0].Value)
}

func TestSenzingMessageJson(test *testing.T) {
	jsonDetail := `{"a":{"b":{"c":{"d":{"e":"f"},"g":{"h":"i"},"j":{}}},"k":{"m":{"n":"o"}},"p":{"q":"r"},"s":{"t":{"u":"v"}}},"w":"x"}`
	jsonString := logger.NewJson(0001, "Bob", "Mary", jsonDetail)
	printRequest(test, jsonString)
	result := &SenzingMessage{}
	err := json.Unmarshal([]byte(jsonString), result)
	testError(test, err)
	printResult(test, result)
	assert.Equal(test, "The person is Bob", result.Text)
	assert.Equal(test, "Bob", result.Details[0].Value)

	aMap, ok := result.Details[2].Value.(map[string]interface{})
	if ok {
		assert.Equal(test, "x", aMap["w"])
	} else {
		assert.FailNow(test, "Map not OK")
	}

}
