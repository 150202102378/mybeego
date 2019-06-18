package models

import (
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id      int
	Name    string
	Profile *Profile `orm:"rel(one)"`
	Post    []*Post  `orm:"reverse(many)"`
}

type Profile struct {
	Id   int
	Age  int16
	User *User `orm:"reverse(one)"`
}

type Post struct {
	Id    int
	Title string
	User  *User `orm:"rel(fk)"`
}

type Tag struct {
	Id    int
	Name  string
	Posts []*Post `orm:"reverse(many)"`
}

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "root:root@/orm_test?charset=utf8")
	orm.RegisterModel(new(User), new(Post), new(Profile), new(Tag))
}
