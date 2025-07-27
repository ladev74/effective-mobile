package postgresClient

const (
	// queryForSaveSubscription inserts a new subscription into the database.
	queryForSaveSubscription = `
	INSERT INTO schema_subscriptions.subscriptions (service_name, price, user_id, start_time, end_time)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// queryForDeleteSubscription deletes a subscription record with the given id from the database.
	queryForDeleteSubscription = `DELETE FROM schema_subscriptions.subscriptions WHERE id = $1`

	// queryForGetSubscriptions selects subscription record with the given id from the database.
	queryForGetSubscriptions = `
	SELECT service_name, price, user_id, start_time, end_time FROM schema_subscriptions.subscriptions WHERE id = $1`

	// queryForListSubscriptions selects all subscription records from the database.
	queryForListSubscriptions = `SELECT * FROM schema_subscriptions.subscriptions`
)
