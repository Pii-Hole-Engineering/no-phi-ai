package gh

import (
	"context"

	"github.com/google/go-github/v58/github"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const LabelCleanPHI string = "Clean scan by no-phi-ai"
const LabelDirtyPHI string = "PHI detected by AI"

func checkResponse(resp *github.Response, err error) error {
	if err != nil {
		return err
	} else if (resp.StatusCode != 200) && (resp.StatusCode != 201) {
		return errors.New("GitHub API server returned status code " + resp.Status)
	}
	return nil
}

func generateLabels() []*github.Label {
	// set defaults for string vars required by github.Label struct
	var (
		cleanPHIlabel = LabelCleanPHI
		cleanPHIdesc  = "No PHI detected in this issue"
		cleanPHIcolor = "56B4E9"
		dirtyPHIlabel = LabelDirtyPHI
		dirtyPHIdesc  = "PHI detected in this issue"
		dirtyPHIcolor = "D55E00"
	)
	// TODO : use config to allow user to override default label names, descriptions, and colors

	return []*github.Label{
		{
			Name:        &cleanPHIlabel,
			Description: &cleanPHIdesc,
			Color:       &cleanPHIcolor,
		},
		{
			Name:        &dirtyPHIlabel,
			Description: &dirtyPHIdesc,
			Color:       &dirtyPHIcolor,
		},
	}
}

func hasLabelMatch(source_labels []*github.Label, label *github.Label) (bool, bool) {
	var nameMatch bool = false
	var descriptionMatch bool = false

	for _, sl := range source_labels {
		if sl.GetName() == label.GetName() {
			nameMatch = true
			if sl.GetDescription() == label.GetDescription() {
				descriptionMatch = true
				break
			}
		}
	}

	return nameMatch, descriptionMatch
}

func hasLabelName(labels []*github.Label, label string) bool {
	for _, l := range labels {
		if l.GetName() == label {
			return true
		}
	}
	return false
}

// listRepoLabels() function returns the list of the current labels for a GitHub repo.
func listRepoLabels(ctx context.Context, client *github.Client, owner string, repo string) (labels []*github.Label, e error) {
	var resp *github.Response

	labels, resp, e = client.Issues.ListLabels(ctx, owner, repo, nil)
	if e = checkResponse(resp, e); e != nil {
		return
	}

	return
}

// setRepoLabels() function creates or updates labels in a GitHub repo to contain the
// full set of labels that we may want to apply to issues in that repo.
func setRepoLabels(ctx context.Context, client *github.Client, owner string, repo string) error {
	// determine which labels already exist in the repo
	repoLabels, err := listRepoLabels(ctx, client, owner, repo)
	if err != nil {
		return err
	}

	// generate the desired-state labels for the repo
	labels := generateLabels()

	// ensure each desired-state label exists in the repo with matching name and description
	for _, label := range labels {
		matched_name, matched_desc := hasLabelMatch(repoLabels, label)
		if !matched_name {
			// create a new label in the repo
			_, resp, err := client.Issues.CreateLabel(ctx, owner, repo, label)
			if err = checkResponse(resp, err); err != nil {
				return err
			}
		} else if matched_name && !matched_desc {
			// update the existing label's description
			label_description := label.GetDescription()
			label_name := label.GetName()
			// get the existing Label object
			current_label, resp, err := client.Issues.GetLabel(ctx, owner, repo, label_name)
			if err = checkResponse(resp, err); err != nil {
				return err
			}
			current_label.Description = &label_description
			// edit the label to update its description
			_, resp, err = client.Issues.EditLabel(ctx, owner, repo, label_name, label)
			if err = checkResponse(resp, err); err != nil {
				return err
			}
		}
	}

	return nil
}

// updateIssueLabels() function updates the labels associated with a GitHub issue,
// including the creation of any labels that don't already exist in the repo.
func updateIssueLabels(ctx context.Context, client *github.Client, owner string, repo string, issueNum int, labelsToAdd, labelsToRemove []string) error {
	if len(labelsToAdd) == 0 && len(labelsToRemove) == 0 {
		return errors.New("cannot apply labels : input label lists are empty")
	}

	var (
		err   error
		issue *github.Issue
		resp  *github.Response
	)
	// create repo labels if they don't exist
	if err = setRepoLabels(ctx, client, owner, repo); err != nil {
		return err
	}

	// get the current details of the issue
	issue, resp, err = client.Issues.Get(ctx, owner, repo, issueNum)
	if err = checkResponse(resp, err); err != nil {
		return err
	}

	newLabels := []string{}
	// check if each label has already been applied to the issue
	for _, label := range labelsToAdd {
		if !hasLabelName(issue.Labels, label) {
			newLabels = append(newLabels, label)
		}
	}
	if len(newLabels) > 0 {
		// add newLabels to the issue
		_, resp, err = client.Issues.AddLabelsToIssue(ctx, owner, repo, issueNum, newLabels)
		if err = checkResponse(resp, err); err != nil {
			return err
		}
	} else {
		log.Ctx(ctx).Debug().Msgf("no new labels to add for %s/%s#%d", owner, repo, issueNum)
	}

	// TODO : remove labelsToRemove from the issue

	return nil
}
