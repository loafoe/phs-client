package client

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
)

const (
	HeaderMAuthorization = "M-Authorization"
)

type SignatureVerifier struct {
	secretKey string
	signer    *APISigner
}

func NewSignatureVerifier(secretKey string) *SignatureVerifier {
	return &SignatureVerifier{
		secretKey: secretKey,
		signer:    NewAPISigner(secretKey),
	}
}

func (v *SignatureVerifier) Init(secretKey string) {
	v.secretKey = secretKey
}

func (v *SignatureVerifier) VerifyRequest(request *http.Request) (bool, error) {
	signature := request.Header.Get(echo.HeaderAuthorization)
	if signature == "" {
		signature = request.Header.Get(HeaderMAuthorization)
	}
	return v.Verify(request.Method, request.RequestURI, signature)
}

func (v *SignatureVerifier) Verify(requestMethod, requestUrl, signature string) (bool, error) {
	spl := strings.Split(signature, ":")
	if len(spl) < 2 {
		// S0
		spl = strings.Split(signature, " ")
		if len(spl) != 2 {
			return false, fmt.Errorf("Malformed S0 signature")
		}
		sig0 := spl[1]
		gen0, _ := v.signer.BuildAuthorizationHeaderValueS0(requestMethod, requestUrl)
		if sig0 == gen0 { // Valid
			return true, nil
		}
		return false, fmt.Errorf("Invalid S0 signature")
	}
	// SX
	switch spl[0] {
	case "S1":
		var timestamp int64
		if n, _ := fmt.Sscanf(spl[1], "%x", &timestamp); n < 1 {
			return false, fmt.Errorf("Malformed S1 signature: invalid timestamp: %s", spl[1])
		}
		now := time.Now().UnixNano() / int64(time.Millisecond)
		if math.Abs(float64(timestamp)-float64(now)) > 30000.0 {
			return false, fmt.Errorf("S1 signature expired")
		}
		sig1, err := v.signer.BuildSigningStringS1(requestMethod, requestUrl, spl[1])
		if err != nil {
			return false, err
		}
		if signature == sig1 {
			return true, nil // Valid
		}
		return false, fmt.Errorf("Invalid signature: %s != %s", signature, sig1)
	default:
		return false, fmt.Errorf("Invalid signature type: %s", spl[0])
	}
	return false, fmt.Errorf("Invalid signature: (%s) (%s) (%s)", spl[0], spl[1], spl[2])
}
