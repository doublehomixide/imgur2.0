package rest

import (
	"bytes"
	"context"
	"fmt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"pictureloader/models"
	"pictureloader/safety/jwt"
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
	jwtUtils := jwt.UtilsJWT{}
	router := api.PathPrefix("/pictures").Subrouter()

	privateRouter := router.PathPrefix("").Subrouter()
	privateRouter.HandleFunc("/create", server.UploadFileHandler).Methods("POST")
	privateRouter.HandleFunc("/my", server.MyPictures).Methods("GET")
	privateRouter.Use(jwtUtils.AuthMiddleware)

	router.HandleFunc("/{imageName}", server.DownloadFileHandler).Methods("GET")
}

// UploadFileHandler Обработчик для загрузки файла
func (s *PictureServer) UploadFileHandler(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	size, err := io.Copy(&buf, file)
	if err != nil {
		log.Printf("Error reading file to buffer: %v", err)
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

	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)
	imageName, err := s.core.Upload(ctx, imageUnit, userID)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		http.Error(w, fmt.Sprintf("Error uploading file: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("File uploaded successfully")
	w.Write([]byte(fmt.Sprintf("File uploaded successfully with image name: %s", imageName)))
}

func (s *PictureServer) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["imageName"]

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	imageURL, err := s.core.Download(ctx, imageName)
	if err != nil {
		log.Printf("Error downloading file from Minio: %v", err)
		http.Error(w, fmt.Sprintf("Error downloading file: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	htmlContent := fmt.Sprintf("<img src=\"%s\" />", imageURL)
	w.Write([]byte(htmlContent))
}

func (s *PictureServer) MyPictures(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	imageURLS, err := s.core.GetAllUserPictures(ctx, userID)
	if err != nil {
		log.Printf("Error fetching pictures: %v", err)
	}
	w.Header().Set("Content-Type", "text/html")

	htmlContent := "<div style=\"display: flex; flex-wrap: wrap; gap: 10px;\">"
	for _, imageURL := range imageURLS {
		htmlContent += fmt.Sprintf("<img src=\"%s\" style=\"max-width: 100px;\" />", imageURL)
	}
	htmlContent += "</div>"
	w.Write([]byte(htmlContent))
}
