// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/albums": {
            "post": {
                "description": "Creates an album for an authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Albums"
                ],
                "summary": "Create an album",
                "parameters": [
                    {
                        "description": "Album name",
                        "name": "Name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AlbumRegister"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/albums/add-image": {
            "post": {
                "description": "Adds an image with imageID to an album with albumID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Albums"
                ],
                "summary": "Add an image to an album",
                "parameters": [
                    {
                        "description": "Data for adding image to album",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/rest.Request"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/albums/my": {
            "get": {
                "description": "Retrieves all albums of the user by their ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Albums"
                ],
                "summary": "Get user albums",
                "responses": {}
            }
        },
        "/albums/{albumID}": {
            "get": {
                "description": "Retrieves all images from the specified album for an authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Albums"
                ],
                "summary": "Get images from an album",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Album ID",
                        "name": "albumID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "delete": {
                "description": "Deletes the specified album along with its associated images for an authenticated user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Albums"
                ],
                "summary": "Delete an album",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Album ID",
                        "name": "albumID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/pictures/create": {
            "post": {
                "description": "This endpoint allows a user to upload an image file.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Image"
                ],
                "summary": "Upload an image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Image file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Image description",
                        "name": "desription",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/pictures/my": {
            "get": {
                "description": "This endpoint allows a user to download his images.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Image"
                ],
                "summary": "Download an image(s)",
                "responses": {}
            }
        },
        "/pictures/{imageURL}": {
            "get": {
                "description": "This endpoint allows a user to download an image file.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Image"
                ],
                "summary": "Download an image",
                "parameters": [
                    {
                        "type": "string",
                        "description": "StorageKey of the image",
                        "name": "imageURL",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "delete": {
                "description": "Delete an image for its url",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Image"
                ],
                "summary": "Delete an image",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Image url",
                        "name": "imageURL",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/users/login": {
            "post": {
                "description": "This endpoint allows a user to log in by providing their username and password. If the credentials are correct, a JWT token will be generated and returned in a cookie for session management.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Login an existing user",
                "parameters": [
                    {
                        "description": "User login credentials (username and password)",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserLogin"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/users/logout": {
            "post": {
                "description": "This endpoint allows a user to log out by deleting the authentication cookie from the client's browser.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Log out a user (delete authentication cookie)",
                "responses": {}
            }
        },
        "/users/register": {
            "post": {
                "description": "This endpoint registers a new user, stores the user in the database, and generates a JWT token for the user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User data for registration",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "models.AlbumRegister": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "models.Image": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "images": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Image"
                    }
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.UserLogin": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "rest.Request": {
            "type": "object",
            "properties": {
                "album_id": {
                    "type": "integer"
                },
                "image_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Imgur 2.0 API",
	Description:      "API для загрузки и просмотра картинок с регистрацией",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
