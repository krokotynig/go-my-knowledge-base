package handler

import (
	"encoding/json"
	"knowledge-base/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Структура для работы со всеми ф-ями handler/question_version.go.
type QuestionVersionHandler struct {
	questionVersionService *service.QuestionVersionService
}

// Фунция для создания объекта типа QuestionVersionHandler.
func NewQuestionVersionHandler(questionVersionService *service.QuestionVersionService) *QuestionVersionHandler {
	return &QuestionVersionHandler{questionVersionService: questionVersionService}
}

// @Summary Get question versions by question ID
// @Description Returns all versions of a specific question by question ID
// @Tags question-versions
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {array} models.QuestionVersion
// @Failure 400 {string} string "Invalid question ID"
// @Router /question-versions/{id} [get]
func (handler *QuestionVersionHandler) GetAllQuestionVersionsByID(w http.ResponseWriter, r *http.Request) {

	//Разбиение пути handler на части.
	vars := mux.Vars(r)
	idStr := vars["id"]

	//Преобразование строк в число.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID вопроса", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	questionVersions, err := handler.questionVersionService.GetAllByID(id)
	if err != nil {
		http.Error(w, "Ошибка получения версий вопроса: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON фомат.
	json.NewEncoder(w).Encode(questionVersions)
}
