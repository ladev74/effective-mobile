package api

// HttpServer defines the configuration parameters for the HTTP server.
type HttpServer struct {
	Host string `env:"HTTP_HOST"`
	Port int    `env:"HTTP_PORT"`
}

// Subscription represents a user's subscription to a service.
// It is used in HTTP requests and responses for CRUDL operation.
type Subscription struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}
