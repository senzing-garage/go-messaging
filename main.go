/*
 */
package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/senzing-garage/go-messaging/messenger"
	"github.com/senzing-garage/go-messaging/parser"
)

var (
	err1       = errors.New("example error")
	idMessages = map[int]string{
		1001: "DEBUG: %s works with %s",
		2001: "INFO: %s works with %s",
		3001: "WARN: %s works with %s",
		4001: "ERROR: %s works with %s",
		5001: "FATAL: %s works with %s",
		6001: "PANIC: %s works with %s",
		7001: "Xxxxx: %s works with %s",
	}
	optionIDMessages = messenger.OptionIDMessages{Value: idMessages}
	reason           = messenger.MessageReason{
		Value: "The reason is...",
	}
)

// ----------------------------------------------------------------------------
// main
// ----------------------------------------------------------------------------

func main() {
	aMessenger, err := messenger.New()
	testError(err, "Error1: %s\n")
	displayMessages("Default messages", aMessenger)

	// Example messages with "text" field.

	optionMessageFields := messenger.OptionMessageFields{
		Value: []string{"id", "text"},
	}
	aMessenger, err = messenger.New(optionMessageFields, optionIDMessages)
	testError(err, "Error2: %s\n")
	displayMessages("Messages with 'text' field", aMessenger)

	// Example messages with "reason" field.

	optionMessageFields = messenger.OptionMessageFields{
		Value: []string{"id", "reason"},
	}
	aMessenger, err = messenger.New(optionMessageFields, optionIDMessages)
	testError(err, "Error3: %s\n")
	displayMessages("Messages with 'reason' field", aMessenger)

	// Example messages with "errors" field.

	optionMessageFields = messenger.OptionMessageFields{
		Value: []string{"id", "errors"},
	}
	aMessenger, err = messenger.New(optionMessageFields, optionIDMessages)
	testError(err, "Error4: %s\n")
	displayMessages("Messages with 'errors' field", aMessenger)

	// Example messages with all fields.

	optionMessageFields = messenger.OptionMessageFields{
		Value: messenger.AllMessageFields,
	}
	aMessenger, err = messenger.New(optionMessageFields, optionIDMessages)
	testError(err, "Error5: %s\n")
	displayMessages("Messages with all fields", aMessenger)

	// Example messages with componentID of 9998.

	optionMessageIDTemplate := messenger.OptionMessageIDTemplate{Value: "XYZ9998%04d"}
	aMessenger, err = messenger.New(optionMessageIDTemplate, optionIDMessages)
	testError(err, "Error6: %s\n")
	displayMessages("Messages with componentID of 9998", aMessenger)

	// Example parsed messages.

	printBanner("Parsed messages")

	message1 := `{"time":"2023-07-11T21:05:51.918625982Z","level":"DEBUG","id":"SZSDK99981001","text":"DEBUG: Bob works with Mary","location":"In main() at main.go:101","errors":["error #1","{\"time\": \"2023-04-10T11:00:20.623748617-04:00\",\"level\": \"TRACE\",\"id\": \"SZSDK99990002\",\"text\": \"A fake error\",\"location\": \"In main() at main.go:36\",\"details\": {\"1\": \"Bob\",\"2\": \"Mary\"}}"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Mary"},{"position":3,"value":"error #1"},{"position":4,"value":"\n\t{\n\t\t\"time\": \"2023-04-10T11:00:20.623748617-04:00\",\n\t\t\"level\": \"TRACE\",\n\t\t\"id\": \"SZSDK99990002\",\n\t\t\"text\": \"A fake error\",\n\t\t\"location\": \"In main() at main.go:36\",\"details\": {\"1\": \"Bob\",\"2\": \"Mary\"}}","valueRaw":{"time":"2023-04-10T11:00:20.623748617-04:00","level":"TRACE","id":"SZSDK99990002","text":"A fake error","location":"In main() at main.go:36","details":{"1":"Bob","2":"Mary"}}}]}	`
	parsedMessage1, err := parser.Parse(message1)
	testError(err, "Error8: %s\n")
	outputf("Parse test 1 - ID: %s; Text: %s\n", parsedMessage1.ID, parsedMessage1.Text)

	// Epilog.

	printBanner("Done")
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func displayMessages(banner string, aMessenger messenger.Messenger) {
	printBanner(banner)

	printJSONMessage(aMessenger.NewJSON(0001, "Bob", "Mary"))
	printJSONMessage(aMessenger.NewJSON(2001, "Bob", "Mary"))
	printJSONMessage(aMessenger.NewJSON(3001, "Bob", "Mary", reason, err1))
	printJSONMessage(aMessenger.NewJSON(4001, "Bob", "Mary", reason, err1))
	outputln()

	outputln(aMessenger.NewSlog(1001, "Bob", "Mary"))
	outputln(aMessenger.NewSlog(2001, "Bob", "Mary"))
	outputln(aMessenger.NewSlog(3001, "Bob", "Mary", reason, err1))
	outputln(aMessenger.NewSlog(4001, "Bob", "Mary", reason, err1))
}

func outputf(format string, message ...any) {
	fmt.Printf(format, message...) //nolint
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}

func printBanner(banner string) {
	outputf("\n%s\n", strings.Repeat("-", 80))
	outputf("-- %s\n", banner)
	outputf("%s\n\n", strings.Repeat("-", 80))
}

func printJSONMessage(message string) {
	outputln(message)
	parsedMessage, err := parser.Parse(message)
	testError(err, "Error7: %s\n")
	outputf("    - Parsed as ID: %s", parsedMessage.ID)

	if len(parsedMessage.Text) > 0 {
		outputf("; Text: %s", parsedMessage.Text)
	}

	if len(parsedMessage.Reason) > 0 {
		outputf("; Reason: %s", parsedMessage.Reason)
	}

	outputf("\n")
}

func testError(err error, stringFormat string) {
	if err != nil {
		outputf(stringFormat, err.Error())
	}
}
