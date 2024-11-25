package environment

import "github.com/go-playground/validator/v10"

type Sqlite struct {
	DatabaseName string `envconfig:"SQLITE_DATABASE" validate:"required" default:"zodo-ccs"`
}

func (s Sqlite) IsValid() bool {
	return validator.New().Struct(s) == nil
}
