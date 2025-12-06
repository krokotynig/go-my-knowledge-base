package handler

import (
	"encoding/json"
	"knowledge-base/internal/models"
	"knowledge-base/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type QuestionHandler struct {
	questionService *service.QuestionService
}

func NewQuestionHandler(questionService *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService}
}

// @Summary Get all questions
// @Description Returns list of all questions
// @Tags questions
// @Produce json
// @Success 200 {array} models.Question
// @Router /questions [get]
func (questionHandler *QuestionHandler) GetAllQuestions(w http.ResponseWriter, r *http.Request) {

	// Вызываем сервисный слой
	questions, err := questionHandler.questionService.GetAll()
	if err != nil {
		http.Error(w, "Ошибка получения вопросов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON
	w.Header().Set("Content-Type", "application/json") // установка заголовка

	// Кодируем результат в JSON и отправляем
	err = json.NewEncoder(w).Encode(questions)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get question by ID
// @Description Returns question by specified ID
// @Tags questions
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} models.Question
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Question not found"
// @Router /questions/{id} [get]
func (questionHandler *QuestionHandler) GetQuestionByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) // "/questions/5"
	idStr := vars["id"] // в шаблоне парамертр называется id

	id, err := strconv.Atoi(idStr) //преобразование строк в число
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	question, err := questionHandler.questionService.GetByID(id)
	if err != nil {
		http.Error(w, "Вопрос не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(question)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete question by ID
// @Description Delete question by ID
// @Tags questions
// @Param id path int true "Question ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Question not found"
// @Router /questions/{id} [delete]
func (questionHandler *QuestionHandler) DeleteQuestionByID(w http.ResponseWriter, r *http.Request) {

	// if r.Method != http.MethodDelete {
	// 	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	// 	return
	// }

	vars := mux.Vars(r) // "/questions/5"
	idStr := vars["id"] // в шаблоне парамертр называется id

	id, err := strconv.Atoi(idStr) //преобразование строк в число
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = questionHandler.questionService.DeleteByID(id)
	if err != nil {
		http.Error(w, "Вопрос не найден", http.StatusNotFound)
		return
	}

	//  Успешный ответ - 204 No connect
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Create new question
// @Description Create a new question with text, tutor_id and edit flag
// @Tags questions
// @Accept json
// @Produce json
// @Param question body models.QuestionsSwaggerRequestBody true "Question data"
// @Success 201 {object} map[string]interface{} "Question created"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /questions [post]
func (questionHandler *QuestionHandler) PostQuestionString(w http.ResponseWriter, r *http.Request) {
	// Проверка метода не нужна - Gorilla Mux уже гарантирует POST

	// Парсим JSON из тела запроса
	var question models.Question

	err := json.NewDecoder(r.Body).Decode(&question)

	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация
	if question.QuestionText == "" {
		http.Error(w, "question_text is required", http.StatusBadRequest)
		return
	}

	// Вызов сервиса
	id, err := questionHandler.questionService.PostString(question.QuestionText, question.TutorID, question.IsEdit)
	if err != nil {
		http.Error(w, "Failed to create question: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Question created successfully",
	})
}

// @Summary Update question
// @Description Update question with text, tutor_id and edit flag
// @Tags questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Param question body models.QuestionsSwaggerRequestBody true "Question data"
// @Success 200 {object} map[string]interface{} "Question updated"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Question not found"
// @Router /questions/{id} [put]
func (questionHandler *QuestionHandler) PutQuestionString(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var question models.Question

	err = json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if question.QuestionText == "" {
		http.Error(w, "question_text is required", http.StatusBadRequest)
		return
	}

	updatedQuestion, err := questionHandler.questionService.PutString(question.QuestionText, question.TutorID, question.IsEdit, id)
	if err != nil {
		http.Error(w, "Вопрос не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"question_text": updatedQuestion.QuestionText,
		"tutor_id":      updatedQuestion.TutorID,
		"created_at":    updatedQuestion.CreatedAt,
		"is_edit":       updatedQuestion.IsEdit,
	})
}
