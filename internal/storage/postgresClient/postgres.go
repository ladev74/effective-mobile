package postgresClient

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"effmob/internal/api"
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

func (ps *PostgresService) SaveSubscription(subscription *api.Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), ps.timeout)
	defer cancel()

	tag, err := ps.pool.Exec(ctx, queryForSaveSubscription,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	)
	if err != nil {
		ps.logger.Error("SaveSubscription: failed to save subscription", zap.Error(err))
		return fmt.Errorf("SaveSubscription: failed to save subscription: %w", err)
	}

	if tag.RowsAffected() == 0 {
		ps.logger.Error("SaveSubscription: no rows affected")
		return fmt.Errorf("SaveSubscription: no rows affected")
	}

	return nil
}

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
