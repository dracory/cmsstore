package cmsstore

import (
	"strings"
	"testing"

	"github.com/dracory/uid"
	"github.com/stretchr/testify/require"
)

func TestGenerateShortIDReturnsLowercaseAndNonEmpty(t *testing.T) {
	generated := GenerateShortID()

	require.NotEmpty(t, generated)
	require.Equal(t, strings.ToLower(generated), generated)

	unshortened, err := uid.UnshortenCrockford(generated)
	require.NoError(t, err)
	require.NotEmpty(t, unshortened)
}

func TestNormalizeIDTrimsWhitespaceAndLowercases(t *testing.T) {
	normalized := NormalizeID("  AbC123XyZ  ")
	require.Equal(t, "abc123xyz", normalized)
}

func TestNormalizeIDEmptyString(t *testing.T) {
	require.Equal(t, "", NormalizeID("   "))
}

func TestIsShortIDReturnsTrueForNineCharIDs(t *testing.T) {
	require.True(t, IsShortID("abc123xyz"))
}

func TestIsShortIDReturnsTrueForTwentyOneCharIDs(t *testing.T) {
	require.True(t, IsShortID("abcdefghijklmnopqrstu"))
}

func TestIsShortIDReturnsFalseForOtherLengths(t *testing.T) {
	require.False(t, IsShortID("abcd"))
}

func TestShortenIDReturnsNormalizedNineCharID(t *testing.T) {
	result := ShortenID("  ABC123XYZ  ")
	require.Equal(t, "abc123xyz", result)
}

func TestShortenIDShortensValidThirtyTwoCharID(t *testing.T) {
	longID := uid.HumanUid()
	require.Len(t, longID, 32)

	short := ShortenID(longID)
	require.Len(t, short, 21)
	require.Equal(t, strings.ToLower(short), short)
}

func TestShortenIDReturnsNormalizedOriginalOnInvalidThirtyTwoCharID(t *testing.T) {
	invalidLong := "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"

	short := ShortenID(invalidLong)
	require.Equal(t, strings.ToLower(invalidLong), short)
}

func TestShortenIDReturnsNormalizedValueForOtherLengths(t *testing.T) {
	result := ShortenID("  SOME-HANDLE ")
	require.Equal(t, "some-handle", result)
}

func TestUnshortenIDReturnsOriginalForNonShortID(t *testing.T) {
	result := UnshortenID("  SOME-HANDLE ")
	require.Equal(t, "some-handle", result)
}

func TestUnshortenIDUnshortensValidTwentyOneCharID(t *testing.T) {
	longID := uid.HumanUid()
	short, err := uid.ShortenCrockford(longID)
	require.NoError(t, err)
	require.Len(t, short, 21)

	unshortened := UnshortenID(strings.ToUpper(short))
	require.Equal(t, longID, unshortened)
}

func TestUnshortenIDReturnsOriginalWhenCrockfordDecodeFails(t *testing.T) {
	invalidShort := "!!!!!!!!!"
	require.Equal(t, invalidShort, UnshortenID(invalidShort))
}

func TestUnshortenIDReturnsOriginalWhenLengthIsNotSupported(t *testing.T) {
	generated := GenerateShortID()
	require.False(t, IsShortID(generated))
	require.Equal(t, generated, UnshortenID(strings.ToUpper(generated)))
}

func TestIsSQLiteReturnsTrueForSQLiteNames(t *testing.T) {
	require.True(t, isSQLite("sqlite"))
	require.True(t, isSQLite("sqlite3"))
	require.True(t, isSQLite("MySQLiteDriver"))
}

func TestIsSQLiteReturnsFalseForNonSQLiteDrivers(t *testing.T) {
	require.False(t, isSQLite("postgres"))
	require.False(t, isSQLite("mysql"))
}
