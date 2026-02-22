package validator

import (
	"context"
	"tickets/internal/domain/errs"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   = validator.New()
	translator ut.Translator
)

func init() {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	translator, _ = uni.GetTranslator("en")
	err := enTranslations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		panic(err)
	}
}

func Validate(ctx context.Context, s any) error {
	err := validate.StructCtx(ctx, s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	msgs := make(map[string]string)
	for _, e := range validationErrors {
		msgs[e.Field()] = e.Translate(translator)
	}

	return errs.ErrInvalidFormat.New(errs.WithMetadata(errs.MetadataErrorsKey, msgs))
}
