// Package auth handles NEPSE API authentication.
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// DefaultTokenTTL is the maximum time a token is considered valid.
// NEPSE tokens expire after ~60 seconds; we refresh at 45s for safety.
const DefaultTokenTTL = 45 * time.Second

// NepseHTTP abstracts the HTTP operations needed for token management.
type NepseHTTP interface {
	// GetTokens performs GET to /api/authenticate/prove and returns the token response.
	GetTokens(ctx context.Context) (*TokenResponse, error)

	// RefreshTokens performs GET to /api/authenticate/refresh-token and returns new tokens.
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenResponse, error)
}

// TokenResponse mirrors the JSON from /api/authenticate/prove.
type TokenResponse struct {
	Salt1        int    `json:"salt1"`
	Salt2        int    `json:"salt2"`
	Salt3        int    `json:"salt3"`
	Salt4        int    `json:"salt4"`
	Salt5        int    `json:"salt5"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ServerTime   int64  `json:"serverTime"`
}

// Manager manages NEPSE auth tokens.
type Manager struct {
	http   NepseHTTP
	parser *tokenParser

	maxUpdatePeriod time.Duration

	mu           sync.RWMutex
	accessToken  string
	refreshToken string
	tokenTS      time.Time
	salts        [5]int

	sf singleflight.Group
}

// NewManager constructs a Manager with embedded WASM parser.
func NewManager(httpClient NepseHTTP) (*Manager, error) {
	parser, err := newTokenParser()
	if err != nil {
		return nil, fmt.Errorf("init wasm parser: %w", err)
	}
	return &Manager{
		http:            httpClient,
		parser:          parser,
		maxUpdatePeriod: DefaultTokenTTL,
	}, nil
}

// Close releases WASM runtime resources.
func (m *Manager) Close(ctx context.Context) error {
	if m.parser != nil {
		return m.parser.close(ctx)
	}
	return nil
}

// AccessToken returns a valid access token, refreshing if needed.
func (m *Manager) AccessToken(ctx context.Context) (string, error) {
	if m.isValid() {
		m.mu.RLock()
		t := m.accessToken
		m.mu.RUnlock()
		return t, nil
	}
	if err := m.update(ctx); err != nil {
		return "", err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.accessToken == "" {
		return "", errors.New("empty access token after update")
	}
	return m.accessToken, nil
}

// GetSalts returns the current salt values.
func (m *Manager) GetSalts(ctx context.Context) ([5]int, error) {
	if !m.isValid() {
		if err := m.update(ctx); err != nil {
			return [5]int{}, err
		}
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.salts, nil
}

// RefreshToken returns a valid refresh token, refreshing if needed.
func (m *Manager) RefreshToken(ctx context.Context) (string, error) {
	if m.isValid() {
		m.mu.RLock()
		t := m.refreshToken
		m.mu.RUnlock()
		return t, nil
	}
	if err := m.update(ctx); err != nil {
		return "", err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.refreshToken == "" {
		return "", errors.New("empty refresh token after update")
	}
	return m.refreshToken, nil
}

// ForceUpdate forces a token refresh.
func (m *Manager) ForceUpdate(ctx context.Context) error {
	return m.update(ctx)
}

func (m *Manager) isValid() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.accessToken == "" || m.tokenTS.IsZero() {
		return false
	}
	return time.Since(m.tokenTS) < m.maxUpdatePeriod
}

func (m *Manager) update(ctx context.Context) error {
	_, err, _ := m.sf.Do("token_update", func() (any, error) {
		if m.isValid() {
			return struct{}{}, nil
		}

		resp, err := m.http.GetTokens(ctx)
		if err != nil {
			return nil, fmt.Errorf("get token: %w", err)
		}

		access, refresh, salts, ts, err := m.parseResponse(*resp)
		if err != nil {
			return nil, err
		}

		m.mu.Lock()
		m.accessToken = access
		m.refreshToken = refresh
		m.salts = salts
		if ts > 0 {
			m.tokenTS = time.Unix(ts, 0)
		} else {
			m.tokenTS = time.Now()
		}
		m.mu.Unlock()

		return struct{}{}, nil
	})
	return err
}

func (m *Manager) parseResponse(tr TokenResponse) (string, string, [5]int, int64, error) {
	salts := [5]int{tr.Salt1, tr.Salt2, tr.Salt3, tr.Salt4, tr.Salt5}

	idx, err := m.parser.indicesFromSalts(salts)
	if err != nil {
		return "", "", salts, 0, fmt.Errorf("wasm parse: %w", err)
	}

	parsedAccess := sliceSkipAt(tr.AccessToken, idx.access...)
	parsedRefresh := sliceSkipAt(tr.RefreshToken, idx.refresh...)

	sec := tr.ServerTime / 1000
	return parsedAccess, parsedRefresh, salts, sec, nil
}

// sliceSkipAt removes characters at specified positions from a string.
func sliceSkipAt(s string, positions ...int) string {
	if len(positions) == 0 {
		return s
	}

	ps := make([]int, len(positions))
	copy(ps, positions)
	sortInts(ps)

	b := []byte(s)
	var out []byte
	prev := 0
	for _, p := range ps {
		if p < 0 || p >= len(b) {
			continue
		}
		out = append(out, b[prev:p]...)
		prev = p + 1
	}
	out = append(out, b[prev:]...)
	return string(out)
}

// sortInts sorts a slice of ints in ascending order.
func sortInts(a []int) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j-1] > a[j]; j-- {
			a[j-1], a[j] = a[j], a[j-1]
		}
	}
}

// SetAuthHeader sets the Authorization header on the request.
func SetAuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Salter "+token)
}
