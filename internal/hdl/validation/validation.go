package validation

import (
	"github.com/JMURv/effective-mobile/internal/dto"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var V *validator.Validate //nolint:gochecknoglobals

func New() {
	V = validator.New()

	V.RegisterStructValidation(
		calculateTotalCostRequestValidation,
		dto.CalculateTotalCostRequest{},
	)
}

func calculateTotalCostRequestValidation(sl validator.StructLevel) {
	req, ok := sl.Current().Interface().(dto.CalculateTotalCostRequest)
	if !ok {
		zap.L().Error("calculate total cost request type cast failed")
		return
	}

	if !req.From.IsZero() && !req.To.IsZero() && req.From.After(req.To) {
		sl.ReportError(
			req.From,
			"from",
			"From",
			"before_to",
			"",
		)
	}
}
