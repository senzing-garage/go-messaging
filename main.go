/*
 */
package main

import (
	"errors"
	"fmt"

	"github.com/senzing/go-messaging/messenger"
	"github.com/senzing/go-messaging/parser"
)

var idMessages = map[int]string{
	0001: "TRACE: %s works with %s",
	1001: "DEBUG: %s works with %s",
	2001: "INFO: %s works with %s",
	3001: "WARN: %s works with %s",
	4001: "ERROR: %s works with %s",
	5001: "FATAL: %s works with %s",
	6001: "PANIC: %s works with %s",
	7001: "Xxxxx: %s works with %s",
}

func main() {

	// Create some fake errors.

	err1 := errors.New("error #1")
	err2 := errors.New(`
	{
		"time": "2023-04-10T11:00:20.623748617-04:00",
		"level": "TRACE",
		"id": "senzing-99990002",
		"text": "A fake error",
		"location": "In main() at main.go:36",
		"details": {
			"1": "Bob",
			"2": "Mary"
		}
	}`)

	// ------------------------------------------------------------------------
	// --- Using a bare message generator
	// ------------------------------------------------------------------------

	// Example from README.md

	aMessenger, _ := messenger.New()
	fmt.Printf("%s\n\n", aMessenger.NewJson(0001, "Bob", "Mary"))
	fmt.Println(aMessenger.NewSlog(1001, "Bob", "Mary"))
	fmt.Println()

	// Create a bare message generator.

	messenger1, err := messenger.New()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	// Print some messages.

	fmt.Printf("%s\n\n", messenger1.NewJson(2001, "Bob", "Mary"))
	fmt.Printf("%s\n\n", messenger1.NewJson(3001, "Bob", "Mary", err1, err2))

	// ------------------------------------------------------------------------
	// --- Using a configured message generator
	// ------------------------------------------------------------------------

	// Create a configured message generator.

	optionSenzingComponentId := &messenger.OptionSenzingComponentId{Value: 9998}
	optionCallerSkip := &messenger.OptionCallerSkip{Value: 2}
	optionIdMessages := &messenger.OptionIdMessages{Value: idMessages}
	messenger2, err := messenger.New(optionSenzingComponentId, optionCallerSkip, optionIdMessages)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	// Print some messages.

	fmt.Printf("%s\n\n", messenger2.NewJson(0001, "Bob", "Mary"))
	fmt.Printf("%s\n\n", messenger2.NewJson(1001, "Bob", "Mary", err1, err2))

	// Parse some messages.

	message1 := messenger2.NewJson(0001, "Bob", "Mary")
	parsedMessage1, err := parser.Parse(message1)
	if err != nil {
		fmt.Printf("Error1: %s\n", err.Error())
	}
	fmt.Printf("Parse test 1:  ID: %s; Text: %s\n", parsedMessage1.ID, parsedMessage1.Text)

	message2 := `{"time":"2023-07-11T21:05:51.918625982Z","level":"DEBUG","id":"senzing-99981001","text":"DEBUG: Bob works with Mary","location":"In main() at main.go:101","errors":["error #1","{\"time\": \"2023-04-10T11:00:20.623748617-04:00\",\"level\": \"TRACE\",\"id\": \"senzing-99990002\",\"text\": \"A fake error\",\"location\": \"In main() at main.go:36\",\"details\": {\"1\": \"Bob\",\"2\": \"Mary\"}}"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Mary"},{"position":3,"value":"error #1"},{"position":4,"value":"\n\t{\n\t\t\"time\": \"2023-04-10T11:00:20.623748617-04:00\",\n\t\t\"level\": \"TRACE\",\n\t\t\"id\": \"senzing-99990002\",\n\t\t\"text\": \"A fake error\",\n\t\t\"location\": \"In main() at main.go:36\",\"details\": {\"1\": \"Bob\",\"2\": \"Mary\"}}","valueRaw":{"time":"2023-04-10T11:00:20.623748617-04:00","level":"TRACE","id":"senzing-99990002","text":"A fake error","location":"In main() at main.go:36","details":{"1":"Bob","2":"Mary"}}}]}	`
	parsedMessage2, err := parser.Parse(message2)
	if err != nil {
		fmt.Printf("Error2: %s\n", err.Error())
	}
	fmt.Printf("Parse test 2:  ID: %s; Text: %s\n", parsedMessage2.ID, parsedMessage2.Text)

	message3 := messenger2.NewJson(2001, "Bob", "Mary", err1, err2)
	parsedMessage3, err := parser.Parse(message3)
	if err != nil {
		fmt.Printf("Error3: %s\n", err.Error())
	}
	fmt.Printf("Parse test 3:  ID: %s; Text: %s\n", parsedMessage3.ID, parsedMessage3.Text)

}
