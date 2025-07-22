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

type UserHandler struct {
	userStore       postgres.UserStore
	submissionStore postgres.SubmissionStore
	problemStore    postgres.ProblemStore
	logger          *log.Logger
}

func NewUserHandler(userStore postgres.UserStore, submissionStore postgres.SubmissionStore, problemStore postgres.ProblemStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore:       userStore,
		submissionStore: submissionStore,
		problemStore:    problemStore,
		logger:          logger,
	}
}

type registerUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin" default:"false"`
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (h *UserHandler) registerUser(w http.ResponseWriter, r *http.Request, is_admin bool) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &postgres.User{
		Username: req.Username,
		IsAdmin:  is_admin,
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		if err == postgres.ErrDuplicateUsername {
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "username already exists"})
			return
		}

		h.logger.Printf("ERROR: registering user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}

// HandleRegisterNormalUser registers a new normal user
// @Summary      Register user
// @Description  Register a new normal user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      registerUserRequest  true  "User registration data"
// @Success      201   {object}  utils.Envelope{user=postgres.User}
// @Failure      400   {object}  utils.Envelope{error=string}
// @Failure      500   {object}  utils.Envelope{error=string}
// @Router       /users/register [post]
func (h *UserHandler) HandleRegisterNormalUser(w http.ResponseWriter, r *http.Request) {
	h.registerUser(w, r, false)
}

// HandleRegisterAdminUser registers a new admin user
// @Summary      Register admin user
// @Description  Register a new admin user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      registerUserRequest  true  "Admin user registration data"
// @Success      201   {object}  utils.Envelope{user=postgres.User}
// @Failure      400   {object}  utils.Envelope{error=string}
// @Failure      500   {object}  utils.Envelope{error=string}
// @Router       /users/register-admin [post]
func (h *UserHandler) HandleRegisterAdminUser(w http.ResponseWriter, r *http.Request) {
	h.registerUser(w, r, true)
}

// HandleGetCurrentUser gets current authenticated user
// @Summary      Get current user
// @Description  Get the current authenticated user's information
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Envelope{user=postgres.User}
// @Failure      401  {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /users/me [get]
func (h *UserHandler) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be authenticated to access this resource"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

type UserStats struct {
	ProblemsSolved int     `json:"problems_solved"`
	SuccessRate    float64 `json:"success_rate"`
}

func (h *UserHandler) HandleGetUserStats(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be authenticated to access this resource"})
		return
	}

	problemsSolved, err := h.submissionStore.GetUserSolvedProblemsCount(user.ID)
	if err != nil {
		h.logger.Printf("ERROR: GetUserSolvedProblemsCount: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	successRate, err := h.submissionStore.GetUserSuccessRate(user.ID)
	if err != nil {
		h.logger.Printf("ERROR: GetUserSuccessRate: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	stats := UserStats{
		ProblemsSolved: problemsSolved,
		SuccessRate:    successRate,
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"stats": stats})
}

// HandleGetAllUsers gets all users
// @Summary      Get all users
// @Description  Get list of all users (admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Envelope{users=[]postgres.User}
// @Failure      401  {object}  utils.Envelope{error=string}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userStore.GetAllUsers()
	if err != nil {
		h.logger.Printf("ERROR: GetAllUsers: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"users": users})
}

// HandleDeleteUser deletes a user by ID
// @Summary      Delete user
// @Description  Delete a user by ID (admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  "No Content"
// @Failure      400  {object}  utils.Envelope{error=string}
// @Failure      401  {object}  utils.Envelope{error=string}
// @Failure      404  {object}  utils.Envelope{error=string}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid user id"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to delete a user"})
		return
	}

	if currentUser.ID == userID {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you cannot delete your own account"})
		return
	}

	err = h.userStore.DeleteUser(userID)
	if err == sql.ErrNoRows {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
		return
	}
	if err != nil {
		h.logger.Printf("ERROR: DeleteUser: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type AdminStats struct {
	TotalProblems    int `json:"total_problems"`
	RegisteredUsers  int `json:"registered_users"`
	TotalSubmissions int `json:"total_submissions"`
	ActiveProblems   int `json:"active_problems"`
}

func (h *UserHandler) HandleGetAdminStats(w http.ResponseWriter, r *http.Request) {
	totalProblems, err := h.problemStore.GetTotalProblemsCount()
	if err != nil {
		h.logger.Printf("ERROR: GetTotalProblemsCount: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	registeredUsers, err := h.userStore.GetTotalUsersCount()
	if err != nil {
		h.logger.Printf("ERROR: GetTotalUsersCount: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	totalSubmissions, err := h.submissionStore.GetTotalSubmissionsCount()
	if err != nil {
		h.logger.Printf("ERROR: GetTotalSubmissionsCount: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	activeProblems, err := h.problemStore.GetActiveProblemsCount()
	if err != nil {
		h.logger.Printf("ERROR: GetActiveProblemsCount: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	stats := AdminStats{
		TotalProblems:    totalProblems,
		RegisteredUsers:  registeredUsers,
		TotalSubmissions: totalSubmissions,
		ActiveProblems:   activeProblems,
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"stats": stats})
}
