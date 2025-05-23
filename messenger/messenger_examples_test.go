package messenger_test

import (
	"fmt"

	"github.com/senzing-garage/go-messaging/messenger"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleBasicMessenger_NewJSON() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/messenger/messenger_examples_test.go
	example, err := messenger.New(getOptionMessageFields())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(example.NewJSON(2001, "Bob", "Jane"))
	//Output: {"level":"INFO","id":"2001","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}
}

func ExampleBasicMessenger_NewSlog() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/messenger/messenger_examples_test.go
	example, err := messenger.New(getOptionMessageFields())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(example.NewSlog(2001, "Bob", "Jane", getOptionMessageFields()))
	//Output: [level INFO id 2001 details [{ 1 string Bob <nil>} { 2 string Jane <nil>}]]
}

func ExampleBasicMessenger_NewSlogLevel() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/messenger/messenger_examples_test.go
	example, err := messenger.New(getOptionMessageFields())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print(example.NewSlogLevel(2001, "Bob", "Jane", getOptionMessageFields()))
	//Output: INFO [level INFO id 2001 details [{ 1 string Bob <nil>} { 2 string Jane <nil>}]]
}
