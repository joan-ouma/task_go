package repository

import (
	"context"
	"time"

	"task-manager-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{
		collection: db.Collection("tasks"),
	}
}

func (r *TaskRepository) CreateTask(task *models.Task) (*models.TaskResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	task.CreatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := &models.TaskResponse{
		ID:        result.InsertedID.(primitive.ObjectID).Hex(),
		Text:      task.Text,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
	}

	return response, nil
}

func (r *TaskRepository) GetAllTasks(filter models.FilterType) ([]models.TaskResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filterQuery bson.M
	switch filter {
	case models.FilterCompleted:
		filterQuery = bson.M{"completed": true}
	case models.FilterPending:
		filterQuery = bson.M{"completed": false}
	default:
		filterQuery = bson.M{}
	}

	cursor, err := r.collection.Find(ctx, filterQuery)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	// Convert to response format
	var response []models.TaskResponse
	for _, task := range tasks {
		response = append(response, models.TaskResponse{
			ID:        task.ID.Hex(),
			Text:      task.Text,
			Completed: task.Completed,
			CreatedAt: task.CreatedAt,
		})
	}

	return response, nil
}

func (r *TaskRepository) GetTaskByID(id string) (*models.TaskResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var task models.Task
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		return nil, err
	}

	response := &models.TaskResponse{
		ID:        task.ID.Hex(),
		Text:      task.Text,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
	}

	return response, nil
}

func (r *TaskRepository) UpdateTask(id string, updateData *models.UpdateTaskRequest) (*models.TaskResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	updateFields := bson.M{}
	if updateData.Text != nil {
		updateFields["text"] = *updateData.Text
	}
	if updateData.Completed != nil {
		updateFields["completed"] = *updateData.Completed
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return nil, err
	}

	// Return updated task
	return r.GetTaskByID(id)
}

func (r *TaskRepository) DeleteTask(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *TaskRepository) DeleteAllTasks() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteMany(ctx, bson.M{})
	return err
}

func (r *TaskRepository) GetTaskStats() (*models.TaskStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get total tasks
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Get completed tasks
	completed, err := r.collection.CountDocuments(ctx, bson.M{"completed": true})
	if err != nil {
		return nil, err
	}

	stats := &models.TaskStats{
		Total:     int(total),
		Completed: int(completed),
		Pending:   int(total) - int(completed),
	}

	return stats, nil
}
