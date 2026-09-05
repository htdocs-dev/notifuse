package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNonOpenersSegmentID(t *testing.T) {
	id := NonOpenersSegmentID("AbCdEf0123456789abcdef0123456789")
	assert.Equal(t, "nonopen_abcdef0123456789abcdef01", id)
	assert.LessOrEqual(t, len(id), 32)
	assert.Regexp(t, "^[a-z0-9_]+$", id)
	assert.Equal(t, "nonopen_short", NonOpenersSegmentID("short"))
}

func TestNonOpenersSegmentTree_IsValid(t *testing.T) {
	tree := NonOpenersSegmentTree("bc1")
	require.NoError(t, tree.Validate())

	leaves := tree.Branch.Leaves
	require.Len(t, leaves, 2)
	assert.Equal(t, "email.sent", leaves[0].Leaf.ContactTimeline.Kind)
	assert.Equal(t, "at_least", leaves[0].Leaf.ContactTimeline.CountOperator)
	assert.Equal(t, 1, leaves[0].Leaf.ContactTimeline.CountValue)
	assert.Equal(t, "email.opened", leaves[1].Leaf.ContactTimeline.Kind)
	assert.Equal(t, "at_most", leaves[1].Leaf.ContactTimeline.CountOperator)
	assert.Equal(t, 0, leaves[1].Leaf.ContactTimeline.CountValue)
	for _, leaf := range leaves {
		assert.Equal(t, "bc1", *leaf.Leaf.ContactTimeline.BroadcastID)
	}
}

func TestResendToNonOpenersRequest_Validate(t *testing.T) {
	assert.Error(t, (&ResendToNonOpenersRequest{}).Validate())
	assert.Error(t, (&ResendToNonOpenersRequest{WorkspaceID: "ws"}).Validate())
	assert.NoError(t, (&ResendToNonOpenersRequest{WorkspaceID: "ws", ID: "b"}).Validate())
}

func TestTruncateName(t *testing.T) {
	assert.Equal(t, "abc", TruncateName("abc"))
	assert.Len(t, TruncateName(strings.Repeat("x", 300)), 255)
}
