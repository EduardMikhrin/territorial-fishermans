package pg

import (
	"database/sql"
	"fmt"
	"github.com/EduardMikhrin/territorial-fishermans/internal/data"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	clientTableName = "client"
	idField         = "id"
)

type ClientQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func NewClientQ(db *pgdb.DB) data.ClientQ {
	return &ClientQ{
		db:  db,
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", clientTableName)),
		upd: sq.Update(clientTableName),
	}
}

func (q ClientQ) New() data.ClientQ {
	return NewClientQ(q.db.Clone())
}

func (q *ClientQ) Get() (*data.Client, error) {
	var result data.Client
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get client")
	}

	return &result, nil
}

func (q *ClientQ) GetAll() ([]data.Client, error) {
	var result []data.Client
	err := q.db.Select(&result, q.sql)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get clients")
	}

	return result, nil
}

func (q *ClientQ) Update(client *data.Client) error {
	clauses := structs.Map(client)
	if err := q.db.Exec(q.upd.SetMap(clauses)); err != nil {
		return errors.Wrap(parsePgError(err), "failed to update client")
	}

	return nil
}

func (q *ClientQ) Insert(value *data.Client) (int64, error) {
	clauses := structs.Map(value)
	var id int64

	if err := q.db.Get(&id, sq.Insert(clientTableName).SetMap(clauses).Suffix("returning id")); err != nil {
		return -1, errors.Wrap(parsePgError(err), "failed to insert client")
	}

	return id, nil
}

func (q *ClientQ) Delete(id int64) error {
	if err := q.db.Exec(sq.Delete(clientTableName).Where(sq.Eq{idField: id})); err != nil {
		return errors.Wrap(err, "failed to delete client")
	}

	return nil
}

func (q *ClientQ) FilterByID(id int64) data.ClientQ {
	q.sql = q.sql.Where(sq.Eq{idField: id})
	q.upd = q.upd.Where(sq.Eq{idField: id})

	return q
}
