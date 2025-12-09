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
// @Router /question-tags/{question_id}/{tag_id} [post]
func (questionTagHandler *QuestionTagHandler) AddTagToQuestion(w http.ResponseWriter, r *http.Request) {

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
	err = questionTagHandler.questionTagService.AddToQuestion(questionID, tagID)
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

// @Summary Get all question-tag relations
// @Description Returns list of all relations between questions and tags
// @Tags question-tags
// @Produce json
// @Success 200 {array} models.QuestionTag
// @Router /question-tags [get]
func (questionTagHandler *QuestionTagHandler) GetAllQuestionTagRelations(w http.ResponseWriter, r *http.Request) {

	// Вызов сервиса для получения всех связей.
	relations, err := questionTagHandler.questionTagService.GetAllRelations()
	if err != nil {
		http.Error(w, "Ошибка получения связей вопрос-тег: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON формат и возвращаем.
	err = json.NewEncoder(w).Encode(relations)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get all question-tag relations by tag ID
// @Description Returns list of all relations between questions and tags by tag ID
// @Tags question-tags
// @Param tag_id path int true "Tag ID"
// @Produce json
// @Success 200 {array} models.QuestionTag
// @Router /question-tags/by-tag/{tag_id} [get]
func (questionTagHandler *QuestionTagHandler) GetAllQuestionTagRelationsByTagID(w http.ResponseWriter, r *http.Request) {

	// Разбиение пути handler на части.
	vars := mux.Vars(r)
	tagIDStr := vars["tag_id"]

	// Преобразование строк в число.
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		http.Error(w, "Неверный ID тега", http.StatusBadRequest)
		return
	}

	// Вызов сервиса для получения всех связей.
	relations, err := questionTagHandler.questionTagService.GetAllRelationsByTagID(tagID)
	if err != nil {
		http.Error(w, "Ошибка получения связей вопрос-тег: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON формат и возвращаем.
	err = json.NewEncoder(w).Encode(relations)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete question-tag relation
// @Description Delete relation between question and tag by their IDs
// @Tags question-tags
// @Param question_id path int true "Question ID"
// @Param tag_id path int true "Tag ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid question ID or tag ID"
// @Failure 404 {string} string "Relation not found"
// @Router /question-tags/{question_id}/{tag_id} [delete]
func (questionTagHandler *QuestionTagHandler) DeleteQuestionTagRelationByID(w http.ResponseWriter, r *http.Request) {

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

	//Возврат кода операции.
	//  Успешный ответ - 204 No connect для удаления.
	err = questionTagHandler.questionTagService.DeleteRelationByID(questionID, tagID)
	if err != nil {
		http.Error(w, "Связь не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
