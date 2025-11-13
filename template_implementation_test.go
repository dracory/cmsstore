package cmsstore

import (
    "testing"

    "github.com/dracory/sb"
    "github.com/dromara/carbon/v2"
    "github.com/stretchr/testify/require"
)

func TestNewTemplateDefaults(t *testing.T) {
    template := NewTemplate()

    require.NotEmpty(t, template.ID(), "ID should be generated")
    require.NotEmpty(t, template.CreatedAt(), "CreatedAt should be set")
    require.NotEmpty(t, template.UpdatedAt(), "UpdatedAt should be set")
    require.Equal(t, "", template.Content(), "Content should default to an empty string")
    require.Equal(t, "", template.Editor(), "Editor should default to an empty string")
    require.Equal(t, "", template.Handle(), "Handle should default to an empty string")
    require.Equal(t, "", template.Memo(), "Memo should default to an empty string")
    require.Equal(t, TEMPLATE_STATUS_DRAFT, template.Status(), "Status should default to draft")
    require.Equal(t, sb.MAX_DATETIME, template.SoftDeletedAt(), "SoftDeletedAt should default to max datetime")
    require.False(t, template.IsSoftDeleted(), "New template should not be soft deleted")

    metas, err := template.Metas()
    require.NoError(t, err)
    require.Empty(t, metas, "Metas should default to an empty map")

    createdCarbon := template.CreatedAtCarbon()
    require.NotNil(t, createdCarbon)
    require.Equal(t, template.CreatedAt(), createdCarbon.ToDateTimeString(carbon.UTC))

    updatedCarbon := template.UpdatedAtCarbon()
    require.NotNil(t, updatedCarbon)
    require.Equal(t, template.UpdatedAt(), updatedCarbon.ToDateTimeString(carbon.UTC))

    softDeletedCarbon := template.SoftDeletedAtCarbon()
    require.NotNil(t, softDeletedCarbon)
    require.True(t, softDeletedCarbon.Gte(carbon.Now(carbon.UTC)), "SoftDeletedAt should be in the future by default")
}

func TestTemplateContentAndEditor(t *testing.T) {
    template := NewTemplate()

    template.SetContent("<div>{{ content }}</div>")
    template.SetEditor("html")

    require.Equal(t, "<div>{{ content }}</div>", template.Content())
    require.Equal(t, "html", template.Editor())
}

func TestTemplateMetasUpsertAndLookup(t *testing.T) {
    template := NewTemplate()

    err := template.SetMetas(map[string]string{"layout": "default"})
    require.NoError(t, err)
    require.Equal(t, "default", template.Meta("layout"))

    err = template.UpsertMetas(map[string]string{"layout": "custom", "author": "cms"})
    require.NoError(t, err)
    require.Equal(t, "custom", template.Meta("layout"))
    require.Equal(t, "cms", template.Meta("author"))
    require.Equal(t, "", template.Meta("missing"))
}

func TestTemplateSoftDeleteBehaviour(t *testing.T) {
    template := NewTemplate()
    require.False(t, template.IsSoftDeleted())

    past := carbon.Now(carbon.UTC).SubHour()
    template.SetSoftDeletedAt(past.ToDateTimeString(carbon.UTC))

    require.True(t, template.IsSoftDeleted(), "Template should be marked as soft deleted when past timestamp is set")
    require.Equal(t, past.ToDateTimeString(carbon.UTC), template.SoftDeletedAt())
}
