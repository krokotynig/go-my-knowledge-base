package handler

import (
	"encoding/json"
	"knowledge-base/internal/service"
	"net/http"
	"strconv" // преобразование строк в число
	"strings" // работа со строками
	// gorilla/mux? Популярная для преоброзования url в число?
)

type TutorHandler struct {
	tutorService *service.Tutor
}

func NewTutorhandler(tutorService *service.Tutor) *TutorHandler {
	return &TutorHandler{tutorService: tutorService}
}

func (tutorHandler *TutorHandler) GetAllTutors(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

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

func (tutorHandler *TutorHandler) GetTutorByID(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path // "/tutors/5"
	parts := strings.Split(path, "/")

	if len(parts) != 3 || parts[1] != "tutors" { // прерывает выполнение, если неверный url
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]              // "1"
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

func (tutorHandler *TutorHandler) DeleteTutorByID(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path // "/tutors/5"
	parts := strings.Split(path, "/")

	if len(parts) != 4 || parts[1] != "tutors" || parts[2] != "delete" { // прерывает выполнение, если неверный url
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	idStr := parts[3]              // "1"
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
