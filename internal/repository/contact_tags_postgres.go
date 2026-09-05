package repository

// Fork patch 7: contact tags stored in domain.ContactTagsField.

import (
	"context"
	"fmt"

	"github.com/Notifuse/notifuse/internal/domain"
	"github.com/lib/pq"
)

// tagsArray reads the field as a JSON array, treating any other shape as empty.
var tagsArray = fmt.Sprintf(
	"(CASE WHEN jsonb_typeof(%s) = 'array' THEN %s ELSE '[]'::jsonb END)",
	domain.ContactTagsField, domain.ContactTagsField)

// addContactTagsSQL sets the field to the sorted union of current and new tags,
// touching only contacts that miss at least one of them.
var addContactTagsSQL = fmt.Sprintf(`
	UPDATE contacts SET %[1]s = (
		SELECT COALESCE(jsonb_agg(DISTINCT t), '[]'::jsonb) FROM (
			SELECT jsonb_array_elements_text(%[2]s) AS t
			UNION SELECT unnest($2::text[])
		) tags
	), updated_at = NOW()
	WHERE email = ANY($1) AND NOT (%[2]s @> to_jsonb($2::text[]))`,
	domain.ContactTagsField, tagsArray)

// removeContactTagsSQL drops the tags from contacts that carry any of them.
var removeContactTagsSQL = fmt.Sprintf(`
	UPDATE contacts SET %[1]s = (
		SELECT COALESCE(jsonb_agg(t), '[]'::jsonb)
		FROM jsonb_array_elements_text(%[1]s) AS t
		WHERE t <> ALL($2::text[])
	), updated_at = NOW()
	WHERE email = ANY($1) AND jsonb_typeof(%[1]s) = 'array' AND %[1]s ?| $2::text[]`,
	domain.ContactTagsField)

func (r *contactRepository) AddContactTags(ctx context.Context, workspaceID string, emails, tags []string) (int, error) {
	return r.execTags(ctx, workspaceID, addContactTagsSQL, emails, tags)
}

func (r *contactRepository) RemoveContactTags(ctx context.Context, workspaceID string, emails, tags []string) (int, error) {
	return r.execTags(ctx, workspaceID, removeContactTagsSQL, emails, tags)
}

func (r *contactRepository) execTags(ctx context.Context, workspaceID, query string, emails, tags []string) (int, error) {
	db, err := r.workspaceRepo.GetConnection(ctx, workspaceID)
	if err != nil {
		return 0, fmt.Errorf("failed to get workspace connection: %w", err)
	}
	res, err := db.ExecContext(ctx, query, pq.Array(emails), pq.Array(tags))
	if err != nil {
		return 0, fmt.Errorf("failed to update contact tags: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to read updated contact count: %w", err)
	}
	return int(n), nil
}
