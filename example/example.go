package main

import (
	"database/sql"
	"github.com/brietsparks/xcrud/cli"
	"github.com/brietsparks/xcrud/data"
)

func main() {
	// replace with your variables
	url := data.MakeUrl(data.Vars{
		Host: "my-db-host.com",
		Name: "database_name",
		User: "postgres",
		Password: "password1234",
		Port: "5432",
	})

	db, _ := sql.Open("postgres", url)
	store, _ := data.NewStore(db, 10)

	// create a user
	createdUser, _ := store.CreateUser(&data.User{FirstName: "Bo", LastName: "Peep"})
	_ = cli.Printed(createdUser) // {"id":1,"firstName":"Bo","lastName":"Peep"}

	// get a user
	retrievedUser, _ := store.GetUserById(createdUser.Id)
	_ = cli.Printed(retrievedUser) // {"id":1,"firstName":"Bo","lastName":"Peep"}

	// update a user
	_ = store.UpdateUser(retrievedUser.Id, &data.User{LastName: "Jackson"}, "LastName")

	// create a group
	createdGroup, _ := store.CreateGroup(&data.Group{Name: "groupA"})
	_ = cli.Printed(createdGroup) // {"id":1,"Name":"groupA"}

	// add a user to a group
	_ = store.LinkGroupToUser(createdGroup.Id, createdUser.Id)

	// get users by group ID
	retrievedUsers, _ := store.GetUsersByGroupId(createdGroup.Id)
	_ = cli.Printed(retrievedUsers) // [{"id":1,"firstName":"Bo","lastName":"Jackson"}]

	// get groups by user ID
	retrievedGroups, _ := store.GetGroupsByUserId(createdUser.Id)
	_ = cli.Printed(retrievedGroups) // [{"id":1,"Name":"groupA"}]

	// remove a user from a group
	_ = store.UnlinkGroupFromUser(createdGroup.Id, createdUser.Id)
}
