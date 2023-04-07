/*
 */
package main

import (
	"errors"
	"fmt"

	"github.com/senzing/go-messaging/appmessage"
)

var (
	productIdentifier int                                    = 9999
	callerSkip        *appmessage.AppMessageOptionCallerSkip = &appmessage.AppMessageOptionCallerSkip{Value: 2}
)

var idMessages = map[int]string{
	2:    "Trace: %s knows %s",
	1001: "Debug: %s knows %s",
	2001: "Info: %s knows %s",
	3001: "Warn: %s knows %s",
	4001: "Error: %s knows %s",
	5001: "Fatal: %s knows %s",
	6001: "Panic: %s knows %s",
	7001: "Xxxxx: %s knows %s",
}

var idStatuses = map[int]string{}

func main() {

	appMessage, err := appmessage.New(productIdentifier, idMessages, idStatuses, callerSkip)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Println(appMessage.NewJson(0002, "Bob", "Mary"))

	err1 := errors.New("error #1")
	err2 := errors.New(`
	{
		"time": "2023-04-07 19:10:21.970756517 +0000 UTC",
		"level": "TRACE",
		"id": "senzing-99990002",
		"text": "Trace: Bob knows Mary",
		"location": "In main() at main.go:36",
		"details": {
			"1": "Bob",
			"2": "Mary"
		}
	}`)

	fmt.Println(appMessage.NewJson(1001, "Bob", "Mary", err1, err2))

	fmt.Printf("\n----- Logging -----------------------------------------------\n\n")

	// ctx := context.TODO()

	// Text logger - long form of construction.

	// textHandler := slog.NewTextHandler(os.Stderr)
	// textLogger := slog.New(textHandler)

	// textOptions := appmessage.HandlerOptions(appmessage.LevelInfoSlog)
	// textHandler := textOptions.NewTextHandler(os.Stderr)
	// textLogger := slog.New(textHandler)

	// JSON logger - short form of construction.

	// jsonHandler := slog.NewJSONHandler(os.Stderr)
	// jsonLogger := slog.New(jsonHandler)

	// jsonLogger := slog.New(appmessage.HandlerOptions(appmessage.LevelInfoSlog).NewJSONHandler(os.Stderr))

	// Initialize message generator.

	// Create a message and details.

	// msg, details := appMessage.NewSlog(2001, "Bob", "Mary")

	// Log the message.

	// textLogger.Info(msg, details...)
	// jsonLogger.Info(msg, details...)

	// Logging with auto-level generation.

	// msg0, level0, details0 := appMessage.NewSlogLevel(0002, "Bob", "Mary")
	// jsonLogger.Log(ctx, level0, msg0, details0...)

	// msg1, level1, details1 := appMessage.NewSlogLevel(1001, "Bob", "Mary")
	// jsonLogger.Log(ctx, level1, msg1, details1...)

	// msg2, level2, details2 := appMessage.NewSlogLevel(2001, "Bob", "Mary")
	// jsonLogger.Log(ctx, level2, msg2, details2...)

	// msg3, level3, details3 := appMessage.NewSlogLevel(3001, "Bob", "Mary")
	// jsonLogger.Log(ctx, level3, msg3, details3...)

	// msg4, level4, details4 := appMessage.NewSlogLevel(4001, "Bob", "Mary")
	// jsonLogger.Log(ctx, level4, msg4, details4...)

	// msg5, level5, details5 := appMessage.NewSlogLevel(5001, "Bob", "Mary")
	// jsonLogger.Log(ctx, level5, msg5, details5...)

	// msg6, level6, details6 := appMessage.NewSlogLevel(6001, "Bob", "Mary")
	// jsonLogger.Log(ctx, level6, msg6, details6...)

	// msg7, level7, details7 := appMessage.NewSlogLevel(7001, "Bob", "Mary")
	// jsonLogger.Log(ctx, level7, msg7, details7...)

	// Shouldn't appear because we're not in TRACE level.

}
