package budget

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Budget type is a group of related financial transactions
type Budget struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	ManagerID   string             `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	BudgetName  string             `bson:"budget_name,omitempty" json:"budget_name,omitempty" validate:"required"`
	BudgetValue float64            `bson:"budget_value,omitempty" json:"budget_value,omitempty" default:"0"` // default doesn't give desired result
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"datetime"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty" validate:"datetime"`
}

// NewBudget type is what's required from the client to create a new Budget
type NewBudget struct {
	ManagerID   string  `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
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
	ManagerID   *string             `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	BudgetName  *string             `bson:"budget_name,omitempty" json:"budget_name,omitempty"`
	BudgetValue *float64            `bson:"budget_value,omitempty" json:"budget_value,omitempty" default:"0"` // default doesn't give desired result
}

// FinancialAccount type is used to track balance record transactions
type FinancialAccount struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AccountName          string             `bson:"account_name,omitempty" json:"account_name,omitempty" validate:"required"`
	CurrentValue         float64            `bson:"current_value,omitempty" json:"current_value,omitempty" validate:"required"`
	FinancialInstitution string             `bson:"financial_institution,omitempty" json:"financial_institution,omitempty" validate:"required"`
	MangerID             string             `bson:"manger_id,omitempty" json:"manger_id,omitempty" validate:"required"`
	CreatedAt            time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"datetime"`
	UpdatedAt            time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty" validate:"datetime"`
}

// NewFinancialAccount type is used to track balance record transactions
type NewFinancialAccount struct {
	AccountName          string  `bson:"account_name,omitempty" json:"account_name,omitempty" validate:"required"`
	CurrentValue         float64 `bson:"current_value,omitempty" json:"current_value,omitempty" validate:"required"`
	FinancialInstitution string  `bson:"financial_institution,omitempty" json:"financial_institution,omitempty" validate:"required"`
	MangerID             string  `bson:"manger_id,omitempty" json:"manger_id,omitempty"`
}

// UpdateFinancialAccount defines what information may be provided to modify an existing Financial Account.
// All fields are optional so clients can send just the fields they want changed.
// It uses pointer fields so we can differentiate between a field that was NOT provided and a field that was provided as explicitly blank.
// Normally we do NOT want to use pointers to basic types but we make exceptions around marshalling/unmarshallying.
type UpdateFinancialAccount struct {
	ID                   *primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ManagerID            *string             `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	AccountName          *string             `bson:"account_name,omitempty" json:"account_name,omitempty"`
	CurrentValue         *float64            `bson:"current_value,omitempty" json:"current_value,omitempty"`
	FinancialInstitution *string             `bson:"financial_institution,omitempty" json:"financial_institution,omitempty"`
}

// Vendor type is a group of vendors that process transactions
type Vendor struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	TransactionIDs []string           `bson:"tranx_id,omitempty" json:"tranx_id,omitempty"`
	VendorName     string             `bson:"vendor_name,omitempty" json:"vendor_name,omitempty"`
	CreatedAt      time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"datetime"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty" validate:"datetime"`
}

// NewVendor type is what's required from the client to create a new vendor
type NewVendor struct {
	TransactionIDs *[]string `bson:"tranx_id,omitempty" json:"tranx_id,omitempty"`
	VendorName     string    `bson:"vendor_name,omitempty" json:"vendor_name,omitempty"`
}

// UpdateVendor defines what information may be provided to modify an existing Vendor.
// All fields are optional so clients can send just the fields they want
// changed.
// It uses pointer fields so we can differentiate between a field that was not provided and a field that was provided as explicitly blank.
// Normally we do not want to use pointers to basic types but we make exceptions around marshalling/unmarshalling.
type UpdateVendor struct {
	TransactionIDs *[]string `bson:"tranx_id,omitempty" json:"tranx_id,omitempty"`
	VendorName     *string   `bson:"vendor_name,omitempty" json:"vendor_name,omitempty"`
}

// Transaction type is used to track and manage financial transactions.
// For example, purchases, investments, etc.
type Transaction struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	BudgetID           string             `bson:"budget_id,omitempty" json:"budget_id,omitempty"`
	CurrencyID         string             `bson:"currency_id,omitempty" json:"currency_id,omitempty"`
	FinancialAccountID []string           `bson:"fin_acc_id,omitempty" json:"fin_acc_id,omitempty"`
	// Occurrence         time.Time            `bson:"occurrence,omitempty" json:"occurrence,omitempty" validate:"datetime"`
	OccurrenceString  string    `bson:"occurrence_string,omitempty" json:"occurrence_string,omitempty"`
	TransactionEvent  string    `bson:"tranx_event,omitempty" json:"tranx_event,omitempty"`
	TransactionCredit float64   `bson:"tranx_credit,omitempty" json:"tranx_credit,omitempty"`
	TransactionDebit  float64   `bson:"tranx_debit,omitempty" json:"tranx_debit,omitempty"`
	VendorID          string    `bson:"vendor_id,omitempty" json:"vendor_id,omitempty"`
	ParticipantID     []string  `bson:"participant_id,omitempty" json:"participant_id,omitempty"`
	CreatedAt         time.Time `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"datetime"`
	UpdatedAt         time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty" validate:"datetime"`
}

// NewTransaction type is what's required from the client to create a new transaction.
type NewTransaction struct {
	BudgetID           string    `bson:"budget_id,omitempty" json:"budget_id,omitempty"`
	CurrencyID         string    `bson:"currency_id,omitempty" json:"currency_id,omitempty"`
	FinancialAccountID *[]string `bson:"fin_acc_id,omitempty" json:"fin_acc_id,omitempty"`
	// Occurrence         *time.Time            `bson:"occurrence,omitempty" json:"occurrence,omitempty" validate:"datetime"`
	OccurrenceString  string    `bson:"occurrence_string,omitempty" json:"occurrence_string,omitempty"`
	TransactionEvent  string    `bson:"tranx_event,omitempty" json:"tranx_event,omitempty"`
	TransactionCredit float64   `bson:"tranx_credit,omitempty" json:"tranx_credit,omitempty"`
	TransactionDebit  float64   `bson:"tranx_debit,omitempty" json:"tranx_debit,omitempty"`
	VendorID          string    `bson:"vendor_id,omitempty" json:"vendor_id,omitempty"`
	ParticipantID     *[]string `bson:"participant_id,omitempty" json:"participant_id,omitempty"`
}

// UpdateTransaction defines what information may be provided to modify an existing Transaction.
// All fields are optional so clients can send just the fields they want
// changed.
// It uses pointer fields so we can differentiate between a field that was not provided and a field that was provided as explicitly blank.
// Normally we do not want to use pointers to basic types but we make exceptions around marshalling/unmarshalling.
type UpdateTransaction struct {
	BudgetID           *string   `bson:"budget_id,omitempty" json:"budget_id,omitempty"`
	CurrencyID         *string   `bson:"currency_id,omitempty" json:"currency_id,omitempty"`
	FinancialAccountID *[]string `bson:"fin_acc_id,omitempty" json:"fin_acc_id,omitempty"`
	// Occurrence         *time.Time            `bson:"occurrence,omitempty" json:"occurrence,omitempty" validate:"datetime"`
	OccurrenceString  *string   `bson:"occurrence_string,omitempty" json:"occurrence_string,omitempty"`
	TransactionEvent  *string   `bson:"tranx_event,omitempty" json:"tranx_event,omitempty"`
	TransactionCredit *float64  `bson:"tranx_credit,omitempty" json:"tranx_credit,omitempty"`
	TransactionDebit  *float64  `bson:"tranx_debit,omitempty" json:"tranx_debit,omitempty"`
	VendorID          *string   `bson:"vendor_id,omitempty" json:"vendor_id,omitempty"`
	ParticipantID     *[]string `bson:"participant_id,omitempty" json:"participant_id,omitempty"`
}

// FilterTransaction type is used to retrieve a filtered list of transactions.
type FilterTransaction struct {
	BudgetID           string `bson:"budget_id,omitempty" json:"budget_id,omitempty"`
	CurrencyID         string `bson:"currency_id,omitempty" json:"currency_id,omitempty"`
	FinancialAccountID string `bson:"fin_acc_id,omitempty" json:"fin_acc_id,omitempty"`
	// Occurrence         time.Time            `bson:"occurrence,omitempty" json:"occurrence,omitempty" validate:"datetime"`
	OccurrenceString  string    `bson:"occurrence_string,omitempty" json:"occurrence_string,omitempty"`
	TransactionEvent  string    `bson:"tranx_event,omitempty" json:"tranx_event,omitempty"`
	TransactionCredit string    `bson:"tranx_credit,omitempty" json:"tranx_credit,omitempty"`
	TransactionDebit  string    `bson:"tranx_debit,omitempty" json:"tranx_debit,omitempty"`
	VendorID          string    `bson:"vendor_id,omitempty" json:"vendor_id,omitempty"`
	ParticipantID     string    `bson:"participant_id,omitempty" json:"participant_id,omitempty"`
	CreatedAt         time.Time `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"datetime"`
	UpdatedAt         time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty" validate:"datetime"`
}

// Currency type is a group of currencies
type Currency struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" validate:"required"`
	CurrencyName string             `bson:"currency,omitempty" json:"currency,omitempty"`
	CurrencyType string             `bson:"curr_type,omitempty" json:"curr_type,omitempty"`
	Symbol       string             `bson:"symbol,omitempty" json:"symbol,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty" validate:"datetime"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty" validate:"datetime"`
}

// NewCurrency type is what's required from the client to create a new Currency
type NewCurrency struct {
	CurrencyName string `bson:"currency,omitempty" json:"currency,omitempty"`
	CurrencyType string `bson:"curr_type,omitempty" json:"curr_type,omitempty"`
	Symbol       string `bson:"symbol,omitempty" json:"symbol,omitempty"`
}

// UpdateCurrency defines what information may be provided to modify an existing Currency.
// All fields are optional so clients can send just the fields they want changed.
// It uses pointer fields so we can differentiate between a field that was not provided and a field that was provided as explicitly blank.
// Normally we do not want to use pointers to basic types but we make exceptions around marshalling/unmarshalling.
type UpdateCurrency struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CurrencyName *string            `bson:"currency,omitempty" json:"currency,omitempty"`
	CurrencyType *string            `bson:"curr_type,omitempty" json:"curr_type,omitempty"`
	Symbol       *string            `bson:"symbol,omitempty" json:"symbol,omitempty"`
}
