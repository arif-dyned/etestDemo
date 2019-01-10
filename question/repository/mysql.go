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

type mysqlQuestionRepository struct {
	Conn *sql.DB
}

// NewMysqlQuestionRepository will create an object that represent the domain.QuestionRepository interface
// This repository stores data into MySQL database
func NewMysqlQuestionRepository(conn *sql.DB) domain.QuestionRepository {
	return &mysqlQuestionRepository{Conn: conn}
}

//fetchData
func (r *mysqlQuestionRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*domain.Question, error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	//load All types
	loadTypes, _ := r.AllType()
	payload := make([]*domain.Question, 0)
	for rows.Next() {
		data := new(domain.Question)

		err := rows.Scan(
			&data.ID,
			&data.OptionMode,
			&data.Instructions,
			&data.Comments,
			&data.Type,
			&data.CreatedAt,
			&data.UpdatedAt,

		)
		typeId, _ := strconv.Atoi(data.Type)
		data.Type = loadTypes[int64(typeId)]
		if err != nil {
			return nil, err
		}
		data.OptionMode = OptionModesToString(data.OptionMode)
		payload = append(payload, data)
	}
	return payload, nil
}

// Fetch returns all Questions
// TODO: implement
func (r *mysqlQuestionRepository) Fetch(ctx context.Context, num int64) ([]*domain.Question, error) {
	query := "Select * From questions limit ?"

	return r.fetch(ctx, query, num)
}

// GetByID returns single Question matched with given ID
// TODO: implement
func (r *mysqlQuestionRepository) GetByID(ctx context.Context, id int) (*domain.Question, error) {

	//query
	query := "SELECT * from questions where id=?"

	//fetch data
	rows, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := &domain.Question{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, domain.ErrNotFound
	}

	return payload, nil

}

// Store saves Question into data storage
// TODO: implement
func (r *mysqlQuestionRepository) Store(Question domain.Question) (int64, error) {
	query := "insert into questions values(?,?,?,?,?,?,?)"

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	// check tag
	var tagsId int
	if r.Conn.QueryRow(fmt.Sprintf("select id from tags where name=?"), Question.TagName).Scan(&tagsId); tagsId == 0 {
		return -1, errors.New("tags not supported")
	}

	// check type
	var typeId int
	if r.Conn.QueryRow(fmt.Sprintf("select id from types where name=?"), Question.Type).Scan(&typeId); typeId == 0 {
		return -1, errors.New("types not supported")
	}

	res, err := stmt.Exec(nil, OptionModes(Question.OptionMode), Question.Instructions, strings.Replace(Question.Comments, "'","\\'", -1), typeId, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	//set question id
	Question.ID, _ = res.LastInsertId()

	//set question tag id
	Question.TagName = strconv.Itoa(tagsId)

	//set question type id
	Question.QuestionSlotType = string(typeId)

	//save question tags
	r.SaveTags(Question)

	//save question slots
	r.SaveSlots(Question)

	return res.LastInsertId()
}

// Save Tags of questions
func (r *mysqlQuestionRepository) SaveTags(Question domain.Question) error {
	query := "insert into question_tags values(?,?,?,?,?)"

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return errors.New(err.Error())
	}
	_, err = stmt.Exec(nil, Question.ID, Question.TagName, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	defer stmt.Close()

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

//Save Slots of question
func (r *mysqlQuestionRepository) SaveSlots(Question domain.Question) error {
	query := "insert into question_slots values(?,?,?,?,?,?,?)"

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return errors.New(err.Error())
	}
	_, err = stmt.Exec(nil, Question.QuestionSlotType, Question.QuestionSlotFile, strings.Replace(Question.QuestionSlotText, "'", "\\'", -1), Question.ID, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	defer stmt.Close()

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// Update existing Question in storage
// TODO: implement
func (r *mysqlQuestionRepository) Update(Question domain.Question) (domain.Question, error) {
	var where []string

	//is set name
	if Question.Type != "" {
		var typeId int
		if r.Conn.QueryRow(fmt.Sprintf("select id from types where name=?"), Question.Type).Scan(&typeId); typeId == 0 {
			return Question, errors.New("types not supported")
		}
		where = append(where, fmt.Sprintf("type_id='%d'", typeId))
	}

	// Option mode
	if Question.OptionMode != "" {
		where = append(where, fmt.Sprintf("option_mode='%d'", OptionModes(Question.OptionMode)))
	}

	// Instructions
	if Question.Instructions != "" {
		where = append(where, fmt.Sprintf("instructions='%s'", strings.Replace(Question.Instructions, "'", "\\'", -1)))
	}

	// Comments
	if Question.Comments != ""{
		where = append(where, fmt.Sprintf("comments='%s'", strings.Replace(Question.Comments, "'", "\\'", -1)))
	}

	// merge where's condition
	mergeQuery := strings.Join(where, ", ")
	if mergeQuery != "" {
		mergeQuery = mergeQuery + ","
	}

	query := fmt.Sprintf("Update questions set %s updated_at=? where id=?", mergeQuery)

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		return Question, errors.New(err.Error())
	}

	res, err := stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), Question.ID)

	if err != nil {
		return Question, errors.New(err.Error())
	}

	// check existing data with ID
	_, err = res.RowsAffected()
	if err != nil {
		return Question, domain.ErrNotFound
	}

	// Update Tags
	//r.UpdateTags(Question)

	// Update Slots
	//r.UpdateSlots(Question)

	defer stmt.Close()

	return Question, nil
}

// Update existing Question Tags in storage
func (r *mysqlQuestionRepository) UpdateTags(Question domain.Question) (error) {
	var where []string

	//is set name
	if Question.Type != "" {
		var typeId int
		if r.Conn.QueryRow(fmt.Sprintf("select id from types where name=?"), Question.Type).Scan(&typeId); typeId == 0 {
			return errors.New("types not supported")
		}
		where = append(where, fmt.Sprintf("type_id='%d'", typeId))
	}

	// is set question tags
	if Question.TagName != "" {
		// check tag
		var tagsId int
		if r.Conn.QueryRow(fmt.Sprintf("select id from tags where name=?"), Question.TagName).Scan(&tagsId); tagsId == 0 {
			return errors.New("tags not supported")
		}
		where = append(where, fmt.Sprintf("tag_id='%d'", tagsId))
	}

	// merge where's condition
	mergeQuery := strings.Join(where, ", ")
	if mergeQuery != "" {
		mergeQuery = mergeQuery + ","
	}
	query := fmt.Sprintf("Update question_tags set %s updated_at=? where id=?", mergeQuery)
	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		return errors.New(err.Error())
	}

	res, err := stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), Question.ID)

	if err != nil {
		return errors.New(err.Error())
	}

	// check existing data with ID
	_, err = res.RowsAffected()
	if err != nil {
		return domain.ErrNotFound
	}

	defer stmt.Close()

	return nil
}

// Update existing Question Slots in storage
func (r *mysqlQuestionRepository) UpdateSlots(Question domain.Question) (error) {
	var where []string

	//is set name
	if Question.QuestionSlotType != "" {
		var typeId int
		if r.Conn.QueryRow(fmt.Sprintf("select id from types where name=?"), Question.QuestionSlotType).Scan(&typeId); typeId == 0 {
			return errors.New("types not supported")
		}
		where = append(where, fmt.Sprintf("type='%d'", typeId))
	}

	// is set file
	if Question.QuestionSlotFile != "" {
		where = append(where, fmt.Sprintf("file='%s'", Question.QuestionSlotFile))
	}

	// is set text
	if Question.QuestionSlotText != "" {
		where = append(where, fmt.Sprintf("file='%s'", strings.Replace(Question.QuestionSlotText, "'", "\\'", -1)))
	}

	// merge where's condition
	mergeQuery := strings.Join(where, ", ")
	if mergeQuery != "" {
		mergeQuery = mergeQuery + ","
	}
	query := fmt.Sprintf("Update question_slots set %s updated_at=? where id=?", mergeQuery)
	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		return errors.New(err.Error())
	}

	res, err := stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), Question.ID)

	if err != nil {
		return errors.New(err.Error())
	}

	// check existing data with ID
	_, err = res.RowsAffected()
	if err != nil {
		return domain.ErrNotFound
	}

	defer stmt.Close()

	return nil
}

// Delete deletes matched Question ID from storage
// TODO: implement
func (r *mysqlQuestionRepository) Delete(id int) error {
	//query
	query := "delete from questions where id=?"

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

// Option Modes
func OptionModes(option_mode string) int {
	switch option_mode {
	case "all":
		return 1
	case "abc":
		return 2
	default:
		return 3
	}
}

func OptionModesToString(option_mode string) string {
	switch option_mode {
	case "1":
		return "all"
	case "2":
		return "abc"
	default:
		return "noa"
	}
}

//ALl Type from storage
func (r *mysqlQuestionRepository) AllType()(map[int64]string, error)  {
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