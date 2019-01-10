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

type mysqlTypesRepository struct {
	Conn *sql.DB
}

// NewMysqlTypesRepository will create an object that represent the domain.TypesRepository interface
// This repository stores data into MySQL database
func NewMysqlTypesRepository(conn *sql.DB) domain.TypesRepository {
	return &mysqlTypesRepository{Conn: conn}
}

//fetchData
func (r *mysqlTypesRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*domain.Types, error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	payload := make([]*domain.Types, 0)
	for rows.Next() {
		data := new(domain.Types)

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

// Fetch returns all Typess
// TODO: implement
func (r *mysqlTypesRepository) Fetch(ctx context.Context, num int64) ([]*domain.Types, error) {
	query := "Select * From types limit ?"

	return r.fetch(ctx, query, num)
}

// GetByID returns single Types matched with given ID
// TODO: implement
func (r *mysqlTypesRepository) GetByID(ctx context.Context, id int) (*domain.Types, error) {

	//query
	query := "SELECT * from types where id=?"

	//fetch data
	rows, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := &domain.Types{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, domain.ErrNotFound
	}

	return payload, nil

}

// Store saves Types into data storage
// TODO: implement
func (r *mysqlTypesRepository) Store(Types domain.Types) (int64, error) {
	query := "insert into types values(?,?,?,?)"

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	res, err := stmt.Exec(nil, Types.Name, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

// Update updates existing Types in storage
// TODO: implement
func (r *mysqlTypesRepository) Update(Types domain.Types) (domain.Types, error) {
	var where []string

	//is set name
	if Types.Name != "" {
		where = append(where, fmt.Sprintf("name='%s'", Types.Name))
	}
	// merge where's condition
	mergeQuery := strings.Join(where, ", ")
	if mergeQuery != "" {
		mergeQuery = mergeQuery + ","
	}
	query := fmt.Sprintf("Update types set %s updated_at=? where id=?", mergeQuery)
	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		return Types, errors.New(err.Error())
	}

	res, err := stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), Types.ID)

	if err != nil {
		return Types, errors.New(err.Error())
	}

	// check existing data with ID
	_, err = res.RowsAffected()
	if err != nil {
		return Types, domain.ErrNotFound
	}

	defer stmt.Close()

	return Types, nil
}

// Delete deletes matched Types ID from storage
// TODO: implement
func (r *mysqlTypesRepository) Delete(id int) error {
	//query
	query := "delete from types where id=?"

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


//ALl Type from storage
func AllType(r *mysqlTypesRepository)(map[int64]string, error)  {
	var cls map[int64]string
	cls = make(map[int64]string)

	rows, err := r.Conn.Query(fmt.Sprintf(`
		SELECT
		id,
		name
		FROM types
	`))

	if err != nil {
		return cls, err
	}

	defer rows.Close()

	for rows.Next() {
		cl := domain.Types{}
		rows.Scan(&cl.ID, &cl.Name)
		cls[cl.ID] = cl.Name
	}
	return cls, nil
}