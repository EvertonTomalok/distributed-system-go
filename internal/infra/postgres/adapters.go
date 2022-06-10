package postgres

import (
	"database/sql"

	log "github.com/sirupsen/logrus"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/entities"
	"golang.org/x/net/context"
)

func New() *Adapter {
	return &Adapter{
		db: db,
	}
}

type Adapter struct {
	db *sql.DB
}

func (a *Adapter) GetAllMethods(ctx context.Context) ([]entities.Method, error) {
	sqlStmt := `
		SELECT id, name, installment FROM methods;
	`

	rows, err := a.db.QueryContext(ctx, sqlStmt)
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
