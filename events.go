package mqueue

// Event is a type that represents the different events that can be published.
type Event string

const (
	EventLoanUpdate        Event = "loan_update"
	EventLoanCreate        Event = "loan_create"
	EventLoanProductUpdate Event = "loan_product_update"
)

// AllEvents is a list of all events that can be published.
var AllEvents = []Event{
	EventLoanUpdate,
	EventLoanCreate,
	EventLoanProductUpdate,
}
