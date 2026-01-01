package requestx

type GetFinanceDataRequest struct {
	Type string `json:"type" validate:"required,oneof=1m 6m 12m all"`
}
