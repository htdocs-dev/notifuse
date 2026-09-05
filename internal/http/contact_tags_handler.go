package http

// Fork patch 7: POST /api/contacts.tag and /api/contacts.untag.

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Notifuse/notifuse/internal/domain"
)

type tagChange func(ctx context.Context, workspaceID string, emails, tags []string) (int, error)

func (h *ContactHandler) handleTag(w http.ResponseWriter, r *http.Request) {
	h.handleTagChange(w, r, h.service.AddContactTags, "Failed to tag contacts")
}

func (h *ContactHandler) handleUntag(w http.ResponseWriter, r *http.Request) {
	h.handleTagChange(w, r, h.service.RemoveContactTags, "Failed to untag contacts")
}

func (h *ContactHandler) handleTagChange(w http.ResponseWriter, r *http.Request, apply tagChange, failMsg string) {
	if r.Method != http.MethodPost {
		WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req domain.TagContactsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	updated, err := apply(r.Context(), req.WorkspaceID, req.Emails, req.Tags)
	if err != nil {
		h.logger.WithField("error", err.Error()).Error(failMsg)
		if writeServiceError(w, err, "You do not have access to this workspace") {
			return
		}
		WriteJSONError(w, failMsg, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "updated": updated})
}
