package handler

import (
	"encoding/json"
	"fmt"
	"knowledge-base/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Структура для работы со всеми ф-ями handler/qustion_tag.go.
type QuestionTagHandler struct {
	questionTagService *service.QuestionTagService
}

// Фунция для создания объекта типа QuestionHandler.
func NewQuestionTagHandler(questionTagService *service.QuestionTagService) *QuestionTagHandler {
	return &QuestionTagHandler{questionTagService: questionTagService}
}

// @Summary Add tag to question
// @Description Associate a tag with a question using IDs from URL path
// @Tags question-tags
// @Produce json
// @Param question_id path int true "Question ID"
// @Param tag_id path int true "Tag ID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Bad request"
// @Router /questions/{question_id}/tags/{tag_id} [post]
func (h *QuestionTagHandler) AddTagToQuestion(w http.ResponseWriter, r *http.Request) {

	// Разбиение пути handler на части.
	vars := mux.Vars(r)
	questionIDStr := vars["question_id"]
	tagIDStr := vars["tag_id"]

	// Преобразование строк в число.
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		http.Error(w, "Неверный ID вопроса", http.StatusBadRequest)
		return
	}

	// Преобразование строк в число.
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		http.Error(w, "Неверный ID тега", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	err = h.questionTagService.AddToQuestion(questionID, tagID)
	if err != nil {
		http.Error(w, "Ошибка добавления тега: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Возврат кода операции.
	w.WriteHeader(http.StatusCreated)

	// Кодируем результат в JSON формат.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"question_id": questionID,
		"tag_id":      tagID,
		"message":     fmt.Sprintf("Тег %d добавлен к вопросу %d", tagID, questionID),
	})
}
