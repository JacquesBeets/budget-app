package models

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Admin    bool   `json:"admin"`
}

func (u *User) Save() error {
	return nil
}

func (u *User) Update() error {
	return nil
}

func (u *User) Delete() error {
	return nil
}
