package command

import (
	_ "encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"git.aemian.com/phs/client"
	"github.com/loafoe/cfutil"
	"github.com/mitchellh/cli"
)

type TestCommand struct {
	Revision          string
	Version           string
	VersionPrerelease string
	Ui                cli.Ui
}

func (uc *TestCommand) Help() string {
	helpText := `
Usage: phsclient test [options]
  Test
Options:
  -secret=stringKey          The shared secret
  -method=GET|POST           Request method 
  -url=http://bla/path?a=b   The path
	`
	return strings.TrimSpace(helpText)
}

func (uc *TestCommand) Synopsis() string {
	return "Test command"
}

func (uc *TestCommand) Run(args []string) int {
	secret := cfutil.Getenv("PRIVATE_API_SHARED_SECRET")
	cmdFlags := flag.NewFlagSet("test", flag.ContinueOnError)
	cmdFlags.Usage = func() { uc.Ui.Output(uc.Help()) }
	method := cmdFlags.String("method", "GET", "Type of request method. Defaults to GET")
	url := cmdFlags.String("url", "", "The URL including query parameters")
	cmdSecret := cmdFlags.String("secret", "", "The shared secret")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	if *cmdSecret != "" {
		secret = *cmdSecret
	}

	fmt.Printf("%s %s\n", *method, *url)
	fmt.Printf("Secret: %s\n", secret)
	signer := client.NewAPISigner(secret)
	signed, err := signer.BuildAuthorizationHeaderValueS1(*method, *url, time.Now())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return 1
	}
	fmt.Printf("Signature: %s\n", signed)

	return 0
}
