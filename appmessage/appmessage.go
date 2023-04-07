package appmessage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// AppMessageImpl is an example type-struct.
type AppMessageImpl struct {
	idMessages          map[int]string // Map message numbers to text format strings
	idStatuses          map[int]string
	messageIdTemplate   string // A string template for fmt.Sprinf()
	callerSkip          int    // Levels of code nexting to skip when calculation location
	sortedIdLevelRanges []int  // The keys of IdLevelRanges in sorted order.
}

type messageErrorsSenzing struct {
	Text interface{} `json:"text,omitempty"` // Text returned by error.Error().
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	LevelTraceInt int = -8
	LevelDebugInt int = -4
	LevelInfoInt  int = 0
	LevelWarnInt  int = 4
	LevelErrorInt int = 8
	LevelFatalInt int = 12
	LevelPanicInt int = 16
)

const (
	LevelDebugSlog = slog.LevelDebug
	LevelErrorSlog = slog.LevelError
	LevelFatalSlog = slog.Level(LevelFatalInt)
	LevelInfoSlog  = slog.LevelInfo
	LevelPanicSlog = slog.Level(LevelPanicInt)
	LevelTraceSlog = slog.Level(LevelTraceInt)
	LevelWarnSlog  = slog.LevelWarn
)

// Strings representing the supported logging levels.
const (
	LevelDebugName = "DEBUG"
	LevelErrorName = "ERROR"
	LevelFatalName = "FATAL"
	LevelInfoName  = "INFO"
	LevelPanicName = "PANIC"
	LevelTraceName = "TRACE"
	LevelWarnName  = "WARN"
)

// Map from string representation to Log level as typed integer.
var TextToLevelMap = map[string]slog.Level{
	LevelDebugName: LevelDebugSlog,
	LevelErrorName: LevelErrorSlog,
	LevelFatalName: LevelFatalSlog,
	LevelInfoName:  LevelInfoSlog,
	LevelPanicName: LevelPanicSlog,
	LevelTraceName: LevelTraceSlog,
	LevelWarnName:  LevelWarnSlog,
}

var LevelToTextMap = map[slog.Level]string{
	LevelDebugSlog: LevelDebugName,
	LevelErrorSlog: LevelErrorName,
	LevelFatalSlog: LevelFatalName,
	LevelInfoSlog:  LevelInfoName,
	LevelPanicSlog: LevelPanicName,
	LevelTraceSlog: LevelTraceName,
	LevelWarnSlog:  LevelWarnName,
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func isJson(unknownString string) bool {
	unknownStringUnescaped, err := strconv.Unquote(unknownString)
	if err != nil {
		unknownStringUnescaped = unknownString
	}
	var jsonString json.RawMessage
	return json.Unmarshal([]byte(unknownStringUnescaped), &jsonString) == nil
}

func jsonAsInterface(unknownString string) interface{} {
	unknownStringUnescaped, err := strconv.Unquote(unknownString)
	if err != nil {
		unknownStringUnescaped = unknownString
	}
	var jsonString json.RawMessage
	json.Unmarshal([]byte(unknownStringUnescaped), &jsonString)
	return jsonString
}

func stringify(unknown interface{}) string {
	// See https://pkg.go.dev/fmt for format strings.
	var result string

	switch value := unknown.(type) {
	case nil:
		result = "<nil>"
	case string:
		result = value
	case int:
		result = fmt.Sprintf("%d", value)
	case float64:
		result = fmt.Sprintf("%g", value)
	case bool:
		result = fmt.Sprintf("%t", value)
	case error:
		result = value.Error()
	default:
		// xType := reflect.TypeOf(unknown)
		// xValue := reflect.ValueOf(unknown)
		// result = fmt.Sprintf("(%s)%#v", xType, xValue)
		result = fmt.Sprintf("%#v", unknown)
	}
	return result
}

func messageDetails(details ...interface{}) interface{} {

	result := make(map[string]interface{})

	// Process different types of details.

	for index, value := range details {
		switch typedValue := value.(type) {
		case nil:
			result[strconv.Itoa(index+1)] = "<nil>"

		case int, float64:
			result[strconv.Itoa(index+1)] = typedValue

		case string:
			if isJson(typedValue) {
				result[strconv.Itoa(index+1)] = jsonAsInterface(typedValue)
			} else {
				result[strconv.Itoa(index+1)] = typedValue
			}

		case bool:
			result[strconv.Itoa(index+1)] = fmt.Sprintf("%t", typedValue)

		case error:
			// do nothing.

		case map[string]string:
			for mapIndex, mapValue := range typedValue {
				mapValueAsString := stringify(mapValue)
				if isJson(mapValueAsString) {
					result[mapIndex] = jsonAsInterface(mapValueAsString)
				} else {
					result[mapIndex] = mapValueAsString
				}
			}

		default:
			valueAsString := stringify(typedValue)
			if isJson(valueAsString) {
				result[strconv.Itoa(index+1)] = jsonAsInterface(valueAsString)
			} else {
				result[strconv.Itoa(index+1)] = valueAsString
			}
		}
	}

	if len(result) == 0 {
		result = nil
	}

	return result
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (appMessage *AppMessageImpl) getKeyValuePairs(appMessageFormat *AppMessageFormat, keys []string) []interface{} {
	var result []interface{} = nil
	keyValueMap := map[string]interface{}{
		"date":     appMessageFormat.Date,
		"time":     appMessageFormat.Time,
		"level":    appMessageFormat.Level,
		"id":       appMessageFormat.Id,
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

func (appMessage *AppMessageImpl) getLevel(messageNumber int) string {
	sortedMessageLevelKeys := appMessage.getSortedIdLevelRanges(IdLevelRangesAsString)
	for _, messageLevelKey := range sortedMessageLevelKeys {
		if messageNumber >= messageLevelKey {
			return IdLevelRangesAsString[messageLevelKey]
		}
	}
	return "UNKNOWN"
}

func (appMessage *AppMessageImpl) getSortedIdLevelRanges(idLevelRanges map[int]string) []int {
	if appMessage.sortedIdLevelRanges == nil {
		appMessage.sortedIdLevelRanges = make([]int, 0, len(idLevelRanges))
		for key := range idLevelRanges {
			appMessage.sortedIdLevelRanges = append(appMessage.sortedIdLevelRanges, key)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(appMessage.sortedIdLevelRanges)))
	}
	return appMessage.sortedIdLevelRanges
}

func (appMessage *AppMessageImpl) populateStructure(messageNumber int, details ...interface{}) *AppMessageFormat {
	now := time.Now()

	var (
		callerSkip int           = 0
		duration   int64         = 0
		errorList  []interface{} = nil
		level      string        = ""
		location   string        = ""
		status     string        = ""
		text       interface{}   = nil
	)

	// Calculate - callerskip

	callerSkip = appMessage.callerSkip

	// Calculate field - level.

	level = appMessage.getLevel(messageNumber)

	// Calculate field - id.

	id := fmt.Sprintf(appMessage.messageIdTemplate, messageNumber)

	// Calculate field - text.

	textTemplate, ok := appMessage.idMessages[messageNumber]
	if ok {
		textRaw := fmt.Sprintf(textTemplate, details...)
		text = strings.Split(textRaw, "%!(")[0]
	}

	// Calculate field - status.

	if appMessage.idStatuses != nil {
		statusCandidate, ok := appMessage.idStatuses[messageNumber]
		if ok {
			status = statusCandidate
		}
	}

	// TODO: Find status in underlying error.
	// See https://github.com/Senzing/go-logging/blob/48487ee9793e94dac4a3e047635ffd40ff9c454e/messagestatus/messagestatus_senzingapi.go#L29-L59

	// Process overrides found in details and filter them out of details.

	var filteredDetails []interface{}
	for _, value := range details {
		switch typedValue := value.(type) {
		case *AppMessageCallerSkip:
			callerSkip = typedValue.Value
		case *AppMessageDuration:
			duration = typedValue.Value
		case *AppMessageId:
			id = typedValue.Value
		case *AppMessageLevel:
			level = typedValue.Value
		case *AppMessageLocation:
			location = typedValue.Value
		case *AppMessageStatus:
			status = typedValue.Value
		case *AppMessageText:
			text = typedValue.Value
		case *AppMessageTimestamp:
			now = typedValue.Value
		case error:
			errorMessage := typedValue.Error()
			var priorError interface{}
			if isJson(errorMessage) {
				priorError = &messageErrorsSenzing{
					Text: jsonAsInterface(errorMessage),
				}
			} else {
				priorError = &messageErrorsSenzing{
					Text: errorMessage,
				}
			}
			errorList = append(errorList, priorError)

			// messageSplits := strings.Split(errorMessage, "|")
			// for key, value := range SenzingApiErrorsMap {
			// 	if messageSplits[0] == key {
			// 		errorMessageList = append(errorMessageList, value)
			// 	}
			// }
			filteredDetails = append(filteredDetails, value)
		case time.Duration:
			duration = typedValue.Nanoseconds()
		default:
			filteredDetails = append(filteredDetails, value)
		}
	}

	// Calculate field - date & time.

	date := fmt.Sprintf("%04d-%02d-%02d", now.UTC().Year(), now.UTC().Month(), now.UTC().Day())
	time := fmt.Sprintf("%02d:%02d:%02d.%09d", now.UTC().Hour(), now.UTC().Minute(), now.Second(), now.Nanosecond())

	// Calculate field - location.
	// See https://pkg.go.dev/runtime#Caller

	pc, file, line, ok := runtime.Caller(callerSkip)
	if ok {
		callingFunction := runtime.FuncForPC(pc)
		runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
		functionName := runtimeFunc.ReplaceAllString(callingFunction.Name(), "$1")
		filename := filepath.Base(file)
		location = fmt.Sprintf("In %s() at %s:%d", functionName, filename, line)
	}

	result := &AppMessageFormat{
		Date:     date,
		Time:     time,
		Level:    level,
		Id:       id,
		Text:     text,
		Status:   status,
		Duration: duration,
		Location: location,
	}

	if len(errorList) > 0 {
		result.Errors = errorList
	}

	if len(filteredDetails) > 0 {
		result.Details = messageDetails(filteredDetails)
	}

	return result
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The NewJson method return a JSON string with the elements of the message.

Input
  - messageNumber: A message identifier which indexes into "idMessages".
  - details: Variadic arguments of any type to be added to the message.

Output
  - A JSON string representing the details formatted by the template identified by the messageNumber.
*/
func (appMessage *AppMessageImpl) NewJson(messageNumber int, details ...interface{}) string {

	appMessageFormat := appMessage.populateStructure(messageNumber, details...)

	// Convert to JSON.

	// Would love to do it this way, but HTML escaping happens.
	// Reported in https://github.com/golang/go/issues/56630
	// result, _ := json.Marshal(messageBuilder)
	// return string(result), err

	// Work-around.

	var resultBytes bytes.Buffer
	enc := json.NewEncoder(&resultBytes)
	enc.SetEscapeHTML(false)
	err := enc.Encode(appMessageFormat)
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(resultBytes.String())
}

/*
The NewSlog method returns a message and list of Key-Value pairs string with the elements of the message.

Input
  - messageNumber: A message identifier which indexes into "idMessages".
  - details: Variadic arguments of any type to be added to the message.

Output
  - A text message
  - A slice of oscillating key-value pairs.
*/
func (appMessage *AppMessageImpl) NewSlog(messageNumber int, details ...interface{}) (string, []interface{}) {
	appMessageFormat := appMessage.populateStructure(messageNumber, details...)
	keys := []string{
		"level",
		"id",
		"status",
		"duration",
		"location",
		"errors",
		"details",
	}
	keyValuePairs := appMessage.getKeyValuePairs(appMessageFormat, keys)
	return appMessageFormat.Text.(string), keyValuePairs
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
func (appMessage *AppMessageImpl) NewSlogLevel(messageNumber int, details ...interface{}) (string, slog.Level, []interface{}) {
	appMessageFormat := appMessage.populateStructure(messageNumber, details...)
	keys := []string{
		"id",
		"status",
		"duration",
		"location",
		"errors",
		"details",
	}

	keyValuePairs := appMessage.getKeyValuePairs(appMessageFormat, keys)

	// Create a slog.Level

	slogLevel, ok := TextToLevelMap[appMessageFormat.Level]
	if !ok {
		slogLevel = LevelPanicSlog
	}

	return appMessageFormat.Text.(string), slogLevel, keyValuePairs
}
