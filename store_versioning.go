package cmsstore

// This file implements versioning operations for the CMS store.
// It provides methods to create, delete, find, list, soft delete,
// and update versioned entities, extending the store struct with
//  versioning capabilities.

import "context"

// VersioningCreate creates a new versioning.
func (store *store) VersioningCreate(ctx context.Context, version VersioningInterface) error {
	return store.versioningStore.VersionCreate(store.toQuerableContext(ctx), version)
}

// VersioningDelete deletes a versioning.
func (store *store) VersioningDelete(ctx context.Context, version VersioningInterface) error {
	return store.versioningStore.VersionDelete(store.toQuerableContext(ctx), version)
}

// VersioningDeleteByID deletes a versioning by ID.
func (store *store) VersioningDeleteByID(ctx context.Context, id string) error {
	return store.versioningStore.VersionDeleteByID(store.toQuerableContext(ctx), id)
}

// VersioningFindByID finds a versioning by ID.
func (store *store) VersioningFindByID(ctx context.Context, versioningID string) (VersioningInterface, error) {
	return store.versioningStore.VersionFindByID(store.toQuerableContext(ctx), versioningID)
}

// VersioningList lists versionings.
func (store *store) VersioningList(ctx context.Context, query VersioningQueryInterface) ([]VersioningInterface, error) {
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
func (store *store) VersioningSoftDelete(ctx context.Context, versioning VersioningInterface) error {
	return store.versioningStore.VersionSoftDelete(store.toQuerableContext(ctx), versioning)
}

// VersioningSoftDeleteByID soft deletes a versioning by ID.
func (store *store) VersioningSoftDeleteByID(ctx context.Context, id string) error {
	return store.versioningStore.VersionSoftDeleteByID(store.toQuerableContext(ctx), id)
}

// VersioningUpdate updates a versioning.
func (store *store) VersioningUpdate(ctx context.Context, version VersioningInterface) error {
	return store.versioningStore.VersionUpdate(store.toQuerableContext(ctx), version)
}
