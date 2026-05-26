package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name"       validate:"required,min=2,max=255"`
	Price       int        `json:"price"              validate:"required,gte=0"`
	UserID      uuid.UUID  `json:"user_id"            validate:"required,uuid"`
	StartDate   time.Time  `json:"start_date"         validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty" validate:"omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name,omitempty" validate:"omitempty,min=2,max=255"`
	Price       *int       `json:"price,omitempty"        validate:"omitempty,gte=0"`
	StartDate   *time.Time `json:"start_date,omitempty"   validate:"omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"     validate:"omitempty"`
}

type CalculateTotalCostRequest struct {
	ServiceName string    `json:"service_name" validate:"required,min=2,max=255"`
	UserID      uuid.UUID `json:"user_id"      validate:"required,uuid"`
	From        time.Time `json:"from"         validate:"required"`
	To          time.Time `json:"to"           validate:"required"`
}

type CalculateTotalCostResponse struct {
	TotalCost int64 `json:"total_cost" validate:"required,min=2,max=255"`
}
