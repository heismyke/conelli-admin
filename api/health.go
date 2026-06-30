package api

import (
	"context"

	"github.com/conelli/admin-backend/http/rs"
	"github.com/conelli/admin-backend/internal/api"
	"github.com/conelli/admin-backend/internal/api/request"
	"github.com/conelli/admin-backend/logic"
	"github.com/gin-gonic/gin"
)

type Health struct {
	svc *logic.CoreLogic
}

func NewHealth(svc *logic.CoreLogic) *Health {
	return &Health{svc: svc}
}

func (a *Health) Definition() api.Definition {
	return api.Definition{
		Name:        "Health APIs",
		Description: "Health check and status endpoints",
		BaseRoute: api.Route{
			Path: "/health",
			Endpoints: api.Endpoints{
				Get: &api.Endpoint{
					Name:         "Health Check",
					Description:  "Check API health status",
					AuthRequired: false,
					Handler:      a.GetHealthStatus,
					ResponseBody: &rs.Status{},
				},
			},
			SubRoutes: []api.Route{
				{
					Path: "/ping",
					Endpoints: api.Endpoints{
						Get: &api.Endpoint{
							Name:         "Ping",
							Description:  "Simple ping endpoint",
							AuthRequired: false,
							Handler:      a.Ping,
							ResponseBody: map[string]string{},
						},
					},
				},
			},
		},
	}
}

func (a *Health) GetHealthStatus(ctx context.Context, req *request.Request) (interface{}, error) {
	return &rs.Status{
		Status:  "healthy",
		Version: "1.0.0",
		Service: "Conelli Admin API",
	}, nil
}

func (a *Health) Ping(ctx context.Context, req *request.Request) (interface{}, error) {
	return map[string]string{"message": "pong"}, nil
}

func (a *Health) RegisterRoutes(engine *gin.Engine) {
	api.RegisterDefinition(engine, a.Definition())
}
