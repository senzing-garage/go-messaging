package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strconv"
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

type theFields struct {
	code            string
	duration        int64
	id              string
	level           string
	location        string
	reason          string
	status          string
	text            string
	callerSkip      int
	errorList       []interface{}
	timeNow         string
	filteredDetails []interface{}
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The NewError method returns an error with a JSON string message.

Input
  - messageNumber: A message identifier which indexes into "idMessages".
  - details: Variadic arguments of any type to be added to the message.

Output
  - An error with a JSON string representing the details formatted by the template identified by the messageNumber.
*/
func (messenger *BasicMessenger) NewError(messageNumber int, details ...interface{}) error {
	return errors.New(messenger.NewJSON(messageNumber, details...)) //nolint
}

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

	result := strings.TrimSpace(resultBytes.String())

	return result
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
func (messenger *BasicMessenger) NewSlogLevel(
	messageNumber int,
	details ...interface{},
) (string, slog.Level, []interface{}) {
	populateDetails := []interface{}{}
	populateDetails = append(populateDetails, details...)
	populateDetails = append(populateDetails, OptionMessageField{Value: "level"})
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
// Private methods
// ----------------------------------------------------------------------------

func (messenger *BasicMessenger) findMessageFields(details ...interface{}) []string {
	var (
		result   []string
		appendix = []string{}
	)

	senzingMessageFields := strings.TrimSpace(strings.ToLower(os.Getenv("SENZING_MESSAGE_FIELDS")))

	if messenger.messageFields == nil {
		messenger.populateMessageFields(senzingMessageFields)
	}

	switch {
	case len(senzingMessageFields) == 0:
		result = messenger.messageFields
	case senzingMessageFields == "all":
		result = AllMessageFields
	default:
		result = []string{}

		messageSplits := strings.Split(senzingMessageFields, ",")
		for _, value := range messageSplits {
			valueTrimmed := strings.TrimSpace(value)
			if slices.Contains(AllMessageFields, valueTrimmed) {
				result = append(result, valueTrimmed)
			}
		}
	}

	for _, value := range details {
		switch typedValue := value.(type) {
		case OptionMessageFields:
			result = typedValue.Value
		case OptionMessageField:
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

// Given a message number, figure out the Level (TRACE, DEBUG, ..., FATAL, PANIC).
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

func (messenger *BasicMessenger) populateMessageFields(senzingMessageFields string) {
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
	actualFields := &theFields{}

	// Calculate fields.

	actualFields.timeNow = time.Now().UTC().Format(time.RFC3339Nano)
	actualFields.callerSkip = messenger.callerSkip
	actualFields.level = messenger.getLevel(messageNumber)
	actualFields.id = fmt.Sprintf(messenger.messageIDTemplate, messageNumber)
	statusCandidate, isOK := messenger.idStatuses[messageNumber]

	if isOK {
		actualFields.status = statusCandidate
	}

	// Construct "text".

	textTemplate, isOK := messenger.idMessages[messageNumber]
	if isOK {
		textRaw := fmt.Sprintf(textTemplate, details...)

		actualFields.text = strings.Split(textRaw, "%!(")[0]
		if isJSON(actualFields.text) {
			textThing := map[string]string{
				"text": fmt.Sprintf("%+v", actualFields.text),
			}
			details = append(details, textThing)
		}
	}

	parseDetails(actualFields, details)

	// Calculate field - location.
	// See https://pkg.go.dev/runtime#Caller

	if actualFields.callerSkip > 0 {
		pc, file, line, ok := runtime.Caller(actualFields.callerSkip)
		if ok {
			callingFunction := runtime.FuncForPC(pc)
			runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
			functionName := runtimeFunc.ReplaceAllString(callingFunction.Name(), "$1")
			filename := filepath.Base(file)
			actualFields.location = fmt.Sprintf("In %s() at %s:%d", functionName, filename, line)
		}
	}

	// Determine fields to print.

	messageFields := messenger.findMessageFields(details...)

	return populateMessageFormat(actualFields, messageFields)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func createBooleanDetail(position int32, value bool) Detail {
	return Detail{
		Position: position,
		Type:     "boolean",
		Value:    interfaceAsString(value),
		ValueRaw: value,
	}
}

func createErrorDetail(position int32, value error) Detail {
	result := Detail{
		Position: position,
		Type:     "error",
		Value:    cleanErrorString(value),
	}
	if isJSON(result.Value) {
		result.ValueRaw = jsonAsInterface(result.Value)
	}

	return result
}

func createFloat64Detail(position int32, value float64) Detail {
	return Detail{
		Position: position,
		Type:     "float",
		Value:    interfaceAsString(value),
		ValueRaw: value,
	}
}

func createIntDetail(position int32, value int) Detail {
	return Detail{
		Position: position,
		Type:     "integer",
		Value:    interfaceAsString(value),
		ValueRaw: value,
	}
}

func createNilDetail(position int32) Detail {
	return Detail{
		Position: position,
		Type:     "nil",
	}
}

func createStringDetail(position int32, value string) Detail {
	result := Detail{
		Position: position,
		Type:     "string",
		Value:    value,
	}
	if isJSON(value) {
		result.ValueRaw = jsonAsInterface(value)
	}

	return result
}

func parseDetails(actualFields *theFields, details []interface{}) {
	for _, value := range details {
		switch typedValue := value.(type) {
		case MessageCode:
			actualFields.code = typedValue.Value
		case MessageDuration:
			actualFields.duration = typedValue.Value
		case MessageID:
			actualFields.id = typedValue.Value
		case MessageLevel:
			actualFields.level = typedValue.Value
		case MessageLocation:
			actualFields.location = typedValue.Value
		case MessageReason:
			actualFields.reason = typedValue.Value
		case MessageStatus:
			actualFields.status = typedValue.Value
		case MessageText:
			actualFields.text = typedValue.Value
		case MessageTime:
			actualFields.timeNow = typedValue.Value.Format(time.RFC3339Nano)
		case OptionCallerSkip:
			actualFields.callerSkip = typedValue.Value
		case error:
			actualFields.errorList = append(actualFields.errorList, cleanErrorString(typedValue))
			actualFields.filteredDetails = append(actualFields.filteredDetails, typedValue)
		case time.Duration:
			actualFields.duration = typedValue.Nanoseconds()
		default:
			actualFields.filteredDetails = append(actualFields.filteredDetails, typedValue)
		}
	}
}

func populateMessageFormat(actualFields *theFields, messageFields []string) *MessageFormat {
	result := &MessageFormat{}
	if slices.Contains(messageFields, "code") {
		result.Code = actualFields.code
	}

	if slices.Contains(messageFields, "details") {
		if len(actualFields.filteredDetails) > 0 {
			result.Details = messageDetails(actualFields.filteredDetails...)
		}
	}

	if slices.Contains(messageFields, "duration") {
		result.Duration = actualFields.duration
	}

	if slices.Contains(messageFields, "errors") {
		if len(actualFields.errorList) > 0 {
			result.Errors = actualFields.errorList
		}
	}

	if slices.Contains(messageFields, "id") {
		result.ID = actualFields.id
	}

	if slices.Contains(messageFields, "level") {
		result.Level = actualFields.level
	}

	if slices.Contains(messageFields, "location") {
		result.Location = actualFields.location
	}

	if slices.Contains(messageFields, "reason") {
		result.Reason = actualFields.reason
	}

	if slices.Contains(messageFields, "status") {
		result.Status = actualFields.status
	}

	if slices.Contains(messageFields, "text") {
		result.Text = actualFields.text
	}

	if slices.Contains(messageFields, "time") {
		result.Time = actualFields.timeNow
	}

	return result
}

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
		result = strconv.Itoa(value)
	case float64:
		result = fmt.Sprintf("%g", value)
	case bool:
		result = strconv.FormatBool(value)
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
		detailPosition := int32(index + 1) //nolint:gosec
		switch typedValue := value.(type) {
		case nil:
			result = append(result, createNilDetail(detailPosition))
		case int:
			result = append(result, createIntDetail(detailPosition, typedValue))
		case float64:
			result = append(result, createFloat64Detail(detailPosition, typedValue))
		case string:
			result = append(result, createStringDetail(detailPosition, typedValue))
		case bool:
			result = append(result, createBooleanDetail(detailPosition, typedValue))
		case error:
			result = append(result, createErrorDetail(detailPosition, typedValue))
		case map[string]string:
			for mapIndex, mapValue := range typedValue {
				detail := Detail{
					Key:      mapIndex,
					Position: detailPosition,
					Type:     "map[string]string",
					Value:    interfaceAsString(mapValue),
				}
				if isJSON(detail.Value) {
					detail.ValueRaw = jsonAsInterface(detail.Value)
				}

				result = append(result, detail)
			}
		case OptionMessageField:
			// Do nothing.
		case OptionMessageFields:
			// Do nothing.
		default:
			detail := Detail{
				Position: detailPosition,
				Type:     fmt.Sprintf("%+v", reflect.TypeOf(value)),
				Value:    interfaceAsString(typedValue),
				ValueRaw: typedValue,
			}
			result = append(result, detail)
		}
	}

	return result
}
