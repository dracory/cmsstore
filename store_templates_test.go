package cmsstore

import (
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	_ "modernc.org/sqlite"
)

func TestStoreTemplateCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_create",
		PageTableName:      "page_table_create",
		SiteTableName:      "site_table_create",
		TemplateTableName:  "template_table_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	template := NewTemplate().SetSiteID("Site1")

	err = store.TemplateCreate(template)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreTemplateFindByHandle(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_find_by_handle",
		PageTableName:      "page_table_find_by_handle",
		SiteTableName:      "site_table_find_by_handle",
		TemplateTableName:  "template_table_find_by_handle",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	template := NewTemplate().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE).
		SetHandle("test-handle")

	err = template.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.TemplateCreate(template)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	templateFound, errFind := store.TemplateFindByHandle(template.Handle())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if templateFound == nil {
		t.Fatal("Template MUST NOT be nil")
	}

	if templateFound.ID() != template.ID() {
		t.Fatal("IDs do not match")
	}

	if templateFound.Status() != template.Status() {
		t.Fatal("Statuses do not match")
	}

	if templateFound.Meta("education_1") != template.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if templateFound.Meta("education_2") != template.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if templateFound.Meta("education_3") != template.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStoreTemplateFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_find_by_id",
		PageTableName:      "page_table_find_by_id",
		SiteTableName:      "site_table_find_by_id",
		TemplateTableName:  "template_table_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	template := NewTemplate().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE)

	err = template.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.TemplateCreate(template)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	templateFound, errFind := store.TemplateFindByID(template.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if templateFound == nil {
		t.Fatal("Template MUST NOT be nil")
	}

	if templateFound.ID() != template.ID() {
		t.Fatal("IDs do not match")
	}

	if templateFound.Status() != template.Status() {
		t.Fatal("Statuses do not match")
	}

	if templateFound.Meta("education_1") != template.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if templateFound.Meta("education_2") != template.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if templateFound.Meta("education_3") != template.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStoreTemplateSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_soft_delete",
		PageTableName:      "page_table_soft_delete",
		SiteTableName:      "site_table_soft_delete",
		TemplateTableName:  "template_table_soft_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	template := NewTemplate().
		SetSiteID("Site1")

	err = store.TemplateCreate(template)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.TemplateSoftDeleteByID(template.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if template.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Template MUST NOT be soft deleted")
	}

	templateFound, errFind := store.TemplateFindByID(template.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if templateFound != nil {
		t.Fatal("Template MUST be nil")
	}
	query := NewTemplateQuery().SetWithSoftDeleted(true)

	query, err = query.SetID(template.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	query, err = query.SetLimit(1)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	templateFindWithSoftDeleted, err := store.TemplateList(query)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(templateFindWithSoftDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(templateFindWithSoftDeleted[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Template MUST be soft deleted", template.SoftDeletedAt())
	}

	if !templateFindWithSoftDeleted[0].IsSoftDeleted() {
		t.Fatal("Template MUST be soft deleted")
	}
}

func TestStoreTemplateDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_delete",
		PageTableName:      "page_table_delete",
		SiteTableName:      "site_table_delete",
		TemplateTableName:  "template_table_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	template := NewTemplate().
		SetSiteID("Site1")

	err = store.TemplateCreate(template)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.TemplateDeleteByID(template.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	query := NewTemplateQuery().SetWithSoftDeleted(true)

	query, err = query.SetID(template.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	query, err = query.SetLimit(1)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	templateFindWithDeleted, err := store.TemplateList(query)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(templateFindWithDeleted) != 0 {
		t.Fatal("Template MUST be deleted, but it is not")
	}
}

func TestStoreTemplateUpdate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_update",
		PageTableName:      "page_table_update",
		SiteTableName:      "site_table_update",
		TemplateTableName:  "template_table_update",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	template := NewTemplate().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE)

	err = store.TemplateCreate(template)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	template.SetStatus(PAGE_STATUS_INACTIVE)

	err = template.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.TemplateUpdate(template)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	templateFound, errFind := store.TemplateFindByID(template.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if templateFound == nil {
		t.Fatal("Template MUST NOT be nil")
	}

	if templateFound.Status() != PAGE_STATUS_INACTIVE {
		t.Fatal("Status MUST be INACTIVE, found: ", templateFound.Status())
	}

	metas, err := templateFound.Metas()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(metas) < 3 {
		t.Fatal("Metas MUST be 3, found: ", len(metas))
	}

	if metas["education_1"] != "Education 1" {
		t.Fatal("Metas do not match")
	}

	if metas["education_2"] != "Education 2" {
		t.Fatal("Metas do not match")
	}

	if metas["education_3"] != "Education 3" {
		t.Fatal("Metas do not match")
	}
}
