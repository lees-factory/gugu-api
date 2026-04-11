package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type UpdateTrackedItemLanguage struct {
	User          auth.RequestUser
	TrackedItemID string
	Language      string
}

func ParseUpdateTrackedItemLanguage(r *http.Request) (UpdateTrackedItemLanguage, error) {
	var body struct {
		Language string `json:"language"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return UpdateTrackedItemLanguage{}, err
	}

	language := strings.ToUpper(strings.TrimSpace(body.Language))
	if language == "" {
		return UpdateTrackedItemLanguage{}, fmt.Errorf("language must not be empty")
	}

	return UpdateTrackedItemLanguage{
		User:          auth.RequestUserFrom(r.Context()),
		TrackedItemID: strings.TrimSpace(chi.URLParam(r, "trackedItemID")),
		Language:      language,
	}, nil
}
