package budget

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Budget type is a group of related financial transactions
type Budget struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	UserID      string             `bson:"userID,omitempty" json:"userID,omitempty"`
	BudgetName  string             `bson:"budget_name,omitempty" json:"budget_name,omitempty" validate:"required"`
	BudgetValue float64            `bson:"budget_value,omitempty" json:"budget_value,omitempty" default:"0"` // default doesn't give desired result
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"datetime"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"datetime"`
}

// NewBudget type is what's required from the client to create a new Budget
type NewBudget struct {
	UserID      string  `bson:"userID,omitempty" json:"userID,omitempty"`
	BudgetName  string  `bson:"budget_name,omitempty" json:"budget_name,omitempty" validate:"required"`
	BudgetValue float64 `bson:"budget_value,omitempty" json:"budget_value,omitempty" default:"0"` // default doesn't give desired result
}

// UpdateBudget defines what information may be provided to modify an existing Budget.
// All fields are optional so clients can send just the fields they want
// changed.
// It uses pointer fields so we can differentiate between a field that was not provided and a field that was provided as explicitly blank.
// Normally we do not want to use pointers to basic types but we make exceptions around marshalling/unmarshalling.
type UpdateBudget struct {
	ID          *primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID      *string             `bson:"userID,omitempty" json:"userID,omitempty"`
	BudgetName  *string             `bson:"budget_name,omitempty" json:"budget_name,omitempty"`
	BudgetValue *float64            `bson:"budget_value,omitempty" json:"budget_value,omitempty" default:"0"` // default doesn't give desired result
}

// FinancialAccount type is used to track balance record transactions
type FinancialAccount struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id, omitempty"`
	AccountName          string             `bson:"account_name,omitempty" json:"account_name,omitempty" validate:"required"`
	CurrentValue         float64            `bson:"current_value,omitempty" json:"current_value,omitempty" validate:"required"`
	FinancialInstitution string             `bson:"financial_institution,omitempty" json:"financial_institution,omitempty" validate:"required"`
	Owner                *string            `bson:"owner,omitempty" json:"owner,omitempty" validate:"required"`
	CreatedAt            time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"datetime"`
	UpdatedAt            time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" validate:"datetime"`
}
