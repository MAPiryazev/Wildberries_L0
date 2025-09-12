package validation

import (
	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterStructValidation(deliveryStructLevelValidation, models.Delivery{})

}

func ValidateOrder(order *models.Order) error {
	return validate.Struct(order)
}

func deliveryStructLevelValidation(sl validator.StructLevel) {
	delivery := sl.Current().Interface().(models.Delivery)

	if delivery.Phone == "" && delivery.Email == "" {
		sl.ReportError(delivery.Phone, "Phone", "phone", "phoneOrEmail", "")
		sl.ReportError(delivery.Email, "Email", "email", "phoneOrEmail", "")
	}
}
