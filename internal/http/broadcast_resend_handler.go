package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Notifuse/notifuse/internal/domain"
	"github.com/Notifuse/notifuse/internal/service"
)

// HandleResendToNonOpeners creates a draft broadcast for the recipients of a
// processed broadcast who did not open it.
func (h *BroadcastHandler) HandleResendToNonOpeners(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.ResendToNonOpenersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithField("error", err.Error()).Error("Failed to decode request body")
		WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	broadcast, err := h.service.ResendToNonOpeners(r.Context(), &req)
	if err != nil {
		if writePermissionError(w, err) {
			return
		}
		if _, ok := err.(*domain.ErrBroadcastNotFound); ok {
			WriteJSONError(w, "Broadcast not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrBroadcastNotProcessed) {
			WriteJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.logger.WithField("error", err.Error()).Error("Failed to create resend broadcast")
		WriteJSONError(w, "Failed to create resend broadcast", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"broadcast": broadcast,
	})
}
