package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"io"
	"log"
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

	router.HandleFunc("/{imageName}", server.DownloadFileHandler).Methods("GET")
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

	claims := r.Context().Value("claims").(jwt2.MapClaims) //извлечение из jwt
	sub := claims["sub"].(float64)
	userID := int(sub)

	imageName, err := s.core.Upload(ctx, imageUnit, userID, imgDesc)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
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
// @Param imageName path string true "Name of the image"
// @Router /pictures/{imageName} [get]
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
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result":"` + imageURL + `"}`))
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
		log.Printf("Error fetching pictures: %v", err)
	}

	jsonResponce, err := json.Marshal(map[string][]string{"result": imageURLS})

	if err != nil {
		log.Printf("Error marshalling pictures: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponce)
}
