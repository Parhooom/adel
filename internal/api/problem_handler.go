package api

import (
	"adel/internal/middleware"
	"adel/internal/service/postgres"
	"adel/internal/utils"
	"database/sql"
	"encoding/json"
	"errors"
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

// HandleGetProblemByID gets a problem by ID
// @Summary      Get problem by ID
// @Description  Get a single problem by its ID
// @Tags         problems
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Problem ID"
// @Success      200  {object}  utils.Envelope{problem=postgres.Problem}
// @Failure      400  {object}  utils.Envelope{error=string}
// @Failure      404  {object}  utils.Envelope{error=string}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Router       /problems/{id} [get]
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

// HandleCreateProblem creates a new problem
// @Summary      Create problem
// @Description  Create a new problem (admin only)
// @Tags         problems
// @Accept       json
// @Produce      json
// @Param        problem  body      postgres.Problem  true  "Problem object"
// @Success      201      {object}  utils.Envelope{problem=postgres.Problem}
// @Failure      400      {object}  utils.Envelope{error=string}
// @Failure      401      {object}  utils.Envelope{error=string}
// @Failure      500      {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /problems [post]
func (ph *ProblemHandler) HandleCreateProblem(w http.ResponseWriter, r *http.Request) {
	var problem postgres.Problem
	err := json.NewDecoder(r.Body).Decode(&problem)
	if err != nil {
		ph.logger.Printf("ERROR: decodingCreateProblem: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to create a problem"})
		return
	}

	problem.UserID = currentUser.ID

	createdProblem, err := ph.problemStore.CreateProblem(&problem)
	if err != nil {
		ph.logger.Printf("ERROR: createProblem: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create problem"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"problem": createdProblem})
}

// HandleDeleteProblemByID deletes a problem by ID
// @Summary      Delete problem
// @Description  Delete a problem by ID (admin only)
// @Tags         problems
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Problem ID"
// @Success      204  "No Content"
// @Failure      400  {object}  utils.Envelope{error=string}
// @Failure      401  {object}  utils.Envelope{error=string}
// @Failure      403  {object}  utils.Envelope{error=string}
// @Failure      404  {object}  utils.Envelope{error=string}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /problems/{id} [delete]
func (ph *ProblemHandler) HandleDeleteProblemByID(w http.ResponseWriter, r *http.Request) {
	problemID, err := utils.ReadIDParam(r)
	if err != nil {
		ph.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid problem id"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to delete a problem"})
		return
	}

	problemOwner, err := ph.problemStore.GetProblemOwner(problemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "problem does not exist"})
			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if problemOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not allowed to delete this problem"})
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

// HandleUpdateProblemByID updates a problem by ID
// @Summary      Update problem
// @Description  Update a problem by ID (admin only)
// @Tags         problems
// @Accept       json
// @Produce      json
// @Param        id       path      int                              true   "Problem ID"
// @Param        problem  body      object{title=string,description=string,difficulty=string,time_limit_ms=int,memory_limit_mb=int,is_active=bool,testcases=[]postgres.TestCase}  true  "Problem update object"
// @Success      200      {object}  utils.Envelope{problem=postgres.Problem}
// @Failure      400      {object}  utils.Envelope{error=string}
// @Failure      401      {object}  utils.Envelope{error=string}
// @Failure      403      {object}  utils.Envelope{error=string}
// @Failure      404      {object}  utils.Envelope{error=string}
// @Failure      500      {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /problems/{id} [put]
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
		TestCases   []postgres.TestCase `json:"testcases"`
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

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to create a problem"})
		return
	}

	if currentUser.ID != existingProblem.UserID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not allowed to update this problem"})
		return
	}

	err = ph.problemStore.UpdateProblem(existingProblem)
	if err != nil {
		ph.logger.Printf("ERROR: updatingProblem: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"problem": existingProblem})
}

// HandleGetAllProblems gets all active problems
// @Summary      Get all problems
// @Description  Get all active problems
// @Tags         problems
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Envelope{problems=[]postgres.Problem}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Router       /problems [get]
func (ph *ProblemHandler) HandleGetAllProblems(w http.ResponseWriter, r *http.Request) {
	problems, err := ph.problemStore.GetAllProblems()
	if err != nil {
		ph.logger.Printf("ERROR: GetAllProblems: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"problems": problems})
}

// HandleGetAllProblemsForAdmin gets all problems for admin
// @Summary      Get all problems for admin
// @Description  Get all problems including inactive ones (admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Envelope{problems=[]postgres.Problem}
// @Failure      401  {object}  utils.Envelope{error=string}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /admin/problems [get]
func (ph *ProblemHandler) HandleGetAllProblemsForAdmin(w http.ResponseWriter, r *http.Request) {
	problems, err := ph.problemStore.GetAllProblemsForAdmin()
	if err != nil {
		ph.logger.Printf("ERROR: GetAllProblemsForAdmin: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"problems": problems})
}
