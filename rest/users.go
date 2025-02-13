package rest

import (
	"context"
	"encoding/json"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"net/http"
	"pictureloader/models"
	"pictureloader/safety/jwtutils"
	"pictureloader/service"
	"time"
)

type Server struct {
	core *service.UserService
}

func NewUserServer(core *service.UserService) *Server {
	return &Server{core: core}
}

func UserRouter(api *mux.Router, server *Server) {
	router := api.PathPrefix("/users").Subrouter()
	router.HandleFunc("/register", server.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", server.LoginUserHandler).Methods("POST")
	router.HandleFunc("/logout", server.LogoutHandler).Methods("POST")

	jwtutils := jwtutils.UtilsJWT{}
	subrouter := router.PathPrefix("/profile").Subrouter()
	subrouter.HandleFunc("", server.DeleteProfile).Methods("DELETE")
	subrouter.HandleFunc("/username", server.ChangeUsername).Methods("PATCH")
	subrouter.Use(jwtutils.AuthMiddleware)
}

// RegisterHandler handles user registration
// @Summary Register a new user
// @Description This endpoint registers a new user, stores the user in the database, and generates a JWT token for the user.
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.UserRegister true "User data for registration"
// @Router /users/register [post]
func (server *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	//получение модели и регистрация ее в бд
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка обработки JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = server.core.RegisterUser(ctx, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//создание jwt
	jwtUtils := jwtutils.UtilsJWT{}
	jwtValue, err := jwtUtils.GenerateToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user-cookie",
		Value:    jwtValue,
		HttpOnly: true,
		Secure:   false, // HTTP & HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 1),
	})
	//ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Registration successful", "cookie":"` + jwtValue + `"}`))
}

// LoginUserHandler handles user login
// @Summary Login an existing user
// @Description This endpoint allows a user to log in by providing their username and password. If the credentials are correct, a JWT token will be generated and returned in a cookie for session management.
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.UserLogin true "User login credentials (username and password)"
// @Router /users/login [post]
func (server *Server) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var userLogin models.UserLogin
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userLogin)
	if err != nil {
		http.Error(w, "Ошибка обработки JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	isCorrect, userID := server.core.LoginUser(ctx, &userLogin)

	if !isCorrect {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect username or password"))
		return
	}
	jwtUtils := jwtutils.UtilsJWT{}
	jwtValue, err := jwtUtils.GenerateToken(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user-cookie",
		Value:    jwtValue,
		HttpOnly: true,
		Secure:   false, // HTTP & HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 1),
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Login successful", "cookie":"` + jwtValue + `"}`))
}

// LogoutHandler handles user logout
// @Summary Log out a user (delete authentication cookie)
// @Description This endpoint allows a user to log out by deleting the authentication cookie from the client's browser.
// @Tags User
// @Accept  json
// @Produce  json
// @Router /users/logout [post]
func (server *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user-cookie",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // HTTP & HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   0,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Logout successful"}`))
}

// DeleteProfile deletes the authenticated user's profile
// @Summary Delete user profile
// @Description This endpoint allows an authenticated user to delete their profile permanently.
// @Tags User
// @Accept json
// @Produce json
// @Router /users/profile [delete]
func (server *Server) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := server.core.DeleteUserByID(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user-cookie",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // HTTP & HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   0,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Delete successful"}`))
}

type usernameReqChange struct {
	Username string `json:"username"`
}

// ChangeUsername allows the user to change their username
// @Summary Change the username of the authenticated user
// @Description This endpoint allows the user to change their username. The new username is passed in the body of the request.
// @Tags User
// @Accept  json
// @Produce  json
// @Param username body usernameReqChange true "New Username"
// @Security BearerAuth
// @Router /users/profile/username [patch]
func (server *Server) ChangeUsername(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	var request usernameReqChange
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "Invalid JSON body"}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := server.core.UpdateUsername(ctx, userID, request.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Change successful"}`))
}
