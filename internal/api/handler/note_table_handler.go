package handler

import (
	"encoding/json"
	"net/http"
	"notes-api/internal/model"
	"notes-api/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type NoteTableHandler struct {
	service service.NoteTableService
}

func NewNoteTableHandler(s service.NoteTableService) *NoteTableHandler {
	return &NoteTableHandler{service: s}
}

func (h *NoteTableHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/{note_id:[0-9]+}/tables", h.CreateTable).Methods("POST")
	r.HandleFunc("/{note_id:[0-9]+}/tables/{table_id:[0-9]+}/rows", h.AddRow).Methods("POST")
}

// CreateTable godoc
// @Summary      Create a table within a note
// @Description  Создает новую вложенную таблицу в существующей заметке
// @Tags         tables
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        note_id path int true "Note ID"
// @Param        table_data body model.CreateNoteTableRequest true "Table Title and Columns"
// @Success      201   {object}  model.NoteTable
// @Failure      400,401,403,404,500 {object} map[string]string
// @Router       /notes/{note_id}/tables [post]
func (h *NoteTableHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Не удалось получить ID пользователя из токена")
		return
	}

	vars := mux.Vars(r)
	noteID, _ := strconv.ParseInt(vars["note_id"], 10, 64)

	var req model.CreateNoteTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	table, err := h.service.CreateTable(&req, noteID, userID)
	if err != nil {
		respondError(w, http.StatusForbidden, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, table)
}

// AddRow godoc
// @Summary      Add a row to a table
// @Description  Добавляет новую строку данных в существующую таблицу
// @Tags         tables
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        note_id path int true "Note ID (для URL)"
// @Param        table_id path int true "Table ID"
// @Param        row_data body model.AddTableRowRequest true "Cell values in correct order"
// @Success      201   {object}  model.TableRow
// @Failure      400,401,403,404,500 {object} map[string]string
// @Router       /notes/{note_id}/tables/{table_id}/rows [post]
func (h *NoteTableHandler) AddRow(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Не удалось получить ID пользователя из токена")
		return
	}

	vars := mux.Vars(r)
	tableID, _ := strconv.ParseInt(vars["table_id"], 10, 64)

	var req model.AddTableRowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	row, err := h.service.AddRow(&req, tableID, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, row)
}
