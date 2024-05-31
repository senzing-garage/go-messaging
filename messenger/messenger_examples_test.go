package messenger

import (
	"fmt"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleImplementation_NewJSON() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/messenger/messenger_examples_test.go
	example, err := New()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(example.NewJSON(2001, "Bob", "Jane", getTimestamp(), getOptionCallerSkip()))
	//Output: {"time":"2000-01-01T00:00:00Z","level":"INFO","id":"SZSDK99992001","location":"In ExampleImplementation_NewJSON() at messenger_examples_test.go:17","details":[{"position":1,"type":"string","value":"Bob"},{"position":2,"type":"string","value":"Jane"}]}
}

func ExampleImplementation_NewSlog() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/messenger/messenger_examples_test.go
	example, err := New()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(example.NewSlog(2001, "Bob", "Jane", getTimestamp(), getOptionCallerSkip()))
	//Output: [id SZSDK99992001 location In NewSlog() at messenger.go:387 details [{ 1 string Bob <nil>} { 2 string Jane <nil>}]]
}

func ExampleImplementation_NewSlogLevel() {
	// For more information, visit https://github.com/senzing-garage/go-messaging/blob/main/messenger/messenger_examples_test.go
	example, err := New()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(example.NewSlogLevel(2001, "Bob", "Jane", getTimestamp(), getOptionCallerSkip()))
	//Output: INFO [id SZSDK99992001 location In ExampleImplementation_NewSlogLevel() at messenger_examples_test.go:37 details [{ 1 string Bob <nil>} { 2 string Jane <nil>}]]
}
