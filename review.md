# Code Review: nepseauth Go Library

**Review Date:** 2025-12-30
**Reviewed Files:** All files in `nepse/` and `auth/` packages
**Severity Levels:** Critical | High | Medium | Low | Nitpick

---

## Executive Summary

This Go library for the NEPSE API is functional but has several areas where it deviates from Go idioms and best practices. The main issues include:

1. **Fat interface anti-pattern** - The `Client` interface is too large
2. **Inconsistent error handling** - Mix of wrapped errors and custom error types
3. **Concurrency safety issues** - Unprotected mutations
4. **Code duplication** - Repetitive methods that could be generalized
5. **Type design issues** - Anonymous structs, stringly-typed configuration

---

## Critical Issues

### 1. Race Condition in `SetTLSVerification` (`http_client.go:258-264`)

```go
func (h *HTTPClient) SetTLSVerification(enabled bool) {
    if transport, ok := h.client.Transport.(*http.Transport); ok {
        transport.TLSClientConfig.InsecureSkipVerify = !enabled
    }
    h.options.TLSVerification = enabled
}
```

**Problem:** This mutates the TLS config and options without any synchronization while requests might be in flight. This is a data race.

**Go Idiom:** Either use `sync.Mutex` to protect mutations, or make the client immutable after creation (preferred for HTTP clients).

**Recommendation:** Remove this method. If TLS config needs to change, create a new client.

---

### 2. Context Misuse in `BatchRequest` (`nepse.go:147-151`)

```go
type BatchRequest struct {
    MaxConcurrency int
    Timeout        context.Context  // WRONG
}
```

**Problem:** A `context.Context` should never be stored in a struct. Contexts are meant to flow through call stacks, not be stored.

**Go Idiom:** Timeouts should be `time.Duration`. Context should be passed as the first parameter to functions.

**Fix:**
```go
type BatchRequest struct {
    MaxConcurrency int
    Timeout        time.Duration
}
```

---

## High Severity Issues

### 3. Fat Interface Anti-Pattern (`client.go:9-78`)

The `Client` interface has 50+ methods. This violates the Go proverb: "The bigger the interface, the weaker the abstraction."

**Go Idiom:** Prefer small, focused interfaces. The standard library's `io.Reader` is one method.

**Recommendation:** Break into smaller interfaces:
```go
type MarketDataReader interface {
    GetMarketSummary(ctx context.Context) (*MarketSummary, error)
    GetMarketStatus(ctx context.Context) (*MarketStatus, error)
    GetNepseIndex(ctx context.Context) (*NepseIndex, error)
}

type SecurityReader interface {
    GetSecurityList(ctx context.Context) ([]Security, error)
    FindSecurity(ctx context.Context, securityID int32) (*Security, error)
}

// Client can embed these if needed
type Client interface {
    MarketDataReader
    SecurityReader
    // ...
    io.Closer
}
```

---

### 4. Inconsistent Error Wrapping (`market_data.go`)

The codebase inconsistently uses `fmt.Errorf` to wrap already-contextualized `NepseError`:

```go
// Line 21 - Double wrapping
return nil, fmt.Errorf("failed to get market summary: %w", err)
// But err is already a NepseError with message and context
```

**Problem:** This creates redundant error chains:
```
"failed to get market summary: nepse internal_error: failed to decode response"
```

**Go Idiom:** Either use custom error types consistently OR use `fmt.Errorf`, but not both. When using custom errors, return them directly.

**Fix:** Since `NepseError` already has context, return it directly:
```go
if err != nil {
    return nil, err  // NepseError already has context
}
```

Or if additional context is needed, ensure `NepseError` doesn't duplicate:
```go
return nil, fmt.Errorf("market summary: %w", err)  // Short prefix only
```

---

### 5. The `Is` Method Implementation is Wrong (`errors.go:42-48`)

```go
func (e *NepseError) Is(target error) bool {
    if target, ok := target.(*NepseError); ok {
        return e.Type == target.Type
    }
    return false
}
```

**Problem:** This only compares `Type`, ignoring the wrapped error. It also shadows the `target` parameter.

**Go Idiom:** `Is` should check if the error matches AND check the wrapped error chain.

**Fix:**
```go
func (e *NepseError) Is(target error) bool {
    t, ok := target.(*NepseError)
    if !ok {
        return false
    }
    return e.Type == t.Type
}
```

Note: For comparing wrapped errors, `errors.Is` will use `Unwrap` automatically.

---

### 6. Stringly-Typed Configuration (`config.go:6-8`)

```go
type Config struct {
    BaseURL      string
    APIEndpoints map[string]string  // PROBLEM
    Headers      map[string]string
}
```

**Problem:** Using `map[string]string` for endpoints is error-prone. Typos in keys won't be caught at compile time:
```go
h.config.APIEndpoints["market_sumary"]  // Typo compiles fine, fails at runtime
```

**Go Idiom:** Use a struct with typed fields:
```go
type Endpoints struct {
    MarketSummary string
    TopGainers    string
    TopLosers     string
    // ...
}

// Usage: h.config.Endpoints.MarketSummary
```

---

### 7. Anonymous Struct Inside Type Definition (`types.go:183-193`)

```go
type MarketDepth struct {
    // ...
    BuyDepth []struct {
        Price    float64 `json:"price"`
        Quantity int64   `json:"quantity"`
        Orders   int32   `json:"orders"`
    } `json:"buyDepth"`
}
```

**Problem:** Anonymous structs cannot be reused, documented separately, or used in function signatures.

**Go Idiom:** Define named types for all public data structures:
```go
type DepthEntry struct {
    Price    float64 `json:"price"`
    Quantity int64   `json:"quantity"`
    Orders   int32   `json:"orders"`
}

type MarketDepth struct {
    SecurityID   int32        `json:"securityId"`
    Symbol       string       `json:"symbol"`
    SecurityName string       `json:"securityName"`
    BuyDepth     []DepthEntry `json:"buyDepth"`
    SellDepth    []DepthEntry `json:"sellDepth"`
}
```

---

## Medium Severity Issues

### 8. Unused Code - Dead Function (`http_client.go:175-179`)

```go
func (h *HTTPClient) getResponseBody(resp *http.Response) (io.ReadCloser, error) {
    return resp.Body, nil
}
```

**Problem:** This function does nothing but return its input. It's dead code.

**Go Idiom:** Delete code that serves no purpose. Comments from "maybe we'll need it later" thinking become lies over time.

**Recommendation:** Remove the function and use `resp.Body` directly.

---

### 9. Function Comments Don't Follow Go Convention (`graphs.go`)

```go
// GetDailyNepseIndexGraph Index Graph Methods
func (h *HTTPClient) GetDailyNepseIndexGraph(...)
```

**Go Idiom:** Comments should be complete sentences starting with the function name:
```go
// GetDailyNepseIndexGraph retrieves the daily NEPSE index graph data.
func (h *HTTPClient) GetDailyNepseIndexGraph(...)
```

This also affects `graphs.go:43`, `graphs.go:148`, and many others.

---

### 10. Massive Code Duplication in Graph Methods (`graphs.go`)

All 17 graph methods follow the exact same pattern:
```go
func (h *HTTPClient) GetDaily*Graph(ctx context.Context) (*GraphResponse, error) {
    var arr []GraphDataPoint
    if err := h.apiRequest(ctx, h.config.APIEndpoints["*_graph"], &arr); err != nil {
        return nil, fmt.Errorf("failed to get daily * graph: %w", err)
    }
    return &GraphResponse{Data: arr}, nil
}
```

**Go Idiom:** DRY (Don't Repeat Yourself). Use a helper:
```go
func (h *HTTPClient) getIndexGraph(ctx context.Context, endpointKey, name string) (*GraphResponse, error) {
    var arr []GraphDataPoint
    if err := h.apiRequest(ctx, h.config.APIEndpoints[endpointKey], &arr); err != nil {
        return nil, fmt.Errorf("failed to get %s graph: %w", name, err)
    }
    return &GraphResponse{Data: arr}, nil
}

func (h *HTTPClient) GetDailyNepseIndexGraph(ctx context.Context) (*GraphResponse, error) {
    return h.getIndexGraph(ctx, "nepse_index_daily_graph", "daily NEPSE index")
}
```

---

### 11. Complex Nested Function Definition (`market_data.go:139-217`)

The `GetSupplyDemand` method defines a recursive closure `decodeWithRetry` inside itself. This is 80+ lines of complex nested code.

**Go Idiom:** Extract complex logic to private methods:
```go
func (h *HTTPClient) GetSupplyDemand(ctx context.Context) ([]SupplyDemandEntry, error) {
    // Simple array first
    var arr []SupplyDemandEntry
    if err := h.apiRequest(ctx, endpoint, &arr); err == nil {
        return arr, nil
    }
    // Fallback
    return h.getSupplyDemandFallback(ctx)
}

func (h *HTTPClient) getSupplyDemandFallback(ctx context.Context) ([]SupplyDemandEntry, error) {
    // ... extracted logic
}
```

---

### 12. Ignoring `bool` Return from `singleflight.Do` (`auth/token.go:158`)

```go
_, err, _ := m.sf.Do("token_update", func() (any, error) {
```

**Problem:** The third return value indicates if the result was shared (from another goroutine). This might be useful for logging/metrics.

**Go Idiom:** If you intentionally ignore a value, consider documenting why. But here it's probably fine.

---

### 13. Manual Sorting Instead of Using `sort` Package (`auth/token.go:218-224`)

```go
for i := 1; i < len(ps); i++ {
    j := i
    for j > 0 && ps[j-1] > ps[j] {
        ps[j-1], ps[j] = ps[j], ps[j-1]
        j--
    }
}
```

**Go Idiom:** Use the standard library:
```go
import "sort"
sort.Ints(ps)
```

---

### 14. Empty String Check for Optional Parameters (`market_data.go:283-285`)

```go
if businessDate != "" {
    endpoint += "?businessDate=" + businessDate + "&size=500"
}
```

**Go Idiom:** Use functional options pattern or a dedicated options struct:
```go
type PriceOptions struct {
    BusinessDate string
    PageSize     int
}

func (h *HTTPClient) GetTodaysPrices(ctx context.Context, opts *PriceOptions) ([]TodayPrice, error)
```

---

### 15. Unused Parameter (`http_client.go:182`)

```go
func (h *HTTPClient) setCommonHeaders(req *http.Request, _ bool) {
```

**Problem:** The second parameter is always ignored.

**Go Idiom:** Remove unused parameters. If it was for future use, you're not gonna need it (YAGNI).

---

## Low Severity Issues

### 16. Inconsistent Indentation (`config.go:53-62`)

The `Headers` map uses different indentation than `APIEndpoints`:
```go
APIEndpoints: map[string]string{
    "price_volume": "/api/...",  // tab-indented
    // ...
},
Headers: map[string]string{
    "User-Agent": "...",  // space-indented
```

**Go Idiom:** Run `gofmt` or `goimports`. Use consistent formatting.

---

### 17. Missing Validation in `DefaultOptions` (`client.go:106-115`)

`DefaultOptions()` returns hardcoded values but doesn't validate them.

**Go Idiom:** Consider using a `Validate()` method:
```go
func (o *Options) Validate() error {
    if o.HTTPTimeout <= 0 {
        return errors.New("HTTPTimeout must be positive")
    }
    if o.MaxRetries < 0 {
        return errors.New("MaxRetries cannot be negative")
    }
    return nil
}
```

---

### 18. Exported But Likely Unused Type (`nepse.go:147-151`)

`BatchRequest` is exported but doesn't appear to be used anywhere in the codebase.

**Go Idiom:** Don't export types until they're needed. Unexported types are easier to change.

---

### 19. Using `int32` for IDs Without Type Alias (`types.go`)

```go
ID int32 `json:"id"`
```

**Go Idiom:** Consider type aliases for domain concepts:
```go
type SecurityID int32
type CompanyID int32
```

This adds type safety - you can't accidentally pass a `CompanyID` where `SecurityID` is expected.

---

### 20. No `context.Context` Timeout/Cancellation Check in Retry Loop (`http_client.go:146-171`)

```go
for attempt := 0; attempt <= h.options.MaxRetries; attempt++ {
    if attempt > 0 {
        time.Sleep(delay)  // Doesn't respect context cancellation
    }
    // ...
}
```

**Go Idiom:** Check context cancellation before sleeping:
```go
select {
case <-ctx.Done():
    return nil, ctx.Err()
case <-time.After(delay):
}
```

---

## Nitpicks

### 21. Variable Naming (`auth/token.go:352`)

```go
pp, err := call5(p.ndx, s1, s2, s4, s3, s5)
```

Using `pp` because `p` is already taken is awkward. Consider `pIdx` or rename the outer `p`.

---

### 22. Magic Numbers (`auth/token.go:75`)

```go
maxUpdatePeriod: 45 * time.Second,
```

**Recommendation:** Define as a constant with documentation:
```go
const (
    // DefaultTokenTTL is the maximum time a token is considered valid.
    // NEPSE tokens expire after ~60 seconds; we refresh at 45s for safety.
    DefaultTokenTTL = 45 * time.Second
)
```

---

### 23. Package-Level Error Variables (`nepse.go:103-118`)

```go
var (
    ErrTokenExpired = NewTokenExpiredError()
    ErrNetworkError = NewNetworkError(nil)
    // ...
)
```

**Recommendation:** These sentinel errors are good but should be documented for use with `errors.Is()`:
```go
// ErrTokenExpired is returned when the access token has expired.
// Use errors.Is(err, ErrTokenExpired) to check.
var ErrTokenExpired = NewTokenExpiredError()
```

---

### 24. Consider `io.Closer` Instead of Custom `Close` (`client.go:77`)

```go
Close(ctx context.Context) error
```

**Go Idiom:** The standard `io.Closer` interface is:
```go
Close() error
```

If context is needed, document why. In this case, the context is passed through but often ignored.

---

## Recommendations Summary

### Quick Wins (Low Effort, High Value)
1. Remove `SetTLSVerification` method (race condition)
2. Fix `BatchRequest.Timeout` type from `context.Context` to `time.Duration`
3. Delete unused `getResponseBody` function
4. Use `sort.Ints()` instead of manual sorting
5. Run `gofmt` to fix indentation

### Medium Effort
1. Extract graph helper method to reduce duplication
2. Extract `GetSupplyDemand` fallback logic to separate method
3. Define named types for `MarketDepth.BuyDepth` and `SellDepth`
4. Add context cancellation check in retry loop

### Larger Refactors
1. Break `Client` interface into smaller, composable interfaces
2. Replace `map[string]string` endpoints with typed struct
3. Standardize error handling (either all `NepseError` or all `fmt.Errorf`, not mixed)
4. Consider functional options pattern for methods with optional parameters

---

## Positive Aspects

The codebase does several things well:

1. **Proper use of `context.Context`** - All public methods accept context as first parameter
2. **Interface-based design** - The `Client` interface allows for mocking in tests
3. **Retry logic with exponential backoff** - Properly implemented
4. **`singleflight` for token refresh** - Prevents thundering herd
5. **Embedded WASM** - Clean approach for the token parsing
6. **Minimal dependencies** - Only 2 external dependencies
7. **Good package documentation** - The `nepse.go` package comment is comprehensive

---

## Conclusion

The library is functional and demonstrates understanding of Go, but would benefit from stricter adherence to Go idioms, particularly around interface design, error handling consistency, and code duplication. The critical race condition in `SetTLSVerification` should be addressed immediately.

**Overall Assessment:** 6/10 - Functional but not idiomatic Go
