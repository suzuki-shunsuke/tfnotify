package github

import (
	"context"
	"log"
	"net/http"

	"github.com/mercari/tfnotify/terraform"
)

// NotifyService handles communication with the notification related
// methods of GitHub API
type NotifyService service

// Notify posts comment optimized for notifications
func (g *NotifyService) Notify(ctx context.Context, body string) (exit int, err error) {
	cfg := g.client.Config
	parser := g.client.Config.Parser
	template := g.client.Config.Template

	result := parser.Parse(body)
	if result.Error != nil {
		return result.ExitCode, result.Error
	}
	if result.Result == "" {
		return result.ExitCode, result.Error
	}

	_, isPlan := parser.(*terraform.PlanParser)
	if isPlan {
		if result.HasDestroy && cfg.WarnDestroy {
			// Notify destroy warning as a new comment before normal plan result
			if err = g.notifyDestoryWarning(ctx, body, result); err != nil {
				return result.ExitCode, err
			}
		}
		if cfg.PR.IsNumber() && cfg.ResultLabels.HasAnyLabelDefined() {
			var (
				labelToAdd string
				labelColor string
			)

			switch {
			case result.HasAddOrUpdateOnly:
				labelToAdd = cfg.ResultLabels.AddOrUpdateLabel
				labelColor = cfg.ResultLabels.AddOrUpdateLabelColor
			case result.HasDestroy:
				labelToAdd = cfg.ResultLabels.DestroyLabel
				labelColor = cfg.ResultLabels.DestroyLabelColor
			case result.HasNoChanges:
				labelToAdd = cfg.ResultLabels.NoChangesLabel
				labelColor = cfg.ResultLabels.NoChangesLabelColor
			case result.HasPlanError:
				labelToAdd = cfg.ResultLabels.PlanErrorLabel
				labelColor = cfg.ResultLabels.PlanErrorLabelColor
			}

			currentLabelColor, err := g.removeResultLabels(ctx, labelToAdd)
			if err != nil {
				log.Printf("[ERROR][tfnotify] remove labels: %v", err)
			}

			if labelToAdd != "" {
				if currentLabelColor == "" {
					labels, _, err := g.client.API.IssuesAddLabels(ctx, cfg.PR.Number, []string{labelToAdd})
					if err != nil {
						log.Printf("[ERROR][tfnotify] add a label %s: %v", labelToAdd, err)
					}
					if labelColor != "" {
						// set the color of label
						for _, label := range labels {
							if labelToAdd == label.GetName() {
								if label.GetColor() != labelColor {
									if _, _, err := g.client.API.IssuesUpdateLabel(ctx, labelToAdd, labelColor); err != nil {
										log.Printf("[ERROR][tfnotify] update a label color(name: %s, color: %s): %v", labelToAdd, labelColor, err)
									}
								}
							}
						}
					}
				} else if labelColor != "" && labelColor != currentLabelColor {
					// set the color of label
					if _, _, err := g.client.API.IssuesUpdateLabel(ctx, labelToAdd, labelColor); err != nil {
						log.Printf("[ERROR][tfnotify] update a label color(name: %s, color: %s): %v", labelToAdd, labelColor, err)
					}
				}
			}
		}
	}

	template.SetValue(terraform.CommonTemplate{
		Title:        cfg.PR.Title,
		Message:      cfg.PR.Message,
		Result:       result.Result,
		Body:         body,
		Link:         cfg.CI,
		UseRawOutput: cfg.UseRawOutput,
		Vars:         cfg.Vars,
	})
	body, err = template.Execute()
	if err != nil {
		return result.ExitCode, err
	}

	value := template.GetValue()

	if cfg.PR.IsNumber() {
		if !cfg.KeepDuplicateComments {
			g.client.Comment.DeleteDuplicates(ctx, value.Title)
		}
	}

	_, isApply := parser.(*terraform.ApplyParser)
	if isApply {
		prNumber, err := g.client.Commits.MergedPRNumber(ctx, cfg.PR.Revision)
		if err == nil {
			cfg.PR.Number = prNumber
		} else if !cfg.PR.IsNumber() {
			commits, err := g.client.Commits.List(ctx, cfg.PR.Revision)
			if err != nil {
				return result.ExitCode, err
			}
			lastRevision, _ := g.client.Commits.lastOne(commits, cfg.PR.Revision)
			cfg.PR.Revision = lastRevision
		}
	}

	return result.ExitCode, g.client.Comment.Post(ctx, body, PostOptions{
		Number:   cfg.PR.Number,
		Revision: cfg.PR.Revision,
	})
}

func (g *NotifyService) notifyDestoryWarning(ctx context.Context, body string, result terraform.ParseResult) error {
	cfg := g.client.Config
	destroyWarningTemplate := g.client.Config.DestroyWarningTemplate
	destroyWarningTemplate.SetValue(terraform.CommonTemplate{
		Title:        cfg.PR.DestroyWarningTitle,
		Message:      cfg.PR.DestroyWarningMessage,
		Result:       result.Result,
		Body:         body,
		Link:         cfg.CI,
		UseRawOutput: cfg.UseRawOutput,
		Vars:         cfg.Vars,
	})
	body, err := destroyWarningTemplate.Execute()
	if err != nil {
		return err
	}

	return g.client.Comment.Post(ctx, body, PostOptions{
		Number:   cfg.PR.Number,
		Revision: cfg.PR.Revision,
	})
}

func (g *NotifyService) removeResultLabels(ctx context.Context, label string) (string, error) {
	cfg := g.client.Config
	labels, _, err := g.client.API.IssuesListLabels(ctx, cfg.PR.Number, nil)
	if err != nil {
		return "", err
	}

	labelColor := ""
	for _, l := range labels {
		labelText := l.GetName()
		if labelText == label {
			labelColor = l.GetColor()
			continue
		}
		if cfg.ResultLabels.IsResultLabel(labelText) {
			resp, err := g.client.API.IssuesRemoveLabel(ctx, cfg.PR.Number, labelText)
			// Ignore 404 errors, which are from the PR not having the label
			if err != nil && resp.StatusCode != http.StatusNotFound {
				return labelColor, err
			}
		}
	}

	return labelColor, nil
}
