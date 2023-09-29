// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Omar Crosby",
            "email": "omar.crosby@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/conferences/:division": {
            "get": {
                "description": "Get a list of conferences",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Conferences"
                ],
                "summary": "Get a list of conferences",
                "parameters": [
                    {
                        "enum": [
                            "all",
                            "di",
                            "dii",
                            "diii",
                            "naia",
                            "njcaa"
                        ],
                        "type": "string",
                        "default": "all",
                        "description": "Specify a division you are interested in",
                        "name": "division",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Conference"
                            }
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Check if the API is up and running",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Health Check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HealthCheckResponse"
                        }
                    }
                }
            }
        },
        "/teams": {
            "get": {
                "description": "Get a list of teams",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Teams"
                ],
                "summary": "Get a list of teams",
                "parameters": [
                    {
                        "enum": [
                            "both",
                            "male",
                            "female"
                        ],
                        "type": "string",
                        "default": "both",
                        "description": "Specify a gender",
                        "name": "gender",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "all",
                            "di",
                            "dii",
                            "diii",
                            "naia",
                            "njcaa"
                        ],
                        "type": "string",
                        "default": "all",
                        "description": "Specify a division you are interested in",
                        "name": "division",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Specify a conference you are interested in",
                        "name": "conference",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The name of the entity you are looking for",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "A partial name of the entity you are looking for",
                        "name": "nameLike",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Specify a target id",
                        "name": "id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Team"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.Error"
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "Get the current version of the API",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Get the API's current version",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.VersionResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.Error": {
            "type": "object",
            "properties": {
                "msg": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "models.Conference": {
            "type": "object",
            "properties": {
                "division": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "models.HealthCheckResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.Team": {
            "type": "object",
            "properties": {
                "conference_id": {
                    "type": "integer"
                },
                "conference_name": {
                    "type": "string"
                },
                "conference_url": {
                    "type": "string"
                },
                "division": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "models.VersionResponse": {
            "type": "object",
            "properties": {
                "version": {
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
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "TopDrawerSoccer API",
	Description:      "This is a simple API providing access to data from TopDrawerSoccer.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
