package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-todolist/internal/domain/task"
	"go-todolist/internal/infrastructure"
)

type TaskHandler struct {
	service task.TaskService
}

func NewTaskHandler(taskService task.TaskService) *TaskHandler {
	return &TaskHandler{service: taskService}
}

func (th *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var req TaskRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(req.Date) == 0 {
		req.Date = time.Now().Format(infrastructure.TimeFormat)
	}

	task, err := req.ToTask()
	if err != nil {
		// Тут наверно светить внутренние ошибки не особо хороший вариант
		log.Println(fmt.Errorf("handler internal error: %w", err))
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	id, err := th.service.AddTask(task)

	if err != nil {
		log.Println(fmt.Errorf("handler internal error: %w", err))
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := struct {
		Id int64 `json:"id"`
	}{
		Id: id,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(fmt.Errorf("handler internal error: %w", err))
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

}

func (th *TaskHandler) NextDate(w http.ResponseWriter, r *http.Request) {
	querry := r.URL.Query()

	now := querry.Get("now")
	date := querry.Get("date")
	repeat := querry.Get("repeat")

	nowTime, err := time.Parse(infrastructure.TimeFormat, now)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "incorrect date format")
		return
	}

	result, err := task.NextDate(nowTime, date, repeat)
	if err != nil {
		log.Println(fmt.Errorf("handler error: %w", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	send := result.Format(infrastructure.TimeFormat)
	w.Write([]byte(send))
}

func (th *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	result, err := th.service.GetAllTasks(50)

	if err != nil {
		log.Println(fmt.Errorf("handler internal error: %w", err))
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := struct {
		Tasks []TaskResponseDTO `json:"tasks"`
	}{
		Tasks: buildTasksResponse(result),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (th *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var req TaskRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "err.Error()")
		return
	}

	if err := req.Validate(); err != nil {
		writeJSONError(w, http.StatusBadRequest, "err.Error()")
		return
	}

	if len(req.Date) == 0 {
		req.Date = time.Now().Format(infrastructure.TimeFormat)
	}

	new, err := req.ToTask()
	if err != nil {
		// Тут наверно светить внутренние ошибки не особо хороший вариант
		log.Println(fmt.Errorf("handler internal error: %w", err))
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := th.service.UpdateTask(new); err != nil {
		if errors.Is(err, task.ErrIDNotFound) {
			writeJSONError(w, http.StatusBadRequest, "id not found")
			return
		} else {
			log.Println(fmt.Errorf("handler internal error: %w", err))
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct{}{})
}

func (th *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	querry := r.URL.Query()

	idQuery := querry.Get("id")

	if len(idQuery) == 0 {
		writeJSONError(w, http.StatusBadRequest, "id field is required")
		return
	}

	id, err := strconv.Atoi(idQuery)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "id is not a number")
		return
	}

	if err := th.service.CompleteTask(id); err != nil {
		if errors.Is(err, task.ErrItemNotFound) {
			writeJSONError(w, http.StatusNotFound, "item not found")
			return
		} else {
			log.Println(fmt.Errorf("handler internal error: %w", err))
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func (th *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	querry := r.URL.Query()

	idQuery := querry.Get("id")

	if len(idQuery) == 0 {
		writeJSONError(w, http.StatusBadRequest, "id field is required")
		return
	}

	id, err := strconv.Atoi(idQuery)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "id is not a number")
		return
	}

	if err := th.service.DeleteTask(id); err != nil {
		if errors.Is(err, task.ErrItemNotFound) {
			writeJSONError(w, http.StatusNotFound, "item not found")
			return
		} else {
			log.Println(fmt.Errorf("handler internal error: %w", err))
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func (th *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	querry := r.URL.Query()

	idQuery := querry.Get("id")

	if len(idQuery) == 0 {
		writeJSONError(w, http.StatusBadRequest, "id field is required")
		return
	}

	id, err := strconv.Atoi(idQuery)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "id is not a number")
		return
	}

	found, err := th.service.GetTask(id)
	if err != nil {
		if errors.Is(err, task.ErrItemNotFound) {
			writeJSONError(w, http.StatusNotFound, "item not found")
			return
		} else {
			log.Println(fmt.Errorf("handler internal error: %w", err))
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewTaskResponseDTO(found))
}

func buildTasksResponse(tasks []task.Task) []TaskResponseDTO {
	result := make([]TaskResponseDTO, len(tasks))
	for i, item := range tasks {
		result[i] = NewTaskResponseDTO(item)
	}
	return result
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	response := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
