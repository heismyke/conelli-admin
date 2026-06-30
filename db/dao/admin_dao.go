package dao

import (
	"github.com/conelli/admin-backend/db/repo"
	"gorm.io/gorm"
)

type AdminDAO struct {
	*BaseDAO
}

func NewAdminDAO(repository *repo.Repo) *AdminDAO {
	return &AdminDAO{BaseDAO: NewBaseDAO(repository)}
}

func (d *AdminDAO) Summary() *repo.AdminSummary {
	return &repo.AdminSummary{
		Service: "Conelli Admin API",
		Status:  "ready",
	}
}

func (d *AdminDAO) DB() *gorm.DB {
	return d.repo.DB()
}
