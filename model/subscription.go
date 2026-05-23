package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID    `json:"id" validate:"required,uuid"`
	ServiceName string       `json:"service_name" validate:"required,max=255"`
	Price       int          `json:"price" validate:"required,gte=0"`
	UserID      uuid.UUID    `json:"user_id" validate:"required,uuid"`
	StartDate   time.Time    `json:"start_date" validate:"required"`
	EndDate     sql.NullTime `json:"end_date"`
}

// SubscriptionSwagger структура для Swagger-документации
// @Description Subscription data
type SubscriptionSwagger struct {
	ID          string `json:"id" swagger:"id,format=uuid"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price" swagger:"minimum=0"`
	UserID      string `json:"user_id" swagger:"format=uuid"`
	StartDate   string `json:"start_date" swagger:"format=date-time"`
	EndDate     string `json:"end_date" swagger:"format=date-time,nullable"`
}

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" validate:"required,max=255"`
	Price       int    `json:"price" validate:"required,gte=0"`
	UserID      string `json:"user_id" validate:"required,uuid"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date"`
}

type UpdateSubscriptionRequest struct {
	ID          string `json:"id" validate:"required,uuid"`
	ServiceName string `json:"service_name" validate:"max=255"`
	Price       int    `json:"price" validate:"omitempty,gte=0"`
	UserID      string `json:"user_id" validate:"omitempty,uuid"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type FilterSubscriptionRequest struct {
	ServiceName string `json:"service_name" validate:"omitempty,max=255"`
	UserID      string `json:"user_id" validate:"omitempty,uuid"`
	StartDate   string `json:"start_date" validate:"omitempty"`
	EndDate     string `json:"end_date" validate:"omitempty"`
}

type SummaResponse struct {
	TotalPrice int `json:"total_price"`
}

// MarshalJSON - кастомная реализация сериализации в JSON
func (s Subscription) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	// Начинаем объект JSON
	buf.WriteString("{")

	// Перебираем поля в нужном порядке
	fields := []struct {
		name  string
		value interface{}
	}{
		{"id", s.ID.String()},
		{"service_name", s.ServiceName},
		{"price", s.Price},
		{"user_id", s.UserID.String()},
		{"start_date", s.StartDate.Format(time.RFC3339)},
	}

	// Добавляем end_date, если валиден
	if s.EndDate.Valid {
		fields = append(fields, struct {
			name  string
			value interface{}
		}{
			name:  "end_date",
			value: s.EndDate.Time.Format(time.RFC3339),
		})
	} else {
		fields = append(fields, struct {
			name  string
			value interface{}
		}{
			name:  "end_date",
			value: nil,
		})
	}

	// Добавляем все поля в нужном порядке
	for i, field := range fields {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("\"")
		buf.WriteString(field.name)
		buf.WriteString("\": ")

		if field.value == nil {
			buf.WriteString("null")
		} else if val, ok := field.value.(string); ok {
			buf.WriteString("\"")
			buf.WriteString(val)
			buf.WriteString("\"")
		} else {
			buf.WriteString(fmt.Sprintf("%v", field.value))
		}
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}
