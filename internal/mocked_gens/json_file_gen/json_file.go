package json_file_gen

import (
	"log"
	"os"
)

// explain:
// example_name_of_project is the name of the project
// example_field is the name of the field
// example_type is the type of the field
// example_adapter is the adapter of the field
// example_repository is the repository of the field
// example_service is the service of the field

/*
How its works:
Generate the json file like this:
"settings": {
    "example_name_of_project": {
        "cmd": {
            "app": {}
        },
    }
This json is equal example_name_of_project/cmd/app

Or like this:
"settings": {
    "example_name_of_project": {
        "internal": {
            "common": {
                "types": {
                    "example_type": [
                        {
                            "file_name": "example_field",
                            "file_ext": ".go",
                            "file_type": "example_type"
                        },
            },
        }
    }
}
This json is equal example_name_of_project/internal/common/types/example_type.go
*/

func GenerateJsonFile() {
	err := os.WriteFile("settings.json", []byte(jsonData), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

const jsonData = `{
    "settings": {
            "example_name_of_project": {
                "cmd": {
                    "app": {}
                },
                "internal": {
                    "config": {},
                    "common": {
                        "types": {
                            "error_with_codes": {},
                            "example_type": [ 
                                {
                                    "file_name": "example_field",
                                    "file_ext": ".go",
                                    "file_type": "example_type"
                                }
                                
                            ]
                        }
                    },
                    "adapter": {
                        "example_adapter": {}
                    },
                    "repository": {
                        "example_repository": {}
                    },
                    "service": {
                        "example_service": {}
                    },
                    "usecase": {
                        "example_usecase": {}
                    },
                    "handler": {
                        "example_handler": {}
                    },
                    "value_object": {
                        "example_value_object": [
                            {
                                "file_name": "example_value_object",
                                "file_ext": ".go",
                                "fields": [
                                    {
                                        "name": "example_field",
                                        "type": "example_type.ExampleType",
                                        "tag": ""
                                    }
                                ]
                            }
                        ]
                    },
                    "model": {
                        "example_model": [
                            {
                                "file_name": "example",
                                "file_ext": ".go",
                                "fields": [
                                    {
                                        "name": "example_field",
                                        "type": "example_type.ExampleType",
                                        "tag": ""
                                    }
                                ]
                            }
                        ]
                    },
                    "data_transfer_object":{
                        "example_dto": {
                            "response": {},
                            "request": {}
                        },
                        "result": {}
                    },
                    "pkg": {
                        "env":{}
                    }
                },
                "migrations": {
                    "postgres": {}
                }
            }
        }
}`
