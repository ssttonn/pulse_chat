package http

import (
	"encoding/json"
	"net/http"
	"pulse/src/api-core/internal/domain"
	"pulse/src/api-core/internal/usecase"
)

type GroupHandler struct {
	usecase usecase.GroupUseCase
}

type createGroupRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type createGroupResponse struct {
	Message string        `json:"message"`
	Group   *domain.Group `json:"group"`
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	creatorID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		http.Error(w, "Server error: Cannot find UserID", http.StatusInternalServerError)
		return
	}

	var req createGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input request", http.StatusBadRequest)
		return
	}

	group, err := h.usecase.CreateNewGroup(r.Context(), creatorID, req.Name, domain.GroupType(req.Type))

	if err != nil {
		http.Error(w, "Cannot create new group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := createGroupResponse{
		Group:   group,
		Message: "Create group successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
