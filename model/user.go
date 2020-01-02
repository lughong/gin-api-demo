package model

// User 结构体
type User struct {
	id       int
	username string
	password string
	age      int
}

// NewUser 获取一个User结构体指针
func NewUser(id int, username, password string, age int) *User {
	return &User{
		id:       id,
		username: username,
		password: password,
		age:      age,
	}
}

func (u *User) GetID() int {
	return u.id
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetAge() int {
	return u.age
}
