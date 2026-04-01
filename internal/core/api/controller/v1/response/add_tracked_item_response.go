package response

import (
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type AddTrackedItemEntry struct {
	TrackedItemID  string `json:"tracked_item_id"`
	AlreadyTracked bool   `json:"already_tracked"`
}

type AddTrackedItems struct {
	Results []AddTrackedItemEntry `json:"results"`
	Total   int                   `json:"total"`
	Added   int                   `json:"added"`
}

func NewAddTrackedItems(results []domaintrackeditem.AddTrackedItemResult) AddTrackedItems {
	entries := make([]AddTrackedItemEntry, 0, len(results))
	added := 0
	for _, r := range results {
		entries = append(entries, AddTrackedItemEntry{
			TrackedItemID:  r.TrackedItem.ID,
			AlreadyTracked: r.AlreadyTracked,
		})
		if !r.AlreadyTracked {
			added++
		}
	}
	return AddTrackedItems{
		Results: entries,
		Total:   len(results),
		Added:   added,
	}
}
