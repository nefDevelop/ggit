package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/markup"
	"github.com/nef734/ggit/internal/osx"
)

const (
	tmplAuthActivate       = "auth/activate"
	tmplAuthActivateEmail  = "auth/activate_email"
	tmplAuthResetPassword  = "auth/reset_passwd"
	tmplAuthRegisterNotify = "auth/register_notify"

	tmplIssueComment = "issue/comment"
	tmplIssueMention = "issue/mention"

	tmplNotifyCollaborator = "notify/collaborator"
)

var (
	tplRender     *template.Template
	tplRenderOnce sync.Once
)

// render renders a mail template with given data.
func render(tpl string, data map[string]any) (string, error) {
	tplRenderOnce.Do(func() {
		tplRender = template.New("mail").Funcs(template.FuncMap{
			"AppName": func() string {
				return conf.App.BrandName
			},
			"AppURL": func() string {
				return conf.Server.ExternalURL
			},
			"Year": func() int {
				return time.Now().Year()
			},
			"Str2HTML": func(raw string) template.HTML {
				return template.HTML(markup.Sanitize(raw))
			},
		})

		customDir := filepath.Join(conf.CustomDir(), "templates", "mail")
		workDir := filepath.Join(conf.WorkDir(), "templates", "mail")

		loadTemplates := func(root string) {
			_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() || !strings.HasSuffix(p, ".tmpl") {
					return nil
				}
				
				rel, _ := filepath.Rel(root, p)
				name := rel[:len(rel)-len(".tmpl")]
				name = filepath.ToSlash(name)

				data, err := os.ReadFile(p)
				if err != nil {
					return nil
				}
				
				_, _ = tplRender.New(name).Parse(string(data))
				return nil
			})
		}

		// Load from work dir first, then override with custom dir
		loadTemplates(workDir)
		if osx.IsDir(customDir) {
			loadTemplates(customDir)
		}
	})

	var buf bytes.Buffer
	if err := tplRender.ExecuteTemplate(&buf, tpl, data); err != nil {
		return "", errors.Wrap(err, "execute template")
	}
	return buf.String(), nil
}

func SendTestMail(email string) error {
	msg, err := newMessage([]string{email}, "Gogs Test Email", "Hello 👋, greeting from Gogs!")
	if err != nil {
		return errors.Wrap(err, "new message")
	}
	return sendMessage(msg)
}

/*
	Setup interfaces of used methods in mail to avoid cycle import.
*/

type User interface {
	ID() int64
	DisplayName() string
	Email() string
	GenerateEmailActivateCode(string) string
}

type Repository interface {
	FullName() string
	HTMLURL() string
	ComposeMetas() map[string]string
}

type Issue interface {
	MailSubject() string
	Content() string
	HTMLURL() string
}

func SendUserMail(u User, tpl, code, subject, info string) error {
	data := map[string]any{
		"Username":          u.DisplayName(),
		"ActiveCodeLives":   conf.Auth.ActivateCodeLives / 60,
		"ResetPwdCodeLives": conf.Auth.ResetPasswordCodeLives / 60,
		"Code":              code,
	}
	body, err := render(tpl, data)
	if err != nil {
		return errors.Wrap(err, "render")
	}

	msg, err := newMessage([]string{u.Email()}, subject, body)
	if err != nil {
		return errors.Wrap(err, "new message")
	}
	msg.info = fmt.Sprintf("UID: %d, %s", u.ID(), info)

	send(msg)
	return nil
}

func SendActivateAccountMail(u User, subject string) error {
	return SendUserMail(u, tmplAuthActivate, u.GenerateEmailActivateCode(u.Email()), subject, "activate account")
}

func SendResetPasswordMail(u User, subject string) error {
	return SendUserMail(u, tmplAuthResetPassword, u.GenerateEmailActivateCode(u.Email()), subject, "reset password")
}

func SendActivateEmailMail(u User, email string, subject string) error {
	data := map[string]any{
		"Username":        u.DisplayName(),
		"ActiveCodeLives": conf.Auth.ActivateCodeLives / 60,
		"Code":            u.GenerateEmailActivateCode(email),
		"Email":           email,
	}
	body, err := render(tmplAuthActivateEmail, data)
	if err != nil {
		return errors.Wrap(err, "render")
	}

	msg, err := newMessage([]string{email}, subject, body)
	if err != nil {
		return errors.Wrap(err, "new message")
	}
	msg.info = fmt.Sprintf("UID: %d, activate email", u.ID())

	send(msg)
	return nil
}

func SendRegisterNotifyMail(u User, subject string) error {
	data := map[string]any{
		"Username": u.DisplayName(),
	}
	body, err := render(tmplAuthRegisterNotify, data)
	if err != nil {
		return errors.Wrap(err, "render")
	}

	msg, err := newMessage([]string{u.Email()}, subject, body)
	if err != nil {
		return errors.Wrap(err, "new message")
	}
	msg.info = fmt.Sprintf("UID: %d, registration notify", u.ID())

	send(msg)
	return nil
}

func SendCollaboratorMail(u, doer User, repo Repository) error {
	subject := fmt.Sprintf("%s added you to %s", doer.DisplayName(), repo.FullName())

	data := map[string]any{
		"Subject":  subject,
		"RepoName": repo.FullName(),
		"Link":     repo.HTMLURL(),
	}
	body, err := render(tmplNotifyCollaborator, data)
	if err != nil {
		return errors.Wrap(err, "render")
	}

	msg, err := newMessage([]string{u.Email()}, subject, body)
	if err != nil {
		return errors.Wrap(err, "new message")
	}
	msg.info = fmt.Sprintf("UID: %d, add collaborator", u.ID())

	send(msg)
	return nil
}

func composeTplData(subject, body, link string) map[string]any {
	data := make(map[string]any, 10)
	data["Subject"] = subject
	data["Body"] = body
	data["Link"] = link
	return data
}

func composeIssueMessage(issue Issue, repo Repository, doer User, tplName string, tos []string, info string) (*message, error) {
	subject := issue.MailSubject()
	body := string(markup.Markdown([]byte(issue.Content()), repo.HTMLURL(), repo.ComposeMetas()))
	data := composeTplData(subject, body, issue.HTMLURL())
	data["Doer"] = doer
	content, err := render(tplName, data)
	if err != nil {
		return nil, errors.Wrapf(err, "render %q", tplName)
	}
	from := (&mail.Address{Name: doer.DisplayName(), Address: conf.Email.FromEmail}).String()
	msg, err := newMessageFrom(tos, from, subject, content)
	if err != nil {
		return nil, errors.Wrap(err, "new message")
	}
	msg.info = fmt.Sprintf("Subject: %s, %s", subject, info)
	return msg, nil
}

// SendIssueCommentMail composes and sends issue comment emails to target receivers.
func SendIssueCommentMail(issue Issue, repo Repository, doer User, tos []string) error {
	if len(tos) == 0 {
		return nil
	}

	msg, err := composeIssueMessage(issue, repo, doer, tmplIssueComment, tos, "issue comment")
	if err != nil {
		return errors.Wrap(err, "compose issue message")
	}
	send(msg)
	return nil
}

// SendIssueMentionMail composes and sends issue mention emails to target receivers.
func SendIssueMentionMail(issue Issue, repo Repository, doer User, tos []string) error {
	if len(tos) == 0 {
		return nil
	}
	msg, err := composeIssueMessage(issue, repo, doer, tmplIssueMention, tos, "issue mention")
	if err != nil {
		return errors.Wrap(err, "compose issue message")
	}
	send(msg)
	return nil
}
