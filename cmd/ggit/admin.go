package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	pikociservice "github.com/nef734/ggit/internal/services/pikoci"
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
			&subcmdPikoCI,
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

var subcmdPikoCI = cli.Command{
	Name:  "pikoci",
	Usage: "Inspect the PikoCI CI/CD integration",
	Commands: []*cli.Command{
		&subcmdPikoCIStatus,
		&subcmdPikoCIBot,
	},
}

var subcmdPikoCIBot = cli.Command{
	Name:  "bot",
	Usage: "Create the CI bot user and a fresh access token",
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
		stringFlag("name, n", "pikoci", "Bot user name"),
		stringFlag("email, e", "pikoci@localhost", "Bot user email"),
		stringFlag("token-name, t", "ggit-cli", "Access token name"),
		stringFlag("repo, r", "", "Add the bot as collaborator of OWNER/NAME"),
	},
	Action: runPikoCIBot,
}

var subcmdPikoCIStatus = cli.Command{
	Name:   "status",
	Usage:  "Show PikoCI configuration, health and linked repositories",
	Action: runPikoCIStatus,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func runPikoCIStatus(_ context.Context, cmd *cli.Command) error {
	if err := conf.Init(configFromLineage(cmd)); err != nil {
		return fmt.Errorf("init configuration: %w", err)
	}
	conf.InitLogging(true)

	if _, err := database.SetEngine(); err != nil {
		return fmt.Errorf("set engine: %w", err)
	}

	fmt.Printf("enabled:  %v\n", conf.PikoCI.Enabled)
	fmt.Printf("mode:     %s\n", conf.PikoCI.Mode)
	fmt.Printf("url:      %s\n", conf.PikoCI.URL)
	if conf.PikoCI.PublicURL != "" {
		fmt.Printf("ui:       %s\n", conf.PikoCI.PublicURL)
	}

	ctx := context.Background()
	linked, err := database.Handle.PikoCIRepos().Count(ctx)
	if err != nil {
		return fmt.Errorf("count pikoci links: %w", err)
	}
	fmt.Printf("linked:   %d\n", linked)

	if !conf.PikoCI.Enabled {
		return nil
	}

	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	health, err := pikociservice.New(conf.PikoCI.URL).Health(healthCtx)
	if err != nil {
		fmt.Printf("health:   unreachable (%v)\n", err)
		return nil
	}
	fmt.Printf("health:   %s (db %s, version %s)\n", health.Status, health.DB, health.Version)
	return nil
}

func runPikoCIBot(ctx context.Context, cmd *cli.Command) error {
	if err := conf.Init(configFromLineage(cmd)); err != nil {
		return fmt.Errorf("init configuration: %w", err)
	}
	conf.InitLogging(true)

	if _, err := database.SetEngine(); err != nil {
		return fmt.Errorf("set engine: %w", err)
	}

	name := strings.TrimSpace(cmd.String("name"))
	bot, err := database.Handle.Users().GetByUsername(ctx, name)
	if err != nil {
		if !database.IsErrUserNotExist(err) {
			return fmt.Errorf("get bot user: %w", err)
		}
		bot, err = database.Handle.Users().Create(ctx, name, cmd.String("email"), database.CreateUserOptions{
			Activated: true,
		})
		if err != nil {
			return fmt.Errorf("create bot user: %w", err)
		}
		fmt.Printf("Usuario %q creado.\n", name)
	} else {
		fmt.Printf("Usuario %q ya existe.\n", name)
	}

	_, raw, err := database.Handle.AccessTokens().Create(ctx, bot.ID, cmd.String("token-name"))
	if err != nil {
		return fmt.Errorf("create access token (prueba con --token-name): %w", err)
	}
	fmt.Printf("Token %q: %s\n", cmd.String("token-name"), raw)
	fmt.Println("Guárdalo ahora: no se vuelve a mostrar.")

	if repoFull := strings.TrimSpace(cmd.String("repo")); repoFull != "" {
		ownerName, repoName, ok := strings.Cut(repoFull, "/")
		if !ok || ownerName == "" || repoName == "" {
			return fmt.Errorf("--repo espera OWNER/NAME, recibido %q", repoFull)
		}

		owner, err := database.Handle.Users().GetByUsername(ctx, ownerName)
		if err != nil {
			return fmt.Errorf("get repository owner: %w", err)
		}
		repo, err := database.Handle.Repositories().GetByName(ctx, owner.ID, repoName)
		if err != nil {
			return fmt.Errorf("get repository: %w", err)
		}
		if err := database.FromDomainRepository(repo).AddCollaborator(database.FromDomainUser(bot)); err != nil {
			return fmt.Errorf("add collaborator: %w", err)
		}
		fmt.Printf("%q es colaborador de %s.\n", name, repoFull)
	}
	return nil
}
