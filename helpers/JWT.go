package helpers

import (
	"errors"
	"net/http"
	"strings"
)

func GetJwtToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("couldn't find authorization header")
	}

	if !strings.HasPrefix(strings.ToLower(header), "bearer") {
		return "", errors.New("authorization header didn't start with 'bearer'")
	}

	return header[7:], nil
}
