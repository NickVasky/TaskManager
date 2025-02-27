package repo

import (
	"database/sql"
	"log"
	"time"
)

type TaskEntity struct {
	Id         int
	UserId     int
	Title      string
	Goal       string
	Measure    string
	Relevance  string
	IsDone     bool
	Deadline   sql.NullTime
	CreatedAt  sql.NullTime
	FinishedAt sql.NullTime
}

func (r *TaskRepo) GetById(id int) (*TaskEntity, error) {
	query := `
		SELECT
			id,
			user_id,
			title,
			goal,
			measure,
			relevance,
			is_done,
			deadline,
			created_at,
			finished_at
		FROM tasks
		WHERE id = $1;
	`
	var t TaskEntity

	err := r.Db.QueryRow(query, id).Scan(
		&t.Id,
		&t.UserId,
		&t.Title,
		&t.Goal,
		&t.Measure,
		&t.Relevance,
		&t.IsDone,
		&t.Deadline,
		&t.CreatedAt,
		&t.FinishedAt)

	if err != nil {
		return nil, err
	}

	log.Printf("Task with id %d was recieved: %v", id, t)

	return &t, nil
}

func (r *TaskRepo) Create(t *TaskEntity) error {
	query := `
	INSERT INTO tasks (
			user_id,
			title,
			goal,
			measure,
			relevance,
			is_done,
			deadline,
			created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`

	t.CreatedAt.Time = time.Now()
	t.CreatedAt.Valid = true

	err := r.Db.QueryRow(
		query,
		t.UserId,
		t.Title,
		t.Goal,
		t.Measure,
		t.Relevance,
		t.IsDone,
		t.Deadline,
		t.CreatedAt).Scan(&t.Id)

	if err != nil {
		return err
	}

	log.Printf("Task \"%s\" created with id %d\n", t.Title, t.Id)

	return nil
}

func (r *TaskRepo) Edit(t *TaskEntity) error {
	query := `
		UPDATE tasks
		SET
			user_id = $2,
			title = $3,
			goal = $4,
			measure = $5,
			relevance = $6,
			is_done = $7,
			deadline = $8,
			created_at = $9,
			finished_at = $10
		WHERE id = $1;
	`

	_, err := r.Db.Exec(
		query,
		t.Id,
		t.UserId,
		t.Title,
		t.Goal,
		t.Measure,
		t.Relevance,
		t.IsDone,
		t.Deadline,
		t.CreatedAt,
		t.FinishedAt)

	if err != nil {
		return err
	}
	log.Printf("Task %d was updated: %v", t.Id, t)

	return nil
}

func (r *TaskRepo) Delete(t *TaskEntity) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1;
	`
	_, err := r.Db.Exec(query, t.Id)

	if err != nil {
		return err
	}

	log.Printf("Task %d was deleted", t.Id)

	return nil
}
