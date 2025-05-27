package pg

import (
	"database/sql"
	"fmt"
	"github.com/EduardMikhrin/territorial-fishermans/internal/data"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

const (
	applicationTableName = "application"
	locationIdField      = "location_id"
	statusIdField        = "status_id"
	fishingDateField     = "fishing_date"
	clientIdField        = "client_id"
	applicationIdField   = "id"
)

type ApplicationQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func NewApplicationQ(db *pgdb.DB) data.ApplicationQ {
	return &ApplicationQ{
		db:  db,
		sql: baseSelect(),
		upd: sq.Update(applicationTableName),
	}
}

func (q *ApplicationQ) New() data.ApplicationQ {
	return NewApplicationQ(q.db.Clone())
}

func (q *ApplicationQ) Get() (*data.Application, error) {
	var result data.Application

	err := q.db.Get(&result, q.sql)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get application")
	}

	return &result, nil
}

func (q *ApplicationQ) GetAll() ([]data.Application, error) {
	var result []data.Application

	err := q.db.Select(&result, q.sql)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get applications")
	}

	return result, nil
}

func (q *ApplicationQ) Update(client *data.Application) error {
	clauses := structs.Map(client)
	if err := q.db.Exec(q.upd.SetMap(clauses)); err != nil {
		return errors.Wrap(parsePgError(err), "failed to update application")
	}

	return nil
}

func (q *ApplicationQ) Insert(value *data.Application) (int64, error) {
	value.CreatedAt = time.Now()
	clauses := structs.Map(value)
	var (
		err error
		id  int64
	)

	if err = q.db.Get(&id, sq.Insert(applicationTableName).SetMap(clauses).Suffix("returning id")); err != nil {
		return -1, errors.Wrap(parsePgError(err), "failed to insert application")
	}

	return id, nil
}

func (q *ApplicationQ) Delete(id uint64) error {
	if err := q.db.Exec(sq.Delete(applicationTableName).Where(sq.Eq{idField: id})); err != nil {
		return errors.Wrap(err, "failed to delete application")
	}

	return nil
}

func (q *ApplicationQ) FilterByID(id int64) data.ApplicationQ {
	q.sql = q.sql.Where(sq.Eq{"a." + idField: id})
	q.upd = q.upd.Where(sq.Eq{"a." + idField: id})

	return q
}

func baseSelect() sq.SelectBuilder {
	return sq.
		Select(
			"a.id",
			"a.fishing_date",
			"a.client_id",
			"a.location_id",
			"l.name          AS location",
			"a.status_id",
			"s.status_name   AS status",
			"a.created_at",
		).
		From(fmt.Sprintf("%s AS a", applicationTableName)).                                     // application AS a
		LeftJoin(fmt.Sprintf("%s AS l ON l.id = a.%s", locationsTableName, locationIdField)).   // JOIN location
		LeftJoin(fmt.Sprintf("%s AS s ON s.id = a.%s", applicationStatusTable, statusIdField)). // JOIN application_status
		PlaceholderFormat(sq.Dollar)
}
