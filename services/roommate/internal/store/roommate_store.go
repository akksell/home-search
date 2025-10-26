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
	pointOfInterestCollectionName = "pointsOfInterest"
	commutesCollectionName        = "commutes"
)

var globalFirestoreInstance *firestore.Client

type PointOfInterestStore struct {
	internalClient *firestore.Client
	projectId      string
	databaseId     string
}

func NewPointOfInterestStore(ctx context.Context, config *config.AppConfig) (*PointOfInterestStore, error) {
	store := &PointOfInterestStore{}

	store.projectId = config.GoogleProjectId
	store.databaseId = config.PointOfInterestStoreDB

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
func (poiStore *PointOfInterestStore) GetPointOfInterest(ctx context.Context, id string) (*models.PointOfInterest, error) {
	if poiStore == nil {
		return nil, fmt.Errorf("Cannot get point of interest: no store instantiated")
	}

	dbCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// TODO: validate that the user has sufficient permission to access the resource before
	// 		returning to client
	pointOfInterestSnap, err := poiStore.internalClient.Collection(pointOfInterestCollectionName).Doc(id).Get(dbCtx)
	if err != nil {
		return nil, err
	} else if !pointOfInterestSnap.Exists() {
		return nil, fmt.Errorf("Cannot get point of interest: point of interest %s does not exist", id)
	}

	var pointOfInterest *models.PointOfInterest
	err = pointOfInterestSnap.DataTo(pointOfInterest)
	if err != nil {
		return nil, err
	}
	return pointOfInterest, nil
}

func (poiStore *PointOfInterestStore) CreatePointOfInterest(ctx context.Context, pointOfInterest *models.PointOfInterest) (string, error) {
	if poiStore == nil {
		return "", fmt.Errorf("Cannot create point of interest: no store instantiated")
	}

	dbCtx, close := context.WithTimeout(ctx, 30*time.Second)
	defer close()

	pointOfInterestSnap, _, err := poiStore.internalClient.Collection(pointOfInterestCollectionName).Add(dbCtx, pointOfInterest)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return pointOfInterestSnap.ID, nil
}

func (poiStore *PointOfInterestStore) GetGroupCommute(ctx context.Context, commuteId string) (*models.GroupCommute, error) {
	if poiStore == nil {
		return nil, fmt.Errorf("Cannot get group commute: no store instantiated")
	}

	dbCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	groupCommuteSnap, err := poiStore.internalClient.Collection(commutesCollectionName).Doc(commuteId).Get(dbCtx)
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

func (poiStore *PointOfInterestStore) CreateGroupCommute(ctx context.Context, commute *models.GroupCommute) (string, error) {
	if poiStore == nil {
		return "", fmt.Errorf("Cannot create group commute: no store instantiated")
	}

	dbCtx, close := context.WithTimeout(ctx, 30*time.Second)
	defer close()

	commuteSnap, _, err := poiStore.internalClient.Collection(commutesCollectionName).Add(dbCtx, commute)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return commuteSnap.ID, nil
}

func (poiStore *PointOfInterestStore) Close() error {
	return poiStore.internalClient.Close()
}
