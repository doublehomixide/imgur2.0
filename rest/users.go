package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"pictureloader/models"
	"pictureloader/safety/jwt"
	"pictureloader/service"
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
}

func (server *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return
	}

	err = server.core.RegisterUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	jwtUtils := jwt.UtilsJWT{}
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
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Registration successful"}`))
}

func (server *Server) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	defer r.Body.Close()

	isExist, userID := server.core.LoginUser(username, password)
	if !isExist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jwtUtils := jwt.UtilsJWT{}
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
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Login successful"}`))
}
