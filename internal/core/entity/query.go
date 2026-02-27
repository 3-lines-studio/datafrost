package entity

type QueryResult struct {
	Columns []string `json:"columns"`
	Rows    [][]any  `json:"rows"`
	Count   int      `json:"count"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}

type Filter struct {
	ID       string `json:"id"`
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type QueryRequest struct {
	Query string `json:"query"`
}
