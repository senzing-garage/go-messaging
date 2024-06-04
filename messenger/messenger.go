package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// SimpleMessenger is an type-struct for an implementation of the MessengerInterface.
type SimpleMessenger struct {
	idMessages          map[int]string // Map message numbers to text format strings
	idStatuses          map[int]string
	messageIDTemplate   string // A string template for fmt.Sprinf()
	callerSkip          int    // Levels of code nexting to skip when calculation location
	sortedIDLevelRanges []int  // The keys of IdLevelRanges in sorted order.
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Strip \t and \n from string.
func cleanTabsAndNewlines(unknownString string) string {
	result := unknownString
	result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, "\t", "")
	return result
}

func cleanErrorString(err error) string {
	return cleanTabsAndNewlines(err.Error())
}

// Determine if string is syntactically JSON.
func isJSON(unknownString string) bool {
	unknownStringUnescaped := cleanTabsAndNewlines(unknownString)
	var jsonRawMessage json.RawMessage
	return json.Unmarshal([]byte(unknownStringUnescaped), &jsonRawMessage) == nil
}

// Cast JSON string into an interface{}.
func jsonAsInterface(unknownString string) interface{} {
	unknownStringUnescaped := cleanTabsAndNewlines(unknownString)
	var jsonRawMessage json.RawMessage
	err := json.Unmarshal([]byte(unknownStringUnescaped), &jsonRawMessage)
	if err != nil {
		panic(err)
	}
	return jsonRawMessage
}

// Cast an interface{} into a string.
func interfaceAsString(unknown interface{}) string {
	// See https://pkg.go.dev/fmt for format strings.
	var result string
	switch value := unknown.(type) {
	case nil:
		result = "<nil>"
	case string:
		if isJSON(value) {
			result = cleanTabsAndNewlines(value)
		} else {
			result = value
		}
	case int:
		result = fmt.Sprintf("%d", value)
	case float64:
		result = fmt.Sprintf("%g", value)
	case bool:
		result = fmt.Sprintf("%t", value)
	case error:
		result = cleanErrorString(value)
	default:
		result = fmt.Sprintf("%#v", unknown)
	}
	return result
}

// Walk through the details to improve their future JSON representation.
func messageDetails(details ...interface{}) []Detail {
	result := []Detail{}

	// Process different types of details.

	for index, value := range details {
		detail := Detail{}
		detail.Position = int32(index + 1)
		switch typedValue := value.(type) {
		case nil:
			detail.Type = "nil"
			result = append(result, detail)
		case int:
			detail.Type = "integer"
			detail.Value = interfaceAsString(typedValue)
			detail.ValueRaw = typedValue
			result = append(result, detail)
		case float64:
			detail.Type = "float"
			detail.Value = interfaceAsString(typedValue)
			detail.ValueRaw = typedValue
			result = append(result, detail)
		case string:
			detail.Type = "string"
			detail.Value = typedValue
			if isJSON(typedValue) {
				detail.ValueRaw = jsonAsInterface(typedValue)
			}
			result = append(result, detail)
		case bool:
			detail.Type = "boolean"
			detail.Value = interfaceAsString(typedValue)
			detail.ValueRaw = typedValue
			result = append(result, detail)
		case error:
			detail.Type = "error"
			detail.Value = cleanErrorString(typedValue)
			if isJSON(detail.Value) {
				detail.ValueRaw = jsonAsInterface(detail.Value)
			}
			result = append(result, detail)
		case map[string]string:
			for mapIndex, mapValue := range typedValue {
				detail := Detail{}
				detail.Position = int32(index + 1)
				detail.Key = mapIndex
				detail.Type = "map[string]string"
				detail.Value = interfaceAsString(mapValue)
				if isJSON(detail.Value) {
					detail.ValueRaw = jsonAsInterface(detail.Value)
				}
				result = append(result, detail)
			}
		default:
			detail.Type = fmt.Sprintf("%+v", reflect.TypeOf(value))
			detail.Value = interfaceAsString(typedValue)
			detail.ValueRaw = typedValue
			result = append(result, detail)
		}
	}

	return result
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

// Create a slice of ["key1", value1, "key2", value2, ...] which has oscillating
// key and values in the slice.
func (messenger *SimpleMessenger) getKeyValuePairs(appMessageFormat *MessageFormat, keys []string) []interface{} {
	var result []interface{}
	keyValueMap := map[string]interface{}{
		"time":     appMessageFormat.Time,
		"level":    appMessageFormat.Level,
		"id":       appMessageFormat.ID,
		"status":   appMessageFormat.Status,
		"duration": appMessageFormat.Duration,
		"location": appMessageFormat.Location,
		"errors":   appMessageFormat.Errors,
		"details":  appMessageFormat.Details,
	}

	// In key order, append values to result.

	for _, key := range keys {
		value, ok := keyValueMap[key]
		if !ok {
			continue
		}
		switch typedValue := value.(type) {
		case string:
			if typedValue != "" {
				result = append(result, key, value)
			}
		case int64:
			if typedValue != 0 {
				result = append(result, key, value)
			}
		default:
			if typedValue != nil {
				result = append(result, key, value)
			}
		}
	}
	return result
}

// Given a message number, figure out the Level (TRACE, DEBUG, ..., FATAL, PANIC)
func (messenger *SimpleMessenger) getLevel(messageNumber int) string {
	sortedMessageLevelKeys := messenger.getSortedIDLevelRanges(IDLevelRangesAsString)
	for _, messageLevelKey := range sortedMessageLevelKeys {
		if messageNumber >= messageLevelKey {
			return IDLevelRangesAsString[messageLevelKey]
		}
	}
	return "UNKNOWN"
}

// Since a map[int]any is not guaranteed to be in order, return an ordered slice of int.
func (messenger *SimpleMessenger) getSortedIDLevelRanges(idLevelRanges map[int]string) []int {
	if messenger.sortedIDLevelRanges == nil {
		messenger.sortedIDLevelRanges = make([]int, 0, len(idLevelRanges))
		for key := range idLevelRanges {
			messenger.sortedIDLevelRanges = append(messenger.sortedIDLevelRanges, key)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(messenger.sortedIDLevelRanges)))
	}
	return messenger.sortedIDLevelRanges
}

// Create a populated MessageFormat.
func (messenger *SimpleMessenger) populateStructure(messageNumber int, details ...interface{}) *MessageFormat {

	var (
		callerSkip int
		duration   int64
		errorList  []interface{}
		level      string
		location   string
		status     string
		text       string
	)

	// Calculate fields.

	timeNow := time.Now().UTC().Format(time.RFC3339Nano)
	callerSkip = messenger.callerSkip
	level = messenger.getLevel(messageNumber)
	id := fmt.Sprintf(messenger.messageIDTemplate, messageNumber)
	statusCandidate, ok := messenger.idStatuses[messageNumber]
	if ok {
		status = statusCandidate
	}

	// Construct "text".

	textTemplate, ok := messenger.idMessages[messageNumber]
	if ok {
		textRaw := fmt.Sprintf(textTemplate, details...)
		text = strings.Split(textRaw, "%!(")[0]
		if isJSON(text) {
			textThing := map[string]string{
				"text": fmt.Sprintf("%+v", text),
			}
			details = append(details, textThing)
		}
	}

	// TODO: Find status in underlying error.
	// See https://github.com/senzing-garage/go-logging/blob/48487ee9793e94dac4a3e047635ffd40ff9c454e/messagestatus/messagestatus_senzingapi.go#L29-L59

	// Process Options found in details and filter them out of details.

	filteredDetails := []interface{}{}
	for _, value := range details {
		switch typedValue := value.(type) {
		case *MessageDuration:
			duration = typedValue.Value
		case *MessageID:
			id = typedValue.Value
		case *MessageLevel:
			level = typedValue.Value
		case *MessageLocation:
			location = typedValue.Value
		case *MessageStatus:
			status = typedValue.Value
		case *MessageText:
			text = typedValue.Value
		case *MessageTime:
			timeNow = typedValue.Value.Format(time.RFC3339Nano)
		case *OptionCallerSkip:
			callerSkip = typedValue.Value
		case error:
			errorList = append(errorList, cleanErrorString(typedValue))
			filteredDetails = append(filteredDetails, typedValue)

			// TODO:
			// messageSplits := strings.Split(errorMessage, "|")
			// for key, value := range SenzingApiErrorsMap {
			// 	if messageSplits[0] == key {
			// 		errorMessageList = append(errorMessageList, value)
			// 	}
			// }
		case time.Duration:
			duration = typedValue.Nanoseconds()
		default:
			filteredDetails = append(filteredDetails, typedValue)
		}
	}

	// Calculate field - location.
	// See https://pkg.go.dev/runtime#Caller

	if callerSkip > 0 {
		pc, file, line, ok := runtime.Caller(callerSkip)
		if ok {
			callingFunction := runtime.FuncForPC(pc)
			runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
			functionName := runtimeFunc.ReplaceAllString(callingFunction.Name(), "$1")
			filename := filepath.Base(file)
			location = fmt.Sprintf("In %s() at %s:%d", functionName, filename, line)
		}
	}

	// Compose result.

	result := &MessageFormat{
		Time:     timeNow,
		Level:    level,
		ID:       id,
		Text:     text,
		Status:   status,
		Duration: duration,
		Location: location,
	}
	if len(errorList) > 0 {
		result.Errors = errorList
	}
	if len(filteredDetails) > 0 {
		result.Details = messageDetails(filteredDetails...)
	}
	return result
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The NewJSON method return a JSON string with the elements of the message.

Input
  - messageNumber: A message identifier which indexes into "idMessages".
  - details: Variadic arguments of any type to be added to the message.

Output
  - A JSON string representing the details formatted by the template identified by the messageNumber.
*/
func (messenger *SimpleMessenger) NewJSON(messageNumber int, details ...interface{}) string {
	messageFormat := messenger.populateStructure(messageNumber, details...)

	// Construct return value.

	// Would love to do it this way, but HTML escaping happens.
	// Reported in https://github.com/golang/go/issues/56630
	// result, _ := json.Marshal(messageBuilder)
	// return string(result), err

	// Work-around.

	var resultBytes bytes.Buffer
	enc := json.NewEncoder(&resultBytes)
	enc.SetEscapeHTML(false)
	err := enc.Encode(messageFormat)
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(resultBytes.String())
}

/*
The NewSlog method returns a message and list of Key-Value pairs string with the elements of the message.
A convenience method for NewSlogLevel(), but without slog.Level returned.

Input
  - messageNumber: A message identifier which indexes into "idMessages".
  - details: Variadic arguments of any type to be added to the message.

Output
  - A text message
  - A slice of oscillating key-value pairs.
*/
func (messenger *SimpleMessenger) NewSlog(messageNumber int, details ...interface{}) (string, []interface{}) {
	message, _, keyValuePairs := messenger.NewSlogLevel(messageNumber, details...)
	return message, keyValuePairs
}

/*
The NewSlogLevel method returns a message. an slog level, and a list of Key-Value pairs string with the elements of the message.

Input
  - messageNumber: A message identifier which indexes into "idMessages".
  - details: Variadic arguments of any type to be added to the message.

Output
  - A text message
  - A message level
  - A slice of oscillating key-value pairs.
*/
func (messenger *SimpleMessenger) NewSlogLevel(messageNumber int, details ...interface{}) (string, slog.Level, []interface{}) {
	messageFormat := messenger.populateStructure(messageNumber, details...)

	// Create a text message.

	message := messageFormat.Text

	// Create a slog.Level message level

	slogLevel, ok := TextToLevelMap[messageFormat.Level]
	if !ok {
		slogLevel = LevelPanicSlog
	}

	// Create a slice of oscillating key-value pairs.

	keys := []string{
		"id",
		"status",
		"duration",
		"location",
		"errors",
		"details",
	}
	keyValuePairs := messenger.getKeyValuePairs(messageFormat, keys)
	return message, slogLevel, keyValuePairs
}
