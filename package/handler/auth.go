package handler

import (
	"encoding/json"
	"net/http"
	filmoteka "vk_restAPI"
	logger "vk_restAPI/logs"
)

// @Summary SignUp
// @Description  create account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body filmoteka.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router       /auth/sign-up [post]
func (h *Handler) handleSignUp(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Sign Up")

	var input filmoteka.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Log.Error("Failed to decaode request body:", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		logger.Log.Error("Failed to create new user:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}
}

type logInInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary LogIn
// @Description  login
// @Tags auth
// @Accept json
// @Produce json
// @Param input body logInInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400 {object} Err
// @Failure 404 {object} Err
// @Failure 500 {object} Err
// @Router       /auth/log-in [post]
func (h *Handler) handleSignIn(w http.ResponseWriter, r *http.Request) {

	logger.Log.Info("Handling Log In")

	var input logInInInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Log.Error("Failed to decaode request body:", err.Error())
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		logger.Log.Error("Failed to generate JWT Token:", err.Error())
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", err.Error())
	}

}
