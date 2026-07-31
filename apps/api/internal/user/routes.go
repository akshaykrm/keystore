package user

import "net/http"

func RegisterRoutes(mux *http.ServeMux, controller *Controller) {
	mux.HandleFunc("GET /users", controller.GetAll)
	mux.HandleFunc("GET /users/{id}", controller.GetByID)
	mux.HandleFunc("PUT /users/{id}", controller.UpdateById)
	mux.HandleFunc("POST /users", controller.Create)
	mux.HandleFunc("DELETE /users/{id}", controller.DeleteById)
}
