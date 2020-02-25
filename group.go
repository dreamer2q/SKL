package skl

import (
	"fmt"
)

func NewGroup() *Group {
	return &Group{}
}

func Groups() []*Group {
	groups := make([]*Group, 0, 10)
	DB.Preload("Users").Find(&groups)
	return groups
}

func (g *Group) BeforeSave() error {
	if g.GroupID == 0 || g.TeacherID == 0 {
		return fmt.Errorf("group should not be empty")
	}
	return nil
}

func (g *Group) AddUser(u *User) *Group {
	g.Users = append(g.Users, u)
	return g
}

func (g *Group) Update() error {
	return DB.Model(&Group{}).Save(g).Error
}

func GetGroup(groupId int64) *Group {
	g := &Group{}
	DB.Preload("Users").First(g, groupId)
	return g
}
