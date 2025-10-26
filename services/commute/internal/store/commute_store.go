package store

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"

	"homesearch.axel.to/commute/config"
	"homesearch.axel.to/commute/internal/models"
)

const (
	commutesCollectionName = "commute"
)

var globalFirestoreInstance *firestore.Client

type CommuteStore struct {
	internalClient *firestore.Client
	projectId      string
	databaseId     string
}

func NewCommuteStore(ctx context.Context, config *config.AppConfig) (*CommuteStore, error) {
	store := &CommuteStore{}

	store.projectId = config.GoogleProjectId
	store.databaseId = config.CommuteStoreDB

	if globalFirestoreInstance == nil {
		newFireStoreInstance, err := firestore.NewClientWithDatabase(ctx, store.projectId, store.databaseId)
		if err != nil {
			return nil, err
		}
		globalFirestoreInstance = newFireStoreInstance
	}
	store.internalClient = globalFirestoreInstance

	return store, nil
}

// TODO: create custom errors
func (cStore *CommuteStore) GetGroupCommute(ctx context.Context, commuteId string) (*models.GroupCommute, error) {
	if cStore == nil {
		return nil, fmt.Errorf("Cannot get group commute: no store instantiated")
	}

	dbCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	groupCommuteSnap, err := cStore.internalClient.Collection(commutesCollectionName).Doc(commuteId).Get(dbCtx)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	} else if !groupCommuteSnap.Exists() {
		return nil, fmt.Errorf("Cannot get group commute: commute does not exist")
	}

	var groupCommute *models.GroupCommute
	err = groupCommuteSnap.DataTo(groupCommute)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return groupCommute, nil
}

func (cStore *CommuteStore) CreateGroupCommute(ctx context.Context, commute *models.GroupCommute) (string, error) {
	if cStore == nil {
		return "", fmt.Errorf("Cannot create group commute: no store instantiated")
	}

	dbCtx, close := context.WithTimeout(ctx, 30*time.Second)
	defer close()

	commuteSnap, _, err := cStore.internalClient.Collection(commutesCollectionName).Add(dbCtx, commute)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return commuteSnap.ID, nil
}

func (cStore *CommuteStore) Close() error {
	return cStore.internalClient.Close()
}
