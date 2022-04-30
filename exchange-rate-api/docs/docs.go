// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "post": {
                "description": "Inserts new exchange rate",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchange-rates"
                ],
                "summary": "InsertSingleExchangeRate",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/all-from-date/{date}": {
            "get": {
                "description": "Returns all exchange rates for the given date",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchange-rates"
                ],
                "summary": "GetAllExchangeRatesFromDate",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Date for which exchange rates should be retrived",
                        "name": "date",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/check": {
            "get": {
                "description": "basic healthcheck",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "healthcheck"
                ],
                "summary": "healthcheck",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/last": {
            "get": {
                "description": "Returns most recent exchange rate source - destinaion currencies",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchange-rates"
                ],
                "summary": "GetLastExchangeRate",
                "parameters": [
                    {
                        "type": "string",
                        "description": "source currency, default is USD",
                        "name": "source",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "destination, currency",
                        "name": "destination",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/range/": {
            "get": {
                "description": "Returns exchange rates for currencies in the time period",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchange-rates"
                ],
                "summary": "GetrangeExchangeRate",
                "parameters": [
                    {
                        "type": "string",
                        "description": "source currency, default is USD",
                        "name": "source",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "destination, currency",
                        "name": "destination",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "From date, inclusive",
                        "name": "from",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Till date, exclusive",
                        "name": "till",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Rate Exchange API",
	Description:      "This is a sample server server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
