package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

func TestToQuestionPlaceholders(t *testing.T) {
	text, args, err := toQuestionPlaceholders(
		"SELECT email FROM contacts WHERE a = $2 AND b = $1 AND c = $2",
		[]interface{}{"one", 2.0})
	require.NoError(t, err)
	require.Equal(t, "SELECT email FROM contacts WHERE a = ? AND b = ? AND c = ?", text)
	require.Equal(t, []interface{}{2.0, "one", 2.0}, args)

	_, _, err = toQuestionPlaceholders("WHERE a = $3", []interface{}{"one"})
	require.Error(t, err)
}

func TestSegmentMembershipExpr(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT id, generated_sql, generated_args FROM segments WHERE id = ANY\(\$1\)`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "generated_sql", "generated_args"}).
			AddRow("vip", "SELECT email FROM contacts WHERE custom_number_1 >= $1", []byte(`[500]`)).
			AddRow("legacy", nil, nil))

	expr, err := segmentMembershipExpr(context.Background(), db, []string{"vip", "legacy"})
	require.NoError(t, err)

	sqlText, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("c.email").From("contacts c").Where(sq.Eq{"c.country": "FR"}).Where(expr).ToSql()
	require.NoError(t, err)
	require.Equal(t,
		"SELECT c.email FROM contacts c WHERE c.country = $1 AND (c.email IN (SELECT email FROM contacts WHERE custom_number_1 >= $2) OR EXISTS (SELECT 1 FROM contact_segments cs WHERE cs.email = c.email AND cs.segment_id = ANY($3)))",
		sqlText)
	require.Len(t, args, 3)
	require.EqualValues(t, 500, args[1])
	require.NoError(t, mock.ExpectationsWereMet())
}
