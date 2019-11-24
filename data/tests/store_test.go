package tests

import (
	"database/sql"
	"flag"
	"github.com/brietsparks/xcrud/data"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"gopkg.in/testfixtures.v2"
	"log"
	"testing"
)

type StoreTestSuite struct {
	suite.Suite
	Store    *data.Store
	fixtures *testfixtures.Context
}

var envPath string

func init() {
	flag.StringVar(&envPath, "env", "", "")
}

func (s *StoreTestSuite) SetupSuite() {
	d := connect(s)

	// load fixtures
	fixtures, err := testfixtures.NewFolder(d, &testfixtures.PostgreSQL{}, "../fixtures")
	if err != nil {
		log.Fatal(err)
	}

	s.fixtures = fixtures

	// create store
	store, err := data.NewStore(d, 10)

	if err != nil {
		s.T().Fatalf("failed to create store: %s", err)
	}

	s.Store = store

	// clear tables
	err = clearTables(d)

	if err != nil {
		s.T().Fatalf("failed to clear table: %s", err)
	}
}

func (s *StoreTestSuite) TearDownSuite() {
	d := connect(s)

	err := clearTables(d)

	if err != nil {
		s.T().Fatalf("failed to clear table: %s", err)
	}
}

func connect(s *StoreTestSuite) *sql.DB {
	if envPath == "" {
		s.T().Fatal("missing variable --env <path to .env file>")
	}

	vars, err := data.LoadEnvVars(envPath)

	if err != nil {
		s.T().Fatalf("failed to load environment variables: %s", err)
	}

	url := data.MakeUrl(vars)
	d, err := sql.Open("postgres", url)

	if err != nil {
		s.T().Fatalf("failed to connect to database: %s", err)
	}

	return d
}

func clearTables(db *sql.DB) error {
	_, err := db.Query(`
		truncate table users cascade;
		truncate table groups cascade;
	`)

	return err
}

func (s *StoreTestSuite) SetupTest() {
	testfixtures.ResetSequencesTo(1)

	if err := s.fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func (s *StoreTestSuite) TestGetUserById() {
	u, _ := s.Store.GetUserById(100)

	expected := &data.User{
		Id:        100,
		FirstName: "A",
		LastName:  "B",
	}
	s.Assert().EqualValues(expected, u)

	u, _ = s.Store.GetUserById(1000)
	s.Assert().Nil(u)
}

func (s *StoreTestSuite) TestUpdateUser() {
	_ = s.Store.UpdateUser(101, &data.User{FirstName: "abc",})

	u, _ := s.Store.GetUserById(101)
	expected := &data.User{
		Id:        101,
		FirstName: "C",
		LastName:  "D",
	}
	s.Assert().EqualValues(expected, u)

	_ = s.Store.UpdateUser(101, &data.User{FirstName: "abc",}, "FirstName")
	u, _ = s.Store.GetUserById(101)
	expected = &data.User{
		Id:        101,
		FirstName: "abc",
		LastName:  "D",
	}
	s.Assert().EqualValues(expected, u)
}

func (s *StoreTestSuite) TestCreateUser() {
	created, _ := s.Store.CreateUser(&data.User{
		FirstName: "foo",
		LastName:  "bar",
	})

	retrieved, _ := s.Store.GetUserById(created.Id)

	s.Assert().EqualValues(created, retrieved)
}

func (s *StoreTestSuite) TestDeleteUser() {
	retrieved, _ := s.Store.GetUserById(102)
	s.Assert().NotNil(retrieved)

	_ = s.Store.DeleteUser(102)
	retrieved, _ = s.Store.GetUserById(102)
	s.Assert().Nil(retrieved)
}

func (s *StoreTestSuite) TestGetGroupById() {
	u, _ := s.Store.GetGroupById(100)

	expected := &data.Group{
		Id:   100,
		Name: "A",
	}
	s.Assert().EqualValues(expected, u)

	u, _ = s.Store.GetGroupById(1000)
	s.Assert().Nil(u)
}

func (s *StoreTestSuite) TestUpdateGroup() {
	_ = s.Store.UpdateGroup(101, &data.Group{Name: "abc",})

	u, _ := s.Store.GetGroupById(101)
	expected := &data.Group{
		Id:   101,
		Name: "B",
	}
	s.Assert().EqualValues(expected, u)

	_ = s.Store.UpdateGroup(101, &data.Group{Name: "abc",}, "Name")
	u, _ = s.Store.GetGroupById(101)
	expected = &data.Group{
		Id:   101,
		Name: "abc",
	}
	s.Assert().EqualValues(expected, u)
}

func (s *StoreTestSuite) TestCreateGroup() {
	created, _ := s.Store.CreateGroup(&data.Group{
		Name: "foo",
	})

	retrieved, _ := s.Store.GetGroupById(created.Id)

	s.Assert().EqualValues(created, retrieved)
}

func (s *StoreTestSuite) TestDeleteGroup() {
	retrieved, _ := s.Store.GetGroupById(102)
	s.Assert().NotNil(retrieved)

	_ = s.Store.DeleteGroup(102)
	retrieved, _ = s.Store.GetGroupById(102)
	s.Assert().Nil(retrieved)
}

func (s *StoreTestSuite) TestGetUsersByGroupId() {
	users, _ := s.Store.GetUsersByGroupId(201)

	expected := []data.User{
		{Id: 201, FirstName: "I", LastName: "J"},
		{Id: 202, FirstName: "K", LastName: "L"},
	}

	s.Assert().Equal(expected, users)
}

func (s *StoreTestSuite) TestGetGroupsByUserId() {
	groups, _ := s.Store.GetGroupsByUserId(202)

	expected := []data.Group{
		{Id: 201, Name: "E"},
		{Id: 202, Name: "F"},
	}

	s.Assert().Equal(expected, groups)
}

func (s *StoreTestSuite) TestLinkGroupToUser() {
	_ = s.Store.LinkGroupToUser(200, 200)

	users, _ := s.Store.GetUsersByGroupId(200)
	expectedUsers := []data.User{{Id: 200, FirstName: "G", LastName: "H"}}
	s.Assert().EqualValues(expectedUsers, users)

	groups, _ := s.Store.GetGroupsByUserId(200)
	expectedGroups := []data.Group{{Id: 200, Name: "D"}}
	s.Assert().EqualValues(expectedGroups, groups)
}

func (s *StoreTestSuite) TestUnlinkGroupFromUser() {
    _ = s.Store.UnlinkGroupFromUser(203, 203)

    users, _ := s.Store.GetUsersByGroupId(203)
	s.Assert().Nil(users)

	groups, _ := s.Store.GetGroupsByUserId(203)
	s.Assert().Nil(groups)
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}
