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

type mysqlTagsRepository struct {
	Conn *sql.DB
}

// NewMysqlTagsRepository will create an object that represent the domain.TagsRepository interface
// This repository stores data into MySQL database
func NewMysqlTagsRepository(conn *sql.DB) domain.TagsRepository {
	return &mysqlTagsRepository{Conn: conn}
}

//fetchData
func (r *mysqlTagsRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*domain.Tags, error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	payload := make([]*domain.Tags, 0)
	for rows.Next() {
		data := new(domain.Tags)

		err := rows.Scan(
			&data.ID,
			&data.Name,
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

// Fetch returns all Tagss
// TODO: implement
func (r *mysqlTagsRepository) Fetch(ctx context.Context, num int64) ([]*domain.Tags, error) {
	query := "Select * From tags limit ?"

	return r.fetch(ctx, query, num)
}

// GetByID returns single Tags matched with given ID
// TODO: implement
func (r *mysqlTagsRepository) GetByID(ctx context.Context, id int) (*domain.Tags, error) {

	//query
	query := "SELECT * from tags where id=?"

	//fetch data
	rows, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := &domain.Tags{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, domain.ErrNotFound
	}

	return payload, nil

}

// Store saves Tags into data storage
// TODO: implement
func (r *mysqlTagsRepository) Store(Tags domain.Tags) (int64, error) {
	query := "insert into tags values(?,?,?,?)"

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	res, err := stmt.Exec(nil, Tags.Name, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

// Update updates existing Tags in storage
// TODO: implement
func (r *mysqlTagsRepository) Update(Tags domain.Tags) (domain.Tags, error) {
	var where []string

	//is set name
	if Tags.Name != "" {
		where = append(where, fmt.Sprintf("name='%s'", Tags.Name))
	}
	// merge where's condition
	mergeQuery := strings.Join(where, ", ")
	if mergeQuery != "" {
		mergeQuery = mergeQuery + ","
	}
	query := fmt.Sprintf("Update tags set %s updated_at=? where id=?", mergeQuery)
	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		return Tags, errors.New(err.Error())
	}

	res, err := stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), Tags.ID)

	if err != nil {
		return Tags, errors.New(err.Error())
	}

	// check existing data with ID
	_, err = res.RowsAffected()
	if err != nil {
		return Tags, domain.ErrNotFound
	}

	defer stmt.Close()

	return Tags, nil
}

// Delete deletes matched Tags ID from storage
// TODO: implement
func (r *mysqlTagsRepository) Delete(id int) error {
	//query
	query := "delete from tags where id=?"

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
