package model

// 定义SQL语句常量，方便统计、查看SQL语句
const (
	selectUserSQL                  = "SELECT * FROM user"
	selectByUsernameAndPasswordSQL = " WHERE username=? and password=?"
	selectUserDetail               = selectUserSQL + selectByUsernameAndPasswordSQL
)

// 定义User结构体
type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Age      int    `json:"age" db:"age"`
}

// 定义ModOption函数类型，用于修改默认值
type ModOption func(u *User)

// NewUser 获取一个User结构体指针
func NewUser(modOptions ...ModOption) *User {
	u := User{
		Id:       0,
		Username: "",
		Password: "",
		Age:      0,
	}

	for _, f := range modOptions {
		f(&u)
	}

	return &u
}

// Find 获取一条user数据
func (u *User) Find() (User, error) {
	var user User
	stmt, err := DB.Conn.Preparex(selectUserDetail)
	if err != nil {
		return user, err
	}

	// 把查找结果解析到User结构体
	if err := stmt.QueryRowx(u.Username, u.Password).StructScan(&user); err != nil {
		return user, err
	}

	return user, nil
}
