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

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	userDB database.UserInterface
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{
		userDB: userDB,
	}
}

// Get JWT godoc
// @Summary Get JWT
// @Description Get JWT
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body dtos.GetJWTInput true "User credentials"
// @Success 200 {object} dtos.GetJwtOutput
// @Failure 400
// @Failure 401
// @failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /users/generate_token [post]
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
		http.Error(w, err.Error(), http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
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

	accessToken := dtos.GetJwtOutput{AccessToken: token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body dtos.CreateUserInput true "User request"
// @Success 201
// @Failure 500 {object} Error
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dtos.CreateUserInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	u, err := entities.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.userDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
