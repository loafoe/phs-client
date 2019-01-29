// This package provides an API client for interacting with the PHS services
package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jeffail/gabs"
)

var (
	TIME_FORMAT = "2006-01-02T15:04:05.000-0700"
)

// The API client
type APIClient struct {
	apiSigner *APISigner
	config    *Config
}

type Config struct {
	Secret string
	Debug  bool
	SkipS0 bool
}

// NewClient creates a new client. It takes a Config struct
// for configuration
func NewClient(config *Config) (*APIClient, error) {
	client := &APIClient{
		apiSigner: NewAPISigner(config.Secret),
		config:    config,
	}
	if os.Getenv("PHS_CLIENT_DEBUG") == "true" {
		client.config.Debug = true
	}
	if os.Getenv("SKIP_S0") == "true" {
		client.config.SkipS0 = true
	}

	return client, nil
}

func (client *APIClient) SignRequest(request *http.Request) {
	client.sign(&request.Header, request.RequestURI, request.Method)
}

func (client *APIClient) sendRestRequest(httpMethod string, uri *url.URL, header *http.Header, body []byte) Response {
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")

	buf := bytes.NewBuffer(body)
	request, err := http.NewRequest(httpMethod, uri.String(), buf)
	if err != nil {
		return Response{
			Body: "error",
		}
	}

	for k, _ := range *header {
		request.Header.Set(k, header.Get(k))
	}
	if client.config.Debug {
		dumped, _ := httputil.DumpRequest(request, true)
		log.Info(string(dumped))
	}

	// Fetch Request
	c := &http.Client{}
	resp, err := c.Do(request)

	if err != nil {
		return Response{
			Body: err.Error(),
		}
	}
	if client.config.Debug {
		dumped, _ := httputil.DumpResponse(resp, false)
		log.Info(string(dumped))
	}

	// Read Response Body
	responseBody, _ := ioutil.ReadAll(resp.Body)

	jsonParsed, err := gabs.ParseJSON([]byte(responseBody))
	if err == nil {
		dhpCode, _ := jsonParsed.Path("responseCode").Data().(string)
		intCode, _ := strconv.Atoi(dhpCode)
		return Response{
			Body:       string(responseBody),
			StatusCode: resp.StatusCode,
			DhpCode:    intCode,
			Response:   resp,
		}
	}
	return Response{
		Body:       string(responseBody),
		StatusCode: resp.StatusCode,
		DhpCode:    0,
		Response:   resp,
	}
}

func (client *APIClient) sendSignedRequest(httpMethod, requestUri string, header *http.Header, body []byte) Response {
	client.sign(header, requestUri, httpMethod)
	uri, _ := url.Parse(requestUri)

	return client.sendRestRequest(httpMethod, uri, header, body)
}

func (client *APIClient) SendRestRequest(httpMethod, requestUri string, header *http.Header, body []byte) Response {
	uri, _ := url.Parse(requestUri)
	return client.sendRestRequest(httpMethod, uri, header, body)
}

// Sends a signed request to the configured service
// It returns a Response struct which contains the parsed response code
// from the service. A full copy of the body is also returned
func (client *APIClient) SendSignedRequest(httpMethod, requestUri string, header *http.Header, body []byte) Response {

	return client.sendSignedRequest(httpMethod, requestUri, header, body)
}

func (client *APIClient) sign(header *http.Header, requestUri, httpMethod string) {
	if !client.config.SkipS0 {
		authHeaderValue, err := client.apiSigner.BuildAuthorizationHeaderValueS0(httpMethod, requestUri)
		if err == nil {
			header.Set("Authorization", "Signature "+authHeaderValue)
		} else {
			log.Error("Failed to sign: ", err.Error())
		}
	}
	authHeaderValueS1, err := client.apiSigner.BuildAuthorizationHeaderValueS1(httpMethod, requestUri, time.Now())
	if err == nil {
		header.Set("M-Authorization", authHeaderValueS1)
	} else {
		log.Error("Failed to sign: ", err.Error())
	}

}
