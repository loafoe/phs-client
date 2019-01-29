package client

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type APISigner struct {
	secretKey string
}

func NewAPISigner(secretKey string) *APISigner {
	return &APISigner{secretKey}
}

func (signer *APISigner) Init(secretKey string) {
	signer.secretKey = secretKey
}

func (signer *APISigner) BuildAuthorizationHeaderValueS1(requestMethod, requestUrl string, now time.Time) (string, error) {
	stime := fmt.Sprintf("%x", now.UnixNano()/int64(time.Millisecond))
	return signer.BuildSigningStringS1(requestMethod, requestUrl, stime)
}

func (signer *APISigner) BuildSigningStringS1(requestMethod, requestUrl, timeString string) (string, error) {
	sig0, err := signer.BuildSigningStringS0(requestMethod, requestUrl)
	if err != nil {
		return "", err
	}
	stringToBeSigned := strings.Join([]string{timeString, sig0}, "|")
	signature := signString([]byte(signer.secretKey), stringToBeSigned)
	return strings.Join([]string{"S1", timeString, signature}, ":"), nil
}

func (signer *APISigner) BuildAuthorizationHeaderValueS0(requestMethod, requestUrl string) (string, error) {
	stringToBeSigned, err := signer.BuildSigningStringS0(requestMethod, requestUrl)
	if err != nil {
		return "", err
	}
	signature := signString([]byte(signer.secretKey), stringToBeSigned)
	return signature, nil
}

func (signer *APISigner) BuildSigningStringS0(requestMethod, requestUrl string) (string, error) {
	parsed, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{parsed.Path, requestMethod, extractSortedValues(parsed.Query())}, "|"), nil
}

func extractSortedValues(values url.Values) string {
	var keys []string
	var vals []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vals = append(vals, k+"="+values.Get(k))
	}
	out := strings.Join(vals, "|")
	return out
}

func signString(signatureKey []byte, stringToBeSigned string) string {
	signatureSlice := hash([]byte(stringToBeSigned), signatureKey)
	return hex.EncodeToString(signatureSlice)
}

func hash(data []byte, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
