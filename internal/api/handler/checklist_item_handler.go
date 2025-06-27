package handler

import (
	"encoding/json"
	"net/http"
	"notes-api/internal/model"
	"notes-api/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type ChecklistItemHandler struct {
	service service.ChecklistItemService
}

func NewChecklistItemHandler(s service.ChecklistItemService) *ChecklistItemHandler {
	return &ChecklistItemHandler{service: s}
}

func (h *ChecklistItemHandler) RegisterRoutes(r *mux.Router) {
	// Эндпоинты будут вложенными в заметки для логичности
	s := r.PathPrefix("/{note_id:[0-9]+}/checklist").Subrouter()
	s.HandleFunc("", h.Create).Methods("POST")
	s.HandleFunc("/{item_id}", h.Update).Methods("PUT")
	s.HandleFunc("/{item_id}", h.Delete).Methods("DELETE")
}

// Create godoc
// @Summary      Create a checklist item
// @Description  Создать новый элемент чек-листа для заметки
// @Tags         checklist
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        note_id path int true "Note ID"
// @Param        item body model.ChecklistItem true "Checklist Item Data (только 'text')"
// @Success      201   {object}  model.ChecklistItem
// @Failure      400,401,404,500 {object} map[string]string
// @Router       /notes/{note_id}/checklist [post]
func (h *ChecklistItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Не удалось получить ID пользователя")
		return
	}

	vars := mux.Vars(r)
	noteID, err := strconv.ParseInt(vars["note_id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Некорректный ID заметки")
		return
	}

	var item model.ChecklistItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}
	item.NoteID = noteID

	if err := h.service.Create(&item, userID); err != nil {
		respondError(w, http.StatusForbidden, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

// Update godoc
// @Summary      Update a checklist item
// @Description  Обновить элемент чек-листа (текст или статус выполнения)
// @Tags         checklist
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        note_id path int true "Note ID"
// @Param        item_id path int true "Checklist Item ID"
// @Param        item body model.ChecklistItem true "Checklist Item Data (только 'text' и 'completed')"
// @Success      200   {object}  model.ChecklistItem
// @Failure      400,401,404,500 {object} map[string]string
// @Router       /notes/{note_id}/checklist/{item_id} [put]
func (h *ChecklistItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Не удалось получить ID пользователя")
		return
	}

	vars := mux.Vars(r)
	itemID, err := strconv.ParseInt(vars["item_id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Некорректный ID элемента")
		return
	}

	var itemData model.ChecklistItem
	if err := json.NewDecoder(r.Body).Decode(&itemData); err != nil {
		respondError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	if err := h.service.Update(&itemData, itemID, userID); err != nil {
		respondError(w, http.StatusForbidden, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, itemData)
}

// Delete godoc
// @Summary      Delete a checklist item
// @Description  Удалить элемент чек-листа
// @Tags         checklist
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        note_id path int true "Note ID"
// @Param        item_id path int true "Checklist Item ID"
// @Success      204
// @Failure      400,401,404,500 {object} map[string]string
// @Router       /notes/{note_id}/checklist/{item_id} [delete]
func (h *ChecklistItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Не удалось получить ID пользователя")
		return
	}

	vars := mux.Vars(r)
	itemID, err := strconv.ParseInt(vars["item_id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Некорректный ID элемента")
		return
	}

	if err := h.service.Delete(itemID, userID); err != nil {
		respondError(w, http.StatusForbidden, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
