package gitlab

import (
	"context"

	"github.com/karlderkaefer/cdk-notifier/config"
	"github.com/xanzy/go-gitlab"
)

const (
	// HeaderPrefix default prefix for comment message
	HeaderPrefix = "## cdk diff for"
)

// NotesService interface for required Gitlab actions with API
type NotesService interface {
	ListMergeRequestNotes(pid interface{}, mergeRequest int, opt *gitlab.ListMergeRequestNotesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Note, *gitlab.Response, error)
	// DeleteComment(ctx context.Context, owner string, repo string, commentID int64) (*gitlab.Response, error)
	// EditComment(ctx context.Context, owner string, repo string, commentID int64, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
	// CreateComment(ctx context.Context, owner string, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
}

// Client Gitlab client configuration
type Client struct {
	Notes  NotesService
	Pid    interface{}
	Client *gitlab.Client

	Token          string
	Owner          string
	Repo           string
	TagID          string
	NoteContent    string
	MergeRequestID int
	DeleteNotes    bool
}

func NewGitlabClient(ctx context.Context, config *config.AppConfig, notesMock NotesService) *Client {
	gitlabClient := &Client{
		Owner:          config.RepoOwner,
		Repo:           config.RepoName,
		TagID:          config.TagID,
		MergeRequestID: config.PullRequest, //TOdo
		DeleteNotes:    config.DeleteComment,
		Token:          config.GithubToken, //TODO
	}

	if notesMock != nil {
		gitlabClient.Notes = notesMock
	} else {
		gitlabClient.Client, _ = gitlab.NewClient(config.GithubToken) //TODO ERROR/naming
		gitlabClient.Notes = gitlabClient.Client.Notes
	}
	return gitlabClient
}

func (gc *Client) Authenticate() {
	gc.Client, _ = gitlab.NewClient(gc.Token)
}

func (gc *Client) ListMergeRequestNotes() ([]*gitlab.Note, error) {
	opt := &gitlab.ListMergeRequestNotesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	}
	notes, _, err := gc.Notes.ListMergeRequestNotes(gc.Pid, gc.MergeRequestID, opt)
	if err != nil {
		return nil, err
	}
	return notes, nil
}
