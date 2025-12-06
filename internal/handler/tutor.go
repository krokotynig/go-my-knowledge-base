package handler

import (
	"encoding/json"
	"knowledge-base/internal/models"
	"knowledge-base/internal/service"
	"net/http"
	"strconv" // преобразование строк в число

	"github.com/gorilla/mux"
	// gorilla/mux? Популярная для преоброзования url в число?
)

type TutorHandler struct {
	tutorService *service.TutorService
}

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
	// Проверяем метод запроса, с gorilla/mux не нухно
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	// 	return
	// }

	// Вызываем сервисный слой
	tutors, err := tutorHandler.tutorService.GetAll()
	if err != nil {
		http.Error(w, "Ошибка получения тьюторов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON
	w.Header().Set("Content-Type", "application/json") // установка заголовка

	// Кодируем результат в JSON и отправляем
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

	vars := mux.Vars(r) // "/tutors/5"
	idStr := vars["id"] // в шаблоне парамертр называется id

	id, err := strconv.Atoi(idStr) //преобразование строк в число
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	tutor, err := tutorHandler.tutorService.GetByID(id)
	if err != nil {
		http.Error(w, "Тьютор не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

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

	// if r.Method != http.MethodDelete {
	// 	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	// 	return
	// }

	vars := mux.Vars(r) // "/tutors/5"
	idStr := vars["id"] // в шаблоне парамертр называется id

	id, err := strconv.Atoi(idStr) //преобразование строк в число
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = tutorHandler.tutorService.DeleteByID(id)
	if err != nil {
		http.Error(w, "Тьютор не найден", http.StatusNotFound)
		return
	}

	//  Успешный ответ - 204 No connect
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
	// Проверка метода не нужна - Gorilla Mux уже гарантирует POST

	// Парсим JSON из тела запроса
	var tutor models.Tutor

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

	// Вызов сервиса
	id, err := tutorHandler.tutorService.PostString(tutor.FullName, tutor.Email)
	if err != nil {
		http.Error(w, "Failed to create tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{ // используем словарь для того, чтобы отдать json (interface{} - что то типа object)
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

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var tutor models.Tutor

	err = json.NewDecoder(r.Body).Decode(&tutor)

	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if tutor.FullName == "" || tutor.Email == "" {
		http.Error(w, "full_name and email are required", http.StatusBadRequest)
		return
	}

	updatedTutor, err := tutorHandler.tutorService.PutString(tutor.FullName, tutor.Email, id)
	if err != nil {
		http.Error(w, "Тьютор не найден", http.StatusNotFound)
		return
	}

	// Ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"full_name": updatedTutor.FullName,
		"email":     updatedTutor.Email,
	})
}
