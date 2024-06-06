package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicMessenger is an type-struct for an implementation of the MessengerInterface.
type BasicMessenger struct {
	callerSkip          int            // Levels of code nexting to skip when calculation location
	idMessages          map[int]string // Map message numbers to text format strings
	idStatuses          map[int]string
	messageFields       []string
	messageIDTemplate   string // A string template for fmt.Sprinf()
	sortedIDLevelRanges []int  // The keys of IdLevelRanges in sorted order.
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
func (messenger *BasicMessenger) NewJSON(messageNumber int, details ...interface{}) string {
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
func (messenger *BasicMessenger) NewSlog(messageNumber int, details ...interface{}) (string, []interface{}) {
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
func (messenger *BasicMessenger) NewSlogLevel(messageNumber int, details ...interface{}) (string, slog.Level, []interface{}) {

	// Add "level" to details temporarily.

	populateDetails := []interface{}{}
	populateDetails = append(populateDetails, details...)
	populateDetails = append(populateDetails, &OptionMessageField{Value: "level"})
	messageFormat := messenger.populateStructure(messageNumber, populateDetails...)

	// Create a text message.

	message := messageFormat.Text

	// Create a slog.Level message level

	slogLevel, ok := TextToLevelMap[messageFormat.Level]
	if !ok {
		slogLevel = LevelPanicSlog
	}

	// Create a slice of oscillating key-value pairs.

	messageFields := messenger.findMessageFields(details...)
	keyValuePairs := messenger.getKeyValuePairs(messageFormat, messageFields)
	return message, slogLevel, keyValuePairs
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
		case *OptionMessageField:
			// Do nothing.
		case *OptionMessageFields:
			// Do nothing.
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

func (messenger *BasicMessenger) findMessageFields(details ...interface{}) []string {
	appendix := []string{}
	if messenger.messageFields == nil {
		messenger.populateMessageFields()
	}
	result := messenger.messageFields
	for _, value := range details {
		switch typedValue := value.(type) {
		case *OptionMessageFields:
			result = typedValue.Value
		case *OptionMessageField:
			appendix = append(appendix, typedValue.Value)
		default:
		}
	}
	result = append(result, appendix...)
	return result
}

// Create a slice of ["key1", value1, "key2", value2, ...] which has oscillating
// key and values in the slice.
func (messenger *BasicMessenger) getKeyValuePairs(appMessageFormat *MessageFormat, keys []string) []interface{} {
	var result []interface{}
	keyValueMap := map[string]interface{}{
		"code":     appMessageFormat.Code,
		"details":  appMessageFormat.Details,
		"duration": appMessageFormat.Duration,
		"errors":   appMessageFormat.Errors,
		"id":       appMessageFormat.ID,
		"level":    appMessageFormat.Level,
		"location": appMessageFormat.Location,
		"reason":   appMessageFormat.Reason,
		"status":   appMessageFormat.Status,
		"time":     appMessageFormat.Time,
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
func (messenger *BasicMessenger) getLevel(messageNumber int) string {
	sortedMessageLevelKeys := messenger.getSortedIDLevelRanges(IDLevelRangesAsString)
	for _, messageLevelKey := range sortedMessageLevelKeys {
		if messageNumber >= messageLevelKey {
			return IDLevelRangesAsString[messageLevelKey]
		}
	}
	return "UNKNOWN"
}

// Since a map[int]any is not guaranteed to be in order, return an ordered slice of int.
func (messenger *BasicMessenger) getSortedIDLevelRanges(idLevelRanges map[int]string) []int {
	if messenger.sortedIDLevelRanges == nil {
		messenger.sortedIDLevelRanges = make([]int, 0, len(idLevelRanges))
		for key := range idLevelRanges {
			messenger.sortedIDLevelRanges = append(messenger.sortedIDLevelRanges, key)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(messenger.sortedIDLevelRanges)))
	}
	return messenger.sortedIDLevelRanges
}

func (messenger *BasicMessenger) populateMessageFields() {
	senzingMessageFields := strings.TrimSpace(strings.ToLower(os.Getenv("SENZING_MESSAGE_FIELDS")))
	switch {
	case len(senzingMessageFields) == 0:
		messenger.messageFields = []string{"id", "text"}
	case senzingMessageFields == "all":
		messenger.messageFields = AllMessageFields
	default:
		messenger.messageFields = []string{}
		messageSplits := strings.Split(senzingMessageFields, ",")
		for _, value := range messageSplits {
			valueTrimmed := strings.TrimSpace(value)
			if slices.Contains(AllMessageFields, valueTrimmed) {
				messenger.messageFields = append(messenger.messageFields, valueTrimmed)
			}
		}
	}
}

// Create a populated MessageFormat.
func (messenger *BasicMessenger) populateStructure(messageNumber int, details ...interface{}) *MessageFormat {

	var (
		callerSkip int
		code       string
		duration   int64
		errorList  []interface{}
		level      string
		location   string
		reason     string
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
		case *MessageCode:
			code = typedValue.Value
		case *MessageDuration:
			duration = typedValue.Value
		case *MessageID:
			id = typedValue.Value
		case *MessageLevel:
			level = typedValue.Value
		case *MessageLocation:
			location = typedValue.Value
		case *MessageReason:
			reason = typedValue.Value
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

	// Determine fields to print.

	messageFields := messenger.findMessageFields(details...)

	// Compose result.

	result := &MessageFormat{}

	if slices.Contains(messageFields, "code") {
		result.Code = code
	}
	if slices.Contains(messageFields, "details") {
		if len(filteredDetails) > 0 {
			result.Details = messageDetails(filteredDetails...)
		}
	}
	if slices.Contains(messageFields, "duration") {
		result.Duration = duration
	}
	if slices.Contains(messageFields, "errors") {
		if len(errorList) > 0 {
			result.Errors = errorList
		}
	}
	if slices.Contains(messageFields, "id") {
		result.ID = id
	}
	if slices.Contains(messageFields, "level") {
		result.Level = level
	}
	if slices.Contains(messageFields, "location") {
		result.Location = location
	}
	if slices.Contains(messageFields, "reason") {
		result.Reason = reason
	}
	if slices.Contains(messageFields, "status") {
		result.Status = status
	}
	if slices.Contains(messageFields, "text") {
		result.Text = text
	}
	if slices.Contains(messageFields, "time") {
		result.Time = timeNow
	}
	return result
}
