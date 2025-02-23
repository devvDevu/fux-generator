package main

import (
	"ca-generator/internal/custom_type_gen"
	"ca-generator/internal/folders_gen"
)

// test
func main() {
	folders := []string{
		"fux/cmd",
		"fux/internal",
		"fux/internal/common",
		"fux/internal/common/types",
		"fux/internal/common/types/user",
	}
	folders_gen.NewFolderGenerator(folders).Generate()

	files := []*custom_type_gen.CustomType{
		custom_type_gen.NewCustomType("fux/internal/common/types/user", "user_id", ".go", "string"),
		custom_type_gen.NewCustomType("fux/internal/common/types/user", "user_name", ".go", "string"),
		custom_type_gen.NewCustomType("fux/internal/common/types/user", "user_email", ".go", "string"),
		custom_type_gen.NewCustomType("fux/internal/common/types/user", "user_password", ".go", "string"),
		custom_type_gen.NewCustomType("fux/internal/common/types/user", "user_created_at", ".go", "time.Time"),
		custom_type_gen.NewCustomType("fux/internal/common/types/user", "user_updated_at", ".go", "time.Time"),
		custom_type_gen.NewCustomType("fux/internal/common/types/user", "user_deleted_at", ".go", "time.Time"),
	}
	custom_type_gen.NewCustomTypeGenerator(files).Generate()
}
