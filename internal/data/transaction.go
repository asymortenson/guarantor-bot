package data

type InMessage struct {
	Message string `json:"message,omitempty"`
	Value   string `json:"value,omitempty"`
}

type Transaction struct {
	InMessage InMessage `json:"in_msg,omitempty"`
}

type Response struct {
	Ok     bool          `json:"ok"`
	Result []Transaction `json:"result,omitempty"`
	Error  string        `json:"error,omitempty"`
	Code   int32         `json:"code,omitempty"`
}
