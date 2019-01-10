package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DynEd/etest/domain"
	"log"
	"strconv"
	"strings"
	"time"
)

type mysqlOptionRepository struct {
	Conn *sql.DB
}

// NewMysqlOptionRepository will create an object that represent the domain.OptionRepository interface
// This repository stores data into MySQL database
func NewMysqlOptionRepository(conn *sql.DB) domain.OptionRepository {
	return &mysqlOptionRepository{Conn: conn}
}

//fetchData
func (r *mysqlOptionRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*domain.Option, error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	payload := make([]*domain.Option, 0)
	for rows.Next() {
		data := new(domain.Option)

		err := rows.Scan(
			&data.ID,
			&data.Type,
			&data.Sequence,
			&data.IsCorrect,
			&data.IsLast,
			&data.File,
			&data.Text,
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

// Fetch returns all Options
// TODO: implement
func (r *mysqlOptionRepository) Fetch(ctx context.Context, num int64) ([]*domain.Option, error) {
	query := "Select * From options limit ?"

	return r.fetch(ctx, query, num)
}

// GetByID returns single Option matched with given ID
// TODO: implement
func (r *mysqlOptionRepository) GetByID(ctx context.Context, id int) (*domain.Option, error) {

	//query
	query := "SELECT * from options where id=?"

	//fetch data
	rows, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}
	
	payload := &domain.Option{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, domain.ErrNotFound
	}

	return payload, nil

}

// Store saves Option into data storage
// TODO: implement
func (r *mysqlOptionRepository) Store(Option domain.Option) (int64, error) {
	query := "insert into options values(?,?,?,?,?,?,?,?,?)"

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	var typeId int
	if r.Conn.QueryRow(fmt.Sprintf("select id from types where name=?"), Option.Type).Scan(&typeId); typeId == 0 {
		return -1, errors.New("types not supported")
	}

	res, err := stmt.Exec(nil, typeId, ConvertStringToBool(Option.Sequence), ConvertStringToBool(Option.IsCorrect), ConvertStringToBool(Option.IsLast), Option.File, Option.Text, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

// Update updates existing Option in storage
// TODO: implement
func (r *mysqlOptionRepository) Update(Option domain.Option) (domain.Option, error) {
	var typeId int
	var where []string

	//define query
	//query := "Update options set type=?, sequence=? , is_correct=?, is_last=?, file=?, text=?, updated_at=? where id=?"

	//is set type
	if Option.Type != "" {
		// select type name of type
		if r.Conn.QueryRow(fmt.Sprintf("select id from types where name=?"), Option.Type).Scan(&typeId); typeId == 0 {
			return Option, errors.New("types not supported")
		}

		id := strconv.Itoa(typeId)
		where = append(where, fmt.Sprintf("type=%s", id))
	}

	//is set sequence
	if Option.Sequence != "" {
		seq, _ := strconv.ParseBool(Option.Sequence)
		Option.Sequence = strconv.FormatBool(seq)

		where = append(where, fmt.Sprintf("sequence=%s", ConvertStringToBool(Option.Sequence)))
	}

	//is set is correct
	if Option.IsCorrect != "" {
		cor, _ := strconv.ParseBool(Option.IsCorrect)
		Option.IsCorrect = strconv.FormatBool(cor)
		where = append(where, fmt.Sprintf("is_correct=%s", ConvertStringToBool(Option.IsCorrect)))
	}

	//is set is last
	if Option.IsLast != "" {
		las, _ := strconv.ParseBool(Option.IsLast)
		Option.IsLast = strconv.FormatBool(las)

		where = append(where, fmt.Sprintf("is_last=%s", ConvertStringToBool(Option.IsLast)))

	}

	//is set file
	if Option.File != "" {
		where = append(where, fmt.Sprintf("file='%s'", Option.File))
	}

	//is set text
	if Option.Text != "" {
		where = append(where, fmt.Sprintf("text='%s'", Option.Text))
	}
	// merge where's condition
	mergeQuery := strings.Join(where, ", ")
	if mergeQuery != "" {
		mergeQuery = mergeQuery + ","
	}
	query := fmt.Sprintf("Update options set %s updated_at=? where id=?", mergeQuery)
	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		return Option, errors.New(err.Error())
	}

	res, err := stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), Option.ID)

	if err != nil {
		return Option, errors.New(err.Error())
	}

	// check existing data with ID
	_, err = res.RowsAffected()
	if err != nil {
		return Option, domain.ErrNotFound
	}

	defer stmt.Close()

	return Option, nil
}

// Delete deletes matched Option ID from storage
// TODO: implement
func (r *mysqlOptionRepository) Delete(id int) error {
	//query
	query := "delete from options where id=?"

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

//Convert Bool to string
func ConvertStringToBool(str string) string {
	if str == "true" {
		return "1"
	}
	return "0"
}
