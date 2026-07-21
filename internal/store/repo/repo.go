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

func (r *Repo) SeedProjects(ctx context.Context) error {
	seed := SeedData()
	data, err := r.AdminData(ctx)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return r.SaveAdminData(ctx, seed)
	}

	upsertProperties(&data, seed.Properties)
	upsertUpdates(&data, seed.Updates)
	upsertMilestones(&data, seed.Milestones)
	upsertDocuments(&data, seed.Documents)
	upsertInvestorProperties(&data, seed.InvestorProperties)
	upsertPayments(&data, seed.Payments)

	return r.SaveAdminData(ctx, data)
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
	if err := normalizeInvestors(&data); err != nil {
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

func (r *Repo) InvestorLogin(ctx context.Context, payload dao.LoginPayload) (*dao.InvestorLoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(payload.Email))
	password := strings.TrimSpace(payload.Password)
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	data, err := r.AdminData(ctx)
	if err != nil {
		return nil, err
	}

	for _, investor := range data.Investors {
		if strings.EqualFold(investor.Email, email) && investor.Status == "active" && validPassword(investor.PasswordHash, password) {
			investor.PasswordHash = ""
			return &dao.InvestorLoginResponse{Investor: investor}, nil
		}
	}

	return nil, errors.New("invalid investor email or password")
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

func normalizeInvestors(data *dao.AdminData) error {
	seen := make(map[string]bool)
	for i := range data.Investors {
		email := strings.ToLower(strings.TrimSpace(data.Investors[i].Email))
		if email == "" {
			return errors.New("investor email is required")
		}
		if seen[email] {
			return errors.New("investor email must be unique")
		}
		seen[email] = true
		data.Investors[i].Email = email

		if strings.TrimSpace(data.Investors[i].PasswordHash) == "" {
			data.Investors[i].PasswordHash = "dev"
			continue
		}

		if shouldHashPassword(data.Investors[i].PasswordHash) {
			hash, err := bcrypt.GenerateFromPassword([]byte(data.Investors[i].PasswordHash), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			data.Investors[i].PasswordHash = string(hash)
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

func upsertProperties(data *dao.AdminData, seeded []dao.Property) {
	index := make(map[string]int, len(data.Properties))
	for i, property := range data.Properties {
		index[property.ID] = i
	}

	for _, property := range seeded {
		if i, ok := index[property.ID]; ok {
			existing := data.Properties[i]
			if existing.CreatedAt != "" {
				property.CreatedAt = existing.CreatedAt
			}
			if property.UpdatedAt == "" {
				property.UpdatedAt = existing.UpdatedAt
			}
			data.Properties[i] = property
			continue
		}
		data.Properties = append(data.Properties, property)
	}
}

func upsertUpdates(data *dao.AdminData, seeded []dao.Update) {
	index := make(map[string]int, len(data.Updates))
	for i, update := range data.Updates {
		index[update.ID] = i
	}
	for _, update := range seeded {
		if i, ok := index[update.ID]; ok {
			data.Updates[i] = update
			continue
		}
		data.Updates = append(data.Updates, update)
	}
}

func upsertMilestones(data *dao.AdminData, seeded []dao.Milestone) {
	index := make(map[string]int, len(data.Milestones))
	for i, milestone := range data.Milestones {
		index[milestone.ID] = i
	}
	for _, milestone := range seeded {
		if i, ok := index[milestone.ID]; ok {
			data.Milestones[i] = milestone
			continue
		}
		data.Milestones = append(data.Milestones, milestone)
	}
}

func upsertDocuments(data *dao.AdminData, seeded []dao.Document) {
	index := make(map[string]int, len(data.Documents))
	for i, document := range data.Documents {
		index[document.ID] = i
	}
	for _, document := range seeded {
		if i, ok := index[document.ID]; ok {
			data.Documents[i] = document
			continue
		}
		data.Documents = append(data.Documents, document)
	}
}

func upsertInvestorProperties(data *dao.AdminData, seeded []dao.InvestorProperty) {
	exists := make(map[string]bool, len(data.InvestorProperties))
	for _, item := range data.InvestorProperties {
		exists[item.InvestorID+"::"+item.PropertyID] = true
	}
	for _, item := range seeded {
		key := item.InvestorID + "::" + item.PropertyID
		if exists[key] {
			continue
		}
		data.InvestorProperties = append(data.InvestorProperties, item)
	}
}

func upsertPayments(data *dao.AdminData, seeded []dao.Payment) {
	index := make(map[string]int, len(data.Payments))
	for i, payment := range data.Payments {
		index[payment.ID] = i
	}
	for _, payment := range seeded {
		if i, ok := index[payment.ID]; ok {
			data.Payments[i] = payment
			continue
		}
		data.Payments = append(data.Payments, payment)
	}
}
