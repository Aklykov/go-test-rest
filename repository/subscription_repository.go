package repository

import (
	"errors"
	"fmt"
	"go-rest-crud/model"

	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// Экспортируемая ошибка для случая, когда запись не найдена
var ErrNotFound = errors.New("subscription not found")

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create добавляет новую подписку
func (r *SubscriptionRepository) Create(sub *model.Subscription) error {
	log.Printf("[REPOSITORY] Creating subscription: %+v", sub)

	// проверка на существование записи с таким сервисом, юзером и интервалом дат
	isExist, err := r.IsExists(sub)
	if isExist {
		log.Printf("[REPOSITORY] Subscription already exists!")
		return fmt.Errorf("Subscription already exists!")
	}

	query := `
		INSERT INTO subscription (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err = r.db.QueryRow(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&sub.ID)
	if err != nil {
		log.Printf("[REPOSITORY] Error creating subscription: %v", err)
		return fmt.Errorf("failed to create subscription")
	}

	log.Printf("[REPOSITORY] Subscription created with id=%s", sub.ID)

	return nil
}

// Update обновляет подписку
func (r *SubscriptionRepository) Update(id uuid.UUID, sub *model.Subscription) error {

	// проверка на существование записи с таким сервисом, юзером и интервалом дат
	isExist, err := r.IsExists(sub)
	if isExist {
		log.Printf("[REPOSITORY] Subscription already exists!")
		return fmt.Errorf("Subscription already exists!")
	}

	query := `
		UPDATE subscription
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
	`
	log.Printf("[REPOSITORY] Updating subscription: id=%s", id)

	_, err = r.db.Exec(query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, id)
	if err != nil {
		log.Printf("[REPOSITORY] Error updating subscription: %v", err)
		return fmt.Errorf("failed to update subscription")
	}

	return nil
}

// GetByID возвращает подписку по ID
func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*model.Subscription, error) {
	log.Printf("[REPOSITORY] Getting subscription by id=%v", id)
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscription
		WHERE id = $1
	`

	var sub model.Subscription
	err := r.db.QueryRow(query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate,
	)
	if err != nil {
		log.Printf("[REPOSITORY] Subscription not found: id=%v", id)
		return nil, fmt.Errorf("Subscription not found")
	}

	log.Printf("[REPOSITORY] Subscription found: id=%v", sub.ID)
	return &sub, nil
}

// GetAll возвращает все подписки с фильтрами
func (r *SubscriptionRepository) GetAll(filter model.FilterSubscriptionRequest) ([]model.Subscription, error) {

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscription
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if filter.ServiceName != "" {
		argCount++
		query += fmt.Sprintf(" AND service_name = $%d", argCount)
		args = append(args, filter.ServiceName)
	}
	if filter.UserID != "" {
		argCount++
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, filter.UserID)
	}
	if filter.StartDate != "" {
		argCount++
		startDate, _ := parseDate(filter.StartDate)
		query += fmt.Sprintf(" AND start_date >= $%d", argCount)
		args = append(args, startDate)
	}
	if filter.EndDate != "" {
		argCount++
		endDate, _ := parseDate(filter.EndDate)
		query += fmt.Sprintf(" AND (end_date <= $%d OR end_date IS NULL)", argCount)
		args = append(args, endDate)
	}

	query += " ORDER BY start_date DESC"
	query += " LIMIT 1000"

	log.Printf("[REPO] Getting subscriptions with filters: %+v", filter)
	log.Printf("[REPO] Getting subscriptions query: %+v", query)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Printf("[REPO] Error querying subscriptions: %v", err)
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			log.Printf("[REPO] Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[REPO] Error iterating subscriptions: %v", err)
		return nil, fmt.Errorf("error iterating subscriptions: %w", err)
	}

	return subscriptions, nil
}

// Delete удаляет подписку
func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
	log.Printf("[REPOSITORY] Delete subscription by id=%v", id)
	query := `DELETE FROM subscription WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("[REPOSITORY] Error delete subscription: %v", err)
		return fmt.Errorf("Failed to delete subscription")
	}

	log.Printf("[REPOSITORY] Success Delete subscription id=: %v", id)

	return nil
}

func (r *SubscriptionRepository) IsExists(sub *model.Subscription) (bool, error) {
	var query string
	var exists bool
	var err error
	if sub.EndDate.Valid {
		query = `
				WITH new_subscription AS (
					SELECT
						$1 AS new_service_name,
						$2::uuid AS new_user_id,
						$3::date AS new_start_date,
						$4::date AS new_end_date
				)
				SELECT
					EXISTS (
						SELECT 1
						FROM subscription s
								JOIN new_subscription ns ON true
						WHERE
							s.service_name = ns.new_service_name AND
							s.user_id = ns.new_user_id AND
							s.id != $5 AND
							(
								(s.end_date IS NOT NULL AND ns.new_start_date <= s.end_date AND ns.new_end_date >= s.start_date)
									OR
								(s.end_date IS NULL AND (ns.new_start_date >= s.start_date OR ns.new_end_date >= s.start_date))
							)
					) AS has_conflict;
			`
		err = r.db.QueryRow(query, sub.ServiceName, sub.UserID, sub.StartDate, sub.EndDate, sub.ID).Scan(&exists)
	} else {
		query = `
				WITH new_subscription AS (
					SELECT
						$1 AS new_service_name,
						$2::uuid AS new_user_id,
						$3::date AS new_start_date
				)
				SELECT
					EXISTS (
						SELECT 1
						FROM subscription s
								JOIN new_subscription ns ON true
						WHERE
							s.service_name = ns.new_service_name AND
							s.user_id = ns.new_user_id AND
							s.id != $4 AND
							(
								(s.end_date IS NOT NULL AND ns.new_start_date <= s.end_date)
									OR
								(s.end_date IS NULL)
							)
					) AS has_conflict;
			`
		err = r.db.QueryRow(query, sub.ServiceName, sub.UserID, sub.StartDate, sub.ID).Scan(&exists)
	}

	if err != nil {
		log.Printf("[REPOSITORY] Error checking existence: %v", err)
		return false, fmt.Errorf("failed to check subscription existence: %w", err)
	}

	return exists, nil
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("01-2006", dateStr)
}
