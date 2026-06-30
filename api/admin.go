package api

import (
	"context"

	"github.com/conelli/admin-backend/http/rs"
	"github.com/conelli/admin-backend/internal/api"
	"github.com/conelli/admin-backend/internal/api/request"
	apierrors "github.com/conelli/admin-backend/internal/errors"
	"github.com/conelli/admin-backend/logic"
	"github.com/gin-gonic/gin"
)

type Admin struct {
	svc *logic.CoreLogic
}

func NewAdmin(svc *logic.CoreLogic) *Admin {
	return &Admin{svc: svc}
}

func (a *Admin) Definition() api.Definition {
	return api.Definition{
		Name:        "Admin APIs",
		Description: "Admin service endpoints",
		BaseRoute: api.Route{
			Path: "/admin",
			Endpoints: api.Endpoints{
				Get: &api.Endpoint{
					Name:         "Admin Summary",
					Description:  "Get admin API summary",
					AuthRequired: false,
					Handler:      a.GetSummary,
					ResponseBody: &rs.AdminSummary{},
				},
			},
			SubRoutes: []api.Route{
				{
					Path: "/data",
					Endpoints: api.Endpoints{
						Get: &api.Endpoint{
							Name:         "Admin Data",
							Description:  "Get admin console data",
							AuthRequired: false,
							Handler:      a.GetData,
						},
						Put: &api.Endpoint{
							Name:         "Save Admin Data",
							Description:  "Persist admin console data",
							AuthRequired: false,
							RequestBody:  &logic.AdminData{},
							Handler:      a.SaveData,
						},
					},
				},
				{
					Path: "/login",
					Endpoints: api.Endpoints{
						Post: &api.Endpoint{
							Name:         "Staff Login",
							Description:  "Validate staff login credentials",
							AuthRequired: false,
							RequestBody:  &logic.LoginPayload{},
							ResponseBody: &logic.LoginResponse{},
							Handler:      a.Login,
						},
					},
				},
			},
		},
	}
}

func (a *Admin) GetSummary(ctx context.Context, req *request.Request) (interface{}, error) {
	return a.svc.Admin.Summary(ctx)
}

func (a *Admin) GetData(ctx context.Context, req *request.Request) (interface{}, error) {
	return a.svc.Admin.Data(ctx)
}

func (a *Admin) SaveData(ctx context.Context, req *request.Request) (interface{}, error) {
	payload, ok := req.Body.(*logic.AdminData)
	if !ok {
		return nil, apierrors.ErrBadRequest
	}
	return a.svc.Admin.SaveData(ctx, *payload)
}

func (a *Admin) Login(ctx context.Context, req *request.Request) (interface{}, error) {
	payload, ok := req.Body.(*logic.LoginPayload)
	if !ok {
		return nil, apierrors.ErrBadRequest
	}
	return a.svc.Admin.Login(ctx, *payload)
}

func (a *Admin) RegisterRoutes(engine *gin.Engine) {
	api.RegisterDefinition(engine, a.Definition())
}
