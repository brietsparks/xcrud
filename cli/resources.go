package cli

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brietsparks/xcrud/data"
	"github.com/urfave/cli"
	"strconv"
)

// NewResourcesCommand returns a resources command tree that can be used by a urfave/cli instance
func NewResourcesCommand(name string, chVars chan data.Vars, logger Logger) cli.Command {
	var vars data.Vars
	var store *data.Store

	// flag values
	var userId int64
	var groupId int64
	var firstName string
	var lastName string
	var groupName string

	return cli.Command{
		Name:  name,
		Usage: "perform operation on data resources",
		Before: func(context *cli.Context) error {
			vars = <-chVars

			url := data.MakeUrl(vars)
			db, err := sql.Open("postgres", url)

			if err != nil {
				return err
			}

			s, err := data.NewStore(db, 10)

			if err != nil {
				return err
			}

			store = s
			return nil
		},
		Subcommands: []cli.Command{
			{
				Name: "user:get",
				Action: func(ctx *cli.Context) error {
					id, err := getIdArg(ctx)

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					user, err := store.GetUserById(id)

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					return Printed(user)
				},
			},
			{
				Name: "user:create",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "FirstName", Destination: &firstName, Required: true},
					cli.StringFlag{Name: "LastName", Destination: &lastName, Required: true},
				},
				Action: func(ctx *cli.Context) error {
					user, err := store.CreateUser(&data.User{
						FirstName: firstName,
						LastName: lastName,
					})

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					return Printed(user)
				},
			},
			{
				Name: "user:update",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "FirstName", Destination: &firstName},
					cli.StringFlag{Name: "LastName", Destination: &lastName},
				},
				Action: func(ctx *cli.Context) error {
					id, err := getIdArg(ctx)

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					fields := getPassedFlagNames(ctx)

					err = store.UpdateUser(id, &data.User{
						FirstName: firstName,
						LastName: lastName,
					}, fields...)

					return err
				},
			},
			{
				Name: "user:delete",
				Action: func(ctx *cli.Context) error {
					id, err := getIdArg(ctx)

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					return store.DeleteUser(id)
				},
			},
			{
				Name: "users:get",
				Flags: []cli.Flag{
					cli.Int64Flag{Name: "GroupId", Destination: &groupId, Required: true},
				},
				Action: func(ctx *cli.Context) error {
					var users []data.User
					var err error

					if ctx.IsSet("GroupId") {
						users, err = store.GetUsersByGroupId(groupId)
					}

					if err != nil {
						return err
					}

					return Printed(users)
				},
			},
			{
				Name: "group:get",
				Action: func(ctx *cli.Context) error {
					id, err := getIdArg(ctx)

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					group, err := store.GetGroupById(id)

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					return Printed(group)
				},
			},
			{
				Name: "group:create",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "Name", Destination: &groupName, Required: true},
				},
				Action: func(ctx *cli.Context) error {
					group, err := store.CreateGroup(&data.Group{Name: groupName})

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					return Printed(group)
				},
			},
			{
				Name: "group:update",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "Name", Destination: &groupName},
				},
				Action: func(ctx *cli.Context) error {
					id, err := getIdArg(ctx)

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					fields := getPassedFlagNames(ctx)

					err = store.UpdateGroup(id, &data.Group{Name: groupName}, fields...)

					return err
				},
			},
			{
				Name: "group:delete",
				Action: func(ctx *cli.Context) error {
					id, err := getIdArg(ctx)

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					return store.DeleteGroup(id)
				},
			},
			{
				Name: "groups:get",
				Flags: []cli.Flag{
					cli.Int64Flag{Name: "UserId", Destination: &userId, Required: true},
				},
				Action: func(ctx *cli.Context) error {
					var groups []data.Group
					var err error

					if ctx.IsSet("UserId") {
						groups, err = store.GetGroupsByUserId(userId)
					}

					if err != nil {
						logger.Error(errors.Unwrap(err))
						return err
					}

					return Printed(groups)
				},
			},
			{
				Name: "group:add-user",
				Flags: []cli.Flag{
					cli.Int64Flag{Name: "GroupId", Destination: &groupId, Required: true},
					cli.Int64Flag{Name: "UserId", Destination: &userId, Required: true},
				},
				Action: func(ctx *cli.Context) error {
					err := store.LinkGroupToUser(groupId, userId)

					if err != nil {
						logger.Error(errors.Unwrap(err))
					}

					return err
				},
			},
			{
				Name: "group:remove-user",
				Flags: []cli.Flag{
					cli.Int64Flag{Name: "GroupId", Destination: &groupId, Required: true},
					cli.Int64Flag{Name: "UserId", Destination: &userId, Required: true},
				},
				Action: func(ctx *cli.Context) error {
					err := store.UnlinkGroupFromUser(groupId, userId)

					if err != nil {
						logger.Error(errors.Unwrap(err))
					}

					return err
				},
			},
		},
	}
}

func getIdArg(ctx *cli.Context) (int64, error) {
	if ctx.NArg() == 0 {
		return 0, fmt.Errorf("missing Id in args")
	}

	s := ctx.Args().Get(0)

	i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("failed to get id from args: %w", err)
	}

	return i, nil
}

func getPassedFlagNames(ctx *cli.Context) []string {
 	fields := make([]string, 0)

	for _, name := range ctx.FlagNames() {
		if ctx.IsSet(name) {
			fields = append(fields, name)
		}
	}

 	return fields
}

func Printed(v interface {}) error {
	j, err := json.Marshal(v)

	if err != nil {
		return fmt.Errorf("failed to convert data to json: %w", err)
	}

	fmt.Println(string(j))

	return nil
}
