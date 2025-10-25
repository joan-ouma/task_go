package handlers

import (
	"net/http"
	"task-manager-api/models"
	"task-manager-api/repository"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	task := &models.Task{
		Text:      req.Text,
		Completed: false,
	}

	createdTask, err := h.repo.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

// GetTasks returns all tasks with optional filtering
func (h *TaskHandler) GetTasks(c *gin.Context) {
	filter := c.DefaultQuery("filter", "all")

	var filterType models.FilterType
	switch filter {
	case "completed":
		filterType = models.FilterCompleted
	case "pending":
		filterType = models.FilterPending
	default:
		filterType = models.FilterAll
	}

	tasks, err := h.repo.GetAllTasks(filterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask returns a single task by ID
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.repo.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask updates a task
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedTask, err := h.repo.UpdateTask(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask deletes a task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// DeleteAllTasks deletes all tasks
func (h *TaskHandler) DeleteAllTasks(c *gin.Context) {
	err := h.repo.DeleteAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete all tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All tasks deleted successfully"})
}

// GetStats returns task statistics
func (h *TaskHandler) GetStats(c *gin.Context) {
	stats, err := h.repo.GetTaskStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ToggleTask toggles the completed status of a task
func (h *TaskHandler) ToggleTask(c *gin.Context) {
	id := c.Param("id")

	// First get the current task
	currentTask, err := h.repo.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Toggle the completed status
	completed := !currentTask.Completed
	updateReq := models.UpdateTaskRequest{
		Completed: &completed,
	}

	updatedTask, err := h.repo.UpdateTask(id, &updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle task"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}
