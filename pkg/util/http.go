package util

import (
	"fmt"
	"io"
	"net/http"
	"slices"
)

func ReadBody(res *http.Response, ignoreStatusCodes []int) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}

	if slices.Contains(ignoreStatusCodes, res.StatusCode) {
		return body, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned error (%d): %s - %s", res.StatusCode, res.Status, string(body))
	}

	return body, nil
}
