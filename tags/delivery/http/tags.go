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

// TagsHandler represent HTTP handler for Tags
type TagsHandler struct {
	TagsUsecase domain.TagsUsecase
}

// NewTagsHandler returns HTTP delivery instance for Tags
func NewTagsHandler(r *mux.Router, u domain.TagsUsecase) {
	h := &TagsHandler{
		TagsUsecase: u,
	}

	r.HandleFunc("/banks/tags", h.Index).Methods("GET").Host(config.HostServerName)
	r.HandleFunc("/banks/tags", h.Create).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/tags", h.Update).Methods("PUT").Host(config.HostServerName)
	r.HandleFunc("/banks/tags/{id:[0-9]+}", h.Delete).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/tags/{id:[0-9]+}", h.GetById).Methods("GET").Host(config.HostServerName)
	// TODO: add more routes
}

// Index returns all Tagss into JSON
// TODO: handle when Tags repository is empty
func (h *TagsHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Tagss, err := h.TagsUsecase.Fetch(r.Context(), 10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Tagss)
}

// Create handles Tags creation
func (h *TagsHandler) Create(w http.ResponseWriter, r *http.Request) {
	Tags := new(domain.Tags)
	errs := binding.Bind(r, Tags)
	if errs.Handle(w) {
		return
	}

	//validate data
	validateData := Tags.Validate(r, errs)
	if validateData.Error() != "" {
		domain.RespondWithError(w, http.StatusInternalServerError, validateData.Error())
		return
	}

	//send to useCase
	res, err := h.TagsUsecase.Store(*Tags)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	Tags.ID = res
	domain.RespondwithJSON(w, http.StatusOK, Tags)
}

// Update handles Existing Tags
func (h *TagsHandler) Update(w http.ResponseWriter, r *http.Request) {
	Tags := new(domain.Tags)
	errs := binding.Bind(r, Tags)
	if errs.Handle(w) {
		return
	}

	//send to useCase
	_, err := h.TagsUsecase.Update(*Tags)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, Tags)
	return
}

// Delete handles a remove tags's
func (h *TagsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID tags
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	err := h.TagsUsecase.Delete(id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, map[string]string{"message": "Delete Successfully"})
	return
}

// GetById handles a get tags by id
func (h *TagsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID tags
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	payload, err := h.TagsUsecase.GetByID(r.Context(), id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, payload)
	return
}

// TODO: implement more handler
