const Ajv = require('ajv');
const ajv = new Ajv();

const schema = {
    type: "object",
    properties: {
        sessionId: { type: "string" },
        userId: { type: "string" },
        events: {
            type: "array",
            items: {
                type: "object",
                required: ["eventName"],
                properties: {
                    eventName: { type: "string" },
                    attributes: {
                        type: "array",
                        items: {
                            type: "object",
                            required: ["name", "value"],
                            properties: {
                                name: { type: "string" },
                                value: { type: "string" }
                            }
                        }
                    }
                }
            }
        }
    },
    required: ["sessionId", "userId", "events"]
};

function validateSchema(data) {
    const validate = ajv.compile(schema);
    const valid = validate(JSON.parse(data));
    return valid;
}

module.exports = { validateSchema };
