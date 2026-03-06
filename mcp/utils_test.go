package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestArgString(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		key      string
		expected string
	}{
		{
			name:     "nil args",
			args:     nil,
			key:      "test",
			expected: "",
		},
		{
			name:     "empty args",
			args:     map[string]any{},
			key:      "test",
			expected: "",
		},
		{
			name:     "key not found",
			args:     map[string]any{"other": "value"},
			key:      "test",
			expected: "",
		},
		{
			name:     "string value",
			args:     map[string]any{"test": "hello"},
			key:      "test",
			expected: "hello",
		},
		{
			name:     "json.Number value",
			args:     map[string]any{"test": json.Number("123")},
			key:      "test",
			expected: "123",
		},
		{
			name:     "float64 value",
			args:     map[string]any{"test": 456.789},
			key:      "test",
			expected: "457",
		},
		{
			name:     "int value",
			args:     map[string]any{"test": 789},
			key:      "test",
			expected: "789",
		},
		{
			name:     "int64 value",
			args:     map[string]any{"test": int64(987)},
			key:      "test",
			expected: "987",
		},
		{
			name:     "bool true value",
			args:     map[string]any{"test": true},
			key:      "test",
			expected: "true",
		},
		{
			name:     "bool false value",
			args:     map[string]any{"test": false},
			key:      "test",
			expected: "false",
		},
		{
			name:     "nil value",
			args:     map[string]any{"test": nil},
			key:      "test",
			expected: "",
		},
		{
			name:     "unsupported type",
			args:     map[string]any{"test": []string{"array"}},
			key:      "test",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := argString(tt.args, tt.key)
			if result != tt.expected {
				t.Errorf("argString(%v, %q) = %q, want %q", tt.args, tt.key, result, tt.expected)
			}
		})
	}
}

func TestArgInt(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		key      string
		expected int64
		found    bool
	}{
		{
			name:     "nil args",
			args:     nil,
			key:      "test",
			expected: 0,
			found:    false,
		},
		{
			name:     "empty args",
			args:     map[string]any{},
			key:      "test",
			expected: 0,
			found:    false,
		},
		{
			name:     "key not found",
			args:     map[string]any{"other": "value"},
			key:      "test",
			expected: 0,
			found:    false,
		},
		{
			name:     "json.Number value",
			args:     map[string]any{"test": json.Number("123")},
			key:      "test",
			expected: 123,
			found:    true,
		},
		{
			name:     "json.Number invalid",
			args:     map[string]any{"test": json.Number("invalid")},
			key:      "test",
			expected: 0,
			found:    false,
		},
		{
			name:     "float64 value",
			args:     map[string]any{"test": 456.789},
			key:      "test",
			expected: 456,
			found:    true,
		},
		{
			name:     "int value",
			args:     map[string]any{"test": 789},
			key:      "test",
			expected: 789,
			found:    true,
		},
		{
			name:     "int64 value",
			args:     map[string]any{"test": int64(987)},
			key:      "test",
			expected: 987,
			found:    true,
		},
		{
			name:     "nil value",
			args:     map[string]any{"test": nil},
			key:      "test",
			expected: 0,
			found:    false,
		},
		{
			name:     "unsupported type",
			args:     map[string]any{"test": "string"},
			key:      "test",
			expected: 0,
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := argInt(tt.args, tt.key)
			if result != tt.expected || found != tt.found {
				t.Errorf("argInt(%v, %q) = (%d, %v), want (%d, %v)", tt.args, tt.key, result, found, tt.expected, tt.found)
			}
		})
	}
}

func TestArgBool(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		key      string
		expected bool
		found    bool
	}{
		{
			name:     "nil args",
			args:     nil,
			key:      "test",
			expected: false,
			found:    false,
		},
		{
			name:     "empty args",
			args:     map[string]any{},
			key:      "test",
			expected: false,
			found:    false,
		},
		{
			name:     "key not found",
			args:     map[string]any{"other": "value"},
			key:      "test",
			expected: false,
			found:    false,
		},
		{
			name:     "bool true value",
			args:     map[string]any{"test": true},
			key:      "test",
			expected: true,
			found:    true,
		},
		{
			name:     "bool false value",
			args:     map[string]any{"test": false},
			key:      "test",
			expected: false,
			found:    true,
		},
		{
			name:     "string true",
			args:     map[string]any{"test": "true"},
			key:      "test",
			expected: true,
			found:    true,
		},
		{
			name:     "string 1",
			args:     map[string]any{"test": "1"},
			key:      "test",
			expected: true,
			found:    true,
		},
		{
			name:     "string yes",
			args:     map[string]any{"test": "yes"},
			key:      "test",
			expected: true,
			found:    true,
		},
		{
			name:     "string false",
			args:     map[string]any{"test": "false"},
			key:      "test",
			expected: false,
			found:    true,
		},
		{
			name:     "string 0",
			args:     map[string]any{"test": "0"},
			key:      "test",
			expected: false,
			found:    true,
		},
		{
			name:     "string no",
			args:     map[string]any{"test": "no"},
			key:      "test",
			expected: false,
			found:    true,
		},
		{
			name:     "string with whitespace",
			args:     map[string]any{"test": "  true  "},
			key:      "test",
			expected: true,
			found:    true,
		},
		{
			name:     "string with mixed case",
			args:     map[string]any{"test": "TRUE"},
			key:      "test",
			expected: true,
			found:    true,
		},
		{
			name:     "string invalid",
			args:     map[string]any{"test": "invalid"},
			key:      "test",
			expected: false,
			found:    false,
		},
		{
			name:     "nil value",
			args:     map[string]any{"test": nil},
			key:      "test",
			expected: false,
			found:    false,
		},
		{
			name:     "unsupported type",
			args:     map[string]any{"test": 123},
			key:      "test",
			expected: false,
			found:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := argBool(tt.args, tt.key)
			if result != tt.expected || found != tt.found {
				t.Errorf("argBool(%v, %q) = (%v, %v), want (%v, %v)", tt.args, tt.key, result, found, tt.expected, tt.found)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		data      any
		expected  string
		expectErr bool
	}{
		{
			name:     "simple string",
			status:   http.StatusOK,
			data:     "hello",
			expected: `"hello"`,
		},
		{
			name:     "simple number",
			status:   http.StatusOK,
			data:     42,
			expected: `42`,
		},
		{
			name:     "simple boolean",
			status:   http.StatusOK,
			data:     true,
			expected: `true`,
		},
		{
			name:     "simple object",
			status:   http.StatusOK,
			data:     map[string]any{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "complex object",
			status:   http.StatusOK,
			data:     map[string]any{"status": "ok", "count": 5, "active": true},
			expected: `{"active":true,"count":5,"status":"ok"}`,
		},
		{
			name:     "null data",
			status:   http.StatusOK,
			data:     nil,
			expected: `null`,
		},
		{
			name:     "custom status",
			status:   http.StatusCreated,
			data:     map[string]any{"id": 123},
			expected: `{"id":123}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			writeJSON(w, tt.status, tt.data)

			// Check status code
			if w.Code != tt.status {
				t.Errorf("writeJSON() status = %d, want %d", w.Code, tt.status)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("writeJSON() Content-Type = %q, want %q", contentType, "application/json")
			}

			// Check body
			body := strings.TrimSpace(w.Body.String())
			if body != tt.expected {
				t.Errorf("writeJSON() body = %q, want %q", body, tt.expected)
			}
		})
	}
}

func TestJSONRPCErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		id       any
		code     int
		message  string
		expected map[string]any
	}{
		{
			name:    "string id",
			id:      "req123",
			code:    -32602,
			message: "Invalid params",
			expected: map[string]any{
				"jsonrpc": "2.0",
				"id":      "req123",
				"error": map[string]any{
					"code":    -32602,
					"message": "Invalid params",
				},
			},
		},
		{
			name:    "number id",
			id:      42,
			code:    -32601,
			message: "Method not found",
			expected: map[string]any{
				"jsonrpc": "2.0",
				"id":      42,
				"error": map[string]any{
					"code":    -32601,
					"message": "Method not found",
				},
			},
		},
		{
			name:    "null id",
			id:      nil,
			code:    -32700,
			message: "Parse error",
			expected: map[string]any{
				"jsonrpc": "2.0",
				"id":      nil,
				"error": map[string]any{
					"code":    -32700,
					"message": "Parse error",
				},
			},
		},
		{
			name:    "empty message",
			id:      "test",
			code:    0,
			message: "",
			expected: map[string]any{
				"jsonrpc": "2.0",
				"id":      "test",
				"error": map[string]any{
					"code":    0,
					"message": "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := jsonRPCErrorResponse(tt.id, tt.code, tt.message)

			// Check structure
			if result["jsonrpc"] != "2.0" {
				t.Errorf("jsonRPCErrorResponse() jsonrpc = %q, want %q", result["jsonrpc"], "2.0")
			}

			if result["id"] != tt.id {
				t.Errorf("jsonRPCErrorResponse() id = %v, want %v", result["id"], tt.id)
			}

			errorVal, ok := result["error"].(map[string]any)
			if !ok {
				t.Errorf("jsonRPCErrorResponse() error is not a map")
				return
			}

			if errorVal["code"] != tt.code {
				t.Errorf("jsonRPCErrorResponse() error.code = %v, want %v", errorVal["code"], tt.code)
			}

			if errorVal["message"] != tt.message {
				t.Errorf("jsonRPCErrorResponse() error.message = %q, want %q", errorVal["message"], tt.message)
			}
		})
	}
}

func TestJSONRPCResultResponse(t *testing.T) {
	tests := []struct {
		name     string
		id       any
		result   any
		expected map[string]any
	}{
		{
			name:   "string id with string result",
			id:     "req123",
			result: "success",
			expected: map[string]any{
				"jsonrpc": "2.0",
				"id":      "req123",
				"result":  "success",
			},
		},
		{
			name:   "number id with number result",
			id:     42,
			result: 123,
			expected: map[string]any{
				"jsonrpc": "2.0",
				"id":      42,
				"result":  123,
			},
		},
		{
			name:   "null id with object result",
			id:     nil,
			result: map[string]any{"status": "ok"},
			expected: map[string]any{
				"jsonrpc": "2.0",
				"id":      nil,
				"result":  map[string]any{"status": "ok"},
			},
		},
		{
			name:   "string id with null result",
			id:     "test",
			result: nil,
			expected: map[string]any{
				"jsonrpc": "2.0",
				"id":      "test",
				"result":  nil,
			},
		},
		{
			name:   "complex result",
			id:     "complex",
			result: map[string]any{"data": []int{1, 2, 3}, "count": 3},
			expected: map[string]any{
				"jsonrpc": "2.0",
				"id":      "complex",
				"result":  map[string]any{"count": 3, "data": []int{1, 2, 3}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := jsonRPCResultResponse(tt.id, tt.result)

			// Check structure
			if result["jsonrpc"] != "2.0" {
				t.Errorf("jsonRPCResultResponse() jsonrpc = %q, want %q", result["jsonrpc"], "2.0")
			}

			if result["id"] != tt.id {
				t.Errorf("jsonRPCResultResponse() id = %v, want %v", result["id"], tt.id)
			}

			// For map comparison, we need to use a different approach
			// since Go doesn't allow direct comparison of map types
			resultJSON, err := json.Marshal(result["result"])
			if err != nil {
				t.Errorf("jsonRPCResultResponse() failed to marshal result: %v", err)
				return
			}
			expectedJSON, err := json.Marshal(tt.result)
			if err != nil {
				t.Errorf("jsonRPCResultResponse() failed to marshal expected: %v", err)
				return
			}

			if string(resultJSON) != string(expectedJSON) {
				t.Errorf("jsonRPCResultResponse() result = %v, want %v", result["result"], tt.result)
			}
		})
	}
}

// TestJSONRPCRequest tests the jsonRPCRequest struct
func TestJSONRPCRequest(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected jsonRPCRequest
		wantErr  bool
	}{
		{
			name:     "valid request",
			jsonData: `{"jsonrpc": "2.0", "id": 1, "method": "test.method", "params": {"key": "value"}}`,
			expected: jsonRPCRequest{
				JSONRPC: "2.0",
				ID:      float64(1),
				Method:  "test.method",
				Params:  json.RawMessage(`{"key": "value"}`),
			},
			wantErr: false,
		},
		{
			name:     "valid request with string id",
			jsonData: `{"jsonrpc": "2.0", "id": "req123", "method": "test.method", "params": {}}`,
			expected: jsonRPCRequest{
				JSONRPC: "2.0",
				ID:      "req123",
				Method:  "test.method",
				Params:  json.RawMessage(`{}`),
			},
			wantErr: false,
		},
		{
			name:     "valid request with null id",
			jsonData: `{"jsonrpc": "2.0", "id": null, "method": "test.method", "params": []}`,
			expected: jsonRPCRequest{
				JSONRPC: "2.0",
				ID:      nil,
				Method:  "test.method",
				Params:  json.RawMessage(`[]`),
			},
			wantErr: false,
		},
		{
			name:     "invalid JSON",
			jsonData: `{"jsonrpc": "2.0", "id": 1, "method": "test.method", "params": {invalid}}`,
			wantErr:  true,
		},
		{
			name:     "missing required fields",
			jsonData: `{"jsonrpc": "2.0"}`,
			expected: jsonRPCRequest{
				JSONRPC: "2.0",
				ID:      nil,
				Method:  "",
				Params:  nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req jsonRPCRequest
			err := json.Unmarshal([]byte(tt.jsonData), &req)

			if (err != nil) != tt.wantErr {
				t.Errorf("jsonRPCRequest unmarshal error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if req.JSONRPC != tt.expected.JSONRPC {
					t.Errorf("jsonRPCRequest.JSONRPC = %q, want %q", req.JSONRPC, tt.expected.JSONRPC)
				}
				if req.ID != tt.expected.ID {
					t.Errorf("jsonRPCRequest.ID = %v, want %v", req.ID, tt.expected.ID)
				}
				if req.Method != tt.expected.Method {
					t.Errorf("jsonRPCRequest.Method = %q, want %q", req.Method, tt.expected.Method)
				}
				if !bytes.Equal(req.Params, tt.expected.Params) {
					t.Errorf("jsonRPCRequest.Params = %v, want %v", req.Params, tt.expected.Params)
				}
			}
		})
	}
}
