package skl

import "fmt"

func Users() []*User {
	users := make([]*User, 0, 10)
	DB.Preload("Groups").Find(&users)
	return users
}

func (u *User) BeforeSave() error {
	if u.QQ == 0 || u.Token == "" || u.UserID == "" {
		return fmt.Errorf("user should not be empty")
	}
	return nil
}

func (u *User) Update() error {
	return DB.Model(&User{}).Save(u).Error
}

func (u *User) AddGroup(g *Group) *User {
	u.Groups = append(u.Groups, g)
	return u
}

func GetUser(userId int64) (u *User) {
	DB.Preload("Groups").Find(u, userId)
	return
}
