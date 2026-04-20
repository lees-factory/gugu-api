package trackeditem

import (
	"context"
	"sort"
	"sync"
	"time"

	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type TrackedItemMemoryRepository struct {
	mu    sync.RWMutex
	byID  map[string]domaintrackeditem.TrackedItem
	index map[string]string
}

func NewRepository() *TrackedItemMemoryRepository {
	return &TrackedItemMemoryRepository{
		byID:  make(map[string]domaintrackeditem.TrackedItem),
		index: make(map[string]string),
	}
}

func (r *TrackedItemMemoryRepository) FindByIDAndUserID(_ context.Context, trackedItemID string, userID string) (*domaintrackeditem.TrackedItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.byID[trackedItemID]
	if !ok || item.UserID != userID {
		return nil, nil
	}
	found := item
	return &found, nil
}

func (r *TrackedItemMemoryRepository) FindByUserIDAndProductID(_ context.Context, userID string, productID string) (*domaintrackeditem.TrackedItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.index[userID+":"+productID]
	if !ok {
		return nil, nil
	}
	item := r.byID[id]
	found := item
	return &found, nil
}

func (r *TrackedItemMemoryRepository) ListByUserID(_ context.Context, userID string) ([]domaintrackeditem.TrackedItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domaintrackeditem.TrackedItem, 0)
	for _, item := range r.byID {
		if item.UserID == userID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (r *TrackedItemMemoryRepository) ListByUserIDWithCursor(_ context.Context, userID string, cursorCreatedAt time.Time, cursorID string, limit int) ([]domaintrackeditem.TrackedItem, error) {
	all, _ := r.ListByUserID(nil, userID)
	sorted := sortDescending(all)

	result := make([]domaintrackeditem.TrackedItem, 0, limit)
	for _, item := range sorted {
		if item.CreatedAt.Before(cursorCreatedAt) || (item.CreatedAt.Equal(cursorCreatedAt) && item.ID < cursorID) {
			result = append(result, item)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (r *TrackedItemMemoryRepository) ListByUserIDFirstPage(_ context.Context, userID string, limit int) ([]domaintrackeditem.TrackedItem, error) {
	all, _ := r.ListByUserID(nil, userID)
	sorted := sortDescending(all)

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted, nil
}

func sortDescending(items []domaintrackeditem.TrackedItem) []domaintrackeditem.TrackedItem {
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID > items[j].ID
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items
}

func (r *TrackedItemMemoryRepository) Create(_ context.Context, trackedItem domaintrackeditem.TrackedItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[trackedItem.ID] = trackedItem
	r.index[trackedItem.UserID+":"+trackedItem.ProductID] = trackedItem.ID
	return nil
}

func (r *TrackedItemMemoryRepository) DeleteByIDAndUserID(_ context.Context, trackedItemID string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.byID[trackedItemID]
	if !ok || item.UserID != userID {
		return nil
	}
	delete(r.byID, trackedItemID)
	delete(r.index, item.UserID+":"+item.ProductID)
	return nil
}

func (r *TrackedItemMemoryRepository) UpdateSKU(_ context.Context, trackedItemID string, userID string, skuID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.byID[trackedItemID]
	if !ok || item.UserID != userID {
		return nil
	}
	item.SKUID = skuID
	r.byID[trackedItemID] = item
	return nil
}

func (r *TrackedItemMemoryRepository) UpdatePreferredLanguage(_ context.Context, trackedItemID string, userID string, preferredLanguage string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.byID[trackedItemID]
	if !ok || item.UserID != userID {
		return nil
	}
	item.PreferredLanguage = preferredLanguage
	r.byID[trackedItemID] = item
	return nil
}

func (r *TrackedItemMemoryRepository) UpdateTrackingScope(_ context.Context, trackedItemID string, userID string, trackingScope string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.byID[trackedItemID]
	if !ok || item.UserID != userID {
		return nil
	}
	item.TrackingScope = trackingScope
	r.byID[trackedItemID] = item
	return nil
}
