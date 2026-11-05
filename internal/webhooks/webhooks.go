// File generated from our OpenAPI spec by Scalar. See README.md for details.

// Package webhooks implements Standard Webhooks signature verification
// (https://github.com/standard-webhooks/standard-webhooks) using only the
// standard library, so generated SDKs verify webhook payloads without any
// third-party module. It is internal because only the SDK's own webhook parse
// methods call it.
package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DefaultTolerance bounds how far a webhook timestamp may drift from "now"
// (in either direction) before verification fails, limiting replay windows.
const DefaultTolerance = 5 * time.Minute

// secretPrefix is the optional Standard Webhooks secret prefix; the remainder
// of the secret is base64.
const secretPrefix = "whsec_"

// signatureVersion marks the HMAC scheme of one entry in the signature header.
const signatureVersion = "v1"

// The reasons a delivery is refused. Every error Verify returns for a refused
// delivery matches exactly one of them, so callers branch with errors.Is. They
// carry the names and messages of the Standard Webhooks Go library's sentinels,
// which SDKs that verified through that library returned wrapped; the SDK
// package re-exports them, because callers cannot import an internal package.
var (
	ErrRequiredHeaders     = errors.New("missing required headers")
	ErrInvalidHeaders      = errors.New("invalid signature headers")
	ErrNoMatchingSignature = errors.New("no matching signature found")
	ErrMessageTooOld       = errors.New("message timestamp too old")
	ErrMessageTooNew       = errors.New("message timestamp too new")
)

// verificationError is a refused delivery: its message says what was wrong in
// the Standard Webhooks library's words, errors.Is matches the sentinel naming
// the reason, and it unwraps to the error that caused it, when there is one.
type verificationError struct {
	message  string
	sentinel error
	cause    error
}

// refused builds the error for a delivery refused for the sentinel's reason,
// worded as the Standard Webhooks library words its own: detail after an
// "unable to verify payload" prefix.
func refused(sentinel error, cause error, detail string) error {
	return &verificationError{message: "unable to verify payload, err: " + detail, sentinel: sentinel, cause: cause}
}

// Error returns the message, worded as the library's.
func (e *verificationError) Error() string { return e.message }

// Is reports whether target is the sentinel naming this refusal, which is how
// errors.Is matches it.
func (e *verificationError) Is(target error) bool { return target == e.sentinel }

// Unwrap returns the error that caused the refusal, if any, so errors.Unwrap
// and errors.As still reach a timestamp's *strconv.NumError as they did when
// it was wrapped with %w.
func (e *verificationError) Unwrap() error { return e.cause }

// The standard header names a delivery carries its id, timestamp and
// signatures in.
const (
	idHeader        = "webhook-id"
	timestampHeader = "webhook-timestamp"
	signatureHeader = "webhook-signature"
)

// Verify checks a Standard Webhooks signature: the HMAC-SHA256 of
// "{id}.{timestamp}.{payload}" keyed with the base64-decoded secret must match
// one of the space-separated "v1,<base64>" entries in the "webhook-signature"
// header, and the "webhook-timestamp" header must be within tolerance of now.
// Comparison is constant-time.
func Verify(payload []byte, headers http.Header, secret string, now time.Time, tolerance time.Duration) error {
	// Decoding the secret first means a misconfigured secret is reported as
	// such, whatever the delivery carries.
	key, err := decodeSecret(secret)
	if err != nil {
		return err
	}

	id := headers.Get(idHeader)
	timestamp := headers.Get(timestampHeader)
	signatures := headers.Get(signatureHeader)
	if id == "" || timestamp == "" || signatures == "" {
		return refused(ErrRequiredHeaders, nil, ErrRequiredHeaders.Error())
	}

	sentAt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		// Worded as the library joins them: the strconv error, which quotes the
		// header's value, then the sentinel's message on its own line.
		return refused(ErrInvalidHeaders, err, "unable to parse timestamp header, err: "+errors.Join(err, ErrInvalidHeaders).Error())
	}
	if err := checkFreshness(time.Unix(sentAt, 0), now, tolerance); err != nil {
		return err
	}

	// The signed timestamp is the parsed integer written back out, so a header
	// sent as "+1614265330" or "01614265330" is checked against the signature
	// of "1614265330".
	expected := sign(key, id, strconv.FormatInt(sentAt, 10), payload)
	if !anySignatureMatches(signatures, expected) {
		return refused(ErrNoMatchingSignature, nil, ErrNoMatchingSignature.Error())
	}
	return nil
}

// checkFreshness rejects a delivery sent more than tolerance before or after
// now, which bounds how long a captured delivery can be replayed. A
// non-positive tolerance means DefaultTolerance.
func checkFreshness(sentAt time.Time, now time.Time, tolerance time.Duration) error {
	if tolerance <= 0 {
		tolerance = DefaultTolerance
	}
	switch {
	case now.Sub(sentAt) > tolerance:
		return refused(ErrMessageTooOld, nil, ErrMessageTooOld.Error())
	case sentAt.Sub(now) > tolerance:
		return refused(ErrMessageTooNew, nil, ErrMessageTooNew.Error())
	}
	return nil
}

// anySignatureMatches reports whether one of the space-separated
// "<version>,<signature>" entries is a v1 signature equal to expected. Entries
// with another version, or without a comma, are skipped so a sender can list
// future schemes beside v1.
func anySignatureMatches(signatures string, expected string) bool {
	for _, entry := range strings.Split(signatures, " ") {
		version, rest, found := strings.Cut(entry, ",")
		if !found || version != signatureVersion {
			continue
		}
		// Anything after a second comma is not part of the signature.
		signature, _, _ := strings.Cut(rest, ",")
		// hmac.Equal is constant-time, so a forged signature cannot be probed
		// byte by byte through timing differences.
		if hmac.Equal([]byte(signature), []byte(expected)) {
			return true
		}
	}
	return false
}

// Sign computes the base64 Standard Webhooks v1 signature (without the "v1,"
// prefix) over "{id}.{timestamp}.{payload}". Exposed so tests can craft valid
// signed payloads.
func Sign(id string, timestamp string, payload []byte, secret string) (string, error) {
	key, err := decodeSecret(secret)
	if err != nil {
		return "", err
	}
	return sign(key, id, timestamp, payload), nil
}

// decodeSecret strips the optional "whsec_" prefix and base64-decodes the
// rest into the HMAC key.
func decodeSecret(secret string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(secret, secretPrefix))
	if err != nil {
		return nil, fmt.Errorf("unable to create webhook, err: %w", err)
	}
	// A bare "whsec_" decodes to an empty key, and an HMAC under an empty key
	// is one anyone can compute, so refuse it rather than accept forgeries.
	if len(key) == 0 {
		return nil, errors.New("unable to create webhook, err: signing secret is empty")
	}
	return key, nil
}

// sign returns the base64 HMAC-SHA256 of "{id}.{timestamp}.{payload}" under key.
func sign(key []byte, id string, timestamp string, payload []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(id))
	mac.Write([]byte("."))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(payload)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
