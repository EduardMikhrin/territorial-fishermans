package data

import "github.com/EduardMikhrin/territorial-fishermans/internal/types"

type ApplicationStatusQ interface {
	New() ApplicationStatusQ
	Get() (*ApplicationStatus, error)
	GetAll() ([]ApplicationStatus, error)
	Insert(data *ApplicationStatus) (int64, error)
	Update(data *ApplicationStatus) error
	Delete(id int64) error
	FilterByID(id int64) ApplicationStatusQ
}
type ApplicationStatus struct {
	Id   int64  `db:"id" structs:"-"`
	Name string `db:"status_name" structs:"status_name"`
}

func ToApplicationStatus(status ApplicationStatus) *types.ApplicationStatus {
	return &types.ApplicationStatus{
		Id:   uint64(status.Id),
		Name: status.Name,
	}
}
