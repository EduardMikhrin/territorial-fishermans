package data

import (
	"github.com/EduardMikhrin/territorial-fishermans/internal/types"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ApplicationQ interface {
	New() ApplicationQ
	Get() (*Application, error)
	GetAll() ([]Application, error)
	Insert(data *Application) (int64, error)
	Update(data *Application) error
	Delete(id uint64) error
	FilterByID(id int64) ApplicationQ
}
type Application struct {
	ID          uint64    `db:"id" structs:"-"`
	FishingDate time.Time `db:"fishing_date" structs:"fishing_date"`
	ClientID    uint64    `db:"client_id" structs:"client_id"`
	LocationId  uint64    `db:"location_id" structs:"location_id"`
	Location    string    `db:"location" structs:"-"`
	Status      string    `db:"status" structs:"-"`
	StatusId    uint64    `db:"status_id" structs:"status_id"`
	CreatedAt   time.Time `db:"created_at" structs:"created_at"`
}

func ToApplication(application *Application) *types.Application {
	return &types.Application{
		Id:          application.ID,
		ClientId:    application.ClientID,
		FishingDate: timestamppb.New(application.FishingDate),
		Location:    application.Location,
		Status:      application.Status,
		CreatedAt:   timestamppb.New(application.CreatedAt),
	}
}
