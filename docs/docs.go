// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "tbd",
        "contact": {
            "name": "Rezoan Tamal",
            "email": "my.full.name.in.lower.case@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/private/customers/profile": {
            "get": {
                "description": "Returns user's profile using access token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customers"
                ],
                "summary": "Get basic profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Set access token here",
                        "name": "authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerSuccessRes"
                        }
                    },
                    "400": {
                        "description": "Invalid request body, or missing required fields.",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "401": {
                        "description": "Unauthorized access attempt.",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "500": {
                        "description": "API sever or db unreachable.",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    }
                }
            },
            "patch": {
                "description": "Update user's basic profile info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customers"
                ],
                "summary": "Update basic profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Set access token here",
                        "name": "authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Some fields are mandatory",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CustomerProfileUpdateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerSuccessRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/private/customers/refresh-token": {
            "get": {
                "description": "Generate new access and refresh tokens using current refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customers"
                ],
                "summary": "Refresh customer's access token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Value of refresh token",
                        "name": "authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.TokenSuccessRes"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/private/customers/verify-token": {
            "get": {
                "description": "VerifyAccessToken lets apps to verify that a provided token is in-fact valid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customers"
                ],
                "summary": "Verify customer's access token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Value of access token",
                        "name": "authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.EmptySuccessRes"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/public/customers/forgot-password": {
            "post": {
                "description": "ForgotPassword uses username and captcha to send otp to user's registered number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customers"
                ],
                "summary": "Reset user's password with otp",
                "parameters": [
                    {
                        "description": "All fields are mandatory",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ForgotPasswordReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.EmptySuccessRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/public/customers/login": {
            "post": {
                "description": "Login uses user defined username and password to authenticate a user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customers"
                ],
                "summary": "Login as a customer",
                "parameters": [
                    {
                        "description": "All fields are mandatory",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.LoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.TokenSuccessRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/public/customers/signup": {
            "post": {
                "description": "Signup a new customer for a valid non-existing phone number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customers"
                ],
                "summary": "Signup a new customer",
                "parameters": [
                    {
                        "description": "All fields are mandatory",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CustomerSignupReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.EmptySuccessRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    }
                }
            }
        },
        "/api/v1/public/customers/verify-signup": {
            "post": {
                "description": "VerifySignUp uses user defined otp and matches it with existing reference in cache to verify a signup",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customers"
                ],
                "summary": "VerifyAccessToken a new customer using otp",
                "parameters": [
                    {
                        "description": "All fields are mandatory",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CustomerSignupVerificationReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.EmptySuccessRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.CustomerErrorRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Customer": {
            "type": "object",
            "properties": {
                "birth_date": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "is_deleted": {
                    "type": "boolean"
                },
                "is_verified": {
                    "type": "boolean"
                },
                "occupation": {
                    "type": "string"
                },
                "organization": {
                    "type": "string"
                },
                "recovery_phone_number": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.CustomerProfileUpdateReq": {
            "type": "object",
            "properties": {
                "birth_date": {
                    "type": "string",
                    "example": "2006-01-02T15:04:05.000Z"
                },
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "gender": {
                    "type": "string",
                    "example": "male/female/other"
                },
                "occupation": {
                    "type": "string"
                },
                "organization": {
                    "type": "string"
                }
            }
        },
        "model.CustomerSignupReq": {
            "type": "object",
            "properties": {
                "captcha_id": {
                    "type": "string"
                },
                "captcha_value": {
                    "type": "string"
                },
                "full_name": {
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
        "model.CustomerSignupVerificationReq": {
            "type": "object",
            "properties": {
                "otp": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.EmptyObject": {
            "type": "object"
        },
        "model.ForgotPasswordReq": {
            "type": "object",
            "properties": {
                "captcha_id": {
                    "type": "string"
                },
                "captcha_value": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.LoginReq": {
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
        "model.Token": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "response.CustomerErrorRes": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.EmptyObject"
                },
                "message": {
                    "type": "string",
                    "example": "failure message"
                },
                "status": {
                    "type": "string",
                    "example": "Status string corresponding to the error"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                },
                "timestamp": {
                    "type": "string",
                    "example": "2006-01-02T15:04:05.000Z"
                }
            }
        },
        "response.CustomerSuccessRes": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.Customer"
                },
                "message": {
                    "type": "string",
                    "example": "success message"
                },
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "success": {
                    "type": "boolean",
                    "example": true
                },
                "timestamp": {
                    "type": "string",
                    "example": "2006-01-02T15:04:05.000Z"
                }
            }
        },
        "response.EmptySuccessRes": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.EmptyObject"
                },
                "message": {
                    "type": "string",
                    "example": "success message"
                },
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                },
                "timestamp": {
                    "type": "string",
                    "example": "2006-01-02T15:04:05.000Z"
                }
            }
        },
        "response.TokenSuccessRes": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.Token"
                },
                "message": {
                    "type": "string",
                    "example": "success message"
                },
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "success": {
                    "type": "boolean",
                    "example": false
                },
                "timestamp": {
                    "type": "string",
                    "example": "2006-01-02T15:04:05.000Z"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "https://auth-iamrz1.cloud.okteto.net",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "auth",
	Description: "This is auth REST api server",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
