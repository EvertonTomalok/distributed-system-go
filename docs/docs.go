// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/methods": {
            "get": {
                "description": "Get available methods",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "methods"
                ],
                "summary": "Get Methods",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.MethodResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Something went wrong"
                    }
                }
            }
        },
        "/orders": {
            "post": {
                "description": "Create order",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "Create Order",
                "parameters": [
                    {
                        "description": "Order to create",
                        "name": "Order",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.OrderRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.OrderResponse"
                        }
                    },
                    "400": {
                        "description": "{'error': 'error description'}"
                    },
                    "404": {
                        "description": "Something went wrong. Try again."
                    },
                    "503": {
                        "description": "Feature Flag is disabled."
                    }
                }
            }
        },
        "/orders/{orderId}": {
            "get": {
                "description": "Get order using its id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "Get Order by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The order id to search",
                        "name": "orderId",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.OrderResponse"
                        }
                    },
                    "404": {
                        "description": "Order not found"
                    },
                    "500": {
                        "description": "Something went wrong"
                    }
                }
            }
        },
        "/orders/{userId}": {
            "get": {
                "description": "Get orders from User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "Get Orders",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The user id to search",
                        "name": "userId",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.OrdersResponse"
                        }
                    },
                    "500": {
                        "description": "Something went wrong"
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.MethodResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "installments": {
                    "type": "integer"
                },
                "method": {
                    "type": "string"
                }
            }
        },
        "dto.OrderRequest": {
            "type": "object",
            "required": [
                "installment",
                "method",
                "user_id",
                "value"
            ],
            "properties": {
                "installment": {
                    "type": "integer"
                },
                "method": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "dto.OrderResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "installment": {
                    "type": "integer"
                },
                "method": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "dto.OrdersResponse": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "orders": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.OrderResponse"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
