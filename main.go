/*
 */
package main

import (
	"fmt"
	"os"

	"github.com/senzing/go-messaging/appmessage"
	"golang.org/x/exp/slog"
)

var productIdentifier int = 9999

var idMessages = map[int]string{
	2001: "%s knows %s",
	3001: "%s knows %s",
	4001: "%s knows %s",
	2:    "%s does not know %s",
}
var callerSkip int = 1

var idStatuses = map[int]string{}

func main() {

	// Initialize message generator.

	appMessage, err := appmessage.New(productIdentifier, idMessages, idStatuses, callerSkip)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// Create a message and details.

	msg, details := appMessage.NewSlog(2001, "Bob", "Mary")

	// Text logging.

	textHandler := slog.NewTextHandler(os.Stderr)
	textLogger := slog.New(textHandler)
	textLogger.Info(msg, details...)

	// JSON logging.

	jsonHandler := slog.NewJSONHandler(os.Stderr)
	jsonLogger := slog.New(jsonHandler)
	jsonLogger.Info(msg, details...)

}
