package handler

import (
	"encoding/json"
	"knowledge-base/internal/models"
	"knowledge-base/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Структура для работы со всеми ф-ями handler/answers.go.
type AnswerHandler struct {
	answerService *service.AnswerService
}

// Фунция для создания объекта типа AnswerHandler.
func NewAnswerHandler(answerService *service.AnswerService) *AnswerHandler {
	return &AnswerHandler{answerService: answerService}
}

// @Summary Get all answers
// @Description Returns list of all answers
// @Tags answers
// @Produce json
// @Success 200 {array} models.Answer
// @Router /answers [get]
func (answerHandler *AnswerHandler) GetAllAnswers(w http.ResponseWriter, r *http.Request) {

	// Вызов сервиса.
	answers, err := answerHandler.answerService.GetAll()
	if err != nil {
		http.Error(w, "Ошибка получения ответов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON фомат и возвращаем.
	err = json.NewEncoder(w).Encode(answers)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get answer by ID
// @Description Returns answer by specified ID
// @Tags answers
// @Produce json
// @Param id path int true "Answer ID"
// @Success 200 {object} models.Answer
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Answer not found"
// @Router /answers/{id} [get]
func (answerHandler *AnswerHandler) GetAnswerByID(w http.ResponseWriter, r *http.Request) {

	//Разбиение пути handler на части.
	vars := mux.Vars(r)
	idStr := vars["id"]

	//Преобразование строк в число.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	answer, err := answerHandler.answerService.GetByID(id)
	if err != nil {
		http.Error(w, "Ответ не найден", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON фомат и возвращаем.
	err = json.NewEncoder(w).Encode(answer)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete answer by ID with version tracking
// @Description Delete answer by ID and mark all answer versions as deleted with tutor who performed deletion
// @Tags answers
// @Param id path int true "Answer ID"
// @Param delete-by path int true "Tutor ID who deleted the answer"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Answer not found"
// @Router /answers/{id}/delete-by-tutor/{delete-by} [delete]
func (answerHandler *AnswerHandler) DeleteAnswerByID(w http.ResponseWriter, r *http.Request) {

	//Разбиение пути handler на части.
	vars := mux.Vars(r)
	idStr := vars["id"]
	deleteByTutorStr := vars["delete-by"]

	//Преобразование строк в число.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	//Преобразование строк в число.
	deleteByTutor, err := strconv.Atoi(deleteByTutorStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	//Возврат кода операции.
	//  Успешный ответ - 204 No connect для удаления.
	err = answerHandler.answerService.DeleteByID(id, deleteByTutor)
	if err != nil {
		http.Error(w, "Вопрос не найден", http.StatusNotFound)
		return
	}

	//Возврат кода операции.
	//  Успешный ответ - 204 No connect для удаления.
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Create new answer and records the version
// @Description Create a new answer with text, tutor_id, question_id. Create a new answer_version with answer_id, answer_text, tutor_id, answer_number
// @Tags answers
// @Accept json
// @Produce json
// @Param answer body models.AnswersSwaggerRequestPostBody true "Answer data"
// @Success 201 {object} map[string]interface{} "Answer created"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /answers [post]
func (answerHandler *AnswerHandler) PostAnswerString(w http.ResponseWriter, r *http.Request) {

	var answer models.Answer

	//Преобразование JSON данных в формат структуры models.Answer.
	err := json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация.
	if answer.AnswersText == "" {
		http.Error(w, "answer_text is required", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	id, err := answerHandler.answerService.PostString(answer.AnswersText, answer.TutorID, answer.QuestionID)
	if err != nil {
		http.Error(w, "Failed to create answer or answer_version: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	//Возврат кода операции.
	w.WriteHeader(http.StatusCreated)

	// Кодируем результат в JSON фомат.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Answer created successfully",
	})
}

// @Summary Update answer and records the version
// @Description Update answer with text, tutor_id, question_id and edit flag. Create a new answer_version with answer_id, answer_text, tutor_id, answer_number
// @Tags answers
// @Accept json
// @Produce json
// @Param id path int true "Answer ID"
// @Param answer body models.AnswersSwaggerRequestPutBody true "Answer data"
// @Success 200 {object} map[string]interface{} "Answer updated"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Answer not found"
// @Router /answers/{id} [put]
func (answerHandler *AnswerHandler) PutAnswerString(w http.ResponseWriter, r *http.Request) {

	//Разбиение пути handler на части.
	vars := mux.Vars(r)
	idStr := vars["id"]

	//Преобразование строк в число.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var answer models.Answer

	//Преобразование JSON данных в формат структуры models.Answer.
	err = json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация.
	if answer.AnswersText == "" {
		http.Error(w, "answer_text is required", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	updatedAnswer, err := answerHandler.answerService.PutString(answer.AnswersText, answer.TutorID, answer.QuestionID, answer.IsEdit, id)
	if err != nil {
		http.Error(w, "Ответ не найден", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	//Возврат кода операции.
	w.WriteHeader(http.StatusOK)

	// Кодируем результат в JSON фомат.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"answer_text": updatedAnswer.AnswersText,
		"tutor_id":    updatedAnswer.TutorID,
		"question_id": updatedAnswer.QuestionID,
		"created_at":  updatedAnswer.CreatedAt,
		"is_edit":     updatedAnswer.IsEdit,
	})
}
