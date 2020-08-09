package viewmodel

// NameVM ...
type NameVM struct {
	ID string `json:"id"`
	EN string `json:"en"`
}

// RedisStringValueVM ...
type RedisStringValueVM struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// LoginVM ....
type LoginVM struct {
	Token       string `json:"token"`
	ExpiredDate string `json:"expired_date"`
}
