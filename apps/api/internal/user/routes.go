package user

import (
	"net/http"

	"github.com/akshaykrm/keystore/apps/api/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, controller *Controller) {
	mux.HandleFunc("POST /users", controller.Create)
	mux.HandleFunc("POST /login", controller.Login)

	mux.HandleFunc("GET /users", auth.IsAuthenticated(controller.GetAll))
	mux.HandleFunc("GET /users/{id}", auth.IsAuthenticated(controller.GetByID))
	mux.HandleFunc("PUT /users/{id}", auth.IsAuthenticated(controller.UpdateById))
	mux.HandleFunc("DELETE /users/{id}", auth.IsAuthenticated(controller.DeleteById))

}
