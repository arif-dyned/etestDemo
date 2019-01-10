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

// QuestionSlotHandler represent HTTP handler for QuestionSlot
type QuestionSlotHandler struct {
	QuestionSlotUsecase domain.QuestionSlotUsecase
}

// NewQuestionSlotHandler returns HTTP delivery instance for QuestionSlot
func NewQuestionSlotHandler(r *mux.Router, u domain.QuestionSlotUsecase) {
	h := &QuestionSlotHandler{
		QuestionSlotUsecase: u,
	}

	r.HandleFunc("/banks/question/{id:[0-9]+}/slot", h.Index).Methods("GET").Host(config.HostServerName)
	r.HandleFunc("/banks/question/slot", h.Create).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/question/slot", h.Update).Methods("PUT").Host(config.HostServerName)
	r.HandleFunc("/banks/question/slot/{id:[0-9]+}", h.Delete).Methods("POST").Host(config.HostServerName)
	r.HandleFunc("/banks/question/{id:[0-9]+}/slot/{slotId:[0-9]+}", h.GetById).Methods("GET").Host(config.HostServerName)
	// TODO: add more routes
}

// Index returns all QuestionSlots into JSON
// TODO: handle when QuestionSlot repository is empty
func (h *QuestionSlotHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	//get ID question
	id, _ := strconv.Atoi(vars["id"])
	QuestionSlots, err := h.QuestionSlotUsecase.Fetch(r.Context(), int64(id))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(QuestionSlots)
}

// Create handles QuestionSlot creation
func (h *QuestionSlotHandler) Create(w http.ResponseWriter, r *http.Request) {
	QuestionSlot := new(domain.QuestionSlot)
	errs := binding.Bind(r, QuestionSlot)
	if errs.Handle(w) {
		return
	}

	//validate data
	validateData := QuestionSlot.Validate(r, errs)
	if validateData.Error() != "" {
		domain.RespondWithError(w, http.StatusInternalServerError, validateData.Error())
		return
	}

	//send to useCase
	res, err := h.QuestionSlotUsecase.Store(*QuestionSlot)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	QuestionSlot.ID = res
	domain.RespondwithJSON(w, http.StatusOK, QuestionSlot)
}

// Update handles Existing QuestionSlot
func (h *QuestionSlotHandler) Update(w http.ResponseWriter, r *http.Request) {
	QuestionSlot := new(domain.QuestionSlot)
	errs := binding.Bind(r, QuestionSlot)
	if errs.Handle(w) {
		return
	}

	//send to useCase
	_, err := h.QuestionSlotUsecase.Update(*QuestionSlot)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, QuestionSlot)
	return
}

// Delete handles a remove questionSlot's
func (h *QuestionSlotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID questionSlot
	id, _ := strconv.Atoi(vars["id"])

	//send to useCase
	err := h.QuestionSlotUsecase.Delete(id)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, map[string]string{"message": "Delete Successfully"})
	return
}

// GetById handles a get questionSlot by id
func (h *QuestionSlotHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//get ID questionSlot
	id, _ := strconv.Atoi(vars["id"])
	slotId, _ := strconv.Atoi(vars["slotId"])

	//send to useCase
	payload, err := h.QuestionSlotUsecase.GetByID(r.Context(), id, slotId)
	if err != nil {
		domain.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	domain.RespondwithJSON(w, http.StatusOK, payload)
	return
}

// TODO: implement more handler
