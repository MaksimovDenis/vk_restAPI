package handler

import (
	"encoding/json"
	"net/http"
	filmoteka "vk_restAPI"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body filmoteka.User true "account info"
// @Success 200 {string} string "token"
// @Failure 400 {object} Err "Bad Request"
// @Failure 404 {object} Err "Not Found"
// @Failure 500 {object} Err "Internal Server Error"
// @Failure default {object} Err "Other Errors"
// @Router /auth/sign-up [post]

func (h *Handler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var input filmoteka.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"id": id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// @Summary SignIn
// @Tags auth
// @Description logIn
// @ID logIn
// @Accept json
// @Produce json
// @Param input body filmoteka.User true "credentials"
// @Success 200 {string} string "token"
// @Failure 400 {object} Err "Bad Request"
// @Failure 404 {object} Err "Not Found"
// @Failure 500 {object} Err "Internal Server Error"
// @Failure default {object} Err "Other Errors"
// @Router /auth/sign-in [post]
func (h *Handler) handleSignIn(w http.ResponseWriter, r *http.Request) {
	var input filmoteka.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
