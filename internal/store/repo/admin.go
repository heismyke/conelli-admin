package repo

import "github.com/conelli/admin-backend/internal/store/dao"

func SeedData() dao.AdminData {
	return dao.AdminData{
		Users: []dao.AdminUser{
			{ID: "user_admin", Name: "Ada Okafor", Email: "admin@conelliengineering.com", PasswordHash: "dev", Role: "ADMIN", CreatedAt: "2026-01-08T09:00:00.000Z"},
			{ID: "user_staff", Name: "Site Office", Email: "staff@conelliengineering.com", PasswordHash: "dev", Role: "STAFF", CreatedAt: "2026-01-10T09:00:00.000Z"},
		},
		Investors: []dao.Investor{
			{ID: "inv_1", Name: "Conelli Partner", Email: "partner@conelliengineering.com", Phone: "+234 800 000 1001", MemberSince: "2022-03-01", Status: "active"},
			{ID: "inv_2", Name: "Infrastructure Holdings Ltd", Email: "ops@infrastructureholdings.ng", Phone: "+234 800 000 1002", MemberSince: "2023-06-15", Status: "active"},
		},
		Properties: []dao.Property{
			{ID: "prop_1", Title: "Real Estate Development, Lifecamp Abuja", Location: "Lifecamp, Abuja", Category: "Real Estate Development", Status: "Under Construction", ProgressPercent: 68, EstCompletionDate: "2026-12-15", CoverImageURL: "/assets/real-estate-lifecamp/lifecamp1.png", Description: "Luxury residential units with sustainable design principles and modern amenities.", CreatedAt: "2026-01-05T09:00:00.000Z", UpdatedAt: "2026-06-20T09:00:00.000Z"},
			{ID: "prop_2", Title: "Residential Development, Durumi Abuja", Location: "Durumi, Abuja", Category: "Residential Development", Status: "Foundation Phase", ProgressPercent: 28, EstCompletionDate: "2027-03-31", CoverImageURL: "/assets/residential-durumi/durumi9.jpg", Description: "Contemporary residential development with high quality finishes for urban living.", CreatedAt: "2026-02-11T09:00:00.000Z", UpdatedAt: "2026-06-18T09:00:00.000Z"},
			{ID: "prop_3", Title: "Drainage and Manhole Cover Project, FCDA", Location: "Abuja, Nigeria", Category: "Infrastructure", Status: "Finishing Works", ProgressPercent: 91, EstCompletionDate: "2026-09-30", CoverImageURL: "/assets/drainage-manhole-cover-fcda/manhole8.jpg", Description: "Drainage channels, modern manhole covers, and water management systems across FCDA districts.", CreatedAt: "2026-01-20T09:00:00.000Z", UpdatedAt: "2026-06-22T09:00:00.000Z"},
		},
		InvestorProperties: []dao.InvestorProperty{
			{InvestorID: "inv_1", PropertyID: "prop_1", InvestmentDate: "2022-03-10"},
			{InvestorID: "inv_1", PropertyID: "prop_2", InvestmentDate: "2023-01-18"},
			{InvestorID: "inv_2", PropertyID: "prop_3", InvestmentDate: "2023-07-04"},
		},
		Updates: []dao.Update{
			{ID: "upd_1", PropertyID: "prop_3", Title: "Urban water management systems enhanced", Body: "Drainage channel checks confirmed stronger runoff control across the active FCDA work areas.", PostedAt: "2026-06-22T12:00:00.000Z", AuthorID: "user_admin"},
			{ID: "upd_2", PropertyID: "prop_1", Title: "Structural works continuing in Lifecamp", Body: "The site team completed the latest quality review and confirmed progress against the construction programme.", PostedAt: "2026-06-20T12:00:00.000Z", AuthorID: "user_staff"},
			{ID: "upd_3", PropertyID: "prop_2", Title: "Foundation material delivery logged", Body: "Blockwork and masonry supply has been received for the next foundation work package.", PostedAt: "2026-06-18T12:00:00.000Z", AuthorID: "user_staff"},
		},
		Milestones: []dao.Milestone{
			{ID: "mile_1", PropertyID: "prop_1", Title: "Structural Works", PlannedDate: "2026-08-30", Status: "in_progress"},
			{ID: "mile_2", PropertyID: "prop_1", Title: "Services Installation", PlannedDate: "2026-10-15", Status: "pending"},
			{ID: "mile_3", PropertyID: "prop_2", Title: "Foundation & Structure", PlannedDate: "2026-09-15", Status: "in_progress"},
			{ID: "mile_4", PropertyID: "prop_3", Title: "Final Inspection", PlannedDate: "2026-08-20", Status: "in_progress"},
		},
		Materials: []dao.Material{
			{ID: "mat_1", PropertyID: "prop_1", MaterialName: "Steel Works Package", Quantity: 30, Unit: "tons", Status: "delivered", UpdatedAt: "2026-06-18T09:00:00.000Z"},
			{ID: "mat_2", PropertyID: "prop_1", MaterialName: "Architectural Finishes", Quantity: 240, Unit: "sqm", Status: "ordered", UpdatedAt: "2026-06-21T09:00:00.000Z"},
			{ID: "mat_3", PropertyID: "prop_2", MaterialName: "Foundation Materials", Quantity: 1, Unit: "lot", Status: "delivered", UpdatedAt: "2026-06-17T09:00:00.000Z"},
			{ID: "mat_4", PropertyID: "prop_3", MaterialName: "Access Cover Package", Quantity: 64, Unit: "units", Status: "installed", UpdatedAt: "2026-06-22T09:00:00.000Z"},
		},
		Documents: []dao.Document{
			{ID: "doc_1", PropertyID: ptr("prop_1"), Title: "Progress Summary - Lifecamp Development", FileURL: "/uploads/lifecamp-summary.pdf", UploadedAt: "2026-06-10T09:00:00.000Z", UploadedBy: "user_admin"},
			{ID: "doc_2", PropertyID: ptr("prop_3"), Title: "Installation Completion Certificate", FileURL: "/uploads/fcda-certificate.pdf", UploadedAt: "2026-06-12T09:00:00.000Z", UploadedBy: "user_staff"},
			{ID: "doc_3", InvestorID: ptr("inv_1"), Title: "Investor Statement - Q2", FileURL: "/uploads/investor-statement-q2.pdf", UploadedAt: "2026-06-15T09:00:00.000Z", UploadedBy: "user_admin"},
		},
	}
}

func ptr(value string) *string {
	return &value
}
