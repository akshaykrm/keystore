package workspace

import (
	"net/http"

	"github.com/akshaykrm/keystore/apps/api/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, controller *Controller) {
	mux.HandleFunc("GET /workspaces", auth.IsAuthenticated(controller.GetAll))
	mux.HandleFunc("GET /workspaces/{id}", auth.IsAuthenticated(controller.GetByID))
	mux.HandleFunc("PUT /workspaces/{id}", auth.IsAuthenticated(controller.UpdateByID))
	mux.HandleFunc("POST /workspaces", auth.IsAuthenticated(controller.Create))
	mux.HandleFunc("DELETE /workspaces/{id}", auth.IsAuthenticated(controller.DeleteById))
}
