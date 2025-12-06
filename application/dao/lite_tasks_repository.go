package dao

import (
	"darbelis.eu/stabas/entities"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "modernc.org/sqlite"
)

type LiteTaskRepository struct {
	db *sql.DB
}

func NewLiteTaskRepository(dbPath string) (*LiteTaskRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &LiteTaskRepository{db: db}
	if err := repo.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (repo *LiteTaskRepository) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT,
		result TEXT,
		sender INTEGER,
		receivers TEXT,
		status INTEGER,
		created_at DATETIME,
		sent_at DATETIME,
		received_at DATETIME,
		executing_at DATETIME,
		finished_at DATETIME,
		closed_at DATETIME,
		deleted BOOLEAN DEFAULT 0,
		version INTEGER DEFAULT 1,
		task_group INTEGER
	);
	`
	_, err := repo.db.Exec(schema)
	return err
}

func (repo *LiteTaskRepository) Close() error {
	return repo.db.Close()
}

func (repo *LiteTaskRepository) scanTask(row *sql.Row) (*entities.Task, error) {
	var task entities.Task
	var receiversJSON string
	var createdAt, sentAt, receivedAt, executingAt, finishedAt, closedAt sql.NullTime

	err := row.Scan(
		&task.Id,
		&task.Message,
		&task.Result,
		&task.Sender,
		&receiversJSON,
		&task.Status,
		&createdAt,
		&sentAt,
		&receivedAt,
		&executingAt,
		&finishedAt,
		&closedAt,
		&task.Deleted,
		&task.Version,
		&task.TaskGroup,
	)

	if err != nil {
		return nil, err
	}

	// Parse receivers JSON
	if receiversJSON != "" {
		if err := json.Unmarshal([]byte(receiversJSON), &task.Receivers); err != nil {
			task.Receivers = []int{}
		}
	} else {
		task.Receivers = []int{}
	}

	// Set nullable time fields
	if createdAt.Valid {
		task.CreatedAt = &createdAt.Time
	}
	if sentAt.Valid {
		task.SentAt = &sentAt.Time
	}
	if receivedAt.Valid {
		task.ReceivedAt = &receivedAt.Time
	}
	if executingAt.Valid {
		task.ExecutingAt = &executingAt.Time
	}
	if finishedAt.Valid {
		task.FinishedAt = &finishedAt.Time
	}
	if closedAt.Valid {
		task.ClosedAt = &closedAt.Time
	}

	task.Children = []*entities.Task{}
	return &task, nil
}

func (repo *LiteTaskRepository) scanTasks(rows *sql.Rows) ([]*entities.Task, error) {
	var tasks []*entities.Task

	for rows.Next() {
		var task entities.Task
		var receiversJSON string
		var createdAt, sentAt, receivedAt, executingAt, finishedAt, closedAt sql.NullTime

		err := rows.Scan(
			&task.Id,
			&task.Message,
			&task.Result,
			&task.Sender,
			&receiversJSON,
			&task.Status,
			&createdAt,
			&sentAt,
			&receivedAt,
			&executingAt,
			&finishedAt,
			&closedAt,
			&task.Deleted,
			&task.Version,
			&task.TaskGroup,
		)

		if err != nil {
			return nil, err
		}

		// Parse receivers JSON
		if receiversJSON != "" {
			if err := json.Unmarshal([]byte(receiversJSON), &task.Receivers); err != nil {
				task.Receivers = []int{}
			}
		} else {
			task.Receivers = []int{}
		}

		// Set nullable time fields
		if createdAt.Valid {
			task.CreatedAt = &createdAt.Time
		}
		if sentAt.Valid {
			task.SentAt = &sentAt.Time
		}
		if receivedAt.Valid {
			task.ReceivedAt = &receivedAt.Time
		}
		if executingAt.Valid {
			task.ExecutingAt = &executingAt.Time
		}
		if finishedAt.Valid {
			task.FinishedAt = &finishedAt.Time
		}
		if closedAt.Valid {
			task.ClosedAt = &closedAt.Time
		}

		task.Children = []*entities.Task{}
		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}

func (repo *LiteTaskRepository) FindById(id int) (*entities.Task, error) {
	query := `
		SELECT id, message, result, sender, receivers, status,
			   created_at, sent_at, received_at, executing_at, finished_at, closed_at,
			   deleted, version, task_group
		FROM tasks
		WHERE id = ?
	`

	row := repo.db.QueryRow(query, id)
	task, err := repo.scanTask(row)
	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	}
	return task, err
}

func (repo *LiteTaskRepository) FindAll() []*entities.Task {
	query := `
		SELECT id, message, result, sender, receivers, status,
			   created_at, sent_at, received_at, executing_at, finished_at, closed_at,
			   deleted, version, task_group
		FROM tasks
		WHERE deleted = 0
	`

	rows, err := repo.db.Query(query)
	if err != nil {
		return []*entities.Task{}
	}
	defer rows.Close()

	tasks, err := repo.scanTasks(rows)
	if err != nil {
		return []*entities.Task{}
	}

	return tasks
}

func (repo *LiteTaskRepository) AddTask(task *entities.Task) int {
	receiversJSON, _ := json.Marshal(task.Receivers)

	query := `
		INSERT INTO tasks (message, result, sender, receivers, status,
						   created_at, sent_at, received_at, executing_at, finished_at, closed_at,
						   deleted, version, task_group)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := repo.db.Exec(query,
		task.Message,
		task.Result,
		task.Sender,
		string(receiversJSON),
		task.Status,
		task.CreatedAt,
		task.SentAt,
		task.ReceivedAt,
		task.ExecutingAt,
		task.FinishedAt,
		task.ClosedAt,
		task.Deleted,
		task.Version,
		task.TaskGroup,
	)

	if err != nil {
		return 0
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0
	}

	task.Id = int(id)

	// If TaskGroup is 0, set it to the task ID
	if task.TaskGroup == 0 {
		task.TaskGroup = task.Id
		repo.db.Exec("UPDATE tasks SET task_group = ? WHERE id = ?", task.TaskGroup, task.Id)
	}

	return task.Id
}

func (repo *LiteTaskRepository) UpdateTask(task *entities.Task) (*entities.Task, error) {
	// First check if task exists
	existingTask, err := repo.FindById(task.Id)
	if err != nil {
		return nil, err
	}

	receiversJSON, _ := json.Marshal(task.Receivers)

	query := `
		UPDATE tasks
		SET message = ?, result = ?, sender = ?, receivers = ?, status = ?,
			created_at = ?, sent_at = ?, received_at = ?, executing_at = ?,
			finished_at = ?, closed_at = ?, version = ?, task_group = ?
		WHERE id = ?
	`

	taskGroup := task.TaskGroup
	if taskGroup == 0 {
		taskGroup = existingTask.TaskGroup
	}

	_, err = repo.db.Exec(query,
		task.Message,
		task.Result,
		task.Sender,
		string(receiversJSON),
		task.Status,
		task.CreatedAt,
		task.SentAt,
		task.ReceivedAt,
		task.ExecutingAt,
		task.FinishedAt,
		task.ClosedAt,
		task.Version,
		taskGroup,
		task.Id,
	)

	if err != nil {
		return nil, err
	}

	return repo.FindById(task.Id)
}

func (repo *LiteTaskRepository) DeleteTask(id int) error {
	query := "UPDATE tasks SET deleted = 1 WHERE id = ?"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (repo *LiteTaskRepository) UpdateTaskWithValidation(task *entities.Task) (*entities.Task, error) {
	existingTask, err := repo.FindById(task.Id)
	if err != nil {
		return nil, err
	}

	if task.Version != existingTask.Version+1 {
		return existingTask, errors.New("Wrong task version")
	}

	return repo.UpdateTask(task)
}

func (repo *LiteTaskRepository) GetCountWithSameGroup(groupId int) int {
	query := "SELECT COUNT(*) FROM tasks WHERE task_group = ? AND deleted = 0"
	var count int
	err := repo.db.QueryRow(query, groupId).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
