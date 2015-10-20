package command

import (
	"testing"

	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
)

func TestValidateCli_BothFileAndSite(t *testing.T) {
	args := cliValidationTestRunner(t, "", "-i", "hello.txt", "somesite.com")
	assert.EqualValues(t, "somesite.com", args[0])
}

func TestValidateCli_MissingInputFileAndSite(t *testing.T) {
	args := cliValidationTestRunner(t, "", "-o", "hello.txt")
	assert.Empty(t, args)
}

func TestValidateCli_OpenReaderError(t *testing.T) {
	args := cliValidationTestRunner(t, "", "-i", "..../...../test")
	assert.Empty(t, args)
}

func TestValidateCli_CreateWriterError(t *testing.T) {
	args := cliValidationTestRunner(t, "", "-o", "..../...../test")
	assert.Empty(t, args)
}

func TestValidateCli_FormatInvalid(t *testing.T) {
	args := cliValidationTestRunner(t, "", "-f", "invalidFormat")
	assert.Empty(t, args)
}

func cliValidationTestRunner(t *testing.T, args ...string) []string {
	var capturedArgs []string

	app := cli.NewApp()
	app.Action = func(c *cli.Context) {
		capturedArgs = c.Args()
		err := validateCLI(c)
		assert.NotNil(t, err)
	}
	app.Flags = AppCliFlags()

	app.Run(args)
	return capturedArgs
}
