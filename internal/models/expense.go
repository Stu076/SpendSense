package models

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type Expense struct {
	bun.BaseModel `bun:"table:expenses" json:"-"`

	ID          int       `bun:",pk,autoincrement"`
	UserID      int       `bun:",notnull"`
	Amount      float64   `bun:",notnull"`
	Category    string    `bun:",notnull"`
	Description string    `bun:",notnull"`
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func CreateExpense(db *bun.DB, expense *Expense) error {
	_, err := db.NewInsert().Model(expense).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}
	return err
}
