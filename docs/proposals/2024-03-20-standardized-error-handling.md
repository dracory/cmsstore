# [Draft] Standardized Error Handling

## Status
**[Draft]** - Basic error handling implemented, enhanced system pending

## Summary
- **Problem**: Error handling could be more consistent with better context and user feedback
- **Solution**: Implement a standardized error handling system with proper error types, logging, and user feedback

## Current Implementation (As-Is)

The CMS Store currently uses standard Go error handling patterns:

**Current Error Handling:**
```go
// Standard Go errors with slog logging
func (frontend *frontend) PageRenderHtmlBySiteAndAlias(...) string {
    page, err := frontend.pageFindBySiteAndAlias(r.Context(), siteID, alias)
    
    if err != nil {
        frontend.logger.Error("PageRenderHtmlBySiteAndAlias: Error finding page", 
            "alias", alias, "error", err)
        return hb.NewDiv().Text("Error loading page").ToHTML()
    }
    
    if page == nil {
        frontend.logger.Warn("PageRenderHtmlBySiteAndAlias: Page not found", 
            "alias", alias)
        return hb.NewDiv().Text("Page with alias '").Text(alias).Text("' not found").ToHTML()
    }
    // ...
}
```

**Files:**
- Various implementation files use standard `error` returns
- `frontend/frontend.go` - Uses `slog.Logger` for structured logging
- `consts.go` - Basic error message constants:
  ```go
  const (
      ERROR_EMPTY_ARRAY     = "array cannot be empty"
      ERROR_EMPTY_STRING    = "string cannot be empty"
      ERROR_NEGATIVE_NUMBER = "number cannot be negative"
  )
  ```

**Current Features:**
- Standard Go `error` interface usage
- Structured logging with `slog.Logger` (Error, Warn, Info levels)
- Context-rich log messages (keys: alias, error, templateID, etc.)
- Graceful degradation (returns HTML error messages to users)
- Error propagation up the call stack

**Current Limitations:**
- No structured error types (just `error` interface)
- No error codes for programmatic handling
- No validation error specifics (field, rule, value)
- No centralized error handling middleware
- No user-friendly error message mapping
- No error metrics/monitoring

## Proposed Enhanced Design (To-Be)

### 1. Structured Error Types

```go
type CMSError struct {
    Code    string
    Message string
    Details map[string]interface{}
    Cause   error
}

type ValidationError struct {
    CMSError
    Field   string
    Value   interface{}
    Rule    string
}

type NotFoundError struct {
    CMSError
    ResourceType string
    Identifier   string
}
```

### 2. Error Codes System

```go
const (
    // Validation errors (1xxx)
    ErrValidation       = "CMS-1000"
    ErrInvalidInput     = "CMS-1001"
    ErrMissingRequired  = "CMS-1002"
    
    // Not found errors (2xxx)
    ErrNotFound         = "CMS-2000"
    ErrPageNotFound     = "CMS-2001"
    ErrTemplateNotFound = "CMS-2002"
)
```

### 3. Centralized Error Handler

```go
func handleError(err error, w http.ResponseWriter, r *http.Request) {
    var response ErrorResponse
    
    switch e := err.(type) {
    case *ValidationError:
        response = ErrorResponse{
            Code:    e.Code,
            Message: e.Message,
            Details: map[string]interface{}{
                "field": e.Field,
                "rule":  e.Rule,
            },
        }
        w.WriteHeader(http.StatusBadRequest)
    case *NotFoundError:
        response = ErrorResponse{
            Code:    e.Code,
            Message: e.Message,
        }
        w.WriteHeader(http.StatusNotFound)
    }
    
    json.NewEncoder(w).Encode(response)
}
```

### 4. User-Friendly Error Messages

```go
var errorMessages = map[string]UserErrorMessage{
    ErrPageNotFound: {
        Title:      "Page Not Found",
        Message:    "The requested page could not be found.",
        Suggestion: "Check the URL or navigate to another page.",
        Action:     "Return to homepage",
    },
}
```

## Implementation Status

| Feature | Status | Notes |
|---------|--------|-------|
| Standard Go errors | Implemented | `error` interface throughout |
| Structured logging | Implemented | `slog.Logger` with structured fields |
| Error context in logs | Implemented | Key-value pairs in log calls |
| Graceful degradation | Implemented | HTML error messages returned |
| Structured error types | Not implemented | No CMSError, ValidationError types |
| Error codes | Not implemented | No CMS-XXXX format codes |
| Error middleware | Not implemented | No centralized handler |
| User message mapping | Not implemented | No errorMessages map |
| Error metrics | Not implemented | No Prometheus error counters |

## Migration Strategy

### Phase 1: Error Types (Backward Compatible)
Create new error types that wrap existing errors:

```go
func NewValidationError(field string, err error) error {
    return &ValidationError{
        CMSError: CMSError{
            Code:    ErrValidation,
            Message: fmt.Sprintf("Validation failed for %s", field),
            Cause:   err,
        },
        Field: field,
    }
}
```

### Phase 2: Gradual Adoption
Update functions incrementally to return structured errors where beneficial.

## Files to Modify (If Implementing)

1. New: `errors.go` - Error type definitions and constants
2. New: `error_handler.go` - Centralized error handling middleware
3. `frontend/frontend.go` - Update to use structured errors where appropriate
4. `admin/` - Add user-friendly error message display
5. New: `error_metrics.go` - Prometheus error counters

## Risks and Mitigations

1. **Migration Complexity**
   - Risk: Difficult to update all error handling
   - Mitigation: Gradual adoption, backward compatible wrappers

2. **Over-Engineering**
   - Risk: Too complex error hierarchy
   - Mitigation: Start simple, add types as needed

3. **Performance**
   - Risk: Structured errors add overhead
   - Mitigation: Benchmark, keep error creation lightweight

4. **User Experience**
   - Risk: Technical errors confuse users
   - Mitigation: Clear mapping to user-friendly messages