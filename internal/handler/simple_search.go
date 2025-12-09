package handler

import (
	"encoding/json"
	"knowledge-base/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

// Структура для работы со всеми ф-ями handler/simple_search.go.
type SimpleSearchHandler struct {
	simpleSearchService *service.SimpleSearchService
}

// Функция для создания объекта типа SimpleSearchHandler.
func NewSimpleSearchHandler(simpleSearchService *service.SimpleSearchService) *SimpleSearchHandler {
	return &SimpleSearchHandler{simpleSearchService: simpleSearchService}
}

// @Summary Search questions by tag name (exact match)
// @Description Search questions by exact tag name
// @Tags search 🔍
// @Produce json
// @Param name path string true "Tag name to search for (exact match)"
// @Success 200 {array} models.Question
// @Failure 400 {string} string "Tag name parameter is required"
// @Router /simple-search/{name} [get]
func (simpleSearchHandler *SimpleSearchHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {

	//Разбиение пути handler на части.
	vars := mux.Vars(r)
	name := vars["name"]

	// Проверка, что параметр name не пустой.
	if name == "" {
		http.Error(w, "Параметр 'name' обязателен для поиска", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	questions, err := simpleSearchHandler.simpleSearchService.SearchLogic(name)
	if err != nil {
		http.Error(w, "Ошибка поиска: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON формат и возвращаем.
	err = json.NewEncoder(w).Encode(questions)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
