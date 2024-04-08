package types

import (
	"errors"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func BlacklistValidation(field string) validation.RuleFunc {
	return func(value interface{}) error {
		val, ok := value.(string)

		if !ok {
			return errors.New("The " + field + " is not a string")
		}

		if val == "" {
			return nil
		}

		match, _ := regexp.MatchString(`^[^'"\[\]<>\{\}]+$`, val)

		if !match {
			return errors.New("The " + field + " contains unsafe characters")
		}

		return nil
	}
}

func DatetimeValidation(field string) validation.RuleFunc {
	return func(value interface{}) error {
		val, ok := value.(string)

		if !ok {
			return errors.New("The " + field + " is not datetime format")
		}

		if len(val) != 19 {
			return errors.New("the " + field + " is not a datetime format")
		}

		_, err := time.Parse("2006-01-02 15:04:05", val)

		if err != nil {
			return errors.New("the " + field + " is not a datetime format")
		}

		return nil
	}
}

func (r GetUserRequest) Validate() interface{} {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Page, is.Digit),
		validation.Field(&r.Limit, is.Digit),
		validation.Field(&r.OrderBy, validation.In("id", "name", "createdAt", "updatedAt")),
		validation.Field(&r.SortBy, validation.In("asc", "desc")),
		validation.Field(&r.Search, validation.By(BlacklistValidation("search"))),
		validation.Field(&r.DisableCalculateTotal, validation.In("true", "false")),
		validation.Field(&r.ID, validation.By(BlacklistValidation("id"))),
	)
}

func (r CreateUserRequest) Validate() interface{} {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.By(BlacklistValidation("name"))),
		validation.Field(&r.Emails, validation.Required),
	)
}

func (r CreateEmailRequest) Validate() interface{} {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email, validation.By(BlacklistValidation("email"))),
	)
}

func (r EditUserRequest) Validate() interface{} {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, validation.By(BlacklistValidation("id"))),
		validation.Field(&r.Name, validation.By(BlacklistValidation("name"))),
	)
}

func (r DeleteUserRequest) Validate() interface{} {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, validation.By(BlacklistValidation("id"))),
	)
}
