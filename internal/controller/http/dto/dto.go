package dto

// Store
type StoreSubscriptionHandlerRequest struct {
	ServiceName string `json:"service_name" validate:"required"`
	Price       string `json:"price" validate:"required"`
	UserId      string `json:"user_id" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date"`
}

type StoreSubscriptionHandlerResponse struct {
	Id string `json:"id"`
}

//--------------------------------------------------------------------------

// Get
type GetSubscriptionHandlerRequest struct {
	Id string `json:"id"`
}

type GetSubscriptionHandlerResponse struct {
	ServiceName string `json:"service_name" validate:"required"`
	Price       string `json:"price" validate:"required"`
	UserId      string `json:"user_id" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date"`
}

//--------------------------------------------------------------------------

// Update
type UpdateSubscriptionHandlerRequest struct {
	ServiceName string `json:"service_name"`
	Price       string `json:"price"`
	UserId      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

//--------------------------------------------------------------------------

// Delete
type DeleteSubscriptionHandlerRequest struct {
	Id string `json:"id"`
}

//--------------------------------------------------------------------------

// List
type ListSubscriptionsHandlerRequest struct {
	// pagination
	Page     int `query:"page" validate:"min=1"`
	PageSize int `query:"page_size" validate:"min=1,max=100"`

	// filters
	ServiceName *string    `query:"service_name"`
	UserID      *string `query:"user_id"`
	Price       *string    `query:"price_min"`
	StartDate   *string    `query:"start_date"` // format: YYYY-MM-DD
	EndDate     *string    `query:"end_date"`   // format: YYYY-MM-DD

	// sort
	SortBy    string `query:"sort_by"`    // field name
	SortOrder string `query:"sort_order"` // asc, desc
}

type ListSubscriptionsHandlerResponse struct {
	Subscriptions []SubscriptionItem `json:"subscriptions"`
	Total         int                `json:"total"`
	Page          int                `json:"page"`
	PageSize      int                `json:"page_size"`
	TotalPages    int                `json:"total_pages"`
}

type SubscriptionItem struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Price       string `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

//--------------------------------------------------------------------------

//Error responce

type ErrorResponse struct {
	Error string `json:"error" example:"message"`
}

//--------------------------------------------------------------------------

// Total cost
type TotalCostHandlerRequest struct {
	UserID      *string `query:"user_id"`
	ServiceName *string `query:"service_name"`
	StartDate   *string `query:"start_date" validate:"required"` // format: YYYY-MM-DD
	EndDate     *string `query:"end_date" validate:"required"`   // format: YYYY-MM-DD
}

type TotalCostHandlerResponse struct {
	TotalCost uint64           `json:"total_cost"`
	Period    Period           `json:"period"`
	Filters   TotalCostFilters `json:"filters"`
}

type Period struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type TotalCostFilters struct {
	UserID      *string `json:"user_id,omitempty"`
	ServiceName *string `json:"service_name,omitempty"`
}
