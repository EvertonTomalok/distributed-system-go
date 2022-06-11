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
			orders (id, value, method_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	now := time.Now()
	_, err := a.Db.Exec(
		sqlStmt,
		order.ID,
		order.Value,
		order.MethodId,
		order.UserId,
		now,
		now,
	)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return order.ID, nil
}
