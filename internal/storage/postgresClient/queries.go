package postgresClient

const (
	// queryForSaveSubscription inserts a new subscription into the database.
	queryForSaveSubscription = `INSERT INTO schema_subscriptions.subscriptions (service_name, price, user_id, start_time, end_time)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	queryForDeleteSubscription = `DELETE FROM schema_subscriptions.subscriptions WHERE id = $1`
)
