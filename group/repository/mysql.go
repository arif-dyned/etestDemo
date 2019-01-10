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

type mysqlGroupsRepository struct {
	Conn *sql.DB
}

// NewMysqlGroupsRepository will create an object that represent the domain.GroupsRepository interface
// This repository stores data into MySQL database
func NewMysqlGroupsRepository(conn *sql.DB) domain.GroupsRepository {
	return &mysqlGroupsRepository{Conn: conn}
}

//fetchData
func (r *mysqlGroupsRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*domain.Groups, error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	payload := make([]*domain.Groups, 0)
	for rows.Next() {
		data := new(domain.Groups)

		err := rows.Scan(
			&data.ID,
			&data.Name,
			&data.IsCore,
			&data.IsRefinement,
			&data.Sequence,
			&data.LevelId,
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

// Fetch returns all Groupss
// TODO: implement
func (r *mysqlGroupsRepository) Fetch(ctx context.Context, num int64) ([]*domain.Groups, error) {
	query := "Select * From groups limit ?"

	return r.fetch(ctx, query, num)
}

// GetByID returns single Groups matched with given ID
// TODO: implement
func (r *mysqlGroupsRepository) GetByID(ctx context.Context, id int) (*domain.Groups, error) {

	//query
	query := "SELECT * from groups where id=?"

	//fetch data
	rows, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	payload := &domain.Groups{}
	if len(rows) > 0 {
		payload = rows[0]
	} else {
		return nil, domain.ErrNotFound
	}

	return payload, nil

}

// Store saves Groups into data storage
// TODO: implement
func (r *mysqlGroupsRepository) Store(Groups domain.Groups) (int64, error) {
	query := "insert into groups values(?,?,?,?,?,?,?,?)"

	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return -1, err
	}

	var levelId int
	if r.Conn.QueryRow(fmt.Sprintf("select id from levels where name=?"), Groups.LevelId).Scan(&levelId); levelId == 0 {
		return -1, errors.New("level not supported")
	}
	res, err := stmt.Exec(
		nil,
		Groups.Name,
		ConvertStringToBool(Groups.IsCore),
		ConvertStringToBool(Groups.IsRefinement),
		Groups.Sequence,
		levelId,
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"))

	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

// Update updates existing Groups in storage
// TODO: implement
func (r *mysqlGroupsRepository) Update(Groups domain.Groups) (domain.Groups, error) {
	var where []string

	//is set name
	if Groups.Name != "" {
		where = append(where, fmt.Sprintf("name='%s'", Groups.Name))
	}

	//is set is core
	if Groups.IsCore != "" {
		where = append(where, fmt.Sprintf("is_core='%s'", ConvertStringToBool(Groups.IsCore)))
	}

	//is set is refinement
	if Groups.IsRefinement != "" {
		where = append(where, fmt.Sprintf("is_refinement='%s'", ConvertStringToBool(Groups.IsRefinement)))
	}

	//is set sequence
	if Groups.Sequence > 0 {
		where = append(where, fmt.Sprintf("sequence='%d'", Groups.Sequence))
	}

	//is set sequence
	if Groups.LevelId != "" {
		var levelId int
		if r.Conn.QueryRow(fmt.Sprintf("select id from levels where name=?"), Groups.LevelId).Scan(&levelId); levelId == 0 {
			return Groups, errors.New("level not supported")
		}
		where = append(where, fmt.Sprintf("level_id='%d'", levelId))
	}

	// merge where's condition
	mergeQuery := strings.Join(where, ", ")
	if mergeQuery != "" {
		mergeQuery = mergeQuery + ","
	}
	query := fmt.Sprintf("Update groups set %s updated_at=? where id=?", mergeQuery)
	stmt, err := r.Conn.Prepare(query)

	if err != nil {
		return Groups, errors.New(err.Error())
	}

	res, err := stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), Groups.ID)

	if err != nil {
		return Groups, errors.New(err.Error())
	}

	// check existing data with ID
	_, err = res.RowsAffected()
	if err != nil {
		return Groups, domain.ErrNotFound
	}

	defer stmt.Close()

	return Groups, nil
}

// Delete deletes matched Groups ID from storage
// TODO: implement
func (r *mysqlGroupsRepository) Delete(id int) error {
	//query
	query := "delete from groups where id=?"

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
