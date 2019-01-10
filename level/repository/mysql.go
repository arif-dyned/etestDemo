package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DynEd/etest/domain"
	"log"
	"strings"
	"time"
)

type mysqlLevelRepository struct {
	Conn *sql.DB
}

// NewMysqlLevelRepository will create an object that represent the domain.LevelRepository interface
// This repository stores data into MySQL database
func NewMysqlLevelRepository(conn *sql.DB) domain.LevelRepository {
	return &mysqlLevelRepository{Conn: conn}
}

//fetchData
func (r *mysqlLevelRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*domain.Level, error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	payload := make([]*domain.Level, 0)
	for rows.Next() {
		data := new(domain.Level)

		err := rows.Scan(
			&data.ID,
			&data.Name,
			&data.Sequence,
			&data.CreatedAt,
			&data.UpdatedAt,

		)
		if err != nil {
			return nil, err
		}
		payload = append(payload, data)
	}
	return payload, nil
}

// Fetch returns all Levels
// TODO: implement
func (r *mysqlLevelRepository) Fetch(ctx context.Context, num int64) ([]*domain.Level, error) {
	query := "Select * From levels limit ?"

	return r.fetch(ctx, query, num)
}

// GetByID returns single Level matched with given ID
// TODO: implement
func (r *mysqlLevelRepository) GetByID(ctx context.Context, id int) (*domain.Level, error) {

	//query
	query := "SELECT * from levels where id=?"

	//fetch data
	rows, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := &domain.Level{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, domain.ErrNotFound
	}

	return payload, nil

}

// Store saves Level into data storage
// TODO: implement
func (r *mysqlLevelRepository) Store(Level domain.Level) (int64, error) {
	query := "insert into levels values(?,?,?,?,?)"

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	res, err := stmt.Exec(nil, Level.Name, Level.Sequence, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

// Update updates existing Level in storage
// TODO: implement
func (r *mysqlLevelRepository) Update(Level domain.Level) (domain.Level, error) {
	var where []string

	//is set sequence
	if Level.Sequence > 0 {
		where = append(where, fmt.Sprintf("sequence=%d", Level.Sequence))
	}

	//is set name
	if Level.Name != "" {
		where = append(where, fmt.Sprintf("name='%s'", Level.Name))
	}
	// merge where's condition
	mergeQuery := strings.Join(where, ", ")
	if mergeQuery != "" {
		mergeQuery = mergeQuery + ","
	}
	query := fmt.Sprintf("Update levels set %s updated_at=? where id=?", mergeQuery)
	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		return Level, errors.New(err.Error())
	}

	res, err := stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), Level.ID)

	if err != nil {
		return Level, errors.New(err.Error())
	}

	// check existing data with ID
	_, err = res.RowsAffected()
	if err != nil {
		return Level, domain.ErrNotFound
	}

	defer stmt.Close()

	return Level, nil
}

// Delete deletes matched Level ID from storage
// TODO: implement
func (r *mysqlLevelRepository) Delete(id int) error {
	//query
	query := "delete from levels where id=?"

	//prepare connect to mysql
	stmt, err := r.Conn.Prepare(query)
	if err != nil {
		return errors.New(err.Error())
	}

	//execute query
	_, err = stmt.Exec(id)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
