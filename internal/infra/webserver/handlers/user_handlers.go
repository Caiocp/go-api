package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/caiocp/go-api/internal/dtos"
	"github.com/caiocp/go-api/internal/entities"
	"github.com/caiocp/go-api/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	userDB database.UserInterface
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{
		userDB: userDB,
	}
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var user dtos.GetJWTInput
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := h.userDB.FindByEmail(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !u.ValidatePassword(user.Password) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	_, token, _ := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dtos.CreateUserInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entities.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.userDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
