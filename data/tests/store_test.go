package tests

import (
	"database/sql"
	"flag"
	"github.com/brietsparks/xcrud/data"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"gopkg.in/testfixtures.v2"
	"testing"
)

type StoreTestSuite struct {
	suite.Suite
	Store    *data.Store
	Fixtures *testfixtures.Context
}

var envPath string

func init() {
	flag.StringVar(&envPath, "env", "", "")
}

func (s *StoreTestSuite) SetupSuite() {
	if envPath == "" {
		s.T().Fatal("missing variable --env <path to .env file>")
	}

	vars, err := data.LoadEnvVars(envPath)

	if err != nil {
		s.T().Fatalf("failed to load environment variables: %s", err)
	}

	store, err := data.NewStore(vars)

	if err != nil {
		s.T().Fatalf("failed to create store: %s", err)
	}

	s.Store = store

	err = clearTables(vars)

	if err != nil {
		s.T().Fatalf("failed to clear table: %s", err)
	}
}

func (s *StoreTestSuite) TearDownSuite() {
	if envPath == "" {
		s.T().Fatal("missing variable --env <path to .env file>")
	}

	vars, err := data.LoadEnvVars(envPath)

	if err != nil {
		s.T().Fatalf("failed to load environment variables: %s", err)
	}

	err = clearTables(vars)

	if err != nil {
		s.T().Fatalf("failed to clear table: %s", err)
	}
}

func clearTables(vars data.Vars) error {
	url := data.MakeUrl(vars)
	db, err := sql.Open("postgres", url)

	if err != nil {
		return err
	}

	_, err = db.Query(`
		truncate table users cascade;
		truncate table groups cascade;
	`)

	return err
}

func (s *StoreTestSuite) TestUser() {
	input := &data.User{
		FirstName: data.RandomString(5),
		LastName:  data.RandomString(5),
	}

	created, err := s.Store.CreateUser(input)

	if err != nil {
		s.T().Errorf("test failed to create user: %w", err)
	}

	retrieved, err := s.Store.GetUserById(created.Id)

	if err != nil {
		s.T().Errorf("test failed to get user: %w", err)
	}

	s.Assert().EqualValues(input, retrieved)

	update := &data.User{FirstName: "asdf"}

	err = s.Store.UpdateUser(created.Id, update, "FirstName")

	if err != nil {
		s.T().Errorf("test failed to update user: %w", err)
	}

	retrieved, err = s.Store.GetUserById(created.Id)

	if err != nil {
		s.T().Errorf("test failed to get user: %w", err)
	}

	s.Assert().Equal(retrieved.FirstName, update.FirstName)
	s.Assert().NotEqual(retrieved.LastName, update.LastName)

	err = s.Store.DeleteUser(created.Id)

	if err != nil {
		s.T().Errorf("test failed to delete user: %w", err)
	}

	retrieved, err = s.Store.GetUserById(created.Id)

	if err != nil {
		s.T().Errorf("test failed to get user: %w", err)
	}

	s.Assert().Nil(retrieved)
}

func (s *StoreTestSuite) TestGroup() {
	input := &data.Group{
		Name: data.RandomString(5),
	}

	created, err := s.Store.CreateGroup(input)

	if err != nil {
		s.T().Errorf("test failed to create group: %w", err)
	}

	retrieved, err := s.Store.GetGroupById(created.Id)

	if err != nil {
		s.T().Errorf("test failed to get group: %w", err)
	}

	s.Assert().EqualValues(input, retrieved)

	update := &data.Group{Name: "asdf"}

	err = s.Store.UpdateGroup(created.Id, update, "Name")

	if err != nil {
		s.T().Errorf("test failed to update group: %w", err)
	}

	retrieved, err = s.Store.GetGroupById(created.Id)

	if err != nil {
		s.T().Errorf("test failed to get group: %w", err)
	}

	s.Assert().Equal(retrieved.Name, update.Name)

	err = s.Store.DeleteGroup(created.Id)

	if err != nil {
		s.T().Errorf("test failed to delete group: %w", err)
	}

	retrieved, err = s.Store.GetGroupById(created.Id)

	if err != nil {
		s.T().Errorf("test failed to get group: %w", err)
	}

	s.Assert().Nil(retrieved)
}

func (s *StoreTestSuite) TestGroupUser() {
	// create the resources
	u1, err := s.Store.CreateUser(&data.User{FirstName: data.RandomString(5), LastName: data.RandomString(5),})
	if err != nil {
		s.T().Errorf("test failed to create user: %w", err)
	}

	u2, err := s.Store.CreateUser(&data.User{FirstName: data.RandomString(5), LastName: data.RandomString(5),})
	if err != nil {
		s.T().Errorf("test failed to create user: %w", err)
	}

	g1, err := s.Store.CreateGroup(&data.Group{Name: data.RandomString(5),})
	if err != nil {
		s.T().Errorf("test failed to create group: %w", err)
	}

	g2, err := s.Store.CreateGroup(&data.Group{Name: data.RandomString(5),})
	if err != nil {
		s.T().Errorf("test failed to create group: %w", err)
	}

	// link the resources
	err = s.Store.LinkGroupToUser(g1.Id, u1.Id)
	if err != nil {
		s.T().Errorf("test failed to link group to user: %w", err)
	}

	err = s.Store.LinkGroupToUser(g2.Id, u1.Id)
	if err != nil {
		s.T().Errorf("test failed to link group to user: %w", err)
	}

	err = s.Store.LinkGroupToUser(g2.Id, u2.Id)
	if err != nil {
		s.T().Errorf("test failed to link group to user: %w", err)
	}

	if err != nil {
		s.T().Errorf("test failed to link group to user: %w", err)
	}

	// verify
	users, err := s.Store.GetUsersByGroupId(g1.Id)
	if err != nil {
		s.T().Errorf("test failed to get users by group id: %w", err)
	}
	s.Assert().EqualValues([]data.User{*u1}, users)

	users, err = s.Store.GetUsersByGroupId(g2.Id)
	if err != nil {
		s.T().Errorf("test failed to get users by group id: %w", err)
	}
	s.Assert().EqualValues([]data.User{*u1, *u2}, users)

	groups, err := s.Store.GetGroupsByUserId(u1.Id)
	if err != nil {
		s.T().Errorf("test failed to get groups by user id: %w", err)
	}
	s.Assert().EqualValues([]data.Group{*g1, *g2}, groups)

	groups, err = s.Store.GetGroupsByUserId(u2.Id)
	if err != nil {
		s.T().Errorf("test failed to get groups by user id: %w", err)
	}
	s.Assert().EqualValues([]data.Group{*g2}, groups)

	// unlink the resource
	err = s.Store.UnlinkGroupFromUser(g2.Id, u1.Id)

	// verify
	users, err = s.Store.GetUsersByGroupId(g2.Id)
	if err != nil {
		s.T().Errorf("test failed to get users by group id: %w", err)
	}
	s.Assert().EqualValues([]data.User{*u2}, users)

	groups, err = s.Store.GetGroupsByUserId(u1.Id)
	if err != nil {
		s.T().Errorf("test failed to get groups by user id: %w", err)
	}
	s.Assert().EqualValues([]data.Group{*g1}, groups)

}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}
