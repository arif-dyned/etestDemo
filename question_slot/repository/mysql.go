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

type mysqlQuestionSlotRepository struct {
	Conn *sql.DB
}

// NewMysqlQuestionSlotRepository will create an object that represent the domain.QuestionSlotRepository interface
// This repository stores data into MySQL database
func NewMysqlQuestionSlotRepository(conn *sql.DB) domain.QuestionSlotRepository {
	return &mysqlQuestionSlotRepository{Conn: conn}
}

//fetchData
func (r *mysqlQuestionSlotRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*domain.QuestionSlot, error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	// load all type

	loadTypes, _ := r.AllType()

	payload := make([]*domain.QuestionSlot, 0)
	for rows.Next() {
		data := new(domain.QuestionSlot)

		err := rows.Scan(
			&data.ID,
			&data.Type,
			&data.File,
			&data.Text,
			&data.QuestionId,
			&data.CreatedAt,
			&data.UpdatedAt,

		)
		if err != nil {
			return nil, err
		}

		typeId, _ := strconv.Atoi(data.Type)
		data.Type = loadTypes[int64(typeId)]
		payload = append(payload, data)
	}
	return payload, nil
}

// Fetch returns all QuestionSlots
// TODO: implement
func (r *mysqlQuestionSlotRepository) Fetch(ctx context.Context, num int64) ([]*domain.QuestionSlot, error) {
	query := "Select * From question_slots where question_id=?"

	return r.fetch(ctx, query, num)
}

// GetByID returns single QuestionSlot matched with given ID
// TODO: implement
func (r *mysqlQuestionSlotRepository) GetByID(ctx context.Context, id, slotId int) (*domain.QuestionSlot, error) {

	//query
	query := "SELECT * from question_slots where question_id=? and id=?"

	//fetch data
	rows, err := r.fetch(ctx, query, id, slotId)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	payload := &domain.QuestionSlot{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, domain.ErrNotFound
	}

	return payload, nil

}

// Store saves QuestionSlot into data storage
// TODO: implement
func (r *mysqlQuestionSlotRepository) Store(QuestionSlot domain.QuestionSlot) (int64, error) {
	query := "insert into question_slots values(?,?,?,?,?,?,?)"

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return -1, err
	}
	// check type
	var typeId int
	if r.Conn.QueryRow(fmt.Sprintf("select id from types where name=?"), QuestionSlot.Type).Scan(&typeId); typeId == 0 {
		return -1, errors.New("types not supported")
	}

	res, err := stmt.Exec(nil, typeId, QuestionSlot.File, strings.Replace(QuestionSlot.Text, "'", "\\'", -1), QuestionSlot.QuestionId, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

// Update updates existing QuestionSlot in storage
// TODO: implement
func (r *mysqlQuestionSlotRepository) Update(QuestionSlot domain.QuestionSlot) (domain.QuestionSlot, error) {
	var where []string

	//is set type
	if QuestionSlot.Type != "" {
		// check type
		var typeId int
		if r.Conn.QueryRow(fmt.Sprintf("select id from types where name=?"), QuestionSlot.Type).Scan(&typeId); typeId == 0 {
			return QuestionSlot, errors.New("types not supported")
		}

		where = append(where, fmt.Sprintf("type='%d'", typeId))
	}

	//is set file
	if QuestionSlot.File != "" {
		where = append(where, fmt.Sprintf("file='%s'", QuestionSlot.File))
	}

	//is set text
	if QuestionSlot.Text != "" {
		where = append(where, fmt.Sprintf("text='%s'", strings.Replace(QuestionSlot.Text, "'", "\\'", -1)))
	}


	// merge where's condition
	mergeQuery := strings.Join(where, ", ")
	if mergeQuery != "" {
		mergeQuery = mergeQuery + ","
	}
	query := fmt.Sprintf("Update question_slots set %s updated_at=? where id=?", mergeQuery)
	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		return QuestionSlot, errors.New(err.Error())
	}

	res, err := stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), QuestionSlot.ID)

	if err != nil {
		return QuestionSlot, errors.New(err.Error())
	}

	// check existing data with ID
	_, err = res.RowsAffected()
	if err != nil {
		return QuestionSlot, domain.ErrNotFound
	}

	defer stmt.Close()

	return QuestionSlot, nil
}

// Delete deletes matched QuestionSlot ID from storage
// TODO: implement
func (r *mysqlQuestionSlotRepository) Delete(id int) error {
	//query
	query := "delete from question_slots where id=?"

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
func (r *mysqlQuestionSlotRepository) AllType()(map[int64]string, error)  {
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