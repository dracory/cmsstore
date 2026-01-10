package cmsstore

// This file implements versioning operations for the CMS store.
// It provides methods to create, delete, find, list, soft delete,
// and update versioned entities, extending the store struct with
//  versioning capabilities.

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dracory/versionstore"
)

type versioningMarshalToInterface interface {
	MarshalToVersioning() (string, error)
}

type versioningDataInterface interface {
	Data() map[string]string
}

func (store *storeImplementation) versioningContentFromEntity(entity any) (string, error) {
	if entity == nil {
		return "", errors.New("entity is nil")
	}

	if v, ok := entity.(versioningMarshalToInterface); ok {
		return v.MarshalToVersioning()
	}

	d, ok := entity.(versioningDataInterface)
	if !ok {
		return "", errors.New("entity does not support versioning")
	}

	versionedData := map[string]string{}
	for k, v := range d.Data() {
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

func (store *storeImplementation) versioningCreateIfChanged(ctx context.Context, entityType string, entityID string, content string) error {
	if !store.VersioningEnabled() {
		return nil
	}

	if store.versioningStore == nil {
		return errors.New("cmsstore: versioning store is nil")
	}

	if entityType == "" {
		return errors.New("cmsstore: entityType is empty")
	}

	if entityID == "" {
		return errors.New("cmsstore: entityID is empty")
	}

	lastVersioningList, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(entityType).
		SetEntityID(entityID).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder("DESC").
		SetLimit(1))
	if err != nil {
		return err
	}

	if len(lastVersioningList) > 0 {
		lastVersioning := lastVersioningList[0]
		if lastVersioning != nil && lastVersioning.Content() == content {
			return nil
		}
	}

	return store.VersioningCreate(ctx, NewVersioning().
		SetEntityID(entityID).
		SetEntityType(entityType).
		SetContent(content))
}

func (store *storeImplementation) versioningTrackEntity(ctx context.Context, entityType string, entityID string, entity any) error {
	if !store.VersioningEnabled() {
		return nil
	}

	content, err := store.versioningContentFromEntity(entity)
	if err != nil {
		return err
	}

	return store.versioningCreateIfChanged(ctx, entityType, entityID, content)
}

// VersioningCreate creates a new versioning.
func (store *storeImplementation) VersioningCreate(ctx context.Context, version VersioningInterface) error {
	return store.versioningStore.VersionCreate(store.toQuerableContext(ctx), version)
}

// VersioningDelete deletes a versioning.
func (store *storeImplementation) VersioningDelete(ctx context.Context, version VersioningInterface) error {
	return store.versioningStore.VersionDelete(store.toQuerableContext(ctx), version)
}

// VersioningDeleteByID deletes a versioning by ID.
func (store *storeImplementation) VersioningDeleteByID(ctx context.Context, id string) error {
	return store.versioningStore.VersionDeleteByID(store.toQuerableContext(ctx), id)
}

// VersioningFindByID finds a versioning by ID.
func (store *storeImplementation) VersioningFindByID(ctx context.Context, versioningID string) (VersioningInterface, error) {
	return store.versioningStore.VersionFindByID(store.toQuerableContext(ctx), versioningID)
}

// VersioningList lists versionings.
func (store *storeImplementation) VersioningList(ctx context.Context, query VersioningQueryInterface) ([]VersioningInterface, error) {
	list, err := store.versioningStore.VersionList(store.toQuerableContext(ctx), query)

	if err != nil {
		return nil, err
	}

	newlist := make([]VersioningInterface, len(list))

	for i, v := range list {
		newlist[i] = v
	}

	return newlist, nil
}

// VersioningSoftDelete soft deletes a versioning.
func (store *storeImplementation) VersioningSoftDelete(ctx context.Context, versioning VersioningInterface) error {
	return store.versioningStore.VersionSoftDelete(store.toQuerableContext(ctx), versioning)
}

// VersioningSoftDeleteByID soft deletes a versioning by ID.
func (store *storeImplementation) VersioningSoftDeleteByID(ctx context.Context, id string) error {
	return store.versioningStore.VersionSoftDeleteByID(store.toQuerableContext(ctx), id)
}

// VersioningUpdate updates a versioning.
func (store *storeImplementation) VersioningUpdate(ctx context.Context, version VersioningInterface) error {
	return store.versioningStore.VersionUpdate(store.toQuerableContext(ctx), version)
}
