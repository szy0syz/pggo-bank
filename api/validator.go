package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/szy0syz/pggo-bank/util"
)

// Func accepts a FieldLevel interface for all validation needs. The return
// value should be true when validation succeeds.
// func(fl FieldLevel) bool
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is supported
		return util.IsSupportedCurrency(currency)
	}
	return false
}
