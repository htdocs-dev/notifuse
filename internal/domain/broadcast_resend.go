package domain

import (
	"fmt"
	"strings"
)

// ResendToNonOpenersRequest asks for a new draft broadcast aimed at the
// recipients of a finished broadcast who never opened it.
type ResendToNonOpenersRequest struct {
	WorkspaceID string `json:"workspace_id"`
	ID          string `json:"id"`
}

// Validate validates the resend-to-non-openers request
func (r *ResendToNonOpenersRequest) Validate() error {
	if r.WorkspaceID == "" {
		return fmt.Errorf("workspace_id is required")
	}
	if r.ID == "" {
		return fmt.Errorf("broadcast id is required")
	}
	return nil
}

// nonOpenersSegmentPrefix + 24 hex chars of the broadcast ID stays within the
// 32-character, [a-z0-9_] segment ID rule.
const nonOpenersSegmentPrefix = "nonopen_"

// NonOpenersSegmentID is the deterministic ID of the "did not open" segment of
// a broadcast, so a second resend of the same broadcast reuses the segment
// instead of creating another one.
func NonOpenersSegmentID(broadcastID string) string {
	id := strings.ToLower(broadcastID)
	if len(id) > 24 {
		id = id[:24]
	}
	return nonOpenersSegmentPrefix + id
}

// NonOpenersSegmentTree matches contacts who were sent the broadcast and never
// opened it. It does not filter on list status: the broadcast audience query
// already skips unsubscribed, bounced and complained rows at send time.
func NonOpenersSegmentTree(broadcastID string) *TreeNode {
	return &TreeNode{
		Kind: "branch",
		Branch: &TreeNodeBranch{
			Operator: "and",
			Leaves: []*TreeNode{
				timelineCountLeaf("email.sent", broadcastID, "at_least", 1),
				timelineCountLeaf("email.opened", broadcastID, "at_most", 0),
			},
		},
	}
}

func timelineCountLeaf(kind, broadcastID, countOperator string, count int) *TreeNode {
	anytime := "anytime"
	return &TreeNode{
		Kind: "leaf",
		Leaf: &TreeNodeLeaf{
			Source: "contact_timeline",
			ContactTimeline: &ContactTimelineCondition{
				Kind:              kind,
				CountOperator:     countOperator,
				CountValue:        count,
				BroadcastID:       &broadcastID,
				TimeframeOperator: &anytime,
			},
		},
	}
}

// TruncateName cuts a display name to the 255-character limit shared by
// broadcast and segment names.
func TruncateName(name string) string {
	const max = 255
	if len(name) <= max {
		return name
	}
	return name[:max]
}
