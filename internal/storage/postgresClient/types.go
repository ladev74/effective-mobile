package postgresClient

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"effmob/internal/api"
)

// DefaultPostgresTimeout defines the default timeout for PostgreSQL operations.
const DefaultPostgresTimeout = 3 * time.Second

var (
	ErrSubscriptionNotFound = fmt.Errorf("Subscription was not found for the specified id")
)

// Config defines the configuration parameters for the PostgresService,
// including credentials and timeout configuration.
type Config struct {
	Host     string        `env:"POSTGRES_HOST"`
	Port     string        `env:"POSTGRES_PORT"`
	User     string        `env:"POSTGRES_USER"`
	Password string        `env:"POSTGRES_PASSWORD"`
	Database string        `env:"POSTGRES_DATABASE"`
	Timeout  time.Duration `env:"POSTGRES_TIMEOUT"`
	MaxConns int           `env:"POSTGRES_MAX_CONNECTIONS"`
	MinConns int           `env:"POSTGRES_MIN_CONNECTIONS"`
}

// PostgresService implements the PostgresClient interface.
// It provides methods for storing and retrieving subscription using a PostgreSQL database.
type PostgresService struct {
	pool    *pgxpool.Pool
	logger  *zap.Logger
	timeout time.Duration
}

// PostgresClient defines an interface for storing and retrieving subscription in a PostgreSQL database.
type PostgresClient interface {
	SaveSubscription(*api.Subscription) (int, error)
	DeleteSubscription(int) error
	GetSubscriptions(int) (*api.Subscription, error)
	Close()
}
