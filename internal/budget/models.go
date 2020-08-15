package budget

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Budget type is a group of related financial transactions
type Budget struct {
	ID          primitive.ObjectID `bson:"_id",omitempty"json:"_id,omitempty" validate:"required"`
	BudgetName  string             `bson:"budget_name,omitempty"json:"budget_name,omitempty"validate:"required"`
	BudgetValue float64            `bson:"budget_value,omitempty"json:"budget_value,omitempty"default:"0"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"datetime"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"datetime"`
}

// NewBudget type is what's required from the client to create a new Budget
type NewBudget struct {
	BudgetName  string  `bson:"budget_name,omitempty"json:"budget_name,omitempty"validate:"required"`
	BudgetValue float64 `bson:"budget_value,omitempty"json:"budget_value,omitempty"default:"0"`
}
