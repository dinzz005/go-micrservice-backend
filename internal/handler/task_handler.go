package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"microservices/internal/db"
	"microservices/internal/response"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TaskHandler handles HTTP requests related to task operations.
// It acts as the entry point between HTTP layer and database layer.
type TaskHandler struct {
	Q *db.Queries
}

// NewTaskHandler initializes and returns a new TaskHandler instance.
// It expects a sqlc-generated Queries object for database interactions.
func NewTaskHandler(q *db.Queries) *TaskHandler {
	return &TaskHandler{Q: q}
}

/*
Create handles the creation of a new task.

Endpoint:
POST /tasks

Expected JSON Body:
{
  "title": "string",
  "description": "string",
  "status_id": 1,
  "user_id": 1,
  "start_time": "RFC3339 timestamp",
  "end_time": "RFC3339 timestamp"
}

Validation Rules:
- title must not be empty
- status_id must be provided
- user_id must be provided
- end_time must be after start_time

Behavior:
- Parses request body
- Validates input
- Inserts task into database via sqlc
- Returns HTTP 201 on success

Future Improvements:
- Extract user_id from JWT instead of request body
- Move validation into service layer
- Add structured error responses
*/
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Request payload structure
	var req struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StatusID    int32     `json:"status_id"`
		UserID      int32     `json:"user_id"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
	}

	// Decode JSON request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
    response.JSON(w,http.StatusBadRequest,  response.APIResponse{
			Status: "Error",
			Message: "invalid request body",
		})
		return
	}

	// -----------------------------
	// Input Validation
	// -----------------------------

	if req.Title == "" {
     response.JSON(w,http.StatusBadRequest,  response.APIResponse{
			Status: "Error",
			Message: "Title is required",
		})
		return
	}

	if req.StatusID == 0 {
		  response.JSON(w,http.StatusBadRequest,  response.APIResponse{
			Status: "Error",
			Message: "Status_id is required",
		})	

		return
	}
	if req.UserID == 0 {
	  response.JSON(w,http.StatusBadRequest,  response.APIResponse{
			Status: "success",
			Message: "Access Denied",
		})	
		return
	}

	if req.EndTime.Before(req.StartTime) {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Status: "Error",
			Message: "end_time must be after start_time",
		})
		return
	}

	// -----------------------------
	// Database Operation
	// -----------------------------

	_, err = h.Q.CreateTask(r.Context(), db.CreateTaskParams{
		Title: req.Title,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
		StatusID:  sql.NullInt32{
			Int32: req.StatusID,
			Valid: true,
		},
		UserID:    sql.NullInt32{
			Int32: req.UserID,
			Valid: true,
		},
		StartTime: sql.NullTime{
			Time: req.StartTime,
			Valid: true,
		} ,
		EndTime: sql.NullTime{
			Time: req.EndTime,
			Valid: true,
		},
	})

	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.APIResponse{
			Status: "Error",
			Message: "something went wrong,failed to create task",
		})
		return
	}

	// -----------------------------
	// Response
	// -----------------------------
  response.JSON(w, http.StatusCreated, response.APIResponse{
		Status: "success",
		Message: "Task created Successfully",
	})
}


// Struct for return the tasks without any sqlc extra values.
type TaskResponse struct {
	ID          int32     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      int32     `json:"user_id"`
	StatusName  string     `json:"status_name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}



/*
GetAll retrieves all tasks.

Endpoint:
GET /tasks

Behavior:
- Fetches all tasks with joined status name
- Returns JSON array of tasks
*/

func(h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Q.ListTasks(r.Context())
	if err != nil {
		response.JSON(w, http.StatusInternalServerError,response.APIResponse{
		 Status: "Error",
		 Message: "Something went wrong, Failed to fetch Tasks",
	 })
		return
	}

  var res []TaskResponse
	for _, t := range tasks {
		res = append(res, TaskResponse{
			ID: t.ID,
			Title: t.Title,
			Description: t.Description.String,
			UserID: t.UserID.Int32,
			StatusName: t.StatusName,
			StartTime: t.StartTime.Time,
			EndTime: t.EndTime.Time,
		})
	}
 
	 response.JSON(w, http.StatusOK,  response.APIResponse{
		 Status: "success",
		 Message: "All the tasks fetched Successfully",
		 Data: res,
	 })
}


/*
GetOne retrieves a single task by ID.

Endpoint:
GET /tasks/{id}

Behavior:
- Extracts task ID from URL
- Fetches task with status
- Returns JSON object
*/

func(h *TaskHandler) GetOne(w http.ResponseWriter, r *http.Request){
	idStr := strings.TrimPrefix( r.URL.Path, "/tasks/")
	id , err := strconv.Atoi(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Status: "Error",
			Message: "Invalid Task ID",
		})
		return
	}

	task, err := h.Q.GetTask(r.Context(), int32(id))
	if err != nil {
		response.JSON(w, http.StatusNotFound, response.APIResponse{
			Status: "Error",
			Message: "task not found",
		})
		return
	}

	res := TaskResponse {
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description.String,
		StatusName:  task.StatusName,
		UserID:      task.UserID.Int32,
		StartTime:   task.StartTime.Time,
		EndTime:     task.EndTime.Time,
		CreatedAt:   task.CreatedAt.Time,
		UpdatedAt:   task.UpdatedAt.Time,
	}


	response.JSON(w, http.StatusOK, response.APIResponse{
		Status: "suucess",
		Message: fmt.Sprintf("Task data for %d fetched Successfully.", id),
		Data: res,
	})

}


/*
Update modifies an existing task.

Endpoint:
PUT /tasks/{id}

Expected JSON Body:
{
  "title": "string",
  "description": "string",
  "status_id": 1,
  "start_time": "...",
  "end_time": "..."
}

Behavior:
- Validates input
- Updates task fields
*/

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id , err := strconv.Atoi(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Status: "Error",
			Message: "Invalid task id",
		})
		return
	}

	var req struct {
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		StatusID       int32     `json:"status_id"`
		StartTime      time.Time `json:"start_time"`
		EndTime        time.Time `json:"end_time"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Status: "Error",
			Message: "invalid request body ",
		})		
		return
	}

		if req.Title == "" {
			response.JSON(w, http.StatusBadRequest, response.APIResponse{
				Status: "Error",
				Message: "title is required.",
			})
		return
	}

	if req.EndTime.Before(req.StartTime) {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Status: "Error",
			Message: "Invalid time range",
		})
		return
	}

	err = h.Q.UpdateTasks(r.Context(), db.UpdateTasksParams{
		ID: int32(id),
		Title: req.Title,
		Description: sql.NullString{
			String: req.Description,
			Valid: req.Description != "",
		},
		StatusID: sql.NullInt32{
			Int32: req.StatusID,
			Valid: true,
		},
		StartTime: sql.NullTime{
			Time: req.StartTime,
			Valid: true,
		},

		EndTime: sql.NullTime{
			Time: req.EndTime,
			Valid: true,
		},

	})

	if err != nil {
		response.JSON(w, http.StatusInternalServerError ,response.APIResponse{
			Status: "Error",
			Message: "failed to update task ",
		})
		return
	}

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Status: "success",
		Message: "Task updated SuccessFully",
	})
}



/*
Delete removes a task by ID.
Endpoint:
DELETE /tasks/{id}
Behavior:
- Deletes task from database
- Returns HTTP 200 on success
*/
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Status: "Error",
			Message: "invalid task id",
		})
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	err = h.Q.DeleteTask(r.Context(), int32(id))
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.APIResponse{
			Status: "Error",
			Message: "Failed to delete Task",
		})
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Status: "success",
		Message: "Task Deleted Successfully",
	})
}
