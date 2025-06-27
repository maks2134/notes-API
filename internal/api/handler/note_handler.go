package handler

import (
	"encoding/json"
	"net/http"
	"notes-api/internal/model"
	"notes-api/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type NoteHandler struct {
	service service.NoteService
}

func NewNoteHandler(s service.NoteService) *NoteHandler {
	return &NoteHandler{service: s}
}

// GetNotes godoc
// @Summary      Get all notes
// @Description  Получить список всех заметок
// @Tags         notes
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   model.Note
// @Failure      500  {object}  map[string]string
// @Router       /notes [get]
func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Не удалось получить ID пользователя из токена", http.StatusInternalServerError)
		return
	}

	notes, err := h.service.GetAllNotes(userID)
	if err != nil {
		http.Error(w, "Не удалось получить заметки", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// CreateNote godoc
// @Summary      Create a new note
// @Description  Создать новую заметку
// @Tags         notes
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        note  body      model.Note  true  "Note Data"
// @Success      201   {object}  model.Note
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes [post]
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Не удалось получить ID пользователя из токена", http.StatusInternalServerError)
		return
	}

	var note model.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	note.UserID = userID
	if err := h.service.CreateNote(&note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// GetNote godoc
// @Summary      Get note by ID
// @Description  Получить заметку по её ID
// @Tags         notes
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path      int  true  "Note ID"
// @Success      200   {object}  model.Note
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /notes/{id} [get]
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Не удалось получить ID пользователя из токена", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	note, err := h.service.GetNoteByID(id, userID)
	if err != nil {
		http.Error(w, "Заметка не найдена или у вас нет к ней доступа", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// UpdateNote godoc
// @Summary      Update an existing note
// @Description  Обновить существующую заметку
// @Tags         notes
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path      int  true  "Note ID"
// @Param        note  body      model.Note  true  "Updated note data"
// @Success      200   {object}  model.Note
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes/{id} [put]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Не удалось получить ID пользователя из токена", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	var note model.Note
	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	note.ID = id

	if err := h.service.UpdateNote(&note, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// DeleteNote godoc
// @Summary      Delete a note
// @Description  Удалить заметку по ID
// @Tags         notes
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Note ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [delete]
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		http.Error(w, "Не удалось получить ID пользователя из токена", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteNote(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
