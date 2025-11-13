package cmsstore

import (
    "testing"

    "github.com/dromara/carbon/v2"
    "github.com/stretchr/testify/require"
    "github.com/dracory/sb"
)

func TestNewTranslationDefaults(t *testing.T) {
    translation := NewTranslation()

    require.NotEmpty(t, translation.ID(), "ID should be generated")
    require.NotEmpty(t, translation.CreatedAt(), "CreatedAt should be set")
    require.NotEmpty(t, translation.UpdatedAt(), "UpdatedAt should be set")
    require.Equal(t, sb.MAX_DATETIME, translation.SoftDeletedAt(), "SoftDeletedAt should default to max datetime")
    require.Equal(t, TEMPLATE_STATUS_DRAFT, translation.Status(), "Status should default to draft")
    require.False(t, translation.IsSoftDeleted(), "New translation should not be marked as soft deleted")

    content, err := translation.Content()
    require.NoError(t, err)
    require.Empty(t, content, "Content should default to an empty map")

    metas, err := translation.Metas()
    require.NoError(t, err)
    require.Empty(t, metas, "Metas should default to an empty map")

    createdCarbon := translation.CreatedAtCarbon()
    require.NotNil(t, createdCarbon)
    require.False(t, createdCarbon.IsZero(), "CreatedAtCarbon should be parseable")

    updatedCarbon := translation.UpdatedAtCarbon()
    require.NotNil(t, updatedCarbon)
    require.False(t, updatedCarbon.IsZero(), "UpdatedAtCarbon should be parseable")

    softDeletedCarbon := translation.SoftDeletedAtCarbon()
    require.NotNil(t, softDeletedCarbon)
    require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)), "SoftDeletedAt should be in the future by default")
}

func TestTranslationContentRoundTrip(t *testing.T) {
    translation := NewTranslation()

    expectedContent := map[string]string{
        "en": "Hello",
        "fr": "Bonjour",
    }

    err := translation.SetContent(expectedContent)
    require.NoError(t, err)

    content, err := translation.Content()
    require.NoError(t, err)
    require.Equal(t, expectedContent, content)
}

func TestTranslationMetasUpsertAndMetaLookup(t *testing.T) {
    translation := NewTranslation()

    err := translation.SetMetas(map[string]string{"locale": "en"})
    require.NoError(t, err)
    require.Equal(t, "en", translation.Meta("locale"))

    err = translation.UpsertMetas(map[string]string{"locale": "fr", "category": "general"})
    require.NoError(t, err)
    require.Equal(t, "fr", translation.Meta("locale"))
    require.Equal(t, "general", translation.Meta("category"))

    require.Equal(t, "", translation.Meta("missing"))
}

func TestTranslationSoftDeleteBehaviour(t *testing.T) {
    translation := NewTranslation()
    require.False(t, translation.IsSoftDeleted())

    past := carbon.Now(carbon.UTC).SubHour()
    translation.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))

    require.True(t, translation.IsSoftDeleted(), "Translation should be marked as soft deleted when past timestamp is set")
    require.Equal(t, past.ToDateTimeString(carbon.UTC), translation.SoftDeletedAt())
}
