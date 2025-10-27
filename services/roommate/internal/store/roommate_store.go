package store

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"

	"homesearch.axel.to/roommate/config"
	"homesearch.axel.to/roommate/internal/store/models"
)

const (
	roommateCollectionName        = "roommate"
	pointOfInterestCollectionName = "pointOfInterest"
	roommateGroupCollectionName   = "roommateGroup"
)

var globalFirestoreInstance *firestore.Client

type RoommateStore struct {
	internalClient *firestore.Client
	projectId      string
	databaseId     string
}

func NewRoommateStore(ctx context.Context, config *config.AppConfig) (*RoommateStore, error) {
	store := &RoommateStore{}

	store.projectId = config.GoogleProjectId
	store.databaseId = config.RoommateStoreInstanceDB

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

func (store *RoommateStore) CreateRoommate(ctx context.Context, roommate *models.Roommate) (string, error) {
	if store == nil {
		return "", fmt.Errorf("Cannot create roommate: no store instantiated")
	}

	newRoommateCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	roommateRef, _, err := store.internalClient.Collection(roommateCollectionName).Add(newRoommateCtx, roommate)
	if err != nil {
		return "", fmt.Errorf("Failed to create roommate: %w", err)
	}
	return roommateRef.ID, nil
}

// TODO: create custom errors
func (store *RoommateStore) GetPointOfInterest(ctx context.Context, id string) (*models.PointOfInterest, error) {
	if store == nil {
		return nil, fmt.Errorf("Cannot get point of interest: no store instantiated")
	}

	dbCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// TODO: validate that the user has sufficient permission to access the resource before
	// 		returning to client
	pointOfInterestSnap, err := store.internalClient.Collection(pointOfInterestCollectionName).Doc(id).Get(dbCtx)
	if err != nil {
		return nil, err
	} else if !pointOfInterestSnap.Exists() {
		return nil, fmt.Errorf("Cannot get point of interest: point of interest %s does not exist", id)
	}

	var pointOfInterest models.PointOfInterest
	err = pointOfInterestSnap.DataTo(&pointOfInterest)
	if err != nil {
		return nil, fmt.Errorf("Failed to get point of interest: %w", err)
	}
	pointOfInterest.ID = pointOfInterestSnap.Ref.ID
	return &pointOfInterest, nil
}

// TODO: create custom errors
func (store *RoommateStore) CreatePointOfInterest(ctx context.Context, pointOfInterest *models.PointOfInterest) (string, error) {
	if store == nil {
		return "", fmt.Errorf("Cannot create point of interest: no store instantiated")
	}

	dbCtx, close := context.WithTimeout(ctx, 30*time.Second)
	defer close()

	pointOfInterestSnap, _, err := store.internalClient.Collection(pointOfInterestCollectionName).Add(dbCtx, pointOfInterest)
	if err != nil {
		return "", fmt.Errorf("Failed to create point of interest: %w", err)
	}
	return pointOfInterestSnap.ID, nil
}

// TODO: create custom errors
func (store *RoommateStore) CreateGroup(ctx context.Context, group *models.RoommateGroup) (string, error) {
	if store == nil {
		return "", fmt.Errorf("Cannot create group: no store instantiated")
	}

	groupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	groupRef, _, err := store.internalClient.Collection(roommateGroupCollectionName).Add(groupCtx, group)
	if err != nil {
		return "", fmt.Errorf("Failed to create group: %w", err)
	}
	return groupRef.ID, nil
}

// TODO: create custom errors
func (store *RoommateStore) GetGroup(ctx context.Context, groupId string) (*models.RoommateGroup, error) {
	if store == nil {
		return nil, fmt.Errorf("Cannot get group: no store instantiated")
	}

	groupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	groupSnap, err := store.internalClient.Collection(roommateGroupCollectionName).Doc(groupId).Get(groupCtx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get group: %w", err)
	}

	var group models.RoommateGroup
	err = groupSnap.DataTo(&group)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert group to model: %w", err)
	}
	group.ID = groupSnap.Ref.ID

	return &group, nil
}

// TODO: handle errors, see if there's a better way to update instead of override like this
// see map[string]interface{} + firestore.MergeAll
func (store *RoommateStore) UpdateGroup(ctx context.Context, group *models.RoommateGroup) (string, error) {
	if store == nil {
		return "", fmt.Errorf("Cannot update group: no store instantiated")
	}
	if group.ID == "" {
		return "", fmt.Errorf("Failed to update group: missing group id")
	}

	existingGroupRef := store.internalClient.Collection(roommateGroupCollectionName).Doc(group.ID)
	if existingGroupRef == nil {
		return group.ID, fmt.Errorf("Failed to update group: could not find group with id %v", group.ID)
	}

	groupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	group.UpdatedAt = time.Now()
	_, err := existingGroupRef.Set(groupCtx, group)
	if err != nil {
		return group.ID, fmt.Errorf("Failed to update group: %w", err)
	}
	return group.ID, nil
}

func (store *RoommateStore) Close() error {
	return store.internalClient.Close()
}
