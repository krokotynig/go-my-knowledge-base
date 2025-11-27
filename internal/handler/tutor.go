package handler

import (
	"encoding/json"
	"knowledge-base/internal/service"
	"net/http"
	"strconv" // преобразование строк в число

	"github.com/gorilla/mux"
	// gorilla/mux? Популярная для преоброзования url в число?
)

type TutorHandler struct {
	tutorService *service.Tutor
}

func NewTutorhandler(tutorService *service.Tutor) *TutorHandler {
	return &TutorHandler{tutorService: tutorService}
}

// @Summary Получить всех тьюторов
// @Description Возвращает список всех тьюторов
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

// @Summary Получить тьютора по ID
// @Description Возвращает тьютора по указанному ID
// @Tags tutors
// @Produce json
// @Param id path int true "ID тьютора"
// @Success 200 {object} models.Tutor
// @Failure 400 {string} string "Неверный ID"
// @Failure 404 {string} string "Тьютор не найден"
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

	//  Успешный ответ - 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
