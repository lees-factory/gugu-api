package response

import (
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type AddTrackedItem struct {
	TrackedItemID  string `json:"tracked_item_id"`
	AlreadyTracked bool   `json:"already_tracked"`
}

func NewAddTrackedItemFromResult(result *domaintrackeditem.AddTrackedItemResult) AddTrackedItem {
	return AddTrackedItem{
		TrackedItemID:  result.TrackedItem.ID,
		AlreadyTracked: result.AlreadyTracked,
	}
}
