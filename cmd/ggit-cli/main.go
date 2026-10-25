// Command ggit-cli es un cliente de la API v2 de Ggit.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nef734/ggit/internal/apiclient"
	apiv2 "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/urfave/cli/v3"
)

const cliVersion = "0.1.0"

func main() {
	cmd := &cli.Command{
		Name:    "ggit-cli",
		Usage:   "Cliente de la API v2 de Ggit",
		Version: cliVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Usage: "URL base del servidor (por defecto GGIT_HOST o el archivo de configuración)"},
			&cli.StringFlag{Name: "token", Usage: "Token de acceso (por defecto GGIT_TOKEN o el archivo de configuración)"},
			&cli.BoolFlag{Name: "json", Usage: "Imprime la respuesta cruda en JSON"},
		},
		Commands: []*cli.Command{&authCommand, &repoCommand, &issueCommand, &prCommand},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "ggit-cli:", err)
		os.Exit(1)
	}
}

var authCommand = cli.Command{
	Name:  "auth",
	Usage: "Comprueba las credenciales contra el servidor",
	Commands: []*cli.Command{
		{
			Name:   "status",
			Usage:  "Muestra el usuario autenticado",
			Action: runAuthStatus,
		},
		{
			Name:   "login",
			Usage:  "Crea un token de acceso y lo guarda en la configuración",
			Action: runAuthLogin,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "username", Usage: "Usuario"},
				&cli.StringFlag{Name: "password", Usage: "Contraseña; usa - para leer de stdin"},
				&cli.StringFlag{Name: "name", Value: "ggit-cli", Usage: "Nombre del token"},
			},
		},
	},
}

var repoCommand = cli.Command{
	Name:  "repo",
	Usage: "Consulta repositorios",
	Commands: []*cli.Command{
		{
			Name:   "list",
			Usage:  "Lista tus repositorios y colaboraciones",
			Action: runRepoList,
		},
		{
			Name:      "view",
			Usage:     "Muestra los metadatos de un repositorio",
			ArgsUsage: "OWNER/NAME",
			Action:    runRepoView,
		},
	},
}

var issueCommand = cli.Command{
	Name:  "issue",
	Usage: "Consulta incidencias",
	Commands: []*cli.Command{
		{
			Name:      "list",
			Usage:     "Lista las incidencias de un repositorio",
			ArgsUsage: "OWNER/NAME",
			Action:    runIssueList,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "state", Value: "open", Usage: "Filtra por estado: open o closed"},
				&cli.IntFlag{Name: "page", Value: 1, Usage: "Página de resultados"},
				&cli.StringFlag{Name: "labels", Usage: "Nombres de etiquetas separados por coma"},
				&cli.IntFlag{Name: "assignee", Usage: "ID del usuario asignado"},
				&cli.IntFlag{Name: "poster", Usage: "ID del autor"},
				&cli.IntFlag{Name: "milestone", Usage: "ID del hito"},
				&cli.StringFlag{Name: "sort", Usage: "Criterio de ordenación"},
			},
		},
		{
			Name:      "view",
			Usage:     "Muestra una incidencia",
			ArgsUsage: "OWNER/NAME NUMBER",
			Action:    runIssueView,
		},
		&subcmdIssueCreate,
		&subcmdIssueComment,
		&subcmdIssueClose,
		&subcmdIssueReopen,
	},
}

func runAuthStatus(ctx context.Context, cmd *cli.Command) error {
	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	user, err := client.AuthStatus(ctx)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(user)
	}

	fmt.Printf("usuario: %s\n", user.Username)
	fmt.Printf("id:      %d\n", user.ID)
	if user.FullName != "" {
		fmt.Printf("nombre:  %s\n", user.FullName)
	}
	fmt.Printf("email:   %s\n", user.Email)
	fmt.Printf("admin:   %s\n", yesNo(user.IsAdmin))
	return nil
}

func runAuthLogin(ctx context.Context, cmd *cli.Command) error {
	cfg, err := apiclient.LoadConfig()
	if err != nil {
		return err
	}

	host := stringValue(cmd, "host")
	if host == "" {
		host = cfg.Host
	}
	baseURL, err := requireHost(host)
	if err != nil {
		return err
	}

	username := strings.TrimSpace(cmd.String("username"))
	if username == "" {
		return fmt.Errorf("falta --username")
	}

	password, err := resolvePassword(cmd)
	if err != nil {
		return err
	}

	// Basic auth basta para /user/settings/tokens: SignInRequired solo mira
	// IsLogged, que también se marca con credenciales básicas.
	client := apiclient.New(baseURL, "").WithBasicAuth(username, password)
	token, err := client.CreateAccessToken(ctx, cmd.String("name"))
	if err != nil {
		return err
	}

	if err := apiclient.SaveConfig(apiclient.Config{Host: baseURL, Token: token}); err != nil {
		return err
	}

	path, err := apiclient.ConfigPath()
	if err != nil {
		return err
	}
	fmt.Printf("token creado para %s y guardado en %s\n", username, path)
	return nil
}

func resolvePassword(cmd *cli.Command) (string, error) {
	password := cmd.String("password")
	if password == "" {
		return "", fmt.Errorf("falta --password (usa --password - para leer de stdin)")
	}
	if password != "-" {
		return password, nil
	}

	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("leer la contraseña de stdin: %w", err)
	}
	password = strings.TrimRight(string(raw), "\r\n")
	if password == "" {
		return "", fmt.Errorf("la contraseña leída de stdin está vacía")
	}
	return password, nil
}

func runRepoList(ctx context.Context, cmd *cli.Command) error {
	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	repos, err := client.ListUserRepos(ctx)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(repos)
	}

	if len(repos) == 0 {
		fmt.Println("no hay repositorios")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVISIBILITY\tISSUES\tPULLS\tUPDATED")
	for _, repo := range repos {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n",
			repo.FullName,
			visibility(repo),
			repo.OpenIssues,
			repo.OpenPulls,
			formatTime(repo.UpdatedAt),
		)
	}
	return w.Flush()
}

func runRepoView(ctx context.Context, cmd *cli.Command) error {
	fullName := cmd.Args().First()
	if fullName == "" {
		return fmt.Errorf("uso: ggit-cli repo view OWNER/NAME")
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	repo, err := client.GetRepo(ctx, fullName)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(repo)
	}

	fmt.Printf("repositorio:  %s\n", repo.FullName)
	fmt.Printf("visibilidad:  %s\n", visibility(repo))
	if repo.Description != "" {
		fmt.Printf("descripción:  %s\n", repo.Description)
	}
	fmt.Printf("rama:         %s\n", repo.DefaultBranch)
	fmt.Printf("url:          %s\n", repo.HTMLURL)
	fmt.Printf("clone https:  %s\n", repo.CloneURL)
	fmt.Printf("clone ssh:    %s\n", repo.SSHURL)
	fmt.Printf("estrellas:    %d\n", repo.NumStars)
	fmt.Printf("forks:        %d\n", repo.NumForks)
	fmt.Printf("incidencias:  %d abiertas\n", repo.OpenIssues)
	fmt.Printf("pulls:        %d abiertos\n", repo.OpenPulls)
	if len(repo.Topics) > 0 {
		fmt.Printf("temas:        %s\n", strings.Join(repo.Topics, ", "))
	}
	fmt.Printf("creado:       %s\n", formatTime(repo.CreatedAt))
	fmt.Printf("actualizado:  %s\n", formatTime(repo.UpdatedAt))
	return nil
}

func runIssueList(ctx context.Context, cmd *cli.Command) error {
	owner, name, err := repoArg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	issues, err := client.ListIssues(ctx, owner, name, apiclient.ListIssuesOptions{
		State:       cmd.String("state"),
		Page:        cmd.Int("page"),
		Labels:      cmd.String("labels"),
		AssigneeID:  int64(cmd.Int("assignee")),
		PosterID:    int64(cmd.Int("poster")),
		MilestoneID: int64(cmd.Int("milestone")),
		Sort:        cmd.String("sort"),
	})
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(issues)
	}

	if len(issues) == 0 {
		fmt.Println("no hay incidencias")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "NUMBER\tSTATE\tTITLE\tUSER\tUPDATED")
	for _, issue := range issues {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			issue.Index,
			issue.State,
			issue.Title,
			userName(issue.Poster),
			formatTime(issue.UpdatedAt),
		)
	}
	return w.Flush()
}

func runIssueView(ctx context.Context, cmd *cli.Command) error {
	owner, name, err := repoArg(cmd)
	if err != nil {
		return err
	}

	index, err := strconv.ParseInt(cmd.Args().Get(1), 10, 64)
	if err != nil || index <= 0 {
		return fmt.Errorf("uso: ggit-cli issue view OWNER/NAME NUMBER")
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	issue, err := client.GetIssue(ctx, owner, name, index)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(issue)
	}

	fmt.Printf("#%d %s\n", issue.Index, issue.Title)
	fmt.Printf("estado:    %s\n", issue.State)
	fmt.Printf("autor:     %s\n", userName(issue.Poster))
	fmt.Printf("asignado:  %s\n", userName(issue.Assignee))
	if issue.Milestone != nil {
		fmt.Printf("hito:      %s\n", issue.Milestone.Title)
	}
	if len(issue.Labels) > 0 {
		names := make([]string, len(issue.Labels))
		for i, label := range issue.Labels {
			names[i] = label.Name
		}
		fmt.Printf("etiquetas: %s\n", strings.Join(names, ", "))
	}
	fmt.Printf("comentarios: %d\n", issue.Comments)
	fmt.Printf("creado:    %s\n", formatTime(issue.CreatedAt))
	fmt.Printf("actualizado: %s\n", formatTime(issue.UpdatedAt))
	fmt.Println()

	if body := strings.TrimSpace(issue.Body); body != "" {
		fmt.Println(body)
	} else {
		fmt.Println("(sin descripción)")
	}
	return nil
}

func newClient(cmd *cli.Command) (*apiclient.Client, error) {
	cfg, err := apiclient.LoadConfig()
	if err != nil {
		return nil, err
	}

	host := stringValue(cmd, "host")
	if host == "" {
		host = cfg.Host
	}
	baseURL, err := requireHost(host)
	if err != nil {
		return nil, err
	}

	token := stringValue(cmd, "token")
	if token == "" {
		token = cfg.Token
	}
	return apiclient.New(baseURL, token), nil
}

func requireHost(host string) (string, error) {
	normalized := apiclient.NormalizeHost(host)
	if normalized == "" {
		return "", fmt.Errorf("falta el host: usa --host, %s o %s", apiclient.EnvHost, configHint())
	}
	return normalized, nil
}

func repoArg(cmd *cli.Command) (owner, name string, err error) {
	fullName := cmd.Args().First()
	if fullName == "" {
		return "", "", fmt.Errorf("uso: %s OWNER/NAME", cmd.FullName())
	}
	return apiclient.SplitRepoName(fullName)
}

// stringValue resuelve una flag declarada en el comando raíz, porque en
// urfave/cli v3 cmd.String solo consulta las flags del propio subcomando.
func stringValue(cmd *cli.Command, name string) string {
	if cmd.IsSet(name) {
		return cmd.String(name)
	}
	for _, parent := range cmd.Lineage() {
		if parent.IsSet(name) {
			return parent.String(name)
		}
	}
	return ""
}

func wantJSON(cmd *cli.Command) bool {
	if cmd.IsSet("json") {
		return cmd.Bool("json")
	}
	for _, parent := range cmd.Lineage() {
		if parent.IsSet("json") {
			return parent.Bool("json")
		}
	}
	return false
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// printRawJSON reindenta una respuesta que el cliente no modela (por ejemplo,
// los modelos internos que devuelven PUT /issues y POST /pulls/{n}/reviews).
func printRawJSON(raw json.RawMessage) error {
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		buf.Write(raw)
	}
	buf.WriteByte('\n')
	_, err := os.Stdout.Write(buf.Bytes())
	return err
}

func configHint() string {
	path, err := apiclient.ConfigPath()
	if err != nil {
		return "el archivo de configuración"
	}
	return path
}

func visibility(repo *apiv2.RepositoryResponse) string {
	if repo.IsPrivate {
		return "private"
	}
	return "public"
}

func userName(user *apiv2.UserResponse) string {
	if user == nil || user.Username == "" {
		return "-"
	}
	return user.Username
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Local().Format("2006-01-02 15:04")
}

func yesNo(v bool) string {
	if v {
		return "sí"
	}
	return "no"
}
