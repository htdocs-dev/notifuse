package domain

// Fork patch 7: contact tags.
//
// Tags are a flat JSON array of strings kept in one custom JSON field, so they
// need no schema change, show up in the API like any custom field, and segments
// match them with the existing `in_array` operator on that field.

import (
	"context"
	"fmt"
	"strings"
)

// ContactTagsField is the contacts column that holds the tags array.
const ContactTagsField = "custom_json_1"

// maxTagBatch caps one tag/untag call; larger jobs page through the API.
const maxTagBatch = 1000

// TagContactsRequest adds or removes tags on a set of contacts.
type TagContactsRequest struct {
	WorkspaceID string   `json:"workspace_id"`
	Emails      []string `json:"emails"`
	Tags        []string `json:"tags"`
}

// Validate normalises emails and tags in place and rejects empty input.
func (r *TagContactsRequest) Validate() error {
	if r.WorkspaceID == "" {
		return fmt.Errorf("workspace_id is required")
	}
	if len(r.Emails) == 0 || len(r.Emails) > maxTagBatch {
		return fmt.Errorf("emails must contain between 1 and %d addresses", maxTagBatch)
	}
	for i, email := range r.Emails {
		r.Emails[i] = NormalizeEmail(email)
		if r.Emails[i] == "" {
			return fmt.Errorf("emails[%d] is empty", i)
		}
	}
	if len(r.Tags) == 0 {
		return fmt.Errorf("tags is required")
	}
	for i, tag := range r.Tags {
		r.Tags[i] = strings.TrimSpace(tag)
		if r.Tags[i] == "" {
			return fmt.Errorf("tags[%d] is empty", i)
		}
	}
	return nil
}

// ContactTagger is the tag/untag surface shared by the service and repository.
type ContactTagger interface {
	// AddContactTags adds tags to the contacts that do not have them all yet
	// and returns how many contacts changed.
	AddContactTags(ctx context.Context, workspaceID string, emails, tags []string) (int, error)
	// RemoveContactTags removes tags from the contacts that carry any of them
	// and returns how many contacts changed.
	RemoveContactTags(ctx context.Context, workspaceID string, emails, tags []string) (int, error)
}
