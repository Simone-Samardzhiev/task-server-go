package task

import (
	"database/sql"
	"github.com/google/uuid"
)

// Repository type interface used to manage task data.
type Repository interface {
	// GetTasksByUserId will get all tasks that belongs to a user.
	GetTasksByUserId(userId *uuid.UUID) ([]Task, error)

	// AddTask will add a new task.
	AddTask(task *Task, userId *uuid.UUID) error

	// UpdateTask will update an existing task
	UpdateTask(task *Task) (bool, error)

	// DeleteTask will delete an existing task.
	DeleteTask(id *uuid.UUID) (bool, error)
}

type PostgresRepository struct {
	database *sql.DB
}

func NewPostgresRepository(database *sql.DB) *PostgresRepository {
	return &PostgresRepository{database}
}

// countTasks will count all tasks that user have.
func (p *PostgresRepository) countTasks(userId *uuid.UUID) (int, error) {
	var count int
	err := p.database.QueryRow("SELECT COUNT(id) FROM tasks WHERE user_id = $1", userId).Scan(&count)
	return count, err
}

func (p *PostgresRepository) GetTasksByUserId(userId *uuid.UUID) ([]Task, error) {
	count, err := p.countTasks(userId)
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0, count)

	rows, err := p.database.Query("SELECT id, name, description, type, due_date, date_completed, date_deleted FROM tasks WHERE user_id = $1", *userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Name, &task.Description, &task.Type, &task.DueDate, &task.DateCompleted, &task.DateDeleted)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (p *PostgresRepository) AddTask(task *Task, userId *uuid.UUID) error {
	_, err := p.database.Exec("INSERT INTO tasks (id, name, description, type, due_date, date_completed, date_deleted, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		task.Id, task.Name, task.Description, task.Type, task.DueDate, task.DateCompleted, task.DateDeleted, *userId,
	)
	return err
}

func (p *PostgresRepository) UpdateTask(task *Task) (bool, error) {
	result, err := p.database.Exec("UPDATE tasks SET name = $1, description = $2, type = $3, due_date = $4, date_completed = $5, date_deleted = $6 WHERE id = $7",
		task.Name, task.Description, task.Type, task.DueDate, task.DateCompleted, task.DateDeleted, task.Id,
	)

	count, err := result.RowsAffected()
	return count > 0, err
}

func (p *PostgresRepository) DeleteTask(id *uuid.UUID) (bool, error) {
	result, err := p.database.Exec("DELETE FROM tasks WHERE id = $1", *id)
	count, err := result.RowsAffected()
	return count > 0, err
}
