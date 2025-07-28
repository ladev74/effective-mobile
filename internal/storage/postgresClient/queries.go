package postgresClient

const (
	// queryForSaveSubscription inserts a new subscription into the database.
	queryForSaveSubscription = `
	INSERT INTO schema_subscriptions.subscriptions (service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// queryForDeleteSubscription deletes a subscription record with the given id from the database.
	queryForDeleteSubscription = `DELETE FROM schema_subscriptions.subscriptions WHERE id = $1`

	// queryForGetSubscription selects subscription record with the given id from the database.
	queryForGetSubscription = `
	SELECT service_name, price, user_id, start_date, end_date FROM schema_subscriptions.subscriptions WHERE id = $1`

	// queryForListSubscriptions selects all subscription records from the database.
	queryForListSubscriptions = `
	SELECT service_name, price, user_id, start_date, end_date FROM schema_subscriptions.subscriptions`

	queryForUpdateSubscription = `
	UPDATE schema_subscriptions.subscriptions SET service_name=$2, price=$3, user_id=$4, start_date=$5, end_date=$6 WHERE id = $1`

	queryForListFilteredSubscriptions = `
	SELECT service_name, price, user_id, start_date, end_date 
	FROM schema_subscriptions.subscriptions WHERE ($1 = '' OR user_id = $1) AND ($2 = '' OR service_name = $2)`
)
