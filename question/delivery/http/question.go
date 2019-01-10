package http

import (
	"encoding/json"
	"github.com/DynEd/etest/domain"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"net/http"
	"strconv"
)

// QuestionHandler represent HTTP handler for Question
type QuestionHandler struct {
	QuestionUsecase domain.QuestionUsecase
}

// NewQuestionHandler returns HTTP delivery instance for Question
func NewQuestionHandler(r *mux.Router, u domain.QuestionUsecase) {
	h := &QuestionHandler{
		QuestionUsecase: u,
	}

	r.HandleFunc("/banks/question", h.Index).Methods("GET")
	r.HandleFunc("/banks/question", h.Create).Methods("POST")
	r.HandleFunc("/banks/question", h.Update).Methods("PUT")
	r.HandleFunc("/banks/question/{id:[0-9]+}", h.Delete).Methods("POST")
	r.HandleFunc("/banks/question/{id:[0-9]+}", h.GetById).Methods("GET")
	// TODO: add more routes
}

// Index returns all Questions into JSON
// TODO: handle when Question repository is empty
func (h *QuestionHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Questions, err := h.QuestionUsecase.Fetch(r.Context(), 10)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Questions)
}

// Create handles Question creation
func (h *QuestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	Question := new(domain.Question)
	errs := binding.Bind(r, Question)
	if errs.Handle(w) {
		return
	}

	//validate data
	validateData := Question.Validate(r, errs)
	if validateData.Error() != "" {
		domain.RespondWithError(w, http.StatusInternalServerError, validateData.Error())
		return
	}

	//send to useCase
	res, err := h.QuestionUsecase.Store(*Question)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	Question.ID = res
	domain.RespondwithJSON(w, http.StatusOK, Question)
}

// Update handles Existing Question
func (h *QuestionHandler) Update(w http.ResponseWriter, r *http.Request) {
	Question := new(domain.Question)
	errs := binding.Bind(r, Question)
	if errs.Handle(w) {
		return
	}

	//send to useCase
	_, err := h.QuestionUsecase.Update(*Question)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, Question)
	return
}

// Delete handles a remove question's
func (h *QuestionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID question
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	err := h.QuestionUsecase.Delete(id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, map[string]string{"message": "Delete Successfully"})
	return
}

// GetById handles a get question by id
func (h *QuestionHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID question
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	payload, err := h.QuestionUsecase.GetByID(r.Context(), id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, payload)
	return
}

// TODO: implement more handler
