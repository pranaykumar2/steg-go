package docs

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "StegGo API",
    "description": "REST API for StegGo - Image Steganography Tool",
    "contact": {
      "name": "pranaykumar2",
      "url": "https://github.com/pranaykumar2/steg-go"
    },
    "license": {
      "name": "MIT"
    },
    "version": "1.0.0"
  },
  "host": "localhost:8080",
  "basePath": "/api",
  "paths": {
    "/health": {
      "get": {
        "produces": ["application/json"],
        "summary": "Health check endpoint",
        "description": "Returns the health status of the API",
        "responses": {
          "200": {
            "description": "API is healthy",
            "schema": {
              "type": "object",
              "properties": {
                "status": {
                  "type": "string",
                  "example": "ok"
                },
                "version": {
                  "type": "string",
                  "example": "1.0.0"
                },
                "time": {
                  "type": "string",
                  "example": "2025-03-24T02:45:00Z"
                }
              }
            }
          }
        }
      }
    },
    "/hide": {
      "post": {
        "consumes": ["multipart/form-data"],
        "produces": ["application/json"],
        "summary": "Hide text in an image",
        "description": "Encrypts and hides a text message inside an image using steganography",
        "parameters": [
          {
            "name": "image",
            "in": "formData",
            "description": "Image file to hide text in (PNG or JPG)",
            "required": true,
            "type": "file"
          },
          {
            "name": "message",
            "in": "formData",
            "description": "Text message to hide",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Message hidden successfully",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": true
                },
                "message": {
                  "type": "string",
                  "example": "Message hidden successfully"
                },
                "data": {
                  "type": "object",
                  "properties": {
                    "key": {
                      "type": "string",
                      "example": "5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d"
                    },
                    "outputFileURL": {
                      "type": "string",
                      "example": "/api/files/stego_abc123.png"
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "Invalid request",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "No image file uploaded"
                }
              }
            }
          },
          "500": {
            "description": "Server error",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "Failed to hide message"
                }
              }
            }
          }
        }
      }
    },
    "/hideFile": {
      "post": {
        "consumes": ["multipart/form-data"],
        "produces": ["application/json"],
        "summary": "Hide a file in an image",
        "description": "Encrypts and hides any file inside an image using steganography",
        "parameters": [
          {
            "name": "image",
            "in": "formData",
            "description": "Cover image file (PNG or JPG)",
            "required": true,
            "type": "file"
          },
          {
            "name": "file",
            "in": "formData",
            "description": "File to hide (PDF, TXT, etc.)",
            "required": true,
            "type": "file"
          }
        ],
        "responses": {
          "200": {
            "description": "File hidden successfully",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": true
                },
                "message": {
                  "type": "string",
                  "example": "File hidden successfully"
                },
                "data": {
                  "type": "object",
                  "properties": {
                    "key": {
                      "type": "string",
                      "example": "5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d5a7b8c9d"
                    },
                    "outputFileURL": {
                      "type": "string",
                      "example": "/api/files/stego_abc123.png"
                    },
                    "fileDetails": {
                      "type": "object",
                      "properties": {
                        "originalName": {
                          "type": "string",
                          "example": "document.pdf"
                        },
                        "fileType": {
                          "type": "string",
                          "example": ".pdf"
                        },
                        "fileSize": {
                          "type": "integer",
                          "example": 12345
                        }
                      }
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "Invalid request",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "No file to hide uploaded"
                }
              }
            }
          },
          "500": {
            "description": "Server error",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "Failed to hide file"
                }
              }
            }
          }
        }
      }
    },
    "/extract": {
      "post": {
        "consumes": ["multipart/form-data"],
        "produces": ["application/json"],
        "summary": "Extract hidden content",
        "description": "Extracts and decrypts hidden content from a steganographic image",
        "parameters": [
          {
            "name": "image",
            "in": "formData",
            "description": "Image containing hidden data",
            "required": true,
            "type": "file"
          },
          {
            "name": "key",
            "in": "formData",
            "description": "Decryption key (64 hexadecimal characters)",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Content extracted successfully",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": true
                },
                "message": {
                  "type": "string",
                  "example": "Content extracted successfully"
                },
                "data": {
                  "type": "object",
                  "properties": {
                    "isFile": {
                      "type": "boolean"
                    },
                    "message": {
                      "type": "string",
                      "example": "This is a secret message."
                    },
                    "fileURL": {
                      "type": "string",
                      "example": "/api/files/document.pdf"
                    },
                    "fileName": {
                      "type": "string",
                      "example": "document.pdf"
                    },
                    "fileType": {
                      "type": "string",
                      "example": ".pdf"
                    },
                    "fileSize": {
                      "type": "integer",
                      "example": 12345
                    },
                    "contentType": {
                      "type": "string",
                      "example": "application/pdf"
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "Invalid request",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "Invalid key format"
                }
              }
            }
          },
          "500": {
            "description": "Server error",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "Failed to extract content"
                }
              }
            }
          }
        }
      }
    },
    "/metadata": {
      "post": {
        "consumes": ["multipart/form-data"],
        "produces": ["application/json"],
        "summary": "Analyze image metadata",
        "description": "Analyzes the metadata of an image including steganography capacity",
        "parameters": [
          {
            "name": "image",
            "in": "formData",
            "description": "Image file to analyze",
            "required": true,
            "type": "file"
          }
        ],
        "responses": {
          "200": {
            "description": "Metadata analysis completed",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": true
                },
                "message": {
                  "type": "string",
                  "example": "Metadata analysis completed"
                },
                "data": {
                  "type": "object",
                  "properties": {
                    "filename": {
                      "type": "string",
                      "example": "image.png"
                    },
                    "fileSize": {
                      "type": "integer",
                      "example": 123456
                    },
                    "fileType": {
                      "type": "string",
                      "example": "PNG"
                    },
                    "mimeType": {
                      "type": "string",
                      "example": "image/png"
                    },
                    "modTime": {
                      "type": "string",
                      "example": "2025-03-24T02:45:00Z"
                    },
                    "imageWidth": {
                      "type": "integer",
                      "example": 1920
                    },
                    "imageHeight": {
                      "type": "integer",
                      "example": 1080
                    },
                    "hasEXIF": {
                      "type": "boolean",
                      "example": false
                    },
                    "privacyRisks": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      },
                      "example": [
                        "Steganography preserves most metadata - consider removing sensitive metadata"
                      ]
                    },
                    "properties": {
                      "type": "object",
                      "additionalProperties": {
                        "type": "string"
                      },
                      "example": {
                        "Image Width": "1920 pixels",
                        "Image Height": "1080 pixels",
                        "Steganography Capacity": "~776.25 KB"
                      }
                    },
                    "steganoCapacity": {
                      "type": "object",
                      "properties": {
                        "bytes": {
                          "type": "integer",
                          "example": 777600
                        },
                        "kilobytes": {
                          "type": "number",
                          "format": "float",
                          "example": 759.38
                        },
                        "megabytes": {
                          "type": "number",
                          "format": "float",
                          "example": 0.742
                        },
                        "text": {
                          "type": "object",
                          "properties": {
                            "characters": {
                              "type": "integer",
                              "example": 777600
                            },
                            "words": {
                              "type": "integer",
                              "example": 155520
                            }
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "Invalid request",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "No image file uploaded"
                }
              }
            }
          },
          "500": {
            "description": "Server error",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "Failed to extract metadata"
                }
              }
            }
          }
        }
      }
    },
    "/files/{filename}": {
      "get": {
        "summary": "Get a file",
        "description": "Retrieve a file generated by the API",
        "parameters": [
          {
            "name": "filename",
            "in": "path",
            "description": "Name of the file to retrieve",
            "required": true,
            "type": "string"
          }
        ],
        "produces": [
          "image/png",
          "image/jpeg",
          "application/pdf",
          "text/plain",
          "application/octet-stream"
        ],
        "responses": {
          "200": {
            "description": "File content",
            "schema": {
              "type": "file"
            }
          },
          "400": {
            "description": "Invalid filename",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "Invalid filename"
                }
              }
            }
          },
          "404": {
            "description": "File not found",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean",
                  "example": false
                },
                "error": {
                  "type": "string",
                  "example": "File not found"
                }
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Response": {
      "type": "object",
      "properties": {
        "success": {
          "type": "boolean"
        },
        "message": {
          "type": "string"
        },
        "data": {
          "type": "object"
        },
        "error": {
          "type": "string"
        }
      }
    }
  }
}`
