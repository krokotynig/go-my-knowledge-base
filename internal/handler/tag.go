package handler

import (
	"encoding/json"
	"knowledge-base/internal/models"
	"knowledge-base/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Структура для работы со всеми ф-ями handler/tag.go.
type TagHandler struct {
	tagService *service.TagService
}

// Фунция для создания объекта типа TagHandler.
func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// @Summary Get all tags
// @Description Returns list of all tags
// @Tags tags
// @Produce json
// @Success 200 {array} models.Tag
// @Router /tags [get]
func (tagHandler *TagHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {

	// Вызов сервиса.
	tags, err := tagHandler.tagService.GetAll()
	if err != nil {
		http.Error(w, "Ошибка получения тегов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON фомат и возвращаем.
	err = json.NewEncoder(w).Encode(tags)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get tag by ID
// @Description Returns tag by specified ID
// @Tags tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} models.Tag
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Tag not found"
// @Router /tags/{id} [get]
func (tagHandler *TagHandler) GetTagByID(w http.ResponseWriter, r *http.Request) {

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
	tag, err := tagHandler.tagService.GetByID(id)
	if err != nil {
		http.Error(w, "Тег не найден", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON фомат и возвращаем.
	err = json.NewEncoder(w).Encode(tag)
	if err != nil {
		http.Error(w, "Ошибка кодирования JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Delete tag by ID
// @Description Delete tag by ID (cascades from questions_tags)
// @Tags tags
// @Param id path int true "Tag ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Tag not found"
// @Router /tags/{id} [delete]
func (tagHandler *TagHandler) DeleteTagByID(w http.ResponseWriter, r *http.Request) {

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
	err = tagHandler.tagService.DeleteByID(id)
	if err != nil {
		http.Error(w, "Тег не найден", http.StatusNotFound)
		return
	}

	//Возврат кода операции.
	//  Успешный ответ - 204 No connect для удаления.
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Create new tag
// @Description Create a new tag with tag name and optional tutor_id
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body models.TagSwaggerRequestBody true "Tag data"
// @Success 201 {object} map[string]interface{} "Tag created"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /tags [post]
func (tagHandler *TagHandler) PostTagString(w http.ResponseWriter, r *http.Request) {

	var tag models.TagSwaggerRequestBody

	//Преобразование JSON данных в формат структуры models.TagSwaggerRequestBody.
	err := json.NewDecoder(r.Body).Decode(&tag)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация.
	if tag.Tag == "" {
		http.Error(w, "tag is required", http.StatusBadRequest)
		return
	}

	// Вызов сервиса.
	id, err := tagHandler.tagService.PostString(tag.Tag, tag.TutorID)
	if err != nil {
		http.Error(w, "Failed to create tag: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок JSON.
	w.Header().Set("Content-Type", "application/json")

	//Возврат кода операции.
	w.WriteHeader(http.StatusCreated)

	// Кодируем результат в JSON фомат.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Tag created successfully",
	})
}
