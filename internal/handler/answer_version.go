package handler

import (
	"encoding/json"
	"knowledge-base/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Структура для работы со всеми ф-ями handler/answer_version.go
type AnswerVersionHandler struct {
	answerVersionService *service.AnswerVersionService
}

// Функция для создания объекта типа AnswerVersionHandler
func NewAnswerVersionHandler(answerVersionService *service.AnswerVersionService) *AnswerVersionHandler {
	return &AnswerVersionHandler{answerVersionService: answerVersionService}
}

// @Summary Get answer versions by answer ID
// @Description Returns all versions of a specific answer by answer ID
// @Tags answer-versions
// @Produce json
// @Param id path int true "Answer ID"
// @Success 200 {array} models.AnswerVersion
// @Failure 400 {string} string "Invalid answer ID"
// @Router /answer-versions/{id} [get]
func (handler *AnswerVersionHandler) GetAllAnswerVersionsByID(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из параметров пути
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Преобразование строки в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID ответа", http.StatusBadRequest)
		return
	}

	// Вызов сервиса для получения версий ответа
	answerVersions, err := handler.answerVersionService.GetAllByID(id)
	if err != nil {
		http.Error(w, "Ошибка получения версий ответа: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON формат
	json.NewEncoder(w).Encode(answerVersions)
}
