// Package docs provides Swagger documentation for the StegGo API
package docs

import (
	"github.com/swaggo/swag"
)

// SwaggerInfo holds exported Swagger Info
var SwaggerInfo = &swag.Spec{
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
