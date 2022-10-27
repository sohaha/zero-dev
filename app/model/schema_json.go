package model

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

const (
	modelSchema = `{
    "type": "object",
    "properties": {
        "name": {
            "title": "模型名称",
            "type": "string"
        },
        "table": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "title": "表名称",
                    "maxLength": 120,
                    "minLength": 1
                },
                "comment": {
                    "type": "string",
                    "title": "表注释"
                }
            },
            "required": [
                "name"
            ],
            "title": "表配置",
            "x-apifox-orders": [
                "name",
                "comment"
            ]
        },
        "columns": {
            "title": "字段配置",
            "allOf": [
                {
                    "type": "object",
                    "properties": {
                        "label": {
                            "type": "string",
                            "title": "字段展示名称"
                        },
                        "name": {
                            "type": "string",
                            "title": "字段名",
                            "description": "表内名称"
                        },
                        "type": {
                            "type": "string",
                            "description": "字段类型",
                            "enum": [
                                "bool",
                                "int",
                                "uint",
                                "float",
                                "string",
                                "time",
                                "text"
                            ],
                            "x-apifox": {
                                "enumDescriptions": {
                                    "text": ""
                                }
                            }
                        },
                        "validations": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "method": {
                                        "type": "string",
                                        "enum": [
                                            "min",
                                            "max",
                                            "minLength",
                                            "maxLength",
                                            "url",
                                            "ip",
                                            "mobile",
                                            "email",
                                            "regex",
                                            "enum",
                                            "json"
                                        ],
                                        "x-apifox": {
                                            "enumDescriptions": {
                                                "enum": "",
                                                "json": ""
                                            }
                                        }
                                    },
                                    "args": {
                                        "type": [
                                            "string",
                                            "number"
                                        ]
                                    }
                                },
                                "required": [
                                    "method"
                                ],
                                "x-apifox-orders": [
                                    "method",
                                    "args"
                                ]
                            },
                            "description": "数据验证"
                        },
                        "nullable": {
                            "type": "boolean",
                            "description": "是否可是为空"
                        },
                        "comment": {
                            "type": "string",
                            "description": "字段注释",
                            "maxLength": 120
                        },
                        "size": {
                            "type": "integer",
                            "description": "字段长度"
                        },
                        "enum": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "value": {
                                        "type": "string",
                                        "description": "值"
                                    },
                                    "label": {
                                        "type": "string",
                                        "description": "展示名"
                                    }
                                },
                                "x-apifox-orders": [
                                    "value",
                                    "label"
                                ],
                                "required": [
                                    "value"
                                ]
                            },
                            "description": "枚举值"
                        },
                        "default": {
                            "oneOf": [
                                {
                                    "type": "string"
                                },
                                {
                                    "type": "integer"
                                },
                                {
                                    "type": "boolean"
                                }
                            ],
                            "description": "默认值"
                        },
                        "unique": {
                            "description": "唯一索引",
                            "anyOf": [
                                {
                                    "type": "string"
                                },
                                {
                                    "type": "boolean"
                                }
                            ]
                        },
                        "index": {
                            "description": "索引",
                            "anyOf": [
                                {
                                    "type": "string"
                                },
                                {
                                    "type": "boolean"
                                }
                            ]
                        },
                        "read_only": {
                            "type": "boolean",
                            "description": "是否创建之后不允许更改",
                            "title": "只读"
                        }
                    },
                    "required": [
                        "type",
                        "label",
                        "name"
                    ],
                    "x-apifox-orders": [
                        "label",
                        "name",
                        "type",
                        "default",
                        "validations",
                        "nullable",
                        "comment",
                        "size",
                        "enum",
                        "unique",
                        "index",
                        "read_only"
                    ],
                    "title": "字段具体配置"
                }
            ]
        },
        "options": {
            "type": "object",
            "properties": {
                "disabled_migrator": {
                    "type": "boolean",
                    "title": "关闭自动建表"
                },
                "soft_deletes": {
                    "type": "boolean",
                    "title": "开启软删除"
                },
                "timestamps": {
                    "type": "boolean",
                    "title": "插入和更新数据时，标记对应操作时间"
                },
                "api": {
                    "oneOf": [
                        {
                            "type": "boolean",
                            "title": ""
                        },
                        {
                            "type": "object",
                            "properties": {
                                "query": {
                                    "oneOf": [
                                        {
                                            "type": "boolean"
                                        },
                                        {
                                            "type": "object",
                                            "properties": {
                                                "public": {
                                                    "type": "boolean"
                                                }
                                            },
                                            "x-apifox-orders": [
                                                "public"
                                            ],
                                            "required": [
                                                "public"
                                            ]
                                        }
                                    ],
                                    "title": "查询单条数据"
                                },
                                "update": {
                                    "oneOf": [
                                        {
                                            "type": "boolean"
                                        },
                                        {
                                            "type": "object",
                                            "properties": {
                                                "public": {
                                                    "type": "boolean"
                                                }
                                            },
                                            "x-apifox-orders": [
                                                "public"
                                            ],
                                            "required": [
                                                "public"
                                            ]
                                        }
                                    ],
                                    "title": "更新数据"
                                },
                                "delete": {
                                    "oneOf": [
                                        {
                                            "type": "boolean"
                                        },
                                        {
                                            "type": "object",
                                            "properties": {
                                                "public": {
                                                    "type": "boolean"
                                                }
                                            },
                                            "x-apifox-orders": [
                                                "public"
                                            ],
                                            "required": [
                                                "public"
                                            ]
                                        }
                                    ],
                                    "title": "删除数据"
                                },
                                "create": {
                                    "oneOf": [
                                        {
                                            "type": "boolean"
                                        },
                                        {
                                            "type": "object",
                                            "properties": {
                                                "public": {
                                                    "type": "boolean"
                                                }
                                            },
                                            "x-apifox-orders": [
                                                "public"
                                            ],
                                            "required": [
                                                "public"
                                            ]
                                        }
                                    ],
                                    "title": "创建数据"
                                },
                                "pages": {
                                    "oneOf": [
                                        {
                                            "type": "boolean"
                                        },
                                        {
                                            "type": "object",
                                            "properties": {
                                                "public": {
                                                    "type": "boolean"
                                                }
                                            },
                                            "x-apifox-orders": [
                                                "public"
                                            ],
                                            "required": [
                                                "public"
                                            ]
                                        }
                                    ],
                                    "title": "分页查询"
                                }
                            },
                            "x-apifox-orders": [
                                "create",
                                "delete",
                                "update",
                                "query",
                                "pages"
                            ]
                        }
                    ],
                    "title": "注册对应 REST API"
                },
                "api_path": {
                    "type": "string",
                    "title": "接口路径,，如果为空则取表名"
                }
            },
            "x-apifox-orders": [
                "api",
                "api_path",
                "disabled_migrator",
                "soft_deletes",
                "timestamps"
            ]
        },
        "values": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "integer",
                        "title": "ID"
                    }
                },
                "x-apifox-orders": [
                    "id"
                ]
            },
            "title": "默认数据"
        }
    },
    "required": [
        "name",
        "table",
        "columns"
    ],
    "x-apifox-orders": [
        "name",
        "table",
        "columns",
        "options",
        "values",
        "relations"
    ],
    "title": ""
}
`
)

var jsonschemaLoader, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(modelSchema))

func ValidateModelSchema(data []byte) error {
	// TODO 先不校验 schema
	return nil
	res, err := jsonschemaLoader.Validate(gojsonschema.NewBytesLoader(data))
	if err != nil {
		return err
	}
	if !res.Valid() {
		return fmt.Errorf("schema is not valid: %s", res.Errors())
	}
	return nil
}
