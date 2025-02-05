package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"io"
	"log/slog"
	"net/http"
	"pictureloader/models"
	"pictureloader/safety/jwtutils"
	"pictureloader/service"
	"time"
)

type PictureServer struct {
	core *service.PictureLoader
}

func PictureNewServer(core *service.PictureLoader) *PictureServer {
	return &PictureServer{core: core}
}

func PictureRouter(api *mux.Router, server *PictureServer) {
	jwtUtils := jwtutils.UtilsJWT{}
	router := api.PathPrefix("/pictures").Subrouter()

	privateRouter := router.PathPrefix("").Subrouter()
	privateRouter.HandleFunc("/create", server.UploadImageHandler).Methods("POST")
	privateRouter.HandleFunc("/my", server.MyPictures).Methods("GET")
	privateRouter.Use(jwtUtils.AuthMiddleware)

	router.HandleFunc("/{imageURL}", server.DownloadFileHandler).Methods("GET")
	router.HandleFunc("/{imageURL}", server.DeleteImageHadler).Methods("DELETE")
}

// UploadImageHandler handles image upload
// @Summary Upload an image
// @Description This endpoint allows a user to upload an image file.
// @Tags Image
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "Image file"
// @Param desription formData string true "Image description"
// @Router /pictures/create [post]
func (s *PictureServer) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	imgDesc := r.FormValue("desription")

	file, _, err := r.FormFile("file")
	if err != nil {
		slog.Error("Error retrieving file", "error", err)
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		slog.Error("Error retrieving file", "error", err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	imageUnit := models.ImageUnit{
		Payload:     bytes.NewReader(buf.Bytes()), // Ensure correct payload
		PayloadName: "uploaded_image.png",
		PayloadSize: size, // Correct file size
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	claims := r.Context().Value("claims").(jwt2.MapClaims) //извлечение из jwt
	sub := claims["sub"].(float64)
	userID := int(sub)

	imageName, err := s.core.Upload(ctx, imageUnit, userID, imgDesc)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading file: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"message":"Picture uploaded successfully.", "picture": "` + imageName + `"}`))
}

// DownloadFileHandler handles image download
// @Summary Download an image
// @Description This endpoint allows a user to download an image file.
// @Tags Image
// @Accept json
// @Produce  json
// @Param imageURL path string true "URL of the image"
// @Router /pictures/{imageURL} [get]
func (s *PictureServer) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["imageURL"]

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	imageURL, description, err := s.core.Download(ctx, imageName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error downloading file: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"img":"` + imageURL + `", "desc": "` + description + `"}`))
}

// MyPictures handles all user images
// @Summary Download an image(s)
// @Description This endpoint allows a user to download his images.
// @Tags Image
// @Accept json
// @Produce  json
// @Router /pictures/my [get]
func (s *PictureServer) MyPictures(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	imageURLS, err := s.core.GetAllUserPictures(ctx, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonResponce, err := json.Marshal(map[string][]string{"result": imageURLS})

	if err != nil {
		slog.Error("Error retrieving file", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponce)
}

// DeleteImageHadler delete image
// @Summary Delete an image
// @Description Delete an image for its url
// @Tags Image
// @Accept json
// @Produce json
// @Param imageURL path string true "Image url"
// @Router /pictures/{imageURL} [delete]
func (s *PictureServer) DeleteImageHadler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["imageURL"]

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := s.core.Delete(ctx, imageName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting image: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Picture deleted successfully."}`))
}
