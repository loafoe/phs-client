package command

import (
	_ "encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"git.aemian.com/phs/client"
	log "github.com/Sirupsen/logrus"
	"github.com/jeffail/gabs"
	"github.com/loafoe/cfutil"
	"github.com/mitchellh/cli"
)

type RequestCommand struct {
	Revision          string
	Version           string
	VersionPrerelease string
	Ui                cli.Ui
}

func (uc *RequestCommand) Help() string {
	helpText := `
Usage: phsclient request [options]
  Request
Options:
  -secret=stringKey             The shared secret
  -method=GET|POST|PUT|DELETE   Request method 
  -body=stringBody              Optional body of request
  -url=http://bla/path?a=b      The path
  -headers="..."                Headers to add. Separate with ;
	`
	return strings.TrimSpace(helpText)
}

func (uc *RequestCommand) Synopsis() string {
	return "Request command"
}

func (uc *RequestCommand) Run(args []string) int {
	var header *http.Header = &http.Header{}
	secret := cfutil.Getenv("PRIVATE_API_SHARED_SECRET")
	cmdFlags := flag.NewFlagSet("test", flag.ContinueOnError)
	cmdFlags.Usage = func() { uc.Ui.Output(uc.Help()) }
	method := cmdFlags.String("method", "GET", "Type of request method. Defaults to GET")
	url := cmdFlags.String("url", "", "The URL including query parameters")
	cmdSecret := cmdFlags.String("secret", "", "The shared secret")
	body := cmdFlags.String("body", "", "Optional JSON body of request")
	headers := cmdFlags.String("headers", "", "Headers to add to request. Separate with ;")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	if *cmdSecret != "" {
		secret = *cmdSecret
	}
	c, err := client.NewClient(&client.Config{Secret: secret})
	if err != nil {
		return 1
	}
	if *url == "" {
		fmt.Println("URL is required")
		return 2
	}

	// Add headers
	splittedHeaders := strings.Split(*headers, ";")
	for _, h := range splittedHeaders {
		kv := strings.Split(h, ":")
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		header.Set(key, val)
	}

	response := c.SendSignedRequest(*method, *url, header, []byte(*body))

	if response.StatusCode < 200 {
		log.Print(string(response.Body))
		for _, e := range response.Errors {
			log.Print(e)
		}
		fmt.Println(client.StatusCodeToString(response.DhpCode))
	}
	jsonParsed, _ := gabs.ParseJSON([]byte(response.Body))
	if jsonParsed != nil {
		fmt.Println(jsonParsed.StringIndent("", "  "))
	} else {
		fmt.Print(response.Body)
	}

	return 0
}
