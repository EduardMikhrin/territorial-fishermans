package pg

import (
	"database/sql"
	"github.com/EduardMikhrin/territorial-fishermans/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type masterQ struct {
	db *pgdb.DB
}

func NewMasterQ(db *pgdb.DB) data.MasterQ {
	return &masterQ{db: db}
}

func (m masterQ) New() data.MasterQ {
	return NewMasterQ(m.db)
}

func (m masterQ) ClientQ() data.ClientQ {
	return NewClientQ(m.db)
}

func (m masterQ) ApplicationQ() data.ApplicationQ {
	return NewApplicationQ(m.db)
}

func (m masterQ) LocationQ() data.LocationQ {
	return NewLocationQ(m.db)
}

func (m masterQ) ApplicationStatusQ() data.ApplicationStatusQ { return NewApplicationStatusQ(m.db) }

func (m masterQ) Transaction(fn func(data interface{}) error, i interface{}) error {
	return m.db.TransactionWithOptions(&sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}, func() error {
		return fn(i)
	})
}
