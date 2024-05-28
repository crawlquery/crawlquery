package dto

type QueryResultPage struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	HtmlLink string `json:"html_link"`
}

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Pages []QueryResultPage `json:"pages"`
}
