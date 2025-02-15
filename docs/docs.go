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
        "/posts": {
            "post": {
                "description": "Creates a new post for the user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Create a new post",
                "parameters": [
                    {
                        "description": "Post",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PostRegister"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/posts/my": {
            "get": {
                "description": "Retrieves all posts of the currently authenticated user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Get all posts of the user",
                "responses": {}
            }
        },
        "/posts/{postID}": {
            "get": {
                "description": "Retrieves images and details from a specific post.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Get post details",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Post ID",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "delete": {
                "description": "Removes a specific post and all its images from the system.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Delete a post",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Post ID",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/posts/{postID}/{imageSK}": {
            "post": {
                "description": "Adds an image to the specified post.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Add an image to a post",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Post ID",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Image Storage Key",
                        "name": "imageSK",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "delete": {
                "description": "Removes a specific image from the specified post.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Posts"
                ],
                "summary": "Delete an image from a post",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Post ID",
                        "name": "postID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Image Storage Key",
                        "name": "imageSK",
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
        "/users/profile": {
            "delete": {
                "description": "This endpoint allows an authenticated user to delete their profile permanently.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Delete user profile",
                "responses": {}
            }
        },
        "/users/profile/password": {
            "patch": {
                "description": "Allows an authenticated user to change their password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Change user password",
                "parameters": [
                    {
                        "description": "New password request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/rest.passwordReqChange"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/users/profile/username": {
            "patch": {
                "description": "This endpoint allows the user to change their username. The new username is passed in the body of the request.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Change the username of the authenticated user",
                "parameters": [
                    {
                        "description": "New Username",
                        "name": "username",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/rest.usernameReqChange"
                        }
                    }
                ],
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
                            "$ref": "#/definitions/models.UserRegister"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "models.PostRegister": {
            "type": "object",
            "properties": {
                "name": {
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
        "models.UserRegister": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "rest.passwordReqChange": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                }
            }
        },
        "rest.usernameReqChange": {
            "type": "object",
            "properties": {
                "username": {
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
