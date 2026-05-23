package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Aklykov/go-test-rest/model"
	"github.com/Aklykov/go-test-rest/repository"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

// парсит дату и возвращает *time.Time (или nil, если endDate пусто)
func parseOptionalDate(date string) (*time.Time, error) {
	if date == "" {
		return nil, nil // возврат nil для необязательного поля
	}
	parsedDate, err := time.Parse("01-2006", date)
	if err != nil {
		return nil, err
	}
	return &parsedDate, nil
}

func (s *SubscriptionService) validateUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return err
	}
	return nil
}

// CreateSubscription создает новую подписку
func (s *SubscriptionService) CreateSubscription(req model.CreateSubscriptionRequest) (*model.Subscription, error) {
	log.Printf("[SERVICE] Creating subscription: %+v", req)

	if req.ServiceName == "" || req.Price < 0 || req.UserID == "" || req.StartDate == "" {
		log.Printf("[SERVICE] Missing required fields")
		return nil, fmt.Errorf("Missing required fields")
	}

	userIdUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Printf("[SERVICE] Invalid id format UUID: %v", err)
		return nil, fmt.Errorf("Invalid user_id format, must be UUID")
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		log.Printf("[SERVICE] Invalid start_date format: %v", err)
		return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
	}

	var endDate sql.NullTime
	if req.EndDate != "" {
		ed, err := parseDate(req.EndDate)
		if err != nil {
			log.Printf("[SERVICE] Invalid end_date format: %v", err)
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
		}
		endDate = sql.NullTime{Time: ed, Valid: true}

		// Validate end_date > start_date
		if endDate.Time.Before(startDate) || endDate.Time.Equal(startDate) {
			log.Printf("[SERVICE] end_date must be after start_date")
			return nil, fmt.Errorf("end_date must be after start_date")
		}
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userIdUUID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	log.Printf("[SERVICE] Subscription: %+v", sub)

	if err := s.repo.Create(sub); err != nil {
		return nil, err
	}

	log.Printf("[SERVICE] Subscription created: id=%s", sub.ID)
	return sub, nil
}

// UpdateSubscription обновляет подписку
func (s *SubscriptionService) UpdateSubscription(req model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	log.Printf("[SERVICE] Update subscription: %+v", req)

	idUUID, err := uuid.Parse(req.ID)
	if err != nil {
		log.Printf("[SERVICE] Invalid id format UUID: %v", err)
		return nil, fmt.Errorf("Invalid user_id format, must be UUID")
	}

	subExist, err := s.repo.GetByID(idUUID)
	if err != nil {
		return nil, err
	}

	sub := &model.Subscription{}
	sub.ID = idUUID

	if req.ServiceName != "" {
		sub.ServiceName = req.ServiceName
	} else {
		sub.ServiceName = subExist.ServiceName
	}

	if req.Price != 0 {
		sub.Price = req.Price
	} else {
		sub.Price = subExist.Price
	}

	if req.UserID != "" {
		userIdUUID, err := uuid.Parse(req.UserID)
		if err != nil {
			log.Printf("[SERVICE] Invalid id format UUID: %v", err)
			return nil, fmt.Errorf("Invalid user_id format, must be UUID")
		}
		sub.UserID = userIdUUID
	} else {
		sub.UserID = subExist.UserID
	}

	if req.StartDate != "" {
		se, err := parseDate(req.StartDate)
		if err != nil {
			log.Printf("[SERVICE] Invalid start_date format: %v", err)
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
		}
		sub.StartDate = se
	} else {
		sub.StartDate = subExist.StartDate
	}
	if req.EndDate != "" {
		ed, err := parseDate(req.EndDate)
		if err != nil {
			log.Printf("[SERVICE] Invalid end_date format: %v", err)
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
		}
		sub.EndDate = sql.NullTime{Time: ed, Valid: true}
	} else {
		sub.EndDate = subExist.EndDate
	}
	// Validate end_date > start_date
	if sub.EndDate.Valid {
		if sub.EndDate.Time.Before(sub.StartDate) || sub.EndDate.Time.Equal(sub.StartDate) {
			log.Printf("[SERVICE] end_date must be after start_date")
			return nil, fmt.Errorf("end_date must be after start_date")
		}
	}

	log.Printf("[SERVICE] Subscription: %+v", sub)

	if err := s.repo.Update(sub.ID, sub); err != nil {
		return nil, err
	}

	log.Printf("[SERVICE] Subscription update: id=%v", req.ID)
	return sub, nil
}

// GetAllSubscriptions получает список подписок с фильтрами
func (s *SubscriptionService) GetAllSubscriptions(filter model.FilterSubscriptionRequest) ([]model.Subscription, error) {
	log.Printf("[SERVICE] Getting all subscriptions with filters: %+v", filter)

	if filter.UserID != "" {
		_, err := uuid.Parse(filter.UserID)
		if err != nil {
			log.Printf("[SERVICE] Invalid id format UUID: %v", err)
			return nil, fmt.Errorf("Invalid user_id format, must be UUID")
		}
	}
	if filter.StartDate != "" {
		_, err := parseDate(filter.StartDate)
		if err != nil {
			log.Printf("[SERVICE] Invalid start_date format: %v", err)
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
		}
	}
	if filter.EndDate != "" {
		_, err := parseDate(filter.EndDate)
		if err != nil {
			log.Printf("[SERVICE] Invalid end_date format: %v", err)
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
		}
	}

	return s.repo.GetAll(filter)
}

// GetAllSubscriptionsSumma
func (s *SubscriptionService) GetAllSubscriptionsSumma(filter model.FilterSubscriptionRequest) (model.SummaResponse, error) {
	log.Printf("[SERVICE] Getting summa subscriptions with filters: %+v", filter)

	summa := model.SummaResponse{TotalPrice: 0}

	_, err := uuid.Parse(filter.UserID)
	if err != nil {
		log.Printf("[SERVICE] Invalid id format UUID: %v", err)
		return summa, fmt.Errorf("Invalid user_id format, must be UUID")
	}
	if filter.StartDate != "" {
		_, err := parseDate(filter.StartDate)
		if err != nil {
			log.Printf("[SERVICE] Invalid start_date format: %v", err)
			return summa, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
		}
	}
	if filter.EndDate != "" {
		_, err := parseDate(filter.EndDate)
		if err != nil {
			log.Printf("[SERVICE] Invalid end_date format: %v", err)
			return summa, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
		}
	}

	subs, err := s.repo.GetAll(filter)
	if err != nil {
		log.Printf("[SERVICE] Invalid get subscriptions: %v", err)
		return summa, fmt.Errorf("Invalid get subscriptions")
	}
	if len(subs) == 0 {
		log.Printf("[SERVICE] No subscriptions found")
		return summa, fmt.Errorf("No subscriptions found")
	}
	for _, sub := range subs {
		summa.TotalPrice += sub.Price
	}

	return summa, nil
}

// GetSubscriptionByID возвращает подписку по ID
func (s *SubscriptionService) GetSubscriptionByID(id string) (*model.Subscription, error) {
	log.Printf("[SERVICE] Get subscription by ID: %+v", id)

	idUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("[SERVICE] Invalid id format UUID: %v", err)
		return nil, fmt.Errorf("Invalid id format, must be UUID")
	}

	subscription, err := s.repo.GetByID(idUUID)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

// DeleteSubscription удаляет подписку по ID
func (s *SubscriptionService) DeleteSubscription(id string) error {
	log.Printf("[SERVICE] Delete subscription: %+v", id)

	idUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("[SERVICE] Invalid id format UUID: %v", err)
		return fmt.Errorf("Invalid id format, must be UUID")
	}

	_, err = s.repo.GetByID(idUUID)
	if err != nil {
		log.Printf("[SERVICE] Subscription not found: %v", err)
		return fmt.Errorf("[SERVICE] Subscription not found")
	}

	err = s.repo.Delete(idUUID)
	if err != nil {
		return err
	}

	return nil
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("01-2006", dateStr)
}
