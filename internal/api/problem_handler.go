package api

import (
	"adel/internal/service/postgres"
	"adel/internal/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type ProblemHandler struct {
	problemStore postgres.ProblemStore
	logger       *log.Logger
}

func NewProblemHandler(problemStore postgres.ProblemStore, logger *log.Logger) *ProblemHandler {
	return &ProblemHandler{
		problemStore: problemStore,
		logger:       logger,
	}
}

func (ph *ProblemHandler) HandleGetProblemByID(w http.ResponseWriter, r *http.Request) {
	problemID, err := utils.ReadIDParam(r)
	if err != nil {
		ph.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid problem id"})
		return
	}

	problem, err := ph.problemStore.GetProblemByID(problemID)
	if problem == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "problem not found"})
		return
	}
	if err != nil {
		ph.logger.Printf("ERROR: GetProblemByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"problem": problem})
}

func (ph *ProblemHandler) HandleCreateProblem(w http.ResponseWriter, r *http.Request) {
	var problem postgres.Problem
	err := json.NewDecoder(r.Body).Decode(&problem)
	if err != nil {
		ph.logger.Printf("ERROR: decodingCreateProblem: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	createdProblem, err := ph.problemStore.CreateProblem(&problem)
	if err != nil {
		ph.logger.Printf("ERROR: createProblem: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create problem"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"problem": createdProblem})
}

func (ph *ProblemHandler) HandleDeleteProblemByID(w http.ResponseWriter, r *http.Request) {
	problemID, err := utils.ReadIDParam(r)
	if err != nil {
		ph.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid problem id"})
		return
	}

	err = ph.problemStore.DeleteProblem(problemID)
	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "problem not found"})
		return
	}
	if err != nil {
		ph.logger.Printf("ERROR: deleteProblemByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ph *ProblemHandler) HandleUpdateProblemByID(w http.ResponseWriter, r *http.Request) {
	problemID, err := utils.ReadIDParam(r)
	if err != nil {
		ph.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid problem id"})
		return
	}

	existingProblem, err := ph.problemStore.GetProblemByID(problemID)
	if err != nil {
		ph.logger.Printf("ERROR: GetProblemByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if existingProblem == nil {
		http.NotFound(w, r)
		return
	}

	var updateProblemRequest struct {
		Title       *string             `json:"title"`
		Description *string             `json:"description"`
		Difficulty  *string             `json:"difficulty"`
		TimeLimit   *int                `json:"time_limit_ms"`
		MemoryLimit *int                `json:"memory_limit_mb"`
		IsActive    *bool               `json:"is_active"`
		TestCases   []postgres.TestCase `json:"test_cases"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateProblemRequest)
	if err != nil {
		ph.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	if updateProblemRequest.Title != nil {
		existingProblem.Title = *updateProblemRequest.Title
	}
	if updateProblemRequest.Description != nil {
		existingProblem.Description = *updateProblemRequest.Description
	}
	if updateProblemRequest.Difficulty != nil {
		existingProblem.Difficulty = *updateProblemRequest.Difficulty
	}
	if updateProblemRequest.TimeLimit != nil {
		existingProblem.TimeLimit = *updateProblemRequest.TimeLimit
	}
	if updateProblemRequest.MemoryLimit != nil {
		existingProblem.MemoryLimit = *updateProblemRequest.MemoryLimit
	}
	if updateProblemRequest.IsActive != nil {
		existingProblem.IsActive = *updateProblemRequest.IsActive
	}
	if updateProblemRequest.TestCases != nil {
		existingProblem.TestCases = updateProblemRequest.TestCases
	}

	err = ph.problemStore.UpdateProblem(existingProblem)
	if err != nil {
		ph.logger.Printf("ERROR: updatingProblem: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"problem": existingProblem})
}
