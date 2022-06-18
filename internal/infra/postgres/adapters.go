package postgres

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	"golang.org/x/net/context"
)

func New(dataBase *sql.DB) *Adapter {
	return &Adapter{
		Db: dataBase,
	}
}

type Adapter struct {
	Db *sql.DB
}

func (a *Adapter) GetAllMethods(ctx context.Context) ([]entities.Method, error) {
	sqlStmt := `
		SELECT id, name, installment FROM methods;
	`

	rows, err := a.Db.QueryContext(ctx, sqlStmt)
	defer rows.Close()

	var methods []entities.Method
	if err != nil {
		log.Info("Any method was found.")
		return methods, nil
	}

	for rows.Next() {
		var method entities.Method
		if err := rows.Scan(&method.ID, &method.Name, &method.Installment); err != nil {
			log.Error("Error scanning method: +v", err)
			return nil, err
		}
		methods = append(methods, method)
	}
	return methods, nil
}

func (a *Adapter) GetMethodByNameAndInstallment(ctx context.Context, methodName string, installment int64) (entities.Method, error) {
	sqlStmt := `
		SELECT id, name, installment FROM methods WHERE name = $1 AND installment = $2;
	`
	method := entities.Method{}
	err := a.Db.QueryRowContext(ctx, sqlStmt, methodName, installment).Scan(&method.ID, &method.Name, &method.Installment)
	if err != nil {
		return entities.Method{}, err
	}
	return method, nil
}

func (a *Adapter) PostOrder(ctx context.Context, order entities.Order) (string, error) {
	sqlStmt := `
		INSERT INTO 
			orders (id, value, method_id, user_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	_, err := a.Db.Exec(
		sqlStmt,
		order.ID,
		order.Value,
		order.MethodId,
		order.UserId,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return order.ID, nil
}

func (a *Adapter) GetOrdersByUserId(ctx context.Context, userId string, offset int64, limit int64) ([]entities.Order, error) {
	sqlStmt := `
		SELECT 
			orders.id, orders.value, orders.method_id, orders.user_id, orders.status, orders.created_at, orders.updated_at, methods.name, methods.installment 
		FROM orders
		INNER JOIN methods ON orders.method_id = methods.id
		WHERE user_id = $1 OFFSET $2 LIMIT $3;
	`

	var orders []entities.Order

	rows, err := a.Db.QueryContext(ctx, sqlStmt, userId, offset, limit)
	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		log.Info("Any method was found.")
		return orders, nil
	}

	for rows.Next() {
		var order entities.Order
		if err := rows.Scan(
			&order.ID,
			&order.Value,
			&order.MethodId,
			&order.UserId,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.Method.Name,
			&order.Method.Installment,
		); err != nil {
			log.Error("Error scanning method: +v", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (a *Adapter) GetOrderById(ctx context.Context, orderId string) (entities.Order, error) {
	sqlStmt := `
		SELECT 
			orders.id, orders.value, orders.method_id, orders.user_id, orders.status, orders.created_at, orders.updated_at,  methods.name, methods.installment 
		FROM orders
		INNER JOIN methods ON orders.method_id = methods.id
		WHERE orders.id = $1;
	`

	var order entities.Order

	err := a.Db.QueryRowContext(ctx, sqlStmt, orderId).Scan(
		&order.ID, &order.Value, &order.MethodId, &order.UserId, &order.Status, &order.CreatedAt, &order.UpdatedAt, &order.Method.Name, &order.Method.Installment,
	)

	if err != nil {
		log.Info("Any method was found.")
		return order, err
	}
	return order, nil
}

func (a *Adapter) UpdateStatusByOrderId(ctx context.Context, orderId string, status string) error {
	sqlStmt := `
		UPDATE orders
		SET status = $2, updated_at = $3
		WHERE id = $1;
	`
	_, err := a.Db.Exec(
		sqlStmt,
		orderId,
		status,
		time.Now(),
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
