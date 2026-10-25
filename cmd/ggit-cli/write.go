package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/nef734/ggit/internal/apiclient"
	"github.com/urfave/cli/v3"
)

var subcmdIssueClose = cli.Command{
	Name:      "close",
	Usage:     "Cierra una incidencia",
	ArgsUsage: "OWNER/NAME NUMBER",
	Action:    runIssueClose,
}

var subcmdIssueReopen = cli.Command{
	Name:      "reopen",
	Usage:     "Reabre una incidencia",
	ArgsUsage: "OWNER/NAME NUMBER",
	Action:    runIssueReopen,
}

var subcmdIssueCreate = cli.Command{
	Name:      "create",
	Usage:     "Crea una incidencia",
	ArgsUsage: "OWNER/NAME",
	Action:    runIssueCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{Name: "title", Usage: "Título de la incidencia"},
		&cli.IntFlag{Name: "assignee", Usage: "ID del usuario asignado"},
		&cli.IntFlag{Name: "milestone", Usage: "ID del hito"},
		&cli.StringFlag{Name: "labels", Usage: "IDs de etiquetas separados por coma"},
	}, bodyFlags()...),
}

var subcmdIssueComment = cli.Command{
	Name:      "comment",
	Usage:     "Añade un comentario a una incidencia",
	ArgsUsage: "OWNER/NAME NUMBER",
	Action:    runIssueComment,
	Flags:     bodyFlags(),
}

func runIssueCreate(ctx context.Context, cmd *cli.Command) error {
	owner, name, err := repoArg(cmd)
	if err != nil {
		return err
	}

	title := strings.TrimSpace(cmd.String("title"))
	if title == "" {
		return fmt.Errorf("falta --title")
	}

	body, err := resolveBody(cmd)
	if err != nil {
		return err
	}

	labels, err := parseIDs(cmd.String("labels"))
	if err != nil {
		return err
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	issue, err := client.CreateIssue(ctx, owner, name, apiclient.CreateIssueOptions{
		Title:       title,
		Content:     body,
		MilestoneID: int64(cmd.Int("milestone")),
		AssigneeID:  int64(cmd.Int("assignee")),
		LabelIDs:    labels,
	})
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printJSON(issue)
	}

	fmt.Printf("#%d %s (%s)\n", issue.Index, issue.Title, issue.State)
	return nil
}

func runIssueClose(ctx context.Context, cmd *cli.Command) error {
	return setIssueClosed(ctx, cmd, true)
}

func runIssueReopen(ctx context.Context, cmd *cli.Command) error {
	return setIssueClosed(ctx, cmd, false)
}

func setIssueClosed(ctx context.Context, cmd *cli.Command, closed bool) error {
	owner, name, index, err := numberArg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient(cmd)
	if err != nil {
		return err
	}

	raw, err := client.SetIssueClosed(ctx, owner, name, index, closed)
	if err != nil {
		return err
	}
	if wantJSON(cmd) {
		return printRawJSON(raw)
	}

	state := "reabierta"
	if closed {
		state = "cerrada"
	}
	fmt.Printf("incidencia #%d %s\n", index, state)
	return nil
}

func runIssueComment(ctx context.Context, cmd *cli.Command) error {
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

func bodyFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "body", Usage: "Texto del mensaje"},
		&cli.StringFlag{Name: "body-file", Usage: "Archivo con el texto; usa - para leer de stdin"},
	}
}

// resolveBody devuelve el texto indicado por --body o --body-file; vacío si no
// se pasó ninguno.
func resolveBody(cmd *cli.Command) (string, error) {
	inline := cmd.String("body")
	file := cmd.String("body-file")

	if inline != "" && file != "" {
		return "", fmt.Errorf("usa --body o --body-file, no ambos")
	}
	if file == "" {
		return inline, nil
	}

	var (
		raw []byte
		err error
	)
	if file == "-" {
		raw, err = io.ReadAll(os.Stdin)
	} else {
		raw, err = os.ReadFile(file)
	}
	if err != nil {
		return "", fmt.Errorf("leer %s: %w", file, err)
	}
	return string(raw), nil
}

func requiredBody(cmd *cli.Command) (string, error) {
	body, err := resolveBody(cmd)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(body) == "" {
		return "", fmt.Errorf("falta el mensaje: usa --body o --body-file")
	}
	return body, nil
}

func numberArg(cmd *cli.Command) (owner, name string, index int64, err error) {
	owner, name, err = repoArg(cmd)
	if err != nil {
		return
	}

	index, err = strconv.ParseInt(cmd.Args().Get(1), 10, 64)
	if err != nil || index <= 0 {
		return "", "", 0, fmt.Errorf("uso: %s OWNER/NAME NUMBER", cmd.FullName())
	}
	return owner, name, index, nil
}

func parseIDs(list string) ([]int64, error) {
	if strings.TrimSpace(list) == "" {
		return nil, nil
	}

	parts := strings.Split(list, ",")
	ids := make([]int64, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("ID inválido %q en --labels", strings.TrimSpace(part))
		}
		ids = append(ids, id)
	}
	return ids, nil
}
