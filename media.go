package cmsstore

import (
	"encoding/json"
	"strings"

	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

type mediaImplementation struct {
	dataobject.DataObject
}

var _ MediaInterface = (*mediaImplementation)(nil)

func NewMedia() MediaInterface {
	o := &mediaImplementation{}
	o.SetID(GenerateShortID())
	o.SetEntityID("")
	o.SetEntityType("")
	o.SetTitle("")
	o.SetDescription("")
	o.SetMemo("")
	o.SetURL("")
	o.SetType("")
	o.SetSize("0")
	o.SetExtension("")
	o.SetSequence("0")
	o.SetStatus(MEDIA_STATUS_DRAFT)
	o.SetHandle("")
	o.SetSiteID("")
	o.SetMetas(map[string]string{})
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(MAX_DATETIME)
	return o
}

func NewMediaFromExistingData(data map[string]string) *mediaImplementation {
	o := &mediaImplementation{}
	o.Hydrate(data)
	return o
}

func (o *mediaImplementation) IsActive() bool {
	return o.Status() == MEDIA_STATUS_ACTIVE
}

func (o *mediaImplementation) IsInactive() bool {
	return o.Status() == MEDIA_STATUS_INACTIVE
}

func (o *mediaImplementation) IsDraft() bool {
	return o.Status() == MEDIA_STATUS_DRAFT
}

func (o *mediaImplementation) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

func (o *mediaImplementation) IsImage() bool {
	return strings.HasPrefix(o.Type(), "image/")
}

func (o *mediaImplementation) IsVideo() bool {
	return strings.HasPrefix(o.Type(), "video/")
}

func (o *mediaImplementation) MarshalToVersioning() (string, error) {
	versionedData := map[string]string{}

	for k, v := range o.Data() {
		if k == COLUMN_CREATED_AT ||
			k == COLUMN_UPDATED_AT ||
			k == COLUMN_SOFT_DELETED_AT {
			continue
		}
		versionedData[k] = v
	}

	b, err := json.Marshal(versionedData)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (o *mediaImplementation) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *mediaImplementation) SetID(id string) MediaInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *mediaImplementation) EntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *mediaImplementation) SetEntityID(entityID string) MediaInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *mediaImplementation) EntityType() string {
	return o.Get(COLUMN_ENTITY_TYPE)
}

func (o *mediaImplementation) SetEntityType(entityType string) MediaInterface {
	o.Set(COLUMN_ENTITY_TYPE, entityType)
	return o
}

func (o *mediaImplementation) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *mediaImplementation) SetTitle(title string) MediaInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *mediaImplementation) Description() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *mediaImplementation) SetDescription(description string) MediaInterface {
	o.Set(COLUMN_DESCRIPTION, description)
	return o
}

func (o *mediaImplementation) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *mediaImplementation) SetMemo(memo string) MediaInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *mediaImplementation) URL() string {
	return o.Get(COLUMN_MEDIA_URL)
}

func (o *mediaImplementation) SetURL(url string) MediaInterface {
	o.Set(COLUMN_MEDIA_URL, url)
	return o
}

func (o *mediaImplementation) Type() string {
	return o.Get(COLUMN_MEDIA_TYPE)
}

func (o *mediaImplementation) SetType(mediaType string) MediaInterface {
	o.Set(COLUMN_MEDIA_TYPE, mediaType)
	return o
}

func (o *mediaImplementation) Size() string {
	return o.Get(COLUMN_FILE_SIZE)
}

func (o *mediaImplementation) SetSize(size string) MediaInterface {
	o.Set(COLUMN_FILE_SIZE, size)
	return o
}

func (o *mediaImplementation) Extension() string {
	return o.Get(COLUMN_FILE_EXTENSION)
}

func (o *mediaImplementation) SetExtension(extension string) MediaInterface {
	o.Set(COLUMN_FILE_EXTENSION, extension)
	return o
}

func (o *mediaImplementation) Sequence() string {
	return o.Get(COLUMN_SEQUENCE)
}

func (o *mediaImplementation) SequenceInt() int {
	s := o.Sequence()
	if s == "" {
		return 0
	}
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}

func (o *mediaImplementation) SetSequence(sequence string) MediaInterface {
	o.Set(COLUMN_SEQUENCE, sequence)
	return o
}

func (o *mediaImplementation) SetSequenceInt(sequence int) MediaInterface {
	o.Set(COLUMN_SEQUENCE, jsonIntToString(sequence))
	return o
}

func (o *mediaImplementation) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *mediaImplementation) SetStatus(status string) MediaInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *mediaImplementation) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

func (o *mediaImplementation) SetHandle(handle string) MediaInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *mediaImplementation) SiteID() string {
	return o.Get(COLUMN_SITE_ID)
}

func (o *mediaImplementation) SetSiteID(siteID string) MediaInterface {
	o.Set(COLUMN_SITE_ID, siteID)
	return o
}

func (o *mediaImplementation) Metas() (map[string]string, error) {
	metasStr := o.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson := map[string]string{}
	errJson := json.Unmarshal([]byte(metasStr), &metasJson)
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return metasJson, nil
}

func (o *mediaImplementation) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *mediaImplementation) SetMeta(name string, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

func (o *mediaImplementation) SetMetas(metas map[string]string) error {
	mapString, err := json.Marshal(metas)
	if err != nil {
		return err
	}

	o.Set(COLUMN_METAS, string(mapString))

	return nil
}

func (o *mediaImplementation) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *mediaImplementation) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *mediaImplementation) SetCreatedAt(createdAt string) MediaInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *mediaImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *mediaImplementation) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *mediaImplementation) SetUpdatedAt(updatedAt string) MediaInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}

func (o *mediaImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}

func (o *mediaImplementation) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *mediaImplementation) SetSoftDeletedAt(softDeletedAt string) MediaInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, softDeletedAt)
	return o
}

func (o *mediaImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt())
}

func jsonIntToString(n int) string {
	if n == 0 {
		return "0"
	}
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
