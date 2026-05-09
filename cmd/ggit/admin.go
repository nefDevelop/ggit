package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/urfave/cli/v3"
)

var (
	adminCommand = cli.Command{
		Name:  "admin",
		Usage: "Perform admin operations on command line",
		Description: `Allow using internal logic of Ggit without hacking into the source code
to make automatic initialization process more smoothly`,
		Commands: []*cli.Command{
			&subcmdCreateUser,
			&subcmdDeleteInactivateUsers,
			&subcmdDeleteRepositoryArchives,
			&subcmdDeleteMissingRepositories,
			&subcmdGitGcRepos,
			&subcmdRewriteAuthorizedKeys,
			&subcmdSyncRepositoryHooks,
			&subcmdReinitMissingRepositories,
			&subcmdUpdateRepositoryCounters,
		},
	}

	subcmdCreateUser = cli.Command{
		Name:   "create-user",
		Usage:  "Create a new user in database",
		Action: runCreateUser,
		Flags: []cli.Flag{
			stringFlag("name", "", "Username"),
			stringFlag("password", "", "User password"),
			stringFlag("email", "", "User email address"),
			boolFlag("admin", "User is an admin"),
			stringFlag("config, c", "", "Custom configuration file path"),
		},
	}

	subcmdDeleteInactivateUsers = cli.Command{
		Name:  "delete-inactive-users",
		Usage: "Delete all inactive accounts",
		Action: adminDashboardOperation(
			func() error { return database.Handle.Users().DeleteInactivated(context.Background()) },
			"All inactivated accounts have been deleted successfully",
		),
		Flags: []cli.Flag{
			stringFlag("config, c", "", "Custom configuration file path"),
		},
	}

	subcmdDeleteRepositoryArchives = cli.Command{
		Name:  "delete-repository-archives",
		Usage: "Delete all repositories archives",
		Action: adminDashboardOperation(
			database.DeleteRepositoryArchives,
			"All repositories archives have been deleted successfully",
		),
		Flags: []cli.Flag{
			stringFlag("config, c", "", "Custom configuration file path"),
		},
	}

	subcmdDeleteMissingRepositories = cli.Command{
		Name:  "delete-missing-repositories",
		Usage: "Delete all repository records that lost Git files",
		Action: adminDashboardOperation(
			database.DeleteMissingRepositories,
			"All repositories archives have been deleted successfully",
		),
		Flags: []cli.Flag{
			stringFlag("config, c", "", "Custom configuration file path"),
		},
	}

	subcmdGitGcRepos = cli.Command{
		Name:  "collect-garbage",
		Usage: "Do garbage collection on repositories",
		Action: adminDashboardOperation(
			database.GitGcRepos,
			"All repositories have done garbage collection successfully",
		),
		Flags: []cli.Flag{
			stringFlag("config, c", "", "Custom configuration file path"),
		},
	}

	subcmdRewriteAuthorizedKeys = cli.Command{
		Name:  "rewrite-authorized-keys",
		Usage: "Rewrite '.ssh/authorized_keys' file (caution: non-Gogs keys will be lost)",
		Action: adminDashboardOperation(
			database.RewriteAuthorizedKeys,
			"All public keys have been rewritten successfully",
		),
		Flags: []cli.Flag{
			stringFlag("config, c", "", "Custom configuration file path"),
		},
	}

	subcmdSyncRepositoryHooks = cli.Command{
		Name:  "resync-hooks",
		Usage: "Resync pre-receive, update and post-receive hooks",
		Action: adminDashboardOperation(
			database.SyncRepositoryHooks,
			"All repositories' pre-receive, update and post-receive hooks have been resynced successfully",
		),
		Flags: []cli.Flag{
			stringFlag("config, c", "", "Custom configuration file path"),
		},
	}

	subcmdReinitMissingRepositories = cli.Command{
		Name:  "reinit-missing-repositories",
		Usage: "Reinitialize all repository records that lost Git files",
		Action: adminDashboardOperation(
			database.ReinitMissingRepositories,
			"All repository records that lost Git files have been reinitialized successfully",
		),
		Flags: []cli.Flag{
			stringFlag("config, c", "", "Custom configuration file path"),
		},
	}

	subcmdUpdateRepositoryCounters = cli.Command{
		Name:  "update-repository-counters",
		Usage: "Update NumCommits and NumBranches for all repositories",
		Action: adminDashboardOperation(
			database.UpdateAllRepositoryCounters,
			"All repository counters have been updated successfully",
		),
		Flags: []cli.Flag{
			stringFlag("config, c", "", "Custom configuration file path"),
		},
	}
)

func runCreateUser(ctx context.Context, cmd *cli.Command) error {
	if !cmd.IsSet("name") {
		return errors.New("Username is not specified")
	} else if !cmd.IsSet("password") {
		return errors.New("Password is not specified")
	} else if !cmd.IsSet("email") {
		return errors.New("Email is not specified")
	}

	err := conf.Init(configFromLineage(cmd))
	if err != nil {
		return fmt.Errorf("init configuration"+": %w", err)
	}
	conf.InitLogging(true)

	if _, err = database.SetEngine(); err != nil {
		return fmt.Errorf("set engine"+": %w", err)
	}

	user, err := database.Handle.Users().Create(
		ctx,
		cmd.String("name"),
		cmd.String("email"),
		database.CreateUserOptions{
			Password:  cmd.String("password"),
			Activated: true,
			Admin:     cmd.Bool("admin"),
		},
	)
	if err != nil {
		if !database.IsErrUserAlreadyExist(err) {
			return fmt.Errorf("create user"+": %w", err)
		}

		// If user already exists, update password and admin status to ensure test credentials are correct
		u, err := database.Handle.Users().GetByUsername(ctx, cmd.String("name"))
		if err != nil {
			return fmt.Errorf("get existing user"+": %w", err)
		}

		pass := cmd.String("password")
		isAdmin := cmd.Bool("admin")
		err = database.Handle.Users().Update(ctx, u.ID, database.UpdateUserOptions{
			Password: &pass,
			IsAdmin:  &isAdmin,
		})
		if err != nil {
			return fmt.Errorf("update existing user"+": %w", err)
		}
		fmt.Printf("User %q already existed, password and admin status have been updated.\n", u.Name)
		return nil
	}

	fmt.Printf("New user %q has been successfully created!\n", user.Name)
	return nil
}

func adminDashboardOperation(operation func() error, successMessage string) func(context.Context, *cli.Command) error {
	return func(_ context.Context, cmd *cli.Command) error {
		err := conf.Init(configFromLineage(cmd))
		if err != nil {
			return fmt.Errorf("init configuration"+": %w", err)
		}
		conf.InitLogging(true)

		if _, err = database.SetEngine(); err != nil {
			return fmt.Errorf("set engine"+": %w", err)
		}

		if err := operation(); err != nil {
			functionName := runtime.FuncForPC(reflect.ValueOf(operation).Pointer()).Name()
			return fmt.Errorf("%s: %v", functionName, err)
		}

		fmt.Printf("%s\n", successMessage)
		return nil
	}
}
