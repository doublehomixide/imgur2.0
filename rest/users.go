package rest

import (
	"encoding/json"
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
}

// RegisterHandler handles user registration
// @Summary Register a new user
// @Description This endpoint registers a new user, stores the user in the database, and generates a JWT token for the user.
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.User true "User data for registration"
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
	err = server.core.RegisterUser(&user)
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

	isCorrect, userID := server.core.LoginUser(&userLogin)

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
