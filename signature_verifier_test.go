package client

import (
	"testing"
	"time"
)

func TestS0Verifier(t *testing.T) {
	secret := "TOP SECRET!"
	signer := NewAPISigner(secret)
	verifier := NewSignatureVerifier(secret)

	method := "GET"
	requestURI := "/api/version"

	s0, err := signer.BuildAuthorizationHeaderValueS0(method, requestURI)
	if err != nil {
		t.Error("S0 signature generation failed:", err)
	}

	valid, err := verifier.Verify(method, requestURI, "Signature "+s0)

	if !valid {
		t.Error("S0 signature failed:", err)
	}

}

func TestS1Verifier(t *testing.T) {
	secret := "TOP SECRET!"
	signer := NewAPISigner(secret)
	verifier := NewSignatureVerifier(secret)

	now := time.Now()
	method := "GET"
	requestURI := "/api/version"

	s1, err := signer.BuildAuthorizationHeaderValueS1(method, requestURI, now)
	if err != nil {
		t.Error("S1 signature generation failed:", err)
	}

	valid, err := verifier.Verify(method, requestURI, s1)

	if !valid {
		t.Error(err)
	}

}
