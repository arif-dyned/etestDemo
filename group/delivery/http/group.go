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

// GroupsHandler represent HTTP handler for Groups
type GroupsHandler struct {
	GroupsUsecase domain.GroupsUsecase
}

// NewGroupsHandler returns HTTP delivery instance for Groups
func NewGroupsHandler(r *mux.Router, u domain.GroupsUsecase) {
	h := &GroupsHandler{
		GroupsUsecase: u,
	}

	r.HandleFunc("/banks/groups", h.Index).Methods("GET").Host(config.HostServerName)
	r.HandleFunc("/banks/groups", h.Create).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/groups", h.Update).Methods("PUT").Host(config.HostServerName)
	r.HandleFunc("/banks/groups/{id:[0-9]+}", h.Delete).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/groups/{id:[0-9]+}", h.GetById).Methods("GET").Host(config.HostServerName)
	// TODO: add more routes
}

// Index returns all Groupss into JSON
// TODO: handle when Groups repository is empty
func (h *GroupsHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Groupss, err := h.GroupsUsecase.Fetch(r.Context(), 10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Groupss)
}

// Create handles Groups creation
func (h *GroupsHandler) Create(w http.ResponseWriter, r *http.Request) {
	Groups := new(domain.Groups)
	errs := binding.Bind(r, Groups)
	if errs.Handle(w) {
		return
	}

	//validate data
	validateData := Groups.Validate(r, errs)
	if validateData.Error() != "" {
		domain.RespondWithError(w, http.StatusInternalServerError, validateData.Error())
		return
	}

	//send to useCase
	res, err := h.GroupsUsecase.Store(*Groups)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Groups.ID = res
	domain.RespondwithJSON(w, http.StatusOK, Groups)
}

// Update handles Existing Groups
func (h *GroupsHandler) Update(w http.ResponseWriter, r *http.Request) {
	Groups := new(domain.Groups)
	errs := binding.Bind(r, Groups)
	if errs.Handle(w) {
		return
	}

	//send to useCase
	_, err := h.GroupsUsecase.Update(*Groups)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, Groups)
	return
}

// Delete handles a remove groups's
func (h *GroupsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID groups
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	err := h.GroupsUsecase.Delete(id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, map[string]string{"message": "Delete Successfully"})
	return
}

// GetById handles a get groups by id
func (h *GroupsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID groups
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	payload, err := h.GroupsUsecase.GetByID(r.Context(), id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, payload)
	return
}

// TODO: implement more handler
