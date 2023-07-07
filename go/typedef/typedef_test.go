package typedef

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-messaging/messenger"
	"github.com/stretchr/testify/assert"
)

var (
	componentId int = 9999
	logger      messenger.MessengerInterface
)

var idMessages = map[int]string{
	1: "bob %s",
}

func testError(test *testing.T, ctx context.Context, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error = nil
	logger, err = messenger.New()
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// --- Test cases
// ----------------------------------------------------------------------------

func TestConfigAddDataSourceResponseTest101(test *testing.T) {

	aMessage := logger.NewJson(0001, "Bob", "Mary")
	test.Log(aMessage)
}
