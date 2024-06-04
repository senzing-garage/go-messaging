/*
 */
package main

import (
	"errors"
	"fmt"

	"github.com/senzing-garage/go-messaging/messenger"
	"github.com/senzing-garage/go-messaging/parser"
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

func testError(err error, stringFormat string) {
	if err != nil {
		fmt.Printf(stringFormat, err.Error())
	}
}

func main() {

	// Create some fake errors.

	err1 := errors.New("error #1")
	err2 := errors.New(`
    {
        "time": "2023-04-10T11:00:20.623748617-04:00",
        "level": "TRACE",
        "id": "SZSDK99990002",
        "text": "A fake error",
        "location": "In main() at main.go:36",
        "errors": ["0027E|Unknown DATA_SOURCE value 'DOESNTEXIST'"],
        "details": [
            {"position": 1, "type": "string", "value": "DoesntExist"},
            {"position": 2, "type": "string", "value": "1070", "valueRaw": 1070},
            {"position": 3, "type": "int64", "value": "-1", "valueRaw": -1},
            {
                "position": 4,
                "type": "g2engine._Ctype_longlong",
                "value": "-2",
                "valueRaw": -2,
            },
            {
                "position": 5,
                "type": "error",
                "value": "0027E|Unknown DATA_SOURCE value 'DOESNTEXIST'",
            },
        ],
    }`)

	// ------------------------------------------------------------------------
	// --- Using a bare message generator
	// ------------------------------------------------------------------------

	// Example from README.md

	optionMessageFields := &messenger.OptionMessageFields{
		Value: []string{"id", "details"},
	}

	aMessenger, _ := messenger.New(optionMessageFields)
	fmt.Printf("%s\n\n", aMessenger.NewJSON(0001, "Bob", "Mary"))
	fmt.Println(aMessenger.NewSlog(1001, "Bob", "Mary"))
	fmt.Println()

	// Create a bare message generator.

	messenger1, err := messenger.New(optionMessageFields)
	testError(err, "Error: %s\n")

	// Print some messages.

	fmt.Printf("%s\n\n", messenger1.NewJSON(2001, "Bob", "Mary"))
	fmt.Printf("%s\n\n", messenger1.NewJSON(3001, "Bob", "Mary", err1, err2))

	// ------------------------------------------------------------------------
	// --- Using a configured message generator
	// ------------------------------------------------------------------------

	// Create a configured message generator.

	optionSenzingComponentID := &messenger.OptionSenzingComponentID{Value: 9998}
	optionCallerSkip := &messenger.OptionCallerSkip{Value: 2}
	optionIDMessages := &messenger.OptionIDMessages{Value: idMessages}
	messenger2, err := messenger.New(optionSenzingComponentID, optionCallerSkip, optionIDMessages)
	testError(err, "Error: %s\n")

	// Print some messages.

	fmt.Printf("%s\n\n", messenger2.NewJSON(0001, "Bob", "Mary"))
	fmt.Printf("%s\n\n", messenger2.NewJSON(1001, "Bob", "Mary", err1, err2))

	// Parse some messages.

	message1 := messenger2.NewJSON(0001, "Bob", "Mary")
	parsedMessage1, err := parser.Parse(message1)
	testError(err, "Error1: %s\n")
	fmt.Printf("Parse test 1:  ID: %s; Text: %s\n", parsedMessage1.ID, parsedMessage1.Text)

	message2 := `{"time":"2023-07-11T21:05:51.918625982Z","level":"DEBUG","id":"SZSDK99981001","text":"DEBUG: Bob works with Mary","location":"In main() at main.go:101","errors":["error #1","{\"time\": \"2023-04-10T11:00:20.623748617-04:00\",\"level\": \"TRACE\",\"id\": \"SZSDK99990002\",\"text\": \"A fake error\",\"location\": \"In main() at main.go:36\",\"details\": {\"1\": \"Bob\",\"2\": \"Mary\"}}"],"details":[{"position":1,"value":"Bob"},{"position":2,"value":"Mary"},{"position":3,"value":"error #1"},{"position":4,"value":"\n\t{\n\t\t\"time\": \"2023-04-10T11:00:20.623748617-04:00\",\n\t\t\"level\": \"TRACE\",\n\t\t\"id\": \"SZSDK99990002\",\n\t\t\"text\": \"A fake error\",\n\t\t\"location\": \"In main() at main.go:36\",\"details\": {\"1\": \"Bob\",\"2\": \"Mary\"}}","valueRaw":{"time":"2023-04-10T11:00:20.623748617-04:00","level":"TRACE","id":"SZSDK99990002","text":"A fake error","location":"In main() at main.go:36","details":{"1":"Bob","2":"Mary"}}}]}	`
	parsedMessage2, err := parser.Parse(message2)
	testError(err, "Error2: %s\n")
	fmt.Printf("Parse test 2:  ID: %s; Text: %s\n", parsedMessage2.ID, parsedMessage2.Text)

	message3 := messenger2.NewJSON(2001, "Bob", "Mary", err1, err2)
	parsedMessage3, err := parser.Parse(message3)
	testError(err, "Error3: %s\n")
	fmt.Printf("Parse test 3:  ID: %s; Text: %s\n", parsedMessage3.ID, parsedMessage3.Text)

}
