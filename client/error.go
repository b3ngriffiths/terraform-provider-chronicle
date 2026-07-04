package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi"
)

type ChronicleAPIError struct {
	Message        string `json:"message"`
	Result         string `json:"result"`
	HTTPStatusCode int
}

func (c *ChronicleAPIError) Error() string {
	return fmt.Sprintf("%s: %s, HTTP status code: %d", c.Result, c.Message, c.HTTPStatusCode)
}

func errorForStatusCode(r *http.Response, err error) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	apiError := &ChronicleAPIError{
		HTTPStatusCode: r.StatusCode,
	}

	// googleapi.CheckResponse has already consumed the response body; it keeps
	// the raw payload in Error.Body, so parse the API's error details from there.
	gError, ok := err.(*googleapi.Error)
	if ok {
		apiError.Message = gError.Message
		_ = json.Unmarshal([]byte(gError.Body), apiError)
	}

	if apiError.Message == "" && ok {
		apiError.Message = gError.Body
	}

	return apiError
}
