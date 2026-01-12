package utils

import (
	"reflect"
	"strings"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Use the existing validator engine from Gin
var validate *validator.Validate

func init() {
	// Initialize validate using Gin's validator engine
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
	} else {
		// Fallback to a new instance if Gin's engine isn't compatible (unlikely)
		validate = validator.New()
	}

	validate.RegisterValidation("solve_strategy", func(fl validator.FieldLevel) bool {
		switch models.SolveStrategy(fl.Field().String()) {
		case models.StrategyBFS, models.StrategyDFS, models.StrategyAStar, models.StrategyGreedy:
			return true
		default:
			return false
		}
	})

	// Register function to get json tag name for fields
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
