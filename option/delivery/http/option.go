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

// OptionHandler represent HTTP handler for Option
type OptionHandler struct {
	OptionUsecase domain.OptionUsecase
}

// NewOptionHandler returns HTTP delivery instance for Option
func NewOptionHandler(r *mux.Router, u domain.OptionUsecase) {
	h := &OptionHandler{
		OptionUsecase: u,
	}

	r.HandleFunc("/banks/options", h.Index).Methods("GET").Host(config.HostServerName)
	r.HandleFunc("/banks/options", h.Create).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/options", h.Update).Methods("PUT").Host(config.HostServerName)
	r.HandleFunc("/banks/options/{id:[0-9]+}", h.Delete).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/options/{id:[0-9]+}", h.GetById).Methods("GET").Host(config.HostServerName)
	// TODO: add more routes
}

// Index returns all Options into JSON
// TODO: handle when Option repository is empty
func (h *OptionHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Options, err := h.OptionUsecase.Fetch(r.Context(), 10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Options)
}

// Create handles Option creation
func (h *OptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	Option := new(domain.Option)
	errs := binding.Bind(r, Option)
	if errs.Handle(w) {
		return
	}

	//validate data
	validateData := Option.Validate(r, errs)
	if validateData.Error() != "" {
		domain.RespondWithError(w, http.StatusInternalServerError, validateData.Error())
		return
	}

	//send to useCase
	res, err := h.OptionUsecase.Store(*Option)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	Option.ID = res
	domain.RespondwithJSON(w, http.StatusOK, Option)
}

// Update handles Existing Option
func (h *OptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	Option := new(domain.Option)
	errs := binding.Bind(r, Option)
	if errs.Handle(w) {
		return
	}

	//send to useCase
	_, err := h.OptionUsecase.Update(*Option)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, Option)
	return
}

// Delete handles a remove option's
func (h *OptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID option
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	err := h.OptionUsecase.Delete(id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, map[string]string{"message": "Delete Successfully"})
	return
}

// GetById handles a get option by id
func (h *OptionHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID option
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	payload, err := h.OptionUsecase.GetByID(r.Context(), id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, payload)
	return
}

// TODO: implement more handler
