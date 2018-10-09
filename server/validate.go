package server

import (
	"errors"

	"github.com/xeipuuv/gojsonschema"
)

var errInvalidPayload = errors.New("invalid payload")

// https://projectborealisgitlab.site/project-borealis/main/pb/blob/code/new-analytics/Plugins/ArcticAnalytics/perf-data.schema.json
var schema = mustCompile(`
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "title": "ArcticAnalytics performance data",
    "description": "Analytics data from the Project Borealis performance test",
    "type": "object",
    "properties": {
        "sessionId": {
            "description": "The unique session identifier for a performance test run.",
            "type": "string"
        },
        "userId": {
            "description": "The random base identifier for a performance test user.",
            "type": "string"
        },
        "events": {
            "description": "The logged analytics events.",
            "type": "array",
            "items": {
                "$ref": "#/definitions/event"
            }
        }
    },
    "required": [ "sessionId", "userId", "events" ],
    "definitions": {
        "event": {
            "description": "An event that was logged by analytics.",
            "type": "object",
            "required": [ "eventName" ],
            "properties": {
                "eventName": {
                    "description": "The key name of this event.",
                    "type": "string"
                },
                "attributes": {
                    "$ref": "#/definitions/attributes"
                }
            }
        },
        "attributes": {
            "description": "The attributes of an event.",
            "type": "array",
            "items": {
                "$ref": "#/definitions/attribute"
            }
        },
        "attribute": {
            "description": "An attribute of an event",
            "type": "object",
            "required": [ "name", "value" ],
            "properties": {
                "name": {
                    "description": "The key name of this attribute.",
                    "type": "string"
                },
                "value": {
                    "description": "The value of the attribute.",
                    "type": "string"
                }
            }
        }
    }
}
`)

func mustCompile(source string) *gojsonschema.Schema {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(source))
	if err != nil {
		panic(err)
	}

	return schema
}

func validate(results []byte) error {
	result, err := schema.Validate(gojsonschema.NewBytesLoader(results))
	if err == nil && result.Valid() {
		return nil
	}

	return errInvalidPayload
}
