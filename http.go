package nepse

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/voidarchive/go-nepse/internal/auth"
)

// nepseClient implements the NEPSE HTTP client with authentication.
type nepseClient struct {
	client      *http.Client
	config      *Config
	authManager *auth.Manager
	options     *Options
}

// newHTTPClient creates a new HTTP client for NEPSE API.
func newHTTPClient(options *Options) (*nepseClient, error) {
	if options == nil {
		options = DefaultOptions()
	}

	if options.Config == nil {
		options.Config = DefaultConfig()
	}

	httpClient := options.HTTPClient
	if httpClient == nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !options.TLSVerification, //nolint:gosec
			},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		}
		httpClient = &http.Client{
			Timeout:   options.HTTPTimeout,
			Transport: transport,
		}
	} else if httpClient.Timeout == 0 {
		httpClient.Timeout = options.HTTPTimeout
	}

	c := &nepseClient{
		client:  httpClient,
		config:  options.Config,
		options: options,
	}

	authManager, err := auth.NewManager(c)
	if err != nil {
		return nil, NewInternalError("failed to create auth manager", err)
	}
	c.authManager = authManager

	return c, nil
}

// GetToken implements auth.NepseHTTP interface.
func (h *nepseClient) GetToken(ctx context.Context) (*auth.TokenResponse, error) {
	url := h.config.BaseURL + "/api/authenticate/prove"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, NewInternalError("failed to create request", err)
	}

	h.setCommonHeaders(req)

	resp, err := h.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, MapHTTPStatusToError(resp.StatusCode, resp.Status)
	}

	var tokenResp auth.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, NewInternalError("failed to decode token response", err)
	}

	return &tokenResp, nil
}

// doRequest performs HTTP request with retry logic.
func (h *nepseClient) doRequest(req *http.Request) (*http.Response, error) {
	var lastErr error
	maxDelay := 30 * time.Second

	for attempt := 0; attempt <= h.options.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := h.options.RetryDelay * time.Duration(1<<uint(attempt-1))
			if delay > maxDelay {
				delay = maxDelay
			}

			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(delay):
			}
		}

		resp, err := h.client.Do(req)
		if err != nil {
			lastErr = NewNetworkError(err)
			continue
		}

		if resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = MapHTTPStatusToError(resp.StatusCode, resp.Status)
			if !lastErr.(*NepseError).IsRetryable() {
				return nil, lastErr
			}
			continue
		}

		return resp, nil
	}

	return nil, lastErr
}

// setCommonHeaders sets common HTTP headers for requests.
func (h *nepseClient) setCommonHeaders(req *http.Request) {
	for key, value := range h.config.Headers {
		if key == "Host" {
			req.Header.Set(key, strings.Replace(h.config.BaseURL, "https://", "", 1))
		} else if key == "Referer" {
			req.Header.Set(key, h.config.BaseURL+"/")
		} else if value != "" {
			req.Header.Set(key, value)
		}
	}

	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Origin", h.config.BaseURL)
}

// apiRequest performs an authenticated API request.
func (h *nepseClient) apiRequest(ctx context.Context, endpoint string, result any) error {
	return h.apiRequestWithRetry(ctx, endpoint, result, 0)
}

// apiRequestWithRetry performs an authenticated API request with token refresh retry.
func (h *nepseClient) apiRequestWithRetry(ctx context.Context, endpoint string, result any, retryCount int) error {
	token, err := h.authManager.AccessToken(ctx)
	if err != nil {
		return NewInternalError("failed to get access token", err)
	}

	url := h.config.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NewInternalError("failed to create request", err)
	}

	auth.SetAuthHeader(req, token)
	req.Header.Set("Content-Type", "application/json")
	h.setCommonHeaders(req)

	resp, err := h.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized && retryCount == 0 {
		if err := h.authManager.ForceUpdate(ctx); err != nil {
			return NewInternalError("failed to refresh token", err)
		}
		return h.apiRequestWithRetry(ctx, endpoint, result, retryCount+1)
	}

	if resp.StatusCode != http.StatusOK {
		return MapHTTPStatusToError(resp.StatusCode, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return NewInternalError("failed to decode response", err)
	}

	return nil
}

// apiRequestRaw performs an authenticated API request and returns raw body.
func (h *nepseClient) apiRequestRaw(ctx context.Context, endpoint string) ([]byte, error) {
	token, err := h.authManager.AccessToken(ctx)
	if err != nil {
		return nil, NewInternalError("failed to get access token", err)
	}

	url := h.config.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, NewInternalError("failed to create request", err)
	}

	auth.SetAuthHeader(req, token)
	req.Header.Set("Content-Type", "application/json")
	h.setCommonHeaders(req)

	resp, err := h.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, MapHTTPStatusToError(resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// GetConfig returns the current configuration.
func (h *nepseClient) GetConfig() *Config {
	return h.config
}

// Close closes the HTTP client and auth manager.
func (h *nepseClient) Close() error {
	if h.authManager != nil {
		return h.authManager.Close()
	}
	return nil
}
