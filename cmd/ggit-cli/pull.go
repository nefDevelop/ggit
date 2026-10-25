package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nef734/ggit/internal/apiclient"
	apiv2 "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/urfave/cli/v3"
)

var prCommand = cli.Command{
	Name:  "pr",
	Usage: "Consulta pull requests",
	Commands: []*cli.Command{
		{
			Name:      "list",
			Usage:     "Lista los pull requests de un repositorio",
			ArgsUsage: "OWNER/NAME",
			Action:    runPRList,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "state", Value: "open", Usage: "Filtra por estado: open o closed"},
				&cli.IntFlag{Name: "page", Value: 1, Usage: "Página de resultados"},
			},
		},
		{
			Name:      "view",
			Usage:     "Muestra un pull request",
			ArgsUsage: "OWNER/NAME NUMBER",
			Action:    runPRView,
		},
		{
			Name:      "reviews",
			Usage:     "Lista las revisiones de un pull request",
			ArgsUsage: "OWNER/NAME NUMBER",
			Action:    runPRReviews,
		},
		{
			Name:      "diff",
			Usage:     "Muestra el diff de un pull request",
			ArgsUsage: "OWNER/NAME NUMBER",
			Action:    runPRDiff,
		},
		&subcmdPRComment,
	},
}

var subcmdPRComment = cli.Command{
	Name:      "comment",
	Usage:     "Añade un comentario a un pull request",
	ArgsUsage: "OWNER/NAME NUMBER",
	Action:    runPRComment,
	Flags:     bodyFlags(),
}

func runPRList(ctx context.Context, cmd *cli.Command) error {
	owner, name, err := repoArg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	pulls, err := client.ListPullRequests(ctx, owner, name, apiclient.ListPullRequestsOptions{
		State: cmd.String("state"),
		Page:  cmd.Int("page"),
	})
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(pulls)
	}

	if len(pulls) == 0 {
		fmt.Println("no hay pull requests")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "NUMBER\tSTATE\tTITLE\tUSER\tBRANCH\tUPDATED")
	for _, pull := range pulls {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			pull.Index,
			pullState(pull),
			pull.Title,
			userName(pull.Poster),
			pull.HeadBranch+" → "+pull.BaseBranch,
			formatTime(pull.UpdatedAt),
		)
	}
	return w.Flush()
}

func runPRView(ctx context.Context, cmd *cli.Command) error {
	owner, name, index, err := numberArg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	pull, err := client.GetPullRequest(ctx, owner, name, index)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(pull)
	}

	fmt.Printf("#%d %s\n", pull.Index, pull.Title)
	fmt.Printf("estado:     %s\n", pullState(pull))
	if pull.Mergeable != nil {
		fmt.Printf("fusionable: %s\n", yesNo(*pull.Mergeable))
	}
	fmt.Printf("autor:      %s\n", userName(pull.Poster))
	fmt.Printf("rama:       %s → %s\n", pull.HeadBranch, pull.BaseBranch)
	fmt.Printf("asignado:   %s\n", userName(pull.Assignee))
	if pull.Milestone != nil {
		fmt.Printf("hito:       %s\n", pull.Milestone.Title)
	}
	if len(pull.Labels) > 0 {
		names := make([]string, len(pull.Labels))
		for i, label := range pull.Labels {
			names[i] = label.Name
		}
		fmt.Printf("etiquetas:  %s\n", strings.Join(names, ", "))
	}
	if pull.CombinedStatus != "" {
		fmt.Printf("checks:     %s\n", pull.CombinedStatus)
	}
	fmt.Printf("comentarios: %d\n", pull.Comments)
	if pull.HasMerged {
		fmt.Printf("fusionado:  %s por %s\n", formatTimePtr(pull.MergedAt), userName(pull.MergedBy))
		if pull.MergedCommitID != nil {
			fmt.Printf("commit:     %s\n", *pull.MergedCommitID)
		}
	}
	fmt.Printf("url:        %s\n", pull.HTMLURL)
	fmt.Printf("creado:     %s\n", formatTime(pull.CreatedAt))
	fmt.Printf("actualizado: %s\n", formatTime(pull.UpdatedAt))
	fmt.Println()

	if body := strings.TrimSpace(pull.Body); body != "" {
		fmt.Println(body)
	} else {
		fmt.Println("(sin descripción)")
	}
	return nil
}

func runPRReviews(ctx context.Context, cmd *cli.Command) error {
	owner, name, index, err := numberArg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	reviews, err := client.ListPullReviews(ctx, owner, name, index)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(reviews)
	}

	if len(reviews) == 0 {
		fmt.Println("no hay revisiones")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "REVIEWER\tTYPE\tCOMMENTS\tCREATED")
	for _, review := range reviews {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
			userName(review.Reviewer),
			reviewTypeName(review.Type),
			len(review.Comments),
			formatTime(review.CreatedAt),
		)
	}
	return w.Flush()
}

func runPRDiff(ctx context.Context, cmd *cli.Command) error {
	owner, name, index, err := numberArg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	diff, err := client.GetPullRequestDiff(ctx, owner, name, index)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(diff)
	}

	if len(diff.Files) == 0 {
		fmt.Println("el diff está vacío")
		return nil
	}

	for i, file := range diff.Files {
		if i > 0 {
			fmt.Println()
		}
		additions, deletions := file.Stats()
		fmt.Printf("=== %s (%s) +%d -%d\n", file.Name, file.TypeName(), additions, deletions)
		for _, section := range file.Sections {
			for _, line := range section.Lines {
				fmt.Println(line.Content)
			}
		}
	}
	return nil
}

func runPRComment(ctx context.Context, cmd *cli.Command) error {
	owner, name, index, err := numberArg(cmd)
	if err != nil {
		return err
	}

	body, err := requiredBody(cmd)
	if err != nil {
		return err
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	// En la API v2 un pull request es una incidencia, así que el comentario va
	// por el mismo endpoint que las incidencias.
	comment, err := client.CreateIssueComment(ctx, owner, name, index, body)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(comment)
	}

	fmt.Printf("comentario creado en #%d\n", index)
	if comment.HTMLURL != "" {
		fmt.Println(comment.HTMLURL)
	}
	return nil
}

func pullState(pull *apiv2.PullRequestResponse) string {
	if pull.HasMerged {
		return "merged"
	}
	return string(pull.State)
}

func reviewTypeName(t int) string {
	switch t {
	case 0:
		return "pending"
	case 1:
		return "approve"
	case 2:
		return "comment"
	case 3:
		return "reject"
	default:
		return "unknown"
	}
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return formatTime(*t)
}
