package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/avast/retry-go"
	"google.golang.org/api/googleapi"
)

// isRetryableError reports whether a failed request may be retried: transport
// errors and 429/5xx responses are retryable, while other API errors (e.g.
// 400/403/404) are permanent and retrying them only delays the failure.
func isRetryableError(err error) bool {
	if !retry.IsRecoverable(err) {
		return false
	}

	var apiErr *ChronicleAPIError
	if errors.As(err, &apiErr) {
		return apiErr.HTTPStatusCode == http.StatusTooManyRequests || apiErr.HTTPStatusCode >= 500
	}

	return true
}

func sendRequest(client *Client, httpClient *http.Client, method, userAgent string, rawurl string, body interface{}) ([]byte, error) {
	if httpClient == nil {
		return nil, fmt.Errorf("no credentials configured for API request to %s: configure the corresponding provider credentials attribute or environment variable", rawurl)
	}

	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	reqHeaders.Set("User-Agent", userAgent)

	var res *http.Response
	err := retry.Do(
		func() error {
			var buf bytes.Buffer
			if body != nil {
				err := json.NewEncoder(&buf).Encode(body)
				if err != nil {
					return retry.Unrecoverable(err)
				}
			}

			u, err := addQueryParams(rawurl, map[string]string{"alt": "json"})
			if err != nil {
				return retry.Unrecoverable(err)
			}

			//nolint:all
			req, err := http.NewRequest(method, u, &buf)
			if err != nil {
				return retry.Unrecoverable(err)
			}

			req.Header = reqHeaders
			//nolint:all
			res, err = httpClient.Do(req)
			if err != nil {
				return err
			}

			if err := googleapi.CheckResponse(res); err != nil {
				googleapi.CloseBody(res)
				return errorForStatusCode(res, err)
			}

			return nil
		}, retry.Attempts(client.requestAttempts), retry.DelayType(retry.BackOffDelay), retry.LastErrorOnly(true),
		retry.RetryIf(isRetryableError), retry.OnRetry(func(n uint, err error) {
			log.Printf("[DEBUG] Retrying request after error: %v", err)
		}),
	)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, fmt.Errorf(`unable to parse server response. This is most likely a terraform problem,
		 please file a bug at https://github.com/form3tech-oss/terraform-provider-chronicle/issues`)
	}

	defer googleapi.CloseBody(res)

	if res.StatusCode == 204 {
		return nil, nil
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func addQueryParams(rawurl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
