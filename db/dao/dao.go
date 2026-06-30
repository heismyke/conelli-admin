package dao

import "github.com/conelli/admin-backend/db/repo"

type BaseDAO struct {
	repo *repo.Repo
}

func NewBaseDAO(repository *repo.Repo) *BaseDAO {
	return &BaseDAO{repo: repository}
}
