package handler

import (
	"encoding/json"
	"knowledge-base/internal/models"
	"knowledge-base/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AnswerHandler struct {
	answerService *service.AnswerService
}

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

	// Вызываем сервисный слой
	answers, err := answerHandler.answerService.GetAll()
	if err != nil {
		http.Error(w, "Ошибка получения ответов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON
	w.Header().Set("Content-Type", "application/json") // установка заголовка

	// Кодируем результат в JSON и отправляем
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

	vars := mux.Vars(r) // "/answers/5"
	idStr := vars["id"] // в шаблоне парамертр называется id

	id, err := strconv.Atoi(idStr) //преобразование строк в число
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	answer, err := answerHandler.answerService.GetByID(id)
	if err != nil {
		http.Error(w, "Ответ не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(answer)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete answer by ID
// @Description Delete answer by ID
// @Tags answers
// @Param id path int true "Answer ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Answer not found"
// @Router /answers/{id} [delete]
func (answerHandler *AnswerHandler) DeleteAnswerByID(w http.ResponseWriter, r *http.Request) {

	// if r.Method != http.MethodDelete {
	// 	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	// 	return
	// }

	vars := mux.Vars(r) // "/answers/5"
	idStr := vars["id"] // в шаблоне парамертр называется id

	id, err := strconv.Atoi(idStr) //преобразование строк в число
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = answerHandler.answerService.DeleteByID(id)
	if err != nil {
		http.Error(w, "Ответ не найден", http.StatusNotFound)
		return
	}

	//  Успешный ответ - 204 No connect
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Create new answer
// @Description Create a new answer with text, tutor_id, question_id and edit flag
// @Tags answers
// @Accept json
// @Produce json
// @Param answer body models.AnswersSwaggerRequestPostBody true "Answer data"
// @Success 201 {object} map[string]interface{} "Answer created"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /answers [post]
func (answerHandler *AnswerHandler) PostAnswerString(w http.ResponseWriter, r *http.Request) {
	// Проверка метода не нужна - Gorilla Mux уже гарантирует POST

	// Парсим JSON из тела запроса
	var answer models.Answer

	err := json.NewDecoder(r.Body).Decode(&answer)

	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация
	if answer.AnswersText == "" {
		http.Error(w, "answer_text is required", http.StatusBadRequest)
		return
	}

	// Вызов сервиса
	id, err := answerHandler.answerService.PostString(answer.AnswersText, answer.TutorID, answer.QuestionID)
	if err != nil {
		http.Error(w, "Failed to create answer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Answer created successfully",
	})
}

// @Summary Update answer
// @Description Update answer with text, tutor_id, question_id and edit flag
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
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var answer models.Answer

	err = json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if answer.AnswersText == "" {
		http.Error(w, "answer_text is required", http.StatusBadRequest)
		return
	}

	updatedAnswer, err := answerHandler.answerService.PutString(answer.AnswersText, answer.TutorID, answer.QuestionID, answer.IsEdit, id)
	if err != nil {
		http.Error(w, "Ответ не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"answer_text": updatedAnswer.AnswersText,
		"tutor_id":    updatedAnswer.TutorID,
		"question_id": updatedAnswer.QuestionID,
		"created_at":  updatedAnswer.CreatedAt,
		"is_edit":     updatedAnswer.IsEdit,
	})
}
