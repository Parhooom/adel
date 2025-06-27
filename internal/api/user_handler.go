package api

import (
	"adel/internal/service/postgres"
	"adel/internal/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type UserHandler struct {
	userStore postgres.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore postgres.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
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

func (h *UserHandler) HandleRegisterNormalUser(w http.ResponseWriter, r *http.Request) {
	h.registerUser(w, r, false)
}

func (h *UserHandler) HandleRegisterAdminUser(w http.ResponseWriter, r *http.Request) {
	h.registerUser(w, r, true)
}
