package products

import (
	"net/http"

	"github.com/akshaykrm/keystore/apps/api/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, controller *Controller) {
	mux.HandleFunc("POST /products", auth.IsAuthenticated(controller.Create))
	mux.HandleFunc("GET /products", auth.IsAuthenticated(controller.GetAll))
	mux.HandleFunc("GET /products/{id}", auth.IsAuthenticated(controller.GetOneByID))
	mux.HandleFunc("PUT /products/{id}", auth.IsAuthenticated(controller.UpdateByID))
	// mux.HandleFunc("DELETE /workspaces/{id}", auth.IsAuthenticated(controller.DeleteById))
}
