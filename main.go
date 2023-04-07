/*
 */
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/senzing/go-messaging/appmessage"
	"golang.org/x/exp/slog"
)

var (
	productIdentifier int                              = 9999
	callerSkip        *appmessage.AppMessageCallerSkip = &appmessage.AppMessageCallerSkip{Value: 2}
)

var idMessages = map[int]string{
	2001: "%s knows %s",
	3001: "%s knows %s",
	4001: "%s knows %s",
	2:    "%s does not know %s",
}

var idStatuses = map[int]string{}

func main() {

	ctx := context.TODO()

	// Initialize message generator.

	appMessage, err := appmessage.New(productIdentifier, idMessages, idStatuses, callerSkip)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	// Create a message and details.

	msg1, details1 := appMessage.NewSlog(2001, "Bob", "Mary")

	// Text logging.

	textHandler := slog.NewTextHandler(os.Stderr)
	textLogger := slog.New(textHandler)
	textLogger.Info(msg1, details1...)

	// JSON logging.

	jsonHandler := slog.NewJSONHandler(os.Stderr)
	jsonLogger := slog.New(jsonHandler)
	jsonLogger.Info(msg1, details1...)

	// Logging with auto-level generation.

	msg2, level2, details2 := appMessage.NewSlogLevel(2001, "Bob", "Mary")
	jsonLogger.Log(ctx, level2, msg2, details2...)

	msg3, level3, details3 := appMessage.NewSlogLevel(3001, "Bob", "Mary")
	jsonLogger.Log(ctx, level3, msg3, details3...)

}
