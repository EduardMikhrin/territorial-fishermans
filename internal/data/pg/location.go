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
	locationsTableName = "location"
	locationNameField  = "name"
)

type LocationQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func NewLocationQ(db *pgdb.DB) data.LocationQ {
	return &LocationQ{
		db:  db,
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", locationsTableName)),
		upd: sq.Update(locationsTableName),
	}
}

func (q *LocationQ) New() data.LocationQ {
	return NewLocationQ(q.db.Clone())
}

func (q *LocationQ) Get() (*data.Location, error) {
	var result data.Location
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(parsePgError(err), "failed to get client")
	}

	return &result, nil
}

func (q *LocationQ) GetAll() ([]data.Location, error) {
	var result []data.Location
	err := q.db.Select(&result, q.sql)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get locations")
	}

	return result, nil
}

func (q *LocationQ) Update(client *data.Location) error {
	clauses := structs.Map(client)
	if err := q.db.Exec(q.upd.SetMap(clauses)); err != nil {
		return errors.Wrap(err, "failed to update client")
	}

	return nil
}

func (q *LocationQ) Insert(value *data.Location) (int64, error) {
	clauses := structs.Map(value)
	var id int64

	if err := q.db.Get(&id, sq.Insert(locationsTableName).SetMap(clauses).Suffix("returning id")); err != nil {
		return -1, errors.Wrap(parsePgError(err), "failed to insert client")
	}

	return id, nil
}

func (q *LocationQ) Delete(id int64) error {
	if err := q.db.Exec(sq.Delete(locationsTableName).Where(sq.Eq{idField: id})); err != nil {
		return errors.Wrap(err, "failed to delete client")
	}

	return nil
}

func (q *LocationQ) FilterByID(id int64) data.LocationQ {
	q.sql = q.sql.Where(sq.Eq{idField: id})
	q.upd = q.upd.Where(sq.Eq{idField: id})

	return q
}
