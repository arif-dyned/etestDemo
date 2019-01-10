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

// LevelHandler represent HTTP handler for Level
type LevelHandler struct {
	LevelUsecase domain.LevelUsecase
}

// NewLevelHandler returns HTTP delivery instance for Level
func NewLevelHandler(r *mux.Router, u domain.LevelUsecase) {
	h := &LevelHandler{
		LevelUsecase: u,
	}

	r.HandleFunc("/banks/levels", h.Index).Methods("GET").Host(config.HostServerName)
	r.HandleFunc("/banks/levels", h.Create).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/levels", h.Update).Methods("PUT").Host(config.HostServerName)
	r.HandleFunc("/banks/levels/{id:[0-9]+}", h.Delete).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/levels/{id:[0-9]+}", h.GetById).Methods("GET").Host(config.HostServerName)
	// TODO: add more routes
}

// Index returns all Levels into JSON
// TODO: handle when Level repository is empty
func (h *LevelHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Levels, err := h.LevelUsecase.Fetch(r.Context(), 10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Levels)
}

// Create handles Level creation
func (h *LevelHandler) Create(w http.ResponseWriter, r *http.Request) {
	Level := new(domain.Level)
	errs := binding.Bind(r, Level)
	if errs.Handle(w) {
		return
	}

	//validate data
	validateData := Level.Validate(r, errs)
	if validateData.Error() != "" {
		domain.RespondWithError(w, http.StatusInternalServerError, validateData.Error())
		return
	}

	//send to useCase
	res, err := h.LevelUsecase.Store(*Level)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	Level.ID = res
	domain.RespondwithJSON(w, http.StatusOK, Level)
}

// Update handles Existing Level
func (h *LevelHandler) Update(w http.ResponseWriter, r *http.Request) {
	Level := new(domain.Level)
	errs := binding.Bind(r, Level)
	if errs.Handle(w) {
		return
	}

	//send to useCase
	_, err := h.LevelUsecase.Update(*Level)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, Level)
	return
}

// Delete handles a remove level's
func (h *LevelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID level
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	err := h.LevelUsecase.Delete(id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, map[string]string{"message": "Delete Successfully"})
	return
}

// GetById handles a get level by id
func (h *LevelHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID level
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	payload, err := h.LevelUsecase.GetByID(r.Context(), id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, payload)
	return
}

// TODO: implement more handler
