package api

import (
	"adel/internal/middleware"
	"adel/internal/service/postgres"
	"adel/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"log"
)

type TokenHandler struct {
	tokenStore postgres.TokenStore
	userStore  postgres.UserStore
	logger     *log.Logger
}

func NewTokenHandler(tokenStore postgres.TokenStore, userStore postgres.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HandleLoginUser authenticates user and returns token
// @Summary      Login user
// @Description  Authenticate user with username and password, returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      loginUserRequest  true  "User login credentials"
// @Success      201          {object}  utils.Envelope{token=token.Token}
// @Failure      400          {object}  utils.Envelope{error=string}
// @Failure      401          {object}  utils.Envelope{error=string}
// @Failure      500          {object}  utils.Envelope{error=string}
// @Router       /users/login [post]
func (h *TokenHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: createTokenRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		} else {
			h.logger.Printf("ERROR: GetUserByUsername: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		}

		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERORR: PasswordHash.Mathes %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if !passwordsDoMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour)
	if err != nil {
		h.logger.Printf("ERORR: Creating Token %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return

	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"token": token})
}

// HandleDeleteToken logs out user by deleting tokens
// @Summary      Logout user
// @Description  Logout user by deleting all their authentication tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Envelope{message=string}
// @Failure      401  {object}  utils.Envelope{error=string}
// @Failure      500  {object}  utils.Envelope{error=string}
// @Security     BearerAuth
// @Router       /users/logout [delete]
func (h *TokenHandler) HandleDeleteToken(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be authenticated to logout"})
		return
	}

	err := h.tokenStore.DeleteAllTokensForUser(user.ID)
	if err != nil {
		h.logger.Printf("ERROR: deleting tokens for user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "logged out successfully"})
}
