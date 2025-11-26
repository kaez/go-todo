package repository

import (
	"database/sql"
	"time"

	"github.com/kaez/go-todo/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(dbPath string) (*TodoRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &TodoRepository{db: db}
	if err := repo.initDB(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *TodoRepository) initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *TodoRepository) Create(req models.CreateTodoRequest) (*models.Todo, error) {
	query := `INSERT INTO todos (title, description) VALUES (?, ?)`
	result, err := r.db.Exec(query, req.Title, req.Description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(int(id))
}

func (r *TodoRepository) GetAll() ([]models.Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM todos ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *TodoRepository) GetByID(id int) (*models.Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = ?`
	var todo models.Todo
	err := r.db.QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) Update(id int, req models.UpdateTodoRequest) (*models.Todo, error) {
	todo, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	query := `UPDATE todos SET title = ?, description = ?, completed = ?, updated_at = ? WHERE id = ?`
	_, err = r.db.Exec(query, todo.Title, todo.Description, todo.Completed, time.Now(), id)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *TodoRepository) Delete(id int) error {
	query := `DELETE FROM todos WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TodoRepository) Close() error {
	return r.db.Close()
}
