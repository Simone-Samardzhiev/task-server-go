package repositories

import (
	"context"
	"database/sql"
	"server/models"
)

// TaskRepository manages tasks data.
type TaskRepository interface {
	// GetTasks will return all tasks of the user.
	GetTasks(ctx context.Context, userId int) ([]models.TaskPayload, error)

	// CheckPriority will check if the task priority is in the database.
	CheckPriority(ctx context.Context, priority string) (bool, error)

	// AddTask will add new task.
	AddTask(ctx context.Context, taskPayload *models.TaskPayload, userId int) error
}

// PostgresTaskRepository is default implementation of [TaskRepository] using postgres database.
type PostgresTaskRepository struct {
	db *sql.DB
}

// countTasks will count all tasks that belongs to a user.
func (r *PostgresTaskRepository) countTasks(ctx context.Context, userId int) (int, error) {
	var count int
	row := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM tasks
                WHERE user_id = $1`,
		userId,
	)

	err := row.Scan(&count)
	return count, err
}

func (r *PostgresTaskRepository) GetTasks(ctx context.Context, userId int) ([]models.TaskPayload, error) {
	count, err := r.countTasks(ctx, userId)
	if err != nil {
		return nil, err
	}

	result := make([]models.TaskPayload, 0, count)

	rows, err := r.db.Query(
		`SELECT id, name, description, priority, date FROM tasks
                WHERE user_id = $1`,
		userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task models.TaskPayload
		err = rows.Scan(&task.Id, &task.Name, &task.Description, &task.Priority, &task.Date)
		if err != nil {
			return nil, err
		}
		result = append(result, task)
	}

	return result, nil
}

func (r *PostgresTaskRepository) CheckPriority(ctx context.Context, priority string) (bool, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM priorities
			WHERE priority = $1`,
		priority,
	)

	var count int
	err := row.Scan(&count)
	return count > 0, err
}

func (r *PostgresTaskRepository) AddTask(ctx context.Context, task *models.TaskPayload, userId int) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tasks (id, name, description, priority, date, user_id)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
		task.Id,
		task.Name,
		task.Description,
		task.Priority,
		&task.Date,
		userId,
	)

	return err
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db}
}
