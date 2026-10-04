package user

import (
	"net/http"

	"github.com/akshaykrm/keystore/apps/api/internal/auth"
)

func RegisterRoutes(mux *http.ServeMux, controller *Controller) {
	mux.HandleFunc("GET /users", auth.IsAuthenticated(controller.GetAll))
	mux.HandleFunc("GET /users/{id}", controller.GetByID)
	mux.HandleFunc("PUT /users/{id}", controller.UpdateById)
	mux.HandleFunc("POST /users", controller.Create)
	mux.HandleFunc("DELETE /users/{id}", controller.DeleteById)

	mux.HandleFunc("POST /login", controller.Login)

}
