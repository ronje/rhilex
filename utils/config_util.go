package utils

import (
	"encoding/json"
	"sync"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate
var once sync.Once = sync.Once{}

func __initValidator() {
	Validator = validator.New()
}

// JSON String to a struct, (can't validate map!!!)
func TransformConfig(s1 []byte, s2 any) error {
	if err := json.Unmarshal(s1, &s2); err != nil {
		return err
	}
	once.Do(func() {
		__initValidator()

	})
	if err := Validator.Struct(s2); err != nil {
		return err
	}
	return nil
}

// Bind config to struct
// config: a Map, s: a struct variable

func BindSourceConfig(config map[string]any, s any) error {
	configBytes, err0 := json.Marshal(&config)
	if err0 != nil {
		return err0
	}
	if err := json.Unmarshal(configBytes, &s); err != nil {
		return err
	}
	once.Do(func() {
		__initValidator()
	})
	if err := Validator.Struct(s); err != nil {
		return err
	}
	return nil
}
