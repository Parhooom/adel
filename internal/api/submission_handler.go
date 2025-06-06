package api

import (
	"adel/internal/service/postgres"
	"adel/internal/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type SubmissionHandler struct {
	submissionStore postgres.SubmissionStore
	logger          *log.Logger
}

func NewSubmissionHandler(submissionStore postgres.SubmissionStore, logger *log.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		submissionStore: submissionStore,
		logger:          logger,
	}
}

func (sh *SubmissionHandler) HandleCreateSubmission(w http.ResponseWriter, r *http.Request) {
	var submission postgres.Submission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		sh.logger.Printf("ERROR: decodingCreateSubmission: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	submission.Status = "pending"
	submission.ExecutionTimeMs = 0
	submission.MemoryUsageMB = 0
	submission.ErrorMessage = ""

	createdSubmission, err := sh.submissionStore.CreateSubmission(&submission)
	if err != nil {
		sh.logger.Printf("ERROR: createSubmission: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create submission"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"submission": createdSubmission})
}

func (sh *SubmissionHandler) HandleGetSubmissionByID(w http.ResponseWriter, r *http.Request) {
	submissionID, err := utils.ReadIDParam(r)
	if err != nil {
		sh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	submission, err := sh.submissionStore.GetSubmissionByID(submissionID)
	if submission == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "submission not found"})
		return
	}
	if err != nil {
		sh.logger.Printf("ERROR: GetSubmissionByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"submission": submission})
}

func (sh *SubmissionHandler) HandleDeleteSubmissionByID(w http.ResponseWriter, r *http.Request) {
	submissionID, err := utils.ReadIDParam(r)
	if err != nil {
		sh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	err = sh.submissionStore.DeleteSubmission(submissionID)
	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "submission not found"})
		return
	}
	if err != nil {
		sh.logger.Printf("ERROR: deleteSubmissionByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (sh *SubmissionHandler) HandleUpdateSubmissionByID(w http.ResponseWriter, r *http.Request) {
	submissionID, err := utils.ReadIDParam(r)
	if err != nil {
		sh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	existingSubmission, err := sh.submissionStore.GetSubmissionByID(submissionID)
	if err != nil {
		sh.logger.Printf("ERROR: GetSubmissionByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if existingSubmission == nil {
		http.NotFound(w, r)
		return
	}

	var updateSubmissionRequest struct {
		Status          *string `json:"status"`
		ExecutionTimeMs *int64  `json:"execution_time_ms"`
		MemoryUsageMB   *int64  `json:"memory_usage_mb"`
		ErrorMessage    *string `json:"error_message"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateSubmissionRequest)
	if err != nil {
		sh.logger.Printf("ERROR: decodingUpdateSubmission: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	if updateSubmissionRequest.Status != nil {
		existingSubmission.Status = *updateSubmissionRequest.Status
	}
	if updateSubmissionRequest.ExecutionTimeMs != nil {
		existingSubmission.ExecutionTimeMs = *updateSubmissionRequest.ExecutionTimeMs
	}
	if updateSubmissionRequest.MemoryUsageMB != nil {
		existingSubmission.MemoryUsageMB = *updateSubmissionRequest.MemoryUsageMB
	}
	if updateSubmissionRequest.ErrorMessage != nil {
		existingSubmission.ErrorMessage = *updateSubmissionRequest.ErrorMessage
	}

	err = sh.submissionStore.UpdateSubmission(existingSubmission)
	if err != nil {
		sh.logger.Printf("ERROR: updateSubmissionByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"submission": existingSubmission})
}
