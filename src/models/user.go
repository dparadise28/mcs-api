package models

type UserCart struct {
	UserId  string `json:"user_id"`
	StoreId string `json:"store_id"`
	CartId  string `json:"cart_id"`
}

type User struct {
	UserName string     `json:"username"`
	Password string     `json:"password"`
	Carts    []UserCart `json:"carts"`
	Email    string     `json:"email"`
}
