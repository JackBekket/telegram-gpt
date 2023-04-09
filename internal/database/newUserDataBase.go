package database

// main database for dialogs, key (int64) is telegram user id
type User struct {
	ID            int64
	Username      string
	Dialog_status int64
	Gpt_key       string
	//gpt_client gpt3.Client
}

var UserMap = make(map[int64]User)

// func NewUserDataBase() *map[int64]*User {
// 	db := make(map[int64]*User)
// 	return &db
// }
