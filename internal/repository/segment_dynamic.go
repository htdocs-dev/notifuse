package repository

// Fork patch 5: dynamic segments.
//
// Broadcast audiences and the contacts search evaluate a segment's stored SQL
// at query time instead of reading contact_segments, so a segment created or
// edited a minute ago sends to the right people without waiting for the build
// task. contact_segments stays the source for automations, the per-contact
// segment chips and the counts on the segments page.

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/Notifuse/notifuse/internal/domain"
	"github.com/lib/pq"
)

var dollarPlaceholder = regexp.MustCompile(`\$(\d+)`)

// toQuestionPlaceholders rewrites `$n` placeholders to squirrel's `?` in the
// order they appear and returns the args reordered to match, so the SQL can be
// nested inside a squirrel builder that renumbers everything at ToSql time.
func toQuestionPlaceholders(sqlText string, args []interface{}) (string, []interface{}, error) {
	ordered := make([]interface{}, 0, len(args))
	var convErr error
	out := dollarPlaceholder.ReplaceAllStringFunc(sqlText, func(m string) string {
		n, err := strconv.Atoi(m[1:])
		if err != nil || n < 1 || n > len(args) {
			convErr = fmt.Errorf("placeholder %s out of range for %d args", m, len(args))
			return m
		}
		ordered = append(ordered, args[n-1])
		return "?"
	})
	if convErr != nil {
		return "", nil, convErr
	}
	return out, ordered, nil
}

// segmentMembershipExpr returns a WHERE expression, for a query whose contacts
// table is aliased `c`, that matches contacts in any of the segments. A segment
// whose SQL was never generated falls back to its contact_segments rows.
func segmentMembershipExpr(ctx context.Context, db *sql.DB, segmentIDs []string) (sq.Sqlizer, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id, generated_sql, generated_args FROM segments WHERE id = ANY($1)`,
		pq.Array(segmentIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to load segment definitions: %w", err)
	}
	defer rows.Close()

	var or sq.Or
	var fallback []string
	for rows.Next() {
		var id string
		var genSQL sql.NullString
		var genArgs domain.JSONArray
		if err := rows.Scan(&id, &genSQL, &genArgs); err != nil {
			return nil, fmt.Errorf("failed to scan segment definition: %w", err)
		}
		if !genSQL.Valid || genSQL.String == "" {
			fallback = append(fallback, id)
			continue
		}
		text, args, err := toQuestionPlaceholders(genSQL.String, genArgs)
		if err != nil {
			return nil, fmt.Errorf("segment %s: %w", id, err)
		}
		or = append(or, sq.Expr("c.email IN ("+text+")", args...))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read segment definitions: %w", err)
	}
	if len(fallback) > 0 {
		or = append(or, sq.Expr(
			"EXISTS (SELECT 1 FROM contact_segments cs WHERE cs.email = c.email AND cs.segment_id = ANY(?))",
			pq.Array(fallback)))
	}
	if len(or) == 0 {
		// None of the requested segments exist: match nobody rather than everybody.
		return sq.Expr("FALSE"), nil
	}
	return or, nil
}
