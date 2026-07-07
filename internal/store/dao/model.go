package dao

import "gorm.io/datatypes"

type AdminUser struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	PasswordHash    string `json:"passwordHash"`
	Role            string `json:"role"`
	ProfileImageURL string `json:"profileImageUrl"`
	CreatedAt       string `json:"createdAt"`
}

type Investor struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	MemberSince  string `json:"memberSince"`
	Status       string `json:"status"`
	PasswordHash string `json:"passwordHash,omitempty"`
}

type Property struct {
	ID                string   `json:"id"`
	Title             string   `json:"title"`
	Location          string   `json:"location"`
	Category          string   `json:"category"`
	Status            string   `json:"status"`
	ProgressPercent   int      `json:"progressPercent"`
	EstCompletionDate string   `json:"estCompletionDate"`
	CoverImageURL     string   `json:"coverImageUrl"`
	Description       string   `json:"description"`
	PublicDescription string   `json:"publicDescription"`
	PublicOverview    string   `json:"publicOverview"`
	Client            string   `json:"client"`
	Year              string   `json:"year"`
	Tags              []string `json:"tags"`
	GalleryImages     []string `json:"galleryImages"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
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

type AdminDataRecord struct {
	ID      uint           `gorm:"primaryKey"`
	Payload datatypes.JSON `gorm:"type:jsonb;not null"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User AdminUser `json:"user"`
}

type InvestorLoginResponse struct {
	Investor Investor `json:"investor"`
}
