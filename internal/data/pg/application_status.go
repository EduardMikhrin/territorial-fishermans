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
	applicationStatusTable = "applicationstatus"
)

type ApplicationStatusQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
	upd sq.UpdateBuilder
}

func NewApplicationStatusQ(db *pgdb.DB) data.ApplicationStatusQ {
	return &ApplicationStatusQ{
		db:  db,
		sql: sq.Select("b.*").From(fmt.Sprintf("%s as b", applicationStatusTable)),
		upd: sq.Update(applicationStatusTable),
	}
}

func (q ApplicationStatusQ) New() data.ApplicationStatusQ {
	return NewApplicationStatusQ(q.db.Clone())
}

func (q *ApplicationStatusQ) Get() (*data.ApplicationStatus, error) {
	var result data.ApplicationStatus
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get application status")
	}

	return &result, nil
}

func (q *ApplicationStatusQ) GetAll() ([]data.ApplicationStatus, error) {
	var result []data.ApplicationStatus
	err := q.db.Select(&result, q.sql)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get application statuses")
	}

	return result, nil
}

func (q *ApplicationStatusQ) Update(data *data.ApplicationStatus) error {
	clauses := structs.Map(data)
	if err := q.db.Exec(q.upd.SetMap(clauses)); err != nil {
		return errors.Wrap(parsePgError(err), "failed to update application status")
	}

	return nil
}

func (q *ApplicationStatusQ) Insert(value *data.ApplicationStatus) (int64, error) {
	clauses := structs.Map(value)
	var id int64

	if err := q.db.Get(&id, sq.Insert(applicationStatusTable).SetMap(clauses).Suffix("returning id")); err != nil {
		return -1, errors.Wrap(parsePgError(err), "failed to insert application status")
	}

	return id, nil
}

func (q *ApplicationStatusQ) Delete(id int64) error {
	if err := q.db.Exec(sq.Delete(applicationStatusTable).Where(sq.Eq{idField: id})); err != nil {
		return errors.Wrap(err, "failed to delete application status")
	}

	return nil
}

func (q *ApplicationStatusQ) FilterByID(id int64) data.ApplicationStatusQ {
	q.sql = q.sql.Where(sq.Eq{idField: id})
	q.upd = q.upd.Where(sq.Eq{idField: id})

	return q
}
