package cmsstore

import (
	"strings"
	"testing"

	"github.com/dracory/uid"
)

func TestGenerateShortIDReturnsLowercaseAndNonEmpty(t *testing.T) {
	generated := GenerateShortID()

	if generated == "" {
		t.Errorf("expected non-empty generated ID")
	}
	if strings.ToLower(generated) != generated {
		t.Errorf("expected lowercase, got %q", generated)
	}

	unshortened, err := uid.UnshortenCrockford(generated)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if unshortened == "" {
		t.Errorf("expected non-empty unshortened ID")
	}
}

func TestNormalizeIDTrimsWhitespaceAndLowercases(t *testing.T) {
	normalized := NormalizeID("  AbC123XyZ  ")
	if normalized != "abc123xyz" {
		t.Errorf("expected %q, got %q", "abc123xyz", normalized)
	}
}

func TestNormalizeIDEmptyString(t *testing.T) {
	if NormalizeID("   ") != "" {
		t.Errorf("expected empty string")
	}
}

func TestIsShortIDReturnsTrueForNineCharIDs(t *testing.T) {
	if !IsShortID("abc123xyz") {
		t.Errorf("expected true for 9 char ID")
	}
}

func TestIsShortIDReturnsTrueForTwentyOneCharIDs(t *testing.T) {
	if !IsShortID("abcdefghijklmnopqrstu") {
		t.Errorf("expected true for 21 char ID")
	}
}

func TestIsShortIDReturnsFalseForOtherLengths(t *testing.T) {
	if IsShortID("abcd") {
		t.Errorf("expected false for 4 char ID")
	}
}

func TestShortenIDReturnsNormalizedNineCharID(t *testing.T) {
	result := ShortenID("  ABC123XYZ  ")
	if result != "abc123xyz" {
		t.Errorf("expected %q, got %q", "abc123xyz", result)
	}
}

func TestShortenIDShortensValidThirtyTwoCharID(t *testing.T) {
	longID := uid.HumanUid()
	if len(longID) != 32 {
		t.Errorf("expected 32 char ID, got %d", len(longID))
	}

	short := ShortenID(longID)
	if len(short) != 21 {
		t.Errorf("expected 21 char short ID, got %d", len(short))
	}
	if strings.ToLower(short) != short {
		t.Errorf("expected lowercase, got %q", short)
	}
}

func TestShortenIDReturnsNormalizedOriginalOnInvalidThirtyTwoCharID(t *testing.T) {
	invalidLong := "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"

	short := ShortenID(invalidLong)
	if short != strings.ToLower(invalidLong) {
		t.Errorf("expected %q, got %q", strings.ToLower(invalidLong), short)
	}
}

func TestShortenIDReturnsNormalizedValueForOtherLengths(t *testing.T) {
	result := ShortenID("  SOME-HANDLE ")
	if result != "some-handle" {
		t.Errorf("expected %q, got %q", "some-handle", result)
	}
}

func TestUnshortenIDReturnsOriginalForNonShortID(t *testing.T) {
	result := UnshortenID("  SOME-HANDLE ")
	if result != "some-handle" {
		t.Errorf("expected %q, got %q", "some-handle", result)
	}
}

func TestUnshortenIDUnshortensValidTwentyOneCharID(t *testing.T) {
	longID := uid.HumanUid()
	short, err := uid.ShortenCrockford(longID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(short) != 21 {
		t.Errorf("expected 21 char short ID, got %d", len(short))
	}

	unshortened := UnshortenID(strings.ToUpper(short))
	if unshortened != longID {
		t.Errorf("expected %q, got %q", longID, unshortened)
	}
}

func TestUnshortenIDReturnsOriginalWhenCrockfordDecodeFails(t *testing.T) {
	invalidShort := "!!!!!!!!!"
	if UnshortenID(invalidShort) != invalidShort {
		t.Errorf("expected original when decode fails")
	}
}

func TestUnshortenIDReturnsOriginalWhenLengthIsNotSupported(t *testing.T) {
	generated := GenerateShortID()
	if IsShortID(generated) {
		t.Errorf("expected generated to NOT be short ID")
	}
	if UnshortenID(strings.ToUpper(generated)) != generated {
		t.Errorf("expected original when length not supported")
	}
}

func TestIsSQLiteReturnsTrueForSQLiteNames(t *testing.T) {
	if !isSQLite("sqlite") {
		t.Errorf("expected true for sqlite")
	}
	if !isSQLite("sqlite3") {
		t.Errorf("expected true for sqlite3")
	}
	if !isSQLite("MySQLiteDriver") {
		t.Errorf("expected true for MySQLiteDriver")
	}
}

func TestIsSQLiteReturnsFalseForNonSQLiteDrivers(t *testing.T) {
	if isSQLite("postgres") {
		t.Errorf("expected false for postgres")
	}
	if isSQLite("mysql") {
		t.Errorf("expected false for mysql")
	}
}
