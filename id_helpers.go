package cmsstore

import (
	"strings"
	"uuid"

	"github.com/dracory/uid"
)

// GenerateShortID generates a new unique ID using the standard library
// UUID v7 (time-ordered). Returns a 36-character lowercase ID
// (e.g., "01952f3a-8c5e-7f3a-9b2d-4e6f8a1c3d5e")
func GenerateShortID() string {
	return uuid.NewV7().String()
}

// NormalizeID normalizes an ID to lowercase for consistent lookups
// Handles both short (9-char, 21-char) and long (32-char) ID formats
func NormalizeID(id string) string {
	return strings.ToLower(strings.TrimSpace(id))
}

// IsShortID checks if an ID appears to be a shortened ID (9 or 21 chars)
// vs a long HumanUid ID (32 chars)
func IsShortID(id string) bool {
	length := len(id)
	return length == 9 || length == 21
}

// ShortenID shortens a long ID to its shortened form
// - 9-char IDs: Return as-is (already shortened TimestampMicro)
// - 32-char IDs: Shorten to 21-char using Crockford
// - Other lengths: Return as-is
func ShortenID(id string) string {
	id = NormalizeID(id)
	length := len(id)

	if length == 9 {
		// Already a shortened TimestampMicro ID
		return id
	}

	if length == 32 {
		// Long HumanUid - shorten to 21 chars
		shortened, err := uid.ShortenCrockford(id)
		if err != nil {
			return id
		}
		return strings.ToLower(shortened)
	}

	// Return as-is for other lengths
	return id
}

// UnshortenID attempts to unshorten a shortened ID back to its original form
// Returns the original ID if it cannot be unshortened
func UnshortenID(id string) string {
	id = NormalizeID(id)

	// Only attempt to unshorten if it looks like a shortened ID
	if !IsShortID(id) {
		return id
	}

	// Try to unshorten using Crockford
	unshortened, err := uid.UnshortenCrockford(id)
	if err == nil && unshortened != "" {
		return unshortened
	}

	// If unshortening fails, return original
	return id
}

// isSQLite checks if the database driver is SQLite (supports both "sqlite" and "sqlite3")
func isSQLite(driverName string) bool {
	return strings.Contains(strings.ToLower(driverName), "sqlite")
}
