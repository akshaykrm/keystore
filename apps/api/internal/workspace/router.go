package workspace

import "net/http"

func RegisterRoutes(mux *http.ServeMux, controller *Controller) {
	mux.HandleFunc("GET /workspaces", controller.GetAll)
	mux.HandleFunc("GET /workspaces/{id}", controller.GetByID)
	mux.HandleFunc("PUT /workspaces/{id}", controller.UpdateByID)
	mux.HandleFunc("POST /workspaces", controller.Create)
	mux.HandleFunc("DELETE /workspaces/{id}", controller.DeleteById)
}
