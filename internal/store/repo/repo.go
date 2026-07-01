package repo

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/conelli/admin-backend/config"
	"github.com/conelli/admin-backend/internal/store"
	"github.com/conelli/admin-backend/internal/store/dao"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	store *store.PostgresDBStore
}

func NewRepo(db *store.PostgresDBStore) *Repo {
	return &Repo{
		store: db,
	}
}

func (r *Repo) db(ctx context.Context) *gorm.DB {
	return r.store.DB(ctx)
}

func (r *Repo) Init() error {
	if err := r.db(context.Background()).AutoMigrate(&dao.AdminDataRecord{}); err != nil {
		return err
	}

	var count int64
	if err := r.db(context.Background()).Model(&dao.AdminDataRecord{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return r.SaveAdminData(context.Background(), SeedData())
	}

	return nil
}

func (r *Repo) AdminSummary(ctx context.Context) map[string]string {
	return map[string]string{
		"service": "Conelli Admin API",
		"status":  "ready",
	}
}

func (r *Repo) AdminData(ctx context.Context) (dao.AdminData, error) {
	var record dao.AdminDataRecord
	if err := r.db(ctx).First(&record, 1).Error; err != nil {
		return dao.AdminData{}, err
	}

	var data dao.AdminData
	if err := json.Unmarshal(record.Payload, &data); err != nil {
		return dao.AdminData{}, err
	}

	return data, nil
}

func (r *Repo) SaveAdminData(ctx context.Context, data dao.AdminData) error {
	if err := normalizeUsers(&data); err != nil {
		return err
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	record := dao.AdminDataRecord{ID: 1, Payload: datatypes.JSON(payload)}
	return r.db(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"payload"}),
	}).Create(&record).Error
}

func (r *Repo) Login(ctx context.Context, payload dao.LoginPayload) (*dao.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(payload.Email))
	password := strings.TrimSpace(payload.Password)
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	if email == strings.ToLower(strings.TrimSpace(config.Envs.ADMIN_EMAIL)) && password == strings.TrimSpace(config.Envs.ADMIN_PASSWORD) {
		return &dao.LoginResponse{
			User: dao.AdminUser{
				ID:    "env_admin",
				Name:  config.Envs.ADMIN_NAME,
				Email: config.Envs.ADMIN_EMAIL,
				Role:  "ADMIN",
			},
		}, nil
	}

	data, err := r.AdminData(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range data.Users {
		if strings.EqualFold(user.Email, email) && validPassword(user.PasswordHash, password) {
			return &dao.LoginResponse{User: user}, nil
		}
	}

	return nil, errors.New("invalid staff email or password")
}

func normalizeUsers(data *dao.AdminData) error {
	seen := make(map[string]bool)
	for i := range data.Users {
		email := strings.ToLower(strings.TrimSpace(data.Users[i].Email))
		if email == "" {
			return errors.New("user email is required")
		}
		if seen[email] {
			return errors.New("user email must be unique")
		}
		seen[email] = true
		data.Users[i].Email = email

		if strings.TrimSpace(data.Users[i].PasswordHash) == "" {
			data.Users[i].PasswordHash = "dev"
			continue
		}

		if shouldHashPassword(data.Users[i].PasswordHash) {
			hash, err := bcrypt.GenerateFromPassword([]byte(data.Users[i].PasswordHash), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			data.Users[i].PasswordHash = string(hash)
		}
	}

	return nil
}

func shouldHashPassword(value string) bool {
	if value == "dev" {
		return false
	}

	return !strings.HasPrefix(value, "$2a$") &&
		!strings.HasPrefix(value, "$2b$") &&
		!strings.HasPrefix(value, "$2y$")
}

func validPassword(hash string, password string) bool {
	if hash == "dev" {
		return strings.TrimSpace(password) != ""
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
