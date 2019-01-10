package http

import (
	"encoding/json"
	"github.com/DynEd/etest/config"
	"github.com/DynEd/etest/domain"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"net/http"
	"strconv"
)

// TypesHandler represent HTTP handler for Types
type TypesHandler struct {
	TypesUsecase domain.TypesUsecase
}

// NewTypesHandler returns HTTP delivery instance for Types
func NewTypesHandler(r *mux.Router, u domain.TypesUsecase) {
	h := &TypesHandler{
		TypesUsecase: u,
	}

	r.HandleFunc("/banks/types", h.Index).Methods("GET").Host(config.HostServerName)
	r.HandleFunc("/banks/types", h.Create).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/types", h.Update).Methods("PUT").Host(config.HostServerName)
	r.HandleFunc("/banks/types/{id:[0-9]+}", h.Delete).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/types/{id:[0-9]+}", h.GetById).Methods("GET").Host(config.HostServerName)
	// TODO: add more routes
}

// Index returns all Typess into JSON
// TODO: handle when Types repository is empty
func (h *TypesHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Typess, err := h.TypesUsecase.Fetch(r.Context(), 10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Typess)
}

// Create handles Types creation
func (h *TypesHandler) Create(w http.ResponseWriter, r *http.Request) {
	Types := new(domain.Types)
	errs := binding.Bind(r, Types)
	if errs.Handle(w) {
		return
	}

	//validate data
	validateData := Types.Validate(r, errs)
	if validateData.Error() != "" {
		domain.RespondWithError(w, http.StatusInternalServerError, validateData.Error())
		return
	}

	//send to useCase
	res, err := h.TypesUsecase.Store(*Types)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	Types.ID = res
	domain.RespondwithJSON(w, http.StatusOK, Types)
}

// Update handles Existing Types
func (h *TypesHandler) Update(w http.ResponseWriter, r *http.Request) {
	Types := new(domain.Types)
	errs := binding.Bind(r, Types)
	if errs.Handle(w) {
		return
	}

	//send to useCase
	_, err := h.TypesUsecase.Update(*Types)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, Types)
	return
}

// Delete handles a remove types's
func (h *TypesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID types
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	err := h.TypesUsecase.Delete(id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, map[string]string{"message": "Delete Successfully"})
	return
}

// GetById handles a get types by id
func (h *TypesHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID types
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	payload, err := h.TypesUsecase.GetByID(r.Context(), id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, payload)
	return
}

// TODO: implement more handler
