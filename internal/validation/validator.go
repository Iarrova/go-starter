package validation

import "github.com/Iarrova/go-starter/internal/apperr"

// Validator collects field errors. Only the first error per field is kept.
type Validator struct {
	fields map[string]string
}

func (v *Validator) Check(ok bool, field string, message string) {
	if ok {
		return
	}
	if v.fields == nil {
		v.fields = map[string]string{}
	}
	if _, exists := v.fields[field]; !exists {
		v.fields[field] = message
	}
}

func (v *Validator) Err() error {
	if len(v.fields) == 0 {
		return nil
	}
	return apperr.InvalidFields(v.fields)
}
