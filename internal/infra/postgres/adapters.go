package postgres

import (
	"database/sql"

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

func (a *Adapter) GetMethodByNameAndInstallment(ctx context.Context, methodName string, installment int) (entities.Method, error) {
	sqlStmt := `
		SELECT id, name, installment FROM methods WHERE name = $1 AND installment = $2;
	`
	_, err := a.Db.QueryContext(ctx, sqlStmt, methodName, installment)
	if err != nil {
		return entities.Method{}, err
	}
	return entities.Method{}, nil
}
