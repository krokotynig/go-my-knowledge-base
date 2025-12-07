package handler

import (
	"encoding/json"
	"knowledge-base/internal/models"
	"knowledge-base/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Структура для работы со всеми ф-ями handler/tutor.go.
type TutorHandler struct {
	tutorService *service.TutorService
}

// Фунция для создания объекта типа TutorHandler.
func NewTutorhandler(tutorService *service.TutorService) *TutorHandler {
	return &TutorHandler{tutorService: tutorService}
}

// @Summary Get all tutors
// @Description Returns list of all tutors
// @Tags tutors
// @Produce json
// @Success 200 {array} models.Tutor
// @Router /tutors [get]
func (tutorHandler *TutorHandler) GetAllTutors(w http.ResponseWriter, r *http.Request) {

	// Вызов сервиса.
	tutors, err := tutorHandler.tutorService.GetAll()
	if err != nil {
		http.Error(w, "Ошибка получения тьюторов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON фомат и возвращаем.
	err = json.NewEncoder(w).Encode(tutors)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get tutor by ID
// @Description Returns tutor by specified ID
// @Tags tutors
// @Produce json
// @Param id path int true "Tutor ID"
// @Success 200 {object} models.Tutor
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Tutor not found"
// @Router /tutors/{id} [get]
func (tutorHandler *TutorHandler) GetTutorByID(w http.ResponseWriter, r *http.Request) {

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
	tutor, err := tutorHandler.tutorService.GetByID(id)
	if err != nil {
		http.Error(w, "Тьютор не найден", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON фомат и возвращаем.
	err = json.NewEncoder(w).Encode(tutor)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete tutor by ID
// @Description Delete tutor by ID
// @Tags tutors
// @Param id path int true "Tutor ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Tutor not found"
// @Router /tutors/{id} [delete]
func (tutorHandler *TutorHandler) DeleteTutorByID(w http.ResponseWriter, r *http.Request) {

	//Разбиение пути handler на части
	vars := mux.Vars(r)
	idStr := vars["id"]

	//Преобразование строк в число.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	err = tutorHandler.tutorService.DeleteByID(id)
	if err != nil {
		http.Error(w, "Тьютор не найден", http.StatusNotFound)
		return
	}

	//Возврат кода операции.
	//  Успешный ответ - 204 No connect для удаления.
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Create new tutor
// @Description Create a new tutor with full name and email
// @Tags tutors
// @Accept json
// @Produce json
// @Param tutor body models.TutorSwaggerRequestBody true "Tutor data"
// @Success 201 {object} map[string]interface{} "Tutor created"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /tutors [post]
func (tutorHandler *TutorHandler) PostTutorString(w http.ResponseWriter, r *http.Request) {

	var tutor models.Tutor

	//Преобразование JSON данных в формат структуры models.Tutor.
	err := json.NewDecoder(r.Body).Decode(&tutor)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация
	if tutor.FullName == "" || tutor.Email == "" {
		http.Error(w, "full_name and email are required", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	id, err := tutorHandler.tutorService.PostString(tutor.FullName, tutor.Email)
	if err != nil {
		http.Error(w, "Failed to create tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	//Возврат кода операции.
	w.WriteHeader(http.StatusCreated)

	// Кодируем результат в JSON фомат.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Tutor created successfully",
	})
}

// @Summary Update tutor
// @Description Update tutor with full name and email
// @Tags tutors
// @Accept json
// @Produce json
// @Param id path int true "Tutor ID"
// @Param tutor body models.TutorSwaggerRequestBody true "Tutor data"
// @Success 200 {object} map[string]string "Tutor updated"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Internal server error"
// @Router /tutors/{id} [put]
func (tutorHandler *TutorHandler) PutTutorString(w http.ResponseWriter, r *http.Request) {

	//Разбиение пути handler на части.
	vars := mux.Vars(r)
	idStr := vars["id"]

	//Преобразование строк в число.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var tutor models.TutorSwaggerRequestBody

	//Преобразование JSON данных в формат структуры models.Tutor.
	err = json.NewDecoder(r.Body).Decode(&tutor)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация.
	if tutor.FullName == "" || tutor.Email == "" {
		http.Error(w, "full_name and email are required", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	updatedTutor, err := tutorHandler.tutorService.PutString(tutor.FullName, tutor.Email, id)
	if err != nil {
		http.Error(w, "Тьютор не найден", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	//Возврат кода операции.
	w.WriteHeader(http.StatusOK)

	// Кодируем результат в JSON фомат.
	json.NewEncoder(w).Encode(map[string]string{
		"full_name": updatedTutor.FullName,
		"email":     updatedTutor.Email,
	})
}
