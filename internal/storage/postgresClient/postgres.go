package postgresClient

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"subscriptions/internal/api"
)

// New creates and returns a new PostgresService instance, applies default timeout if not set,
// establishes a connection pool, and runs the migration located at migrationsPath.
func New(ctx context.Context, config *Config, logger *zap.Logger, migrationsPath string) (*PostgresService, error) {
	if config.Timeout == 0 {
		config.Timeout = DefaultPostgresTimeout
	}

	url := buildURL(config)
	dsn := buildDSN(config)

	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}

	err = upMigration(url, migrationsPath)
	if err != nil {
		return nil, err
	}

	return &PostgresService{
		pool:    pool,
		logger:  logger,
		timeout: config.Timeout,
	}, nil
}

// SaveSubscription inserts the given subscription into the database and returns its generated ID.
func (ps *PostgresService) SaveSubscription(subscription *api.Subscription) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ps.timeout)
	defer cancel()

	var id int

	err := ps.pool.QueryRow(ctx, queryForSaveSubscription,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	).Scan(&id)
	if err != nil {
		ps.logger.Error("SaveSubscription: failed to save subscription", zap.Error(err))
		return 0, fmt.Errorf("SaveSubscription: failed to save subscription: %w", err)
	}

	return id, nil
}

// DeleteSubscription deletes a subscription by specified id.
func (ps *PostgresService) DeleteSubscription(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), ps.timeout)
	defer cancel()

	tag, err := ps.pool.Exec(ctx, queryForDeleteSubscription, id)
	if err != nil {
		ps.logger.Error("DeleteSubscription: failed to delete subscription", zap.Error(err), zap.Int("id", id))
		return fmt.Errorf("DeleteSubscription: failed to delete subscription: %w", err)
	}

	if tag.RowsAffected() == 0 {
		ps.logger.Error(ErrSubscriptionNotFound.Error())
		return ErrSubscriptionNotFound
	}

	return nil
}

// GetSubscription return a stored subscription by specified id.
// If there was no subscription, returns ErrSubscriptionNotFound.
func (ps *PostgresService) GetSubscription(id int) (*api.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ps.timeout)
	defer cancel()

	res := &api.Subscription{}

	err := ps.pool.QueryRow(ctx, queryForGetSubscription, id).Scan(
		&res.ServiceName,
		&res.Price,
		&res.UserID,
		&res.StartDate,
		&res.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ps.logger.Error(ErrSubscriptionNotFound.Error())
			return nil, ErrSubscriptionNotFound
		}
		ps.logger.Error("GetSubscriptions: failed to retrieve subscription", zap.Error(err))
		return nil, fmt.Errorf("GetSubscriptions: failed to retrieve subscription: %w", err)
	}

	return res, nil
}

// ListSubscriptions returns all stored subscriptions.
// If there were no subscriptions, returns ErrSubscriptionNotFound.
func (ps *PostgresService) ListSubscriptions() ([]*api.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ps.timeout)
	defer cancel()

	var res []*api.Subscription

	rows, err := ps.pool.Query(ctx, queryForListSubscriptions)
	defer rows.Close()
	if err != nil {
		ps.logger.Error("ListSubscription: failed to retrieve subscriptions", zap.Error(err))
		return nil, fmt.Errorf("ListSubscription: failed to retrieve subscriptions: %w", err)
	}

	for rows.Next() {
		var serviceName, userId, startDate string
		var endDate sql.NullString
		var price int

		err = rows.Scan(&serviceName, &price, &userId, &startDate, &endDate)
		if err != nil {
			ps.logger.Error("ListSubscription: failed to fetch subscriptions", zap.Error(err))
			return nil, fmt.Errorf("ListSubscription: failed to fetch subscriptions: %w", err)
		}

		subscription := &api.Subscription{
			ServiceName: serviceName,
			Price:       price,
			UserID:      userId,
			StartDate:   startDate,
		}

		if endDate.Valid {
			subscription.EndDate = endDate.String
		}

		res = append(res, subscription)
	}

	if rows.Err() != nil {
		ps.logger.Error("ListSubscription: rows error", zap.Error(rows.Err()))
		return nil, fmt.Errorf("ListSubscription: rows error: %w", rows.Err())
	}

	if len(res) == 0 {
		return nil, ErrSubscriptionNotFound
	}

	return res, nil
}

// UpdateSubscription updates specified record by given id.
func (ps *PostgresService) UpdateSubscription(id int, subscription *api.Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), ps.timeout)
	defer cancel()

	tag, err := ps.pool.Exec(ctx, queryForUpdateSubscription,
		id, subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate, subscription.EndDate,
	)
	if err != nil {
		ps.logger.Error("UpdateSubscription: failed to update subscription", zap.Error(err))
		return fmt.Errorf("UpdateSubscription: failed to update subscription: %w", err)
	}

	if tag.RowsAffected() == 0 {
		ps.logger.Error(ErrSubscriptionNotFound.Error())
		return ErrSubscriptionNotFound
	}

	return nil
}

// ListFilteredSubscriptions retrieves a list of subscriptions filtered by user_id and/or service_name.
func (ps *PostgresService) ListFilteredSubscriptions(userID string, serviceName string) ([]*api.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ps.timeout)
	defer cancel()

	rows, err := ps.pool.Query(ctx, queryForListFilteredSubscriptions, userID, serviceName)
	if err != nil {
		return nil, fmt.Errorf("ListFilteredSubscriptions: %w", err)
	}
	defer rows.Close()

	var res []*api.Subscription
	for rows.Next() {
		var subscription api.Subscription
		var endDate sql.NullString
		err = rows.Scan(
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&endDate,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			subscription.EndDate = endDate.String
		}
		res = append(res, &subscription)
	}
	return res, nil
}

// Close closes a connections pool.
func (ps *PostgresService) Close() {
	ps.pool.Close()
}

// buildURL creates a PostgreSQL URL by specified parameters on Config, for perform migrations.
func buildURL(config *Config) string {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	return url
}

// buildDSN creates a PostgreSQL DSN by specified parameters on Config, for perform pool connection.
func buildDSN(config *Config) string {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_max_conns=%d pool_min_conns=%d",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.MaxConns,
		config.MinConns,
	)

	return dsn
}

// upMigration applies database migrations using the specified file path and connection URL.
func upMigration(url string, path string) error {
	migration, err := migrate.New(path, url)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	return nil
}
