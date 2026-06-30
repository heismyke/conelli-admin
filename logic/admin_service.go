package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/conelli/admin-backend/db/dao"
	"github.com/conelli/admin-backend/db/repo"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

type AdminUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	Role         string `json:"role"`
	CreatedAt    string `json:"createdAt"`
}

type Investor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	MemberSince string `json:"memberSince"`
	Status      string `json:"status"`
}

type Property struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	Location          string `json:"location"`
	Category          string `json:"category"`
	Status            string `json:"status"`
	ProgressPercent   int    `json:"progressPercent"`
	EstCompletionDate string `json:"estCompletionDate"`
	CoverImageURL     string `json:"coverImageUrl"`
	Description       string `json:"description"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type InvestorProperty struct {
	InvestorID       string   `json:"investorId"`
	PropertyID       string   `json:"propertyId"`
	InvestmentAmount *float64 `json:"investmentAmount"`
	InvestmentDate   string   `json:"investmentDate"`
}

type Update struct {
	ID         string `json:"id"`
	PropertyID string `json:"propertyId"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	PostedAt   string `json:"postedAt"`
	AuthorID   string `json:"authorId"`
}

type Milestone struct {
	ID            string  `json:"id"`
	PropertyID    string  `json:"propertyId"`
	Title         string  `json:"title"`
	PlannedDate   string  `json:"plannedDate"`
	CompletedDate *string `json:"completedDate"`
	Status        string  `json:"status"`
}

type Material struct {
	ID           string  `json:"id"`
	PropertyID   string  `json:"propertyId"`
	MaterialName string  `json:"materialName"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	Status       string  `json:"status"`
	UpdatedAt    string  `json:"updatedAt"`
}

type Document struct {
	ID         string  `json:"id"`
	PropertyID *string `json:"propertyId"`
	InvestorID *string `json:"investorId"`
	Title      string  `json:"title"`
	FileURL    string  `json:"fileUrl"`
	UploadedAt string  `json:"uploadedAt"`
	UploadedBy string  `json:"uploadedBy"`
}

type AdminData struct {
	Users              []AdminUser        `json:"users"`
	Investors          []Investor         `json:"investors"`
	Properties         []Property         `json:"properties"`
	InvestorProperties []InvestorProperty `json:"investorProperties"`
	Updates            []Update           `json:"updates"`
	Milestones         []Milestone        `json:"milestones"`
	Materials          []Material         `json:"materials"`
	Documents          []Document         `json:"documents"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User AdminUser `json:"user"`
}

type AdminDataRecord struct {
	ID      uint           `gorm:"primaryKey"`
	Payload datatypes.JSON `gorm:"type:jsonb;not null"`
}

type AdminService struct {
	adminDAO *dao.AdminDAO
	mu       sync.RWMutex
	data     AdminData
}

func NewAdminService(adminDAO *dao.AdminDAO) *AdminService {
	service := &AdminService{
		adminDAO: adminDAO,
		data:     seedData(),
	}
	if err := service.load(); err != nil {
		if err := service.seedLocked(); err != nil {
			panic(fmt.Sprintf("initialize admin data: %v", err))
		}
	}
	return service
}

func (s *AdminService) Summary(ctx context.Context) (*repo.AdminSummary, error) {
	return s.adminDAO.Summary(), nil
}

func (s *AdminService) Data(ctx context.Context) (AdminData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data, nil
}

func (s *AdminService) SaveData(ctx context.Context, payload AdminData) (AdminData, error) {
	if len(payload.Users) == 0 {
		return AdminData{}, errors.New("at least one staff user is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = payload
	if err := s.saveLocked(); err != nil {
		return AdminData{}, err
	}
	return s.data, nil
}

func (s *AdminService) Login(ctx context.Context, payload LoginPayload) (*LoginResponse, error) {
	if strings.TrimSpace(payload.Email) == "" || strings.TrimSpace(payload.Password) == "" {
		return nil, errors.New("email and password are required")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.data.Users {
		if strings.EqualFold(user.Email, payload.Email) {
			return &LoginResponse{User: user}, nil
		}
	}
	return nil, errors.New("invalid staff email or password")
}

func (s *AdminService) load() error {
	if err := s.adminDAO.DB().AutoMigrate(&AdminDataRecord{}); err != nil {
		return fmt.Errorf("migrate admin data table: %w", err)
	}

	var record AdminDataRecord
	if err := s.adminDAO.DB().First(&record, 1).Error; err != nil {
		return err
	}

	var data AdminData
	if err := json.Unmarshal(record.Payload, &data); err != nil {
		return err
	}
	s.data = data
	return nil
}

func (s *AdminService) saveLocked() error {
	payload, err := json.Marshal(s.data)
	if err != nil {
		return fmt.Errorf("encode admin data: %w", err)
	}

	record := AdminDataRecord{ID: 1, Payload: datatypes.JSON(payload)}
	if err := s.adminDAO.DB().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"payload"}),
	}).Create(&record).Error; err != nil {
		return fmt.Errorf("save admin data: %w", err)
	}
	return nil
}

func (s *AdminService) seedLocked() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = seedData()
	return s.saveLocked()
}

func seedData() AdminData {
	return AdminData{
		Users: []AdminUser{
			{ID: "user_admin", Name: "Ada Okafor", Email: "admin@conelliengineering.com", PasswordHash: "dev", Role: "ADMIN", CreatedAt: "2026-01-08T09:00:00.000Z"},
			{ID: "user_staff", Name: "Site Office", Email: "staff@conelliengineering.com", PasswordHash: "dev", Role: "STAFF", CreatedAt: "2026-01-10T09:00:00.000Z"},
		},
		Investors: []Investor{
			{ID: "inv_1", Name: "Conelli Partner", Email: "partner@conelliengineering.com", Phone: "+234 800 000 1001", MemberSince: "2022-03-01", Status: "active"},
			{ID: "inv_2", Name: "Infrastructure Holdings Ltd", Email: "ops@infrastructureholdings.ng", Phone: "+234 800 000 1002", MemberSince: "2023-06-15", Status: "active"},
		},
		Properties: []Property{
			{ID: "prop_1", Title: "Real Estate Development, Lifecamp Abuja", Location: "Lifecamp, Abuja", Category: "Real Estate Development", Status: "Under Construction", ProgressPercent: 68, EstCompletionDate: "2026-12-15", CoverImageURL: "/assets/real-estate-lifecamp/lifecamp1.png", Description: "Luxury residential units with sustainable design principles and modern amenities.", CreatedAt: "2026-01-05T09:00:00.000Z", UpdatedAt: "2026-06-20T09:00:00.000Z"},
			{ID: "prop_2", Title: "Residential Development, Durumi Abuja", Location: "Durumi, Abuja", Category: "Residential Development", Status: "Foundation Phase", ProgressPercent: 28, EstCompletionDate: "2027-03-31", CoverImageURL: "/assets/residential-durumi/durumi9.jpg", Description: "Contemporary residential development with high quality finishes for urban living.", CreatedAt: "2026-02-11T09:00:00.000Z", UpdatedAt: "2026-06-18T09:00:00.000Z"},
			{ID: "prop_3", Title: "Drainage and Manhole Cover Project, FCDA", Location: "Abuja, Nigeria", Category: "Infrastructure", Status: "Finishing Works", ProgressPercent: 91, EstCompletionDate: "2026-09-30", CoverImageURL: "/assets/drainage-manhole-cover-fcda/manhole8.jpg", Description: "Drainage channels, modern manhole covers, and water management systems across FCDA districts.", CreatedAt: "2026-01-20T09:00:00.000Z", UpdatedAt: "2026-06-22T09:00:00.000Z"},
		},
		InvestorProperties: []InvestorProperty{
			{InvestorID: "inv_1", PropertyID: "prop_1", InvestmentDate: "2022-03-10"},
			{InvestorID: "inv_1", PropertyID: "prop_2", InvestmentDate: "2023-01-18"},
			{InvestorID: "inv_2", PropertyID: "prop_3", InvestmentDate: "2023-07-04"},
		},
		Updates: []Update{
			{ID: "upd_1", PropertyID: "prop_3", Title: "Urban water management systems enhanced", Body: "Drainage channel checks confirmed stronger runoff control across the active FCDA work areas.", PostedAt: "2026-06-22T12:00:00.000Z", AuthorID: "user_admin"},
			{ID: "upd_2", PropertyID: "prop_1", Title: "Structural works continuing in Lifecamp", Body: "The site team completed the latest quality review and confirmed progress against the construction programme.", PostedAt: "2026-06-20T12:00:00.000Z", AuthorID: "user_staff"},
			{ID: "upd_3", PropertyID: "prop_2", Title: "Foundation material delivery logged", Body: "Blockwork and masonry supply has been received for the next foundation work package.", PostedAt: "2026-06-18T12:00:00.000Z", AuthorID: "user_staff"},
		},
		Milestones: []Milestone{
			{ID: "mile_1", PropertyID: "prop_1", Title: "Structural Works", PlannedDate: "2026-08-30", Status: "in_progress"},
			{ID: "mile_2", PropertyID: "prop_1", Title: "Services Installation", PlannedDate: "2026-10-15", Status: "pending"},
			{ID: "mile_3", PropertyID: "prop_2", Title: "Foundation & Structure", PlannedDate: "2026-09-15", Status: "in_progress"},
			{ID: "mile_4", PropertyID: "prop_3", Title: "Final Inspection", PlannedDate: "2026-08-20", Status: "in_progress"},
		},
		Materials: []Material{
			{ID: "mat_1", PropertyID: "prop_1", MaterialName: "Steel Works Package", Quantity: 30, Unit: "tons", Status: "delivered", UpdatedAt: "2026-06-18T09:00:00.000Z"},
			{ID: "mat_2", PropertyID: "prop_1", MaterialName: "Architectural Finishes", Quantity: 240, Unit: "sqm", Status: "ordered", UpdatedAt: "2026-06-21T09:00:00.000Z"},
			{ID: "mat_3", PropertyID: "prop_2", MaterialName: "Foundation Materials", Quantity: 1, Unit: "lot", Status: "delivered", UpdatedAt: "2026-06-17T09:00:00.000Z"},
			{ID: "mat_4", PropertyID: "prop_3", MaterialName: "Access Cover Package", Quantity: 64, Unit: "units", Status: "installed", UpdatedAt: "2026-06-22T09:00:00.000Z"},
		},
		Documents: []Document{
			{ID: "doc_1", PropertyID: ptr("prop_1"), Title: "Progress Summary - Lifecamp Development", FileURL: "/uploads/lifecamp-summary.pdf", UploadedAt: "2026-06-10T09:00:00.000Z", UploadedBy: "user_admin"},
			{ID: "doc_2", PropertyID: ptr("prop_3"), Title: "Installation Completion Certificate", FileURL: "/uploads/fcda-certificate.pdf", UploadedAt: "2026-06-12T09:00:00.000Z", UploadedBy: "user_staff"},
			{ID: "doc_3", InvestorID: ptr("inv_1"), Title: "Investor Statement - Q2", FileURL: "/uploads/investor-statement-q2.pdf", UploadedAt: "2026-06-15T09:00:00.000Z", UploadedBy: "user_admin"},
		},
	}
}

func ptr(value string) *string {
	return &value
}
