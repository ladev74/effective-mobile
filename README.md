<details><summary>
basePath: /
definitions:
  api.Subscription:
    properties:
      end_date:
        type: string
      price:
        type: integer
      service_name:
        type: string
      start_date:
        type: string
      user_id:
        type: string
    type: object
  handlers.response:
    properties:
      data: {}
      message:
        type: string
      status:
        type: string
    type: object
host: localhost:8081
info:
  contact: {}
  description: This is a service for managing subscriptions.
  title: Subscriptions API
  version: "1.0"
paths:
  /subscriptions:
    get:
      description: Returns a list of all subscriptions stored in the database.
      produces:
      - application/json
      responses:
        "200":
          description: OK
          schema:
            items:
              $ref: '#/definitions/api.Subscription'
            type: array
        "404":
          description: Not Found
          schema:
            $ref: '#/definitions/handlers.response'
        "500":
          description: Internal Server Error
          schema:
            $ref: '#/definitions/handlers.response'
      summary: List all subscriptions
      tags:
      - subscriptions
    post:
      consumes:
      - application/json
      description: Adds a new subscription for a user
      parameters:
      - description: Subscription data
        in: body
        name: subscription
        required: true
        schema:
          $ref: '#/definitions/api.Subscription'
      produces:
      - application/json
      responses:
        "201":
          description: Created
          schema:
            $ref: '#/definitions/handlers.response'
        "400":
          description: Bad Request
          schema:
            $ref: '#/definitions/handlers.response'
        "500":
          description: Internal Server Error
          schema:
            $ref: '#/definitions/handlers.response'
      summary: Add a new subscription
      tags:
      - subscriptions
  /subscriptions/{id}:
    delete:
      description: Deletes a subscription by ID
      parameters:
      - description: Subscription ID
        in: path
        name: id
        required: true
        type: integer
      produces:
      - application/json
      responses:
        "200":
          description: OK
          schema:
            $ref: '#/definitions/handlers.response'
        "400":
          description: Bad Request
          schema:
            $ref: '#/definitions/handlers.response'
        "404":
          description: Not Found
          schema:
            $ref: '#/definitions/handlers.response'
        "500":
          description: Internal Server Error
          schema:
            $ref: '#/definitions/handlers.response'
      summary: Delete a subscription
      tags:
      - subscriptions
    get:
      description: Returns subscription details for the given subscription ID.
      parameters:
      - description: Subscription ID
        in: path
        name: id
        required: true
        type: integer
      produces:
      - application/json
      responses:
        "200":
          description: OK
          schema:
            $ref: '#/definitions/api.Subscription'
        "400":
          description: Bad Request
          schema:
            $ref: '#/definitions/handlers.response'
        "404":
          description: Not Found
          schema:
            $ref: '#/definitions/handlers.response'
        "500":
          description: Internal Server Error
          schema:
            $ref: '#/definitions/handlers.response'
      summary: Get subscription by ID
      tags:
      - subscriptions
    put:
      consumes:
      - application/json
      description: Updates subscription data for the given subscription ID.
      parameters:
      - description: Subscription ID
        in: path
        name: id
        required: true
        type: integer
      - description: Subscription data to update
        in: body
        name: subscription
        required: true
        schema:
          $ref: '#/definitions/api.Subscription'
      produces:
      - application/json
      responses:
        "200":
          description: OK
          schema:
            $ref: '#/definitions/handlers.response'
        "400":
          description: Bad Request
          schema:
            $ref: '#/definitions/handlers.response'
        "500":
          description: Internal Server Error
          schema:
            $ref: '#/definitions/handlers.response'
      summary: Update a subscription by ID
      tags:
      - subscriptions
  /subscriptions/total-price:
    get:
      description: Calculates the total price for subscriptions filtered by user_id
        and/or service_name during the specified date range.
      parameters:
      - description: User ID filter
        in: query
        name: user_id
        type: string
      - description: Service Name filter
        in: query
        name: service_name
        type: string
      - description: Start date in MM-YYYY format
        in: query
        name: start_date
        required: true
        type: string
      - description: End date in MM-YYYY format
        in: query
        name: end_date
        required: true
        type: string
      produces:
      - application/json
      responses:
        "200":
          description: Total price
          schema:
            type: integer
        "400":
          description: Bad Request
          schema:
            $ref: '#/definitions/handlers.response'
        "500":
          description: Internal Server Error
          schema:
            $ref: '#/definitions/handlers.response'
      summary: Calculate total price of subscriptions
      tags:
      - subscriptions
swagger: "2.0"

</summary></details>