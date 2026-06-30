package repo

import (
	"github.com/conelli/admin-backend/db"
	"gorm.io/gorm"
)

type Repo struct {
	client *db.Client
}

func New(client *db.Client) *Repo {
	return &Repo{client: client}
}

func (r *Repo) Client() *db.Client {
	return r.client
}

func (r *Repo) DB() *gorm.DB {
	return r.client.DB()
}
