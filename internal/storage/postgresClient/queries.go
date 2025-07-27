package postgresClient

const (
	// queryForSaveSubscription inserts a new subscription into the database.
	queryForSaveSubscription = `
	INSERT INTO schema_subscriptions.subscriptions (service_name, price, user_id, start_time, end_time)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// queryForDeleteSubscription deletes a subscription record with the given id from the database.
	queryForDeleteSubscription = `DELETE FROM schema_subscriptions.subscriptions WHERE id = $1`

	// queryForGetSubscription selects subscription record with the given id from the database.
	queryForGetSubscription = `
	SELECT service_name, price, user_id, start_time, end_time FROM schema_subscriptions.subscriptions WHERE id = $1`

	// queryForListSubscriptions selects all subscription records from the database.
	queryForListSubscriptions = `
	SELECT service_name, price, user_id, start_time, end_time FROM schema_subscriptions.subscriptions`

	queryForUpdateSubscription = `
	UPDATE schema_subscriptions.subscriptions SET service_name=$2, price=$3, user_id=$4, start_time=$5, end_time=$6 WHERE id = $1`
)
