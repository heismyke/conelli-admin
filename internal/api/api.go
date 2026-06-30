package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/conelli/admin-backend/internal/api/request"
	"github.com/conelli/admin-backend/internal/api/response"
	appErrors "github.com/conelli/admin-backend/internal/errors"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx context.Context, req *request.Request) (interface{}, error)

type Definition struct {
	Name        string
	Description string
	BaseRoute   Route
}

type Route struct {
	Path      string
	Endpoints Endpoints
	SubRoutes []Route
}

type Endpoints struct {
	Get    *Endpoint
	Post   *Endpoint
	Put    *Endpoint
	Patch  *Endpoint
	Delete *Endpoint
}

type Endpoint struct {
	Name         string
	Description  string
	AuthRequired bool
	RequestBody  interface{}
	ResponseBody interface{}
	Handler      HandlerFunc
}

func RegisterDefinition(engine *gin.Engine, definition Definition) {
	registerRoute(engine.Group(definition.BaseRoute.Path), definition.BaseRoute)
}

func registerRoute(group *gin.RouterGroup, route Route) {
	registerEndpoint(group, http.MethodGet, route.Endpoints.Get)
	registerEndpoint(group, http.MethodPost, route.Endpoints.Post)
	registerEndpoint(group, http.MethodPut, route.Endpoints.Put)
	registerEndpoint(group, http.MethodPatch, route.Endpoints.Patch)
	registerEndpoint(group, http.MethodDelete, route.Endpoints.Delete)

	for _, subRoute := range route.SubRoutes {
		registerRoute(group.Group(subRoute.Path), subRoute)
	}
}

func registerEndpoint(group *gin.RouterGroup, method string, endpoint *Endpoint) {
	if endpoint == nil {
		return
	}

	handler := func(c *gin.Context) {
		req := &request.Request{Context: c}
		if endpoint.RequestBody != nil {
			req.Body = endpoint.RequestBody
			if err := c.ShouldBindJSON(req.Body); err != nil {
				c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
				return
			}
		}

		result, err := endpoint.Handler(c.Request.Context(), req)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, appErrors.ErrBadRequest) {
				status = http.StatusBadRequest
			}
			if errors.Is(err, appErrors.ErrUnauthorized) {
				status = http.StatusUnauthorized
			}
			if errors.Is(err, appErrors.ErrForbidden) {
				status = http.StatusForbidden
			}
			if errors.Is(err, appErrors.ErrNotFound) {
				status = http.StatusNotFound
			}
			c.JSON(status, response.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}

	switch method {
	case http.MethodGet:
		group.GET("", handler)
	case http.MethodPost:
		group.POST("", handler)
	case http.MethodPut:
		group.PUT("", handler)
	case http.MethodPatch:
		group.PATCH("", handler)
	case http.MethodDelete:
		group.DELETE("", handler)
	}
}
