package rest

import (
	"encoding/json"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"net/http"
	"pictureloader/models"
	"pictureloader/safety/jwtutils"
	"pictureloader/service"
	"strconv"
)

type AlbumServer struct {
	service service.AlbumService
}

func NewAlbumServer(service service.AlbumService) *AlbumServer {
	return &AlbumServer{service: service}
}

func AlbumRouter(api *mux.Router, server *AlbumServer) {
	jwtUtils := jwtutils.UtilsJWT{}
	router := api.PathPrefix("/albums").Subrouter()
	router.HandleFunc("/", server.CreateAlbumHandler).Methods("POST")
	router.HandleFunc("/my", server.GetMyAlbums).Methods("GET")
	router.HandleFunc("/{albumID}", server.GetAlbum).Methods("GET")
	router.HandleFunc("/add-image", server.AddImageToAlbum).Methods("POST")
	router.HandleFunc("/{albumID}", server.DeleteAlbum).Methods("DELETE")
	router.Use(jwtUtils.AuthMiddleware)
}

// CreateAlbumHandler creates a new album for the user.
// @Summary Create an album
// @Description Creates an album for an authenticated user
// @Tags Albums
// @Accept  json
// @Produce  json
// @Param Name body models.AlbumRegister true "Album name"
// @Router /albums [post]
func (as *AlbumServer) CreateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	var album models.Album
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&album)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	album.UserID = userID

	err = as.service.CreateAlbum(&album)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{status:"Album created"}`))
}

// GetAlbum retrieves images from a specific album.
// @Summary Get images from an album
// @Description Retrieves all images from the specified album for an authenticated user
// @Tags Albums
// @Accept  json
// @Produce  json
// @Param albumID path int true "Album ID"
// @Router /albums/{albumID} [get]
func (as *AlbumServer) GetAlbum(w http.ResponseWriter, r *http.Request) {
	albumIDstr := mux.Vars(r)["albumID"]
	albumID, err := strconv.Atoi(albumIDstr)
	result, err := as.service.GetAlbum(albumID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(result)
}

type Request struct {
	AlbumID int    `json:"album_id"`
	ImageID string `json:"image_id"`
}

// AddImageToAlbum adds an image to an album.
// @Summary Add an image to an album
// @Description Adds an image with imageID to an album with albumID
// @Tags Albums
// @Accept  json
// @Produce  json
// @Param request body Request true "Data for adding image to album"
// @Router /albums/add-image [post]
func (as *AlbumServer) AddImageToAlbum(w http.ResponseWriter, r *http.Request) {
	var req Request
	json.NewDecoder(r.Body).Decode(&req)
	err := as.service.AppendImageToAlbum(req.AlbumID, req.ImageID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"status":"Image added"}`))
}

// GetMyAlbums get my albums handler
// @Summary Get user albums
// @Description Retrieves all albums of the user by their ID.
// @Tags Albums
// @Accept json
// @Produce json
// @Router /albums/my [get]
func (as *AlbumServer) GetMyAlbums(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	result, err := as.service.GetUserAlbums(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(result)
}

// DeleteAlbum removes a specific album and its images from the system.
// @Summary Delete an album
// @Description Deletes the specified album along with its associated images for an authenticated user.
// @Tags Albums
// @Accept  json
// @Produce  json
// @Param albumID path int true "Album ID"
// @Router /albums/{albumID} [delete]
func (as *AlbumServer) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	albumIDstr := mux.Vars(r)["albumID"]
	albumID, err := strconv.Atoi(albumIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = as.service.DeleteAlbum(albumID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"status":"Album deleted"}`))
}
