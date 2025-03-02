# [Draft] Standardized Error Handling

## Summary
- **Problem**: Error handling across the CMS is inconsistent, making it difficult to debug issues and provide good user feedback
- **Solution**: Implement a standardized error handling system with proper error types, logging, and user feedback

## Background

Current error handling has several issues:
- Inconsistent error types and messages
- Mixed logging levels
- Unclear error recovery paths
- Limited context in error messages
- No standardized way to present errors to users
- Difficult to track error patterns

## Detailed Design

### 1. Error Types Hierarchy

```go
// Base error type
type CMSError struct {
    Code    string
    Message string
    Details map[string]interface{}
    Cause   error
}

// Specific error types
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

type PermissionError struct {
    CMSError
    RequiredPermission string
    UserPermissions   []string
}

type ProcessingError struct {
    CMSError
    Stage    string
    Context  map[string]interface{}
}
```

### 2. Error Codes System

```go
const (
    // Validation errors (1xxx)
    ErrValidation           = "CMS-1000"
    ErrInvalidInput        = "CMS-1001"
    ErrMissingRequired     = "CMS-1002"
    
    // Not found errors (2xxx)
    ErrNotFound            = "CMS-2000"
    ErrPageNotFound        = "CMS-2001"
    ErrTemplateNotFound    = "CMS-2002"
    ErrBlockNotFound       = "CMS-2003"
    
    // Permission errors (3xxx)
    ErrPermissionDenied    = "CMS-3000"
    ErrUnauthorized        = "CMS-3001"
    
    // Processing errors (4xxx)
    ErrProcessing          = "CMS-4000"
    ErrTemplateProcessing  = "CMS-4001"
    ErrBlockProcessing     = "CMS-4002"
    ErrCacheError          = "CMS-4003"
)
```

### 3. Error Creation Helpers

```go
// Error factory functions
func NewValidationError(field string, value interface{}, rule string) error {
    return &ValidationError{
        CMSError: CMSError{
            Code:    ErrValidation,
            Message: fmt.Sprintf("Validation failed for %s", field),
        },
        Field: field,
        Value: value,
        Rule:  rule,
    }
}

// Error wrapping with context
func WrapError(err error, code string, message string) error {
    return &CMSError{
        Code:    code,
        Message: message,
        Cause:   err,
    }
}
```

### 4. Standardized Error Handling

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
            Details: map[string]interface{}{
                "resourceType": e.ResourceType,
                "identifier":   e.Identifier,
            },
        }
        w.WriteHeader(http.StatusNotFound)
        
    // ... other error types
    }
    
    json.NewEncoder(w).Encode(response)
}
```

### 5. Logging Strategy

```go
type ErrorLogger struct {
    logger *slog.Logger
}

func (l *ErrorLogger) LogError(err error, r *http.Request) {
    // Extract error details
    code := "UNKNOWN"
    msg := err.Error()
    details := map[string]interface{}{}
    
    if cmsErr, ok := err.(*CMSError); ok {
        code = cmsErr.Code
        details = cmsErr.Details
    }
    
    // Log with context
    l.logger.Error("Error occurred",
        slog.String("code", code),
        slog.String("message", msg),
        slog.String("url", r.URL.String()),
        slog.String("method", r.Method),
        slog.Any("details", details),
    )
}
```

### 6. User Feedback System

```go
type UserErrorMessage struct {
    Title       string
    Message     string
    Suggestion  string
    Action      string
}

var errorMessages = map[string]UserErrorMessage{
    ErrPageNotFound: {
        Title:      "Page Not Found",
        Message:    "The requested page could not be found.",
        Suggestion: "Check the URL or navigate to another page.",
        Action:     "Return to homepage",
    },
    // ... other messages
}
```

## Alternatives Considered

1. **Simple Error Strings**
   - Pros: Simple implementation
   - Cons: Limited context, harder to handle systematically
   - Rejected: Need structured error handling

2. **Third-party Error Package**
   - Pros: Ready-made solution
   - Cons: Additional dependency, less control
   - Rejected: Need custom implementation for CMS-specific needs

3. **HTTP-only Error Handling**
   - Pros: Simpler model
   - Cons: Not suitable for all error types
   - Rejected: Need broader error handling

## Implementation Plan

1. Phase 1: Core Error Types (1 week)
   - Implement error hierarchy
   - Add error codes
   - Create helper functions

2. Phase 2: Error Handling (1 week)
   - Implement error handlers
   - Add logging system
   - Create user feedback system

3. Phase 3: Integration (2 weeks)
   - Update existing code
   - Add error documentation
   - Create examples

4. Phase 4: Testing (1 week)
   - Unit tests
   - Integration tests
   - Error handling scenarios

## Risks and Mitigations

1. **Migration Complexity**
   - Risk: Difficult to update all error handling
   - Mitigation: Gradual migration, tooling support

2. **Performance Impact**
   - Risk: Additional overhead from structured errors
   - Mitigation: Benchmark critical paths, optimize

3. **Error Proliferation**
   - Risk: Too many error types
   - Mitigation: Regular review, consolidation

4. **User Experience**
   - Risk: Technical errors confuse users
   - Mitigation: Clear user messages, actionable feedback 