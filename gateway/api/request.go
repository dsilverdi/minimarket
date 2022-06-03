package api

type UserRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CommentRequestBody struct {
	ProductID int    `json:"product_id"`
	Message   string `json:"message"`
	ReplyID   int    `json:"reply_id,omitempty"`
}
