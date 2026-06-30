package logic

import (
	"github.com/conelli/admin-backend/config"
	"github.com/conelli/admin-backend/db/dao"
	"github.com/conelli/admin-backend/db/repo"
)

type CoreLogic struct {
	cfg *config.Config

	adminDAO *dao.AdminDAO

	Admin *AdminService
}

func New(cfg *config.Config, repository *repo.Repo) *CoreLogic {
	adminDAO := dao.NewAdminDAO(repository)
	adminService := NewAdminService(adminDAO)

	return &CoreLogic{
		cfg:      cfg,
		adminDAO: adminDAO,
		Admin:    adminService,
	}
}
