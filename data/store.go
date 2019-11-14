package data

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	"gopkg.in/go-playground/validator.v9"
)

type Store struct {
	db       *dbr.Session
	validate *validator.Validate
}

func NewStore(vars Vars) (*Store, error) {
	url := MakeUrl(vars)

	// create a connection (e.g. "postgres", "mysql", or "sqlite3")
	conn, _ := dbr.Open("postgres", url, nil)
	conn.SetMaxOpenConns(10)

	// create a session for each business unit of execution (e.g. a web request or goworkers job)
	sess := conn.NewSession(nil)

	// create a tx from sessions
	_, err := sess.Begin()

	if err != nil {
		return nil, fmt.Errorf("unable to create data store: %w", err)
	}

	v := validator.New()

	return &Store{
		db:       sess,
		validate: v,
	}, nil
}

func (s *Store) CreateUser(u *User) (*User, error) {
	if err := s.validate.Struct(u); err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	columns := []string{"first_name", "last_name",}
	id, err := s.create("users", u, columns)

	if id != nil {
		u.Id = id.(int64)
	}

	return u, err
}

func (s *Store) UpdateUser(id int64, u *User, fields ...string) error {
	if err := s.validate.StructPartial(u, fields...); err != nil {
		return fmt.Errorf("invalid user data: %w", err)
	}

	return s.update("users", id, fields,
		set{"FirstName", "first_name", u.FirstName},
		set{"LastName", "last_name", u.LastName},
	)
}

func (s *Store) GetUserById(id int64) (*User, error) {
	u := &User{}
	retrieved, err := s.getById("users", id, u)

	if retrieved == nil {
		return nil, err
	}

	return retrieved.(*User), err
}

func (s *Store) DeleteUser(id int64) error {
    _, err := s.db.DeleteFrom("users").Where("id = ?", id).Exec()

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

    return nil
}

func (s *Store) CreateGroup(g *Group) (*Group, error) {
	if err := s.validate.Struct(g); err != nil {
		return nil, fmt.Errorf("invalid group data: %w", err)
	}

	columns := []string{"name",}
	id, err := s.create("groups", g, columns)
	g.Id = id.(int64)

	return g, err
}

func (s *Store) UpdateGroup(id int64, g *Group, fields ...string) error {
	if err := s.validate.StructPartial(g, fields...); err != nil {
		return fmt.Errorf("invalid group data: %w", err)
	}

	return s.update("groups", id, fields,
		set{"Name", "name", g.Name},
	)
}

func (s *Store) GetGroupById(id int64) (*Group, error) {
	g := &Group{}
	retrieved, err := s.getById("groups", id, g)

	if retrieved == nil {
		return nil, err
	}

	return retrieved.(*Group), err
}

func (s *Store) DeleteGroup(id int64) error {
	_, err := s.db.DeleteFrom("groups").Where("id = ?", id).Exec()

	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}

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
		return nil, fmt.Errorf("failed to get users by groupId: %v", err)
	}

	return users, nil
}

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
		return nil, fmt.Errorf("failed to get groups by userId: %v", err)
	}

	return groups, nil
}

func (s *Store) LinkGroupToUser(groupId int64, userId int64) error {
	_, err := s.db.
		InsertInto("groups_users").
		Pair("group_id", groupId).
		Pair("user_id", userId).
		Exec()

	return err
}


func (s *Store) UnlinkGroupFromUser(groupId int64, userId int64) error {
	_, err := s.db.
		DeleteFrom("groups_users").
		Where("group_id = ? and user_id = ?", groupId, userId).
		Exec()

	return err
}
