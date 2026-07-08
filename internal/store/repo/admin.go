package repo

import "github.com/conelli/admin-backend/internal/store/dao"

func SeedData() dao.AdminData {
	return dao.AdminData{
		Users: []dao.AdminUser{
			{ID: "user_admin", Name: "Ada Okafor", Email: "admin@conelliengineering.com", PasswordHash: "dev", Role: "ADMIN", CreatedAt: "2026-01-08T09:00:00.000Z"},
			{ID: "user_staff", Name: "Site Office", Email: "staff@conelliengineering.com", PasswordHash: "dev", Role: "STAFF", CreatedAt: "2026-01-10T09:00:00.000Z"},
		},
		Investors: []dao.Investor{
			{ID: "inv_1", Name: "Conelli Partner", Email: "partner@conelliengineering.com", Phone: "+234 800 000 1001", MemberSince: "2022-03-01", Status: "active", PasswordHash: "dev"},
			{ID: "inv_2", Name: "Infrastructure Holdings Ltd", Email: "ops@infrastructureholdings.ng", Phone: "+234 800 000 1002", MemberSince: "2023-06-15", Status: "active", PasswordHash: "dev"},
		},
		Properties: []dao.Property{
			{
				ID: "prop_3", Title: "Drainage and Manhole Cover Project, FCDA", Location: "Abuja, Nigeria", Category: "Infrastructure", Status: "Finishing Works", ProgressPercent: 91, EstCompletionDate: "2026-09-30", CoverImageURL: "/assets/drainage-manhole-cover-fcda/manhole8.jpg",
				Description:       "Drainage channels, modern manhole covers, and water management systems across FCDA districts.",
				PublicDescription: "Comprehensive drainage infrastructure and manhole cover installation project for the Federal Capital Development Authority, enhancing urban water management systems.",
				PublicOverview:    "This large-scale infrastructure project involved the design, installation, and implementation of a comprehensive drainage system for the Federal Capital Development Authority. The project included the installation of modern manhole covers, drainage channels, and water management systems to improve urban flood control and water runoff management across multiple districts.",
				Client:            "Federal Capital Development Authority", Year: "2023", Tags: []string{"Infrastructure", "Drainage", "FCDA"},
				GalleryImages: []string{"/assets/drainage-manhole-cover-fcda/manhole1.jpg", "/assets/drainage-manhole-cover-fcda/manhole2.jpg", "/assets/drainage-manhole-cover-fcda/manhole3.jpg", "/assets/drainage-manhole-cover-fcda/manhole4.jpg", "/assets/drainage-manhole-cover-fcda/manhole5.jpg", "/assets/drainage-manhole-cover-fcda/manhole6.jpg", "/assets/drainage-manhole-cover-fcda/manhole7.jpg", "/assets/drainage-manhole-cover-fcda/manhole8.jpg", "/assets/drainage-manhole-cover-fcda/manhole9.jpg"},
				CreatedAt:     "2026-01-20T09:00:00.000Z", UpdatedAt: "2026-06-22T09:00:00.000Z",
			},
			{
				ID: "prop_1", Title: "Real Estate Development, Lifecamp Abuja", Location: "Lifecamp, Abuja", Category: "Real Estate Development", Status: "Under Construction", ProgressPercent: 68, EstCompletionDate: "2026-12-15", CoverImageURL: "/assets/real-estate-lifecamp/lifecamp1.png",
				Description:       "Luxury residential units with sustainable design principles and modern amenities.",
				PublicDescription: "Premium real estate development project in Lifecamp, Abuja, featuring modern residential and commercial properties designed for contemporary urban living.",
				PublicOverview:    "A residential development project in the prestigious Lifecamp area, featuring luxury residential units. The project incorporates sustainable design principles, modern amenities, and quality construction standards.",
				Client:            "Private Developer", Year: "2023", Tags: []string{"Real Estate", "Development", "Lifecamp"},
				GalleryImages: []string{"/assets/real-estate-lifecamp/lifecamp1.png", "/assets/real-estate-lifecamp/lifecamp2.png", "/assets/real-estate-lifecamp/lifecamp3.png", "/assets/real-estate-lifecamp/lifecamp4.png", "/assets/real-estate-lifecamp/lifecamp5.png", "/assets/real-estate-lifecamp/lifecamp6.jpg"},
				CreatedAt:     "2026-01-05T09:00:00.000Z", UpdatedAt: "2026-06-20T09:00:00.000Z",
			},
			{
				ID: "prop_2", Title: "Residential Development, Durumi Abuja", Location: "Durumi, Abuja", Category: "Residential Development", Status: "Foundation Phase", ProgressPercent: 28, EstCompletionDate: "2027-03-31", CoverImageURL: "/assets/residential-durumi/durumi9.jpg",
				Description:       "Contemporary residential development with high quality finishes for urban living.",
				PublicDescription: "Modern residential development project in Durumi, Abuja, featuring contemporary design and quality construction for comfortable urban living.",
				PublicOverview:    "A premium residential development featuring modern architectural design, sustainable building practices, and high-quality finishes. The project includes multiple residential units with contemporary amenities designed for comfortable urban living.",
				Client:            "Private Developer", Year: "2023", Tags: []string{"Real Estate", "Residential", "Development"},
				GalleryImages: []string{"/assets/residential-durumi/durumi1.jpg", "/assets/residential-durumi/durumi2.jpg", "/assets/residential-durumi/durumi8.jpg", "/assets/residential-durumi/durumi9.jpg"},
				CreatedAt:     "2026-02-11T09:00:00.000Z", UpdatedAt: "2026-06-18T09:00:00.000Z",
			},
			{
				ID: "prop_4", Title: "Mosque, FIRS Abuja", Location: "Abuja, Nigeria", Category: "Religious Facility", Status: "Completed", ProgressPercent: 100, EstCompletionDate: "2023-12-15", CoverImageURL: "/assets/mosque-firs/worship_center1.jpg",
				Description:       "Design and construction of a modern worship center for the Federal Inland Revenue Service.",
				PublicDescription: "Design and construction of a modern worship center for the Federal Inland Revenue Service, creating a serene and functional space for spiritual activities.",
				PublicOverview:    "A beautifully designed worship center featuring modern architecture, excellent acoustics, and comfortable seating arrangements. The project included specialized lighting, sound systems, and climate control for optimal worship experience. Mosque Construction No 14 Sokode St, FIRS",
				Client:            "Federal Inland Revenue Service", Year: "2023", Tags: []string{"Religious Facility", "Construction", "FIRS"},
				GalleryImages: []string{"/assets/mosque-firs/worship_center1.jpg", "/assets/mosque-firs/worship_center2.jpg", "/assets/mosque-firs/worship_center3.jpg", "/assets/mosque-firs/worship_center4.jpg", "/assets/mosque-firs/worship_center5.jpg", "/assets/mosque-firs/worship_center6.jpg", "/assets/mosque-firs/worship_center7.jpg", "/assets/mosque-firs/worship_center8.jpg", "/assets/mosque-firs/worship_center9.jpg", "/assets/mosque-firs/worship_center10.jpg", "/assets/mosque-firs/worship_center11.jpg", "/assets/mosque-firs/worship_center12.jpg", "/assets/mosque-firs/worship_center13.jpg", "/assets/mosque-firs/worship_center14.jpg"},
				CreatedAt:     "2026-01-22T09:00:00.000Z", UpdatedAt: "2026-06-10T09:00:00.000Z",
			},
			{
				ID: "prop_5", Title: "School Remodelling, GTC Danbatta", Location: "Danbatta, Kano", Category: "Education", Status: "Completed", ProgressPercent: 100, EstCompletionDate: "2022-12-15", CoverImageURL: "/assets/school-remodelling-danbatta/school4.jpg",
				Description:       "Comprehensive remodelling of Government Technical College Danbatta.",
				PublicDescription: "Comprehensive remodelling of Government Technical College Danbatta, upgrading facilities to provide enhanced learning environments for students.",
				PublicOverview:    "Complete renovation and modernization of Government Technical College Danbatta, including classroom upgrades, laboratory renovations, administrative block remodelling, and infrastructure improvements to create a conducive learning environment.",
				Client:            "Government Technical College", Year: "2022", Tags: []string{"Education", "Remodelling", "Renovation"},
				GalleryImages: []string{"/assets/school-remodelling-danbatta/school8.jpg", "/assets/school-remodelling-danbatta/school7.jpg", "/assets/school-remodelling-danbatta/school1.jpg", "/assets/school-remodelling-danbatta/school4.jpg", "/assets/school-remodelling-danbatta/school5.jpg", "/assets/school-remodelling-danbatta/school6.jpg", "/assets/school-remodelling-danbatta/school2.jpg", "/assets/school-remodelling-danbatta/school3.jpg"},
				CreatedAt:     "2026-01-24T09:00:00.000Z", UpdatedAt: "2026-06-09T09:00:00.000Z",
			},
			{
				ID: "prop_6", Title: "Concrete Road Construction, Baze University", Location: "Abuja, Nigeria", Category: "Infrastructure", Status: "Completed", ProgressPercent: 100, EstCompletionDate: "2022-10-15", CoverImageURL: "/assets/concrete-road-baze/construct1.jpg",
				Description:       "High-quality concrete road construction project at Baze University.",
				PublicDescription: "High-quality concrete road construction project at Baze University, providing durable and sustainable transportation infrastructure for the campus community.",
				PublicOverview:    "A comprehensive road construction project that transformed the internal road network of Baze University. The project utilized high-grade concrete and modern construction techniques to create durable, weather-resistant roads that enhance campus accessibility and aesthetics.",
				Client:            "Baze University", Year: "2022", Tags: []string{"Road Construction", "Education", "Infrastructure"},
				GalleryImages: []string{"/assets/concrete-road-baze/construct1.jpg", "/assets/concrete-road-baze/construct2.jpg", "/assets/concrete-road-baze/construct3.jpg", "/assets/concrete-road-baze/construct4.jpg", "/assets/concrete-road-baze/construct5.jpg", "/assets/concrete-road-baze/construct6.jpg", "/assets/concrete-road-baze/construct7.jpg", "/assets/concrete-road-baze/construct8.jpg", "/assets/concrete-road-baze/construct9.jpg", "/assets/concrete-road-baze/construct10.jpg", "/assets/concrete-road-baze/construct11.jpg", "/assets/concrete-road-baze/construct12.jpg", "/assets/concrete-road-baze/construct13jpg.jpg", "/assets/concrete-road-baze/construct14.jpg", "/assets/concrete-road-baze/construct15.jpg", "/assets/concrete-road-baze/construct16.jpg", "/assets/concrete-road-baze/construct17.jpg", "/assets/concrete-road-baze/construct18.jpg", "/assets/concrete-road-baze/construct19.jpg"},
				CreatedAt:     "2026-01-26T09:00:00.000Z", UpdatedAt: "2026-06-08T09:00:00.000Z",
			},
			{
				ID: "prop_7", Title: "Sports Court Construction, Baze University", Location: "Abuja, Nigeria", Category: "Sports Facility", Status: "Completed", ProgressPercent: 100, EstCompletionDate: "2022-11-15", CoverImageURL: "/assets/sports-court-baze/sport1.jpg",
				Description:       "Professional sports court construction at Baze University.",
				PublicDescription: "Professional sports court construction at Baze University, providing state-of-the-art recreational facilities for students and staff.",
				PublicOverview:    "Construction of modern sports courts including basketball, volleyball, and tennis facilities. The project incorporated professional-grade surfaces, proper drainage systems, and lighting for extended playing hours.",
				Client:            "Baze University", Year: "2022", Tags: []string{"Sports Facility", "Construction", "Education"},
				GalleryImages: []string{"/assets/sports-court-baze/sport1.jpg", "/assets/sports-court-baze/sport2.jpg", "/assets/sports-court-baze/sport3.jpg", "/assets/sports-court-baze/sport4.jpg", "/assets/sports-court-baze/sport5.jpg", "/assets/sports-court-baze/sport6.jpg", "/assets/sports-court-baze/sport7.jpg", "/assets/sports-court-baze/sport8.jpg", "/assets/sports-court-baze/sport9.jpg", "/assets/sports-court-baze/sport10.jpg"},
				CreatedAt:     "2026-01-28T09:00:00.000Z", UpdatedAt: "2026-06-07T09:00:00.000Z",
			},
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
