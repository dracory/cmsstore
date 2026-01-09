package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// argString extracts a string value from args map
func argString(args map[string]any, key string) string {
	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		return fmt.Sprintf("%.0f", t)
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

// argInt extracts an integer value from args map
func argInt(args map[string]any, key string) (int64, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case json.Number:
		i64, err := t.Int64()
		if err != nil {
			return 0, false
		}
		return i64, true
	case float64:
		return int64(t), true
	case int:
		return int64(t), true
	case int64:
		return t, true
	default:
		return 0, false
	}
}

// argBool extracts a boolean value from args map
func argBool(args map[string]any, key string) (bool, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return false, false
	}
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		vv := strings.TrimSpace(strings.ToLower(t))
		if vv == "true" || vv == "1" || vv == "yes" {
			return true, true
		}
		if vv == "false" || vv == "0" || vv == "no" {
			return false, true
		}
		return false, false
	default:
		return false, false
	}
}

// jsonRPCRequest represents a JSON-RPC 2.0 request
type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// jsonRPCErrorResponse creates a JSON-RPC error response
func jsonRPCErrorResponse(id any, code int, message string) map[string]any {
	return map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	}
}

// jsonRPCResultResponse creates a JSON-RPC result response
func jsonRPCResultResponse(id any, result any) map[string]any {
	return map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}
}
