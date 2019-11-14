package data

type User struct {
	Id        int64  `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"firstName" validate:"required,lte=100"`
	LastName  string `db:"last_name" json:"lastName" validate:"required,lte=100"`
}

type Group struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name" validate:"required,lte=100"`
}
