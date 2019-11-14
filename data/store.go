package data

import (
	"database/sql"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"gopkg.in/go-playground/validator.v9"
)

type Store struct {
	db       *dbr.Session
	validate *validator.Validate
}

func NewStore(d *sql.DB, maxConn int) (*Store, error) {
	conn := &dbr.Connection{
		DB: d,
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect: dialect.PostgreSQL,
	}

	conn.SetMaxOpenConns(maxConn)
	sess := conn.NewSession(nil)
	_, err := sess.Begin()

	if err != nil {
		return nil, NewError(err, "unable to create data store")
	}

	v := validator.New()

	return &Store{
		db:       sess,
		validate: v,
	}, nil
}

// CreateUser creates a new user
func (s *Store) CreateUser(u *User) (*User, error) {
	if err := s.validate.Struct(u); err != nil {
		return nil, NewError(err, "invalid user data")
	}

	columns := []string{"first_name", "last_name",}
	id, err := s.create("users", u, columns)

	if err != nil {
		return nil, NewError(err, "failed to create user")
	}

	if id != nil {
		u.Id = id.(int64)
	}

	return u, nil
}

// UpdateUser updates an existing user.
// The variadic "fields" arg should contain the field names that should be updated
func (s *Store) UpdateUser(id int64, u *User, fields ...string) error {
	if err := s.validate.StructPartial(u, fields...); err != nil {
		return NewError(err, "invalid user data")
	}

	err :=  s.update("users", id, fields,
		set{"FirstName", "first_name", u.FirstName},
		set{"LastName", "last_name", u.LastName},
	)

	if err != nil {
		return NewError(err, "failed to update user")
	}

	return nil
}

// GetUserById gets a user by ID
func (s *Store) GetUserById(id int64) (*User, error) {
	u := &User{}
	retrieved, err := s.getById("users", id, u)

	if err != nil {
		return nil, NewError(err, "failed to get user by id")
	}

	if retrieved == nil {
		return nil, nil
	}

	return retrieved.(*User), err
}

// DeleteUser deletes a user
func (s *Store) DeleteUser(id int64) error {
    _, err := s.db.DeleteFrom("users").Where("id = ?", id).Exec()

	if err != nil {
		return NewError(err, "failed to delete user")
	}

    return nil
}

// CreateGroup creates a new group
func (s *Store) CreateGroup(g *Group) (*Group, error) {
	if err := s.validate.Struct(g); err != nil {
		return nil, NewError(err, "invalid group data")
	}

	columns := []string{"name",}
	id, err := s.create("groups", g, columns)

	if err != nil {
		return nil, NewError(err, "failed to create group")
	}

	if id != nil {
		g.Id = id.(int64)
	}

	return g, nil
}

// UpdateGroup updates an existing group
// The variadic "fields" arg should contain the field names that should be updated
func (s *Store) UpdateGroup(id int64, g *Group, fields ...string) error {
	if err := s.validate.StructPartial(g, fields...); err != nil {
		return NewError(err, "invalid group data")
	}

	err := s.update("groups", id, fields,
		set{"Name", "name", g.Name},
	)

	if err != nil {
		return NewError(err, "failed to update group")
	}

	return nil
}

// GetGroupById gets a group by ID
func (s *Store) GetGroupById(id int64) (*Group, error) {
	g := &Group{}
	retrieved, err := s.getById("groups", id, g)

	if err != nil {
		return nil, NewError(err, "failed to get group by id")
	}

	if retrieved == nil {
		return nil, err
	}

	return retrieved.(*Group), err
}

// DeleteUser deletes a group
func (s *Store) DeleteGroup(id int64) error {
	_, err := s.db.DeleteFrom("groups").Where("id = ?", id).Exec()

	if err != nil {
		return NewError(err, "failed to delete group")
	}

	return nil
}

// GetUsersByGroupId returns an array of users that belong to a group
func (s *Store) GetUsersByGroupId(groupId int64) ([]User, error) {
	var users []User

	_, err := s.selectJunction(s.db, groupId, junction{
		table1: "users",
		table2: "groups",
		junctionTable: "groups_users",
		junctionFk1: "user_id",
		junctionFk2: "group_id",
	}).Load(&users)

	if err != nil {
		return nil, NewError(err, "failed to get users by groupId")
	}

	return users, nil
}

// GetUsersByGroupId returns an array of groups that contain a user
func (s *Store) GetGroupsByUserId(userId int64) ([]Group, error) {
    var groups []Group

	_, err := s.selectJunction(s.db, userId, junction{
		table1: "groups",
		table2: "users",
		junctionTable: "groups_users",
		junctionFk1: "group_id",
		junctionFk2: "user_id",
	}).Load(&groups)

	if err != nil {
		return nil, NewError(err, "failed to get groups by userId")
	}

	return groups, nil
}

// LinkGroupToUser links a group to a user
func (s *Store) LinkGroupToUser(groupId int64, userId int64) error {
	_, err := s.db.
		InsertInto("groups_users").
		Pair("group_id", groupId).
		Pair("user_id", userId).
		Exec()

	if err != nil {
		return NewError(err, "failed to link group to user")
	}

	return nil
}

// UnlinkGroupFromUser unlinks a group from a user
func (s *Store) UnlinkGroupFromUser(groupId int64, userId int64) error {
	_, err := s.db.
		DeleteFrom("groups_users").
		Where("group_id = ? and user_id = ?", groupId, userId).
		Exec()

	if err != nil {
		return NewError(err, "failed to unlink group from user")
	}

	return nil
}
