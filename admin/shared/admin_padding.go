package shared

import "net/http"

func PaddingTopPx(r *http.Request) int {
	return readPadding(r, KeyPaddingTopPx)
}

func PaddingRightPx(r *http.Request) int {
	return readPadding(r, KeyPaddingRightPx)
}

func PaddingBottomPx(r *http.Request) int {
	return readPadding(r, KeyPaddingBottomPx)
}

func PaddingLeftPx(r *http.Request) int {
	return readPadding(r, KeyPaddingLeftPx)
}

func readPadding(r *http.Request, key string) int {
	value := r.Context().Value(key)
	if value == nil {
		return 2
	}

	padding, ok := value.(int)
	if !ok {
		return 2
	}

	if padding == 0 {
		return 2
	}

	return padding
}
