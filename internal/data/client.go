package data

import "github.com/EduardMikhrin/territorial-fishermans/internal/types"

type ClientQ interface {
	New() ClientQ
	Get() (*Client, error)
	GetAll() ([]Client, error)
	Insert(data *Client) (int64, error)
	Update(data *Client) error
	Delete(id int64) error
	FilterByID(id int64) ClientQ
}

type Client struct {
	ID        uint64 `db:"id" structs:"-"`
	FirstName string `db:"first_name" structs:"first_name"`
	LastName  string `db:"last_name" structs:"last_name"`
	Photo     string `db:"photo" structs:"photo"`
	Contact   string `db:"contact" structs:"contact"`
}

func ToClient(c *Client) *types.Client {
	return &types.Client{
		Id:      c.ID,
		Name:    c.FirstName,
		Surname: c.LastName,
		Photo:   c.Photo,
		Contact: c.Contact,
	}
}
