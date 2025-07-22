package api

import (
	"adel/internal/middleware"
	"adel/internal/service/postgres"
	"adel/internal/service/rabbitmq"
	"adel/internal/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type SubmissionHandler struct {
	submissionStore postgres.SubmissionStore
	problemStore    postgres.ProblemStore
	rabbitmqClient  *rabbitmq.RabbitMQClient
	logger          *log.Logger
}

func NewSubmissionHandler(submissionStore postgres.SubmissionStore, problemStore postgres.ProblemStore, rabbitmqClient *rabbitmq.RabbitMQClient, logger *log.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		submissionStore: submissionStore,
		problemStore:    problemStore,
		rabbitmqClient:  rabbitmqClient,
		logger:          logger,
	}
}

// HandleCreateSubmission creates a new submission
// @Summary      Create submission
// @Description  Create a new code submission for a problem
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        submission  body      postgres.Submission  true  "Submission object"
// @Success      201         {object}  utils.Envelope{submission=postgres.Submission}
// @Failure      400         {object}  utils.Envelope{error=string}
// @Failure      401         {object}  utils.Envelope{error=string}
// @Failure      500         {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /submissions [post]
func (sh *SubmissionHandler) HandleCreateSubmission(w http.ResponseWriter, r *http.Request) {
	var submission postgres.Submission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		sh.logger.Printf("ERROR: decodingCreateSubmission: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to create a submission"})
		return
	}

	submission.UserID = currentUser.ID

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

	problem, err := sh.problemStore.GetProblemByID(createdSubmission.ProblemID)
	if err != nil {
		sh.logger.Printf("ERROR: getProblemByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to get problem"})
		return
	}

	job := &rabbitmq.SubmissionJob{
		Submission: *createdSubmission,
		Problem:    *problem,
	}

	err = sh.rabbitmqClient.PublishSubmissionJob(job)
	if err != nil {
		sh.logger.Printf("ERROR: publishSubmissionJob: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to publish submission job"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"submission": createdSubmission})
}

// HandleGetSubmissionByID gets a submission by ID
// @Summary      Get submission by ID
// @Description  Get a single submission by its ID
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Submission ID"
// @Success      200  {object}  utils.Envelope{submission=postgres.Submission}
// @Failure      400  {object}  utils.Envelope{error=string}
// @Failure      401  {object}  utils.Envelope{error=string}
// @Failure      404  {object}  utils.Envelope{error=string}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /submissions/{id} [get]
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

// HandleDeleteSubmissionByID deletes a submission by ID
// @Summary      Delete submission
// @Description  Delete a submission by ID
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Submission ID"
// @Success      204  "No Content"
// @Failure      400  {object}  utils.Envelope{error=string}
// @Failure      401  {object}  utils.Envelope{error=string}
// @Failure      404  {object}  utils.Envelope{error=string}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /submissions/{id} [delete]
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

// HandleUpdateSubmissionByID updates a submission by ID
// @Summary      Update submission
// @Description  Update submission results (used by judge system)
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        id          path      int                                                                                     true   "Submission ID"
// @Param        submission  body      object{status=string,execution_time_ms=int64,memory_usage_mb=int64,error_message=string}  true  "Submission update object"
// @Success      200         {object}  utils.Envelope{submission=postgres.Submission}
// @Failure      400         {object}  utils.Envelope{error=string}
// @Failure      404         {object}  utils.Envelope{error=string}
// @Failure      500         {object}  utils.Envelope{error=string}
// @Router       /submissions/{id} [put]
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

// HandleGetSubmissionsByUserID gets submissions for current user
// @Summary      Get user submissions
// @Description  Get all submissions for the current authenticated user
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Envelope{submissions=[]postgres.Submission}
// @Failure      401  {object}  utils.Envelope{error=string}
// @Failure      404  {object}  utils.Envelope{error=string}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /submissions/user [get]
func (sh *SubmissionHandler) HandleGetSubmissionsByUserID(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to get submissions"})
		return
	}

	submissions, err := sh.submissionStore.GetSubmissionsByUserID(currentUser.ID)
	if submissions == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "submission not found"})
		return
	}
	if err != nil {
		sh.logger.Printf("ERROR: GetSubmissionByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"submissions": submissions})
}
