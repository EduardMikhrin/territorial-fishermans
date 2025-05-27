package data

import "github.com/EduardMikhrin/territorial-fishermans/internal/types"

type LocationQ interface {
	New() LocationQ
	Get() (*Location, error)
	GetAll() ([]Location, error)
	Insert(data *Location) (int64, error)
	Update(data *Location) error
	Delete(id int64) error
	FilterByID(id int64) LocationQ
}

type Location struct {
	ID          int64  `db:"id" structs:"-"`
	Name        string `db:"name" structs:"name"`
	Description string `db:"description" structs:"description"`
	Photo       string `db:"photo" structs:"photo"`
}

func ToLocation(l *Location) *types.Location {
	return &types.Location{
		Id:          uint64(l.ID),
		Name:        l.Name,
		Description: l.Description,
		Photo:       l.Photo,
	}
}
