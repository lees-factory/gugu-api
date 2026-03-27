package response

type LoadHotProductsResponse struct {
	TotalFetched   int `json:"total_fetched"`
	NewlyCreated   int `json:"newly_created"`
	AlreadyExisted int `json:"already_existed"`
}
