package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"classifier.buhtigexa.net/internal/cache"
	"errors"
)

// Classifier represents a classifier in our system
// Uses sql.Null types for optional fields to properly handle NULL values from DB
type Classifier struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description,omitempty"` // omitempty hides null values in JSON
	IsActive    sql.NullBool   `json:"is_active,omitempty"`  // omitempty hides null values in JSON
	CreatedAt   time.Time      `json:"created_at"`
}

type ClassifierModel struct {
	DB        *sql.DB
	cache     *cache.Cache
	countStmt *sql.Stmt
	listStmt  *sql.Stmt
}

func NewClassifierModel(db *sql.DB) (*ClassifierModel, error) {
	countStmt, err := db.Prepare("SELECT COUNT(*) FROM classifiers")
	if err != nil {
		return nil, err
	}

	listStmt, err := db.Prepare(`
		SELECT id, name, description, is_active, created_at 
		FROM classifiers 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`)
	if err != nil {
		countStmt.Close()
		return nil, err
	}

	return &ClassifierModel{
		DB:        db,
		cache:     cache.New(),
		countStmt: countStmt,
		listStmt:  listStmt,
	}, nil
}

// Close releases all resources (cache and prepared statements)
func (m *ClassifierModel) Close() error {
	// Primero cerramos la caché y su goroutine
	if err := m.cache.Close(); err != nil {
		return fmt.Errorf("error closing cache: %w", err)
	}

	// Luego cerramos los prepared statements
	return m.CloseStatements()
}

// CloseStatements releases the prepared statements
func (m *ClassifierModel) CloseStatements() error {
	if err := m.countStmt.Close(); err != nil {
		return fmt.Errorf("error closing count statement: %w", err)
	}
	if err := m.listStmt.Close(); err != nil {
		return fmt.Errorf("error closing list statement: %w", err)
	}
	return nil
}

func (m *ClassifierModel) Insert(name string, description string, isActive *bool) (int64, error) {
	// Re piola query para insertar un classifier con los campos nuevos
	query := `INSERT INTO classifiers (name, description, is_active) VALUES (?, ?, ?)`
	
	// Manejamos los campos nullables con mucho cuidado, viste
	var descriptionSQL sql.NullString
	if description != "" {
		descriptionSQL = sql.NullString{String: description, Valid: true}
	}
	
	var isActiveSQL sql.NullBool
	if isActive != nil {
		isActiveSQL = sql.NullBool{Bool: *isActive, Valid: true}
	}
	
	result, err := m.DB.Exec(query, name, descriptionSQL, isActiveSQL)
	if err != nil {
		// Uh, something went wrong with the DB, que quilombo!
		return 0, err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		// Che, couldn't get the ID, alto bardo!
		return 0, err
	}

	// Tenemos que invalidar el cache porque hay data nueva
	// Si no hacemos esto, everything gets desynchronized viste
	m.cache.Delete("classifiers:list")
	return id, nil
}

func (m *ClassifierModel) Get(id int64) (*Classifier, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("classifier:%d", id)
	if cached, ok := m.cache.Get(cacheKey); ok {
		return cached.(*Classifier), nil
	}

	query := `SELECT id, name, description, is_active, created_at FROM classifiers WHERE id = ?`
	c := getClassifier() // Get from pool
	err := m.DB.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description, &c.IsActive, &c.CreatedAt)
	if err != nil {
		putClassifier(c) // Return to pool on error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	// Cache the result for 5 minutes
	m.cache.Set(cacheKey, c, 5*time.Minute)
	return c, nil
}

type ListClassifiersOptions struct {
	Page     int
	PageSize int
}

func (m *ClassifierModel) List(opts ListClassifiersOptions) ([]*Classifier, int, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 20
	}

	// Try to get from cache first
	cacheKey := makeCacheKey("classifiers", "list", strconv.Itoa(opts.Page), strconv.Itoa(opts.PageSize))
	if cached, ok := m.cache.Get(cacheKey); ok {
		return cached.([]*Classifier), 0, nil
	}

	// Get total count using a prepared statement
	var total int
	if err := m.countStmt.QueryRow().Scan(&total); err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (opts.Page - 1) * opts.PageSize

	// Use prepared statement and preallocate slice with exact capacity
	rows, err := m.listStmt.Query(opts.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Preallocate slice with exact capacity needed
	classifiers := make([]*Classifier, 0, opts.PageSize)
	
	// Reuse classifier objects from pool
	for rows.Next() {
		c := getClassifier()
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			// Return objects to pool on error
			for _, cls := range classifiers {
				putClassifier(cls)
			}
			putClassifier(c)
			return nil, 0, err
		}
		classifiers = append(classifiers, c)
	}
	
	if err = rows.Err(); err != nil {
		// Return objects to pool on error
		for _, c := range classifiers {
			putClassifier(c)
		}
		return nil, 0, err
	}

	// Cache the result for 1 minute since this data changes more frequently
	m.cache.Set(cacheKey, classifiers, 1*time.Minute)
	return classifiers, total, nil
}
