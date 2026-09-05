package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Notifuse/notifuse/internal/domain"
)

// SetSegmentService injects the segment service used by ResendToNonOpeners.
// Kept out of the constructor so existing callers are untouched.
func (s *BroadcastService) SetSegmentService(segmentService domain.SegmentService) {
	s.segmentService = segmentService
}

// ErrBroadcastNotProcessed is returned when a resend is asked for a broadcast
// that has not finished sending.
var ErrBroadcastNotProcessed = errors.New("broadcast has not finished sending")

// ResendToNonOpeners creates (or refreshes) the "did not open" segment of a
// processed broadcast and returns a new draft broadcast aimed at it. The draft
// keeps the list, the winning (or only) template, the UTM parameters and the
// data feed of the original. Nothing is sent: the operator schedules the draft.
func (s *BroadcastService) ResendToNonOpeners(ctx context.Context, request *domain.ResendToNonOpenersRequest) (*domain.Broadcast, error) {
	ctx, _, userWorkspace, err := s.authService.AuthenticateUserForWorkspace(ctx, request.WorkspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate user: %w", err)
	}
	if !userWorkspace.HasPermission(domain.PermissionResourceBroadcasts, domain.PermissionTypeWrite) {
		return nil, domain.NewPermissionError(
			domain.PermissionResourceBroadcasts,
			domain.PermissionTypeWrite,
			"Insufficient permissions: write access to broadcasts required",
		)
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}
	if s.segmentService == nil {
		return nil, fmt.Errorf("resend to non-openers is not available: segment service not configured")
	}

	original, err := s.repo.GetBroadcast(ctx, request.WorkspaceID, request.ID)
	if err != nil {
		return nil, err
	}
	if original.Status != domain.BroadcastStatusProcessed {
		return nil, fmt.Errorf("%w: status is %s", ErrBroadcastNotProcessed, original.Status)
	}

	templateID := ""
	if original.WinningTemplate != nil && *original.WinningTemplate != "" {
		templateID = *original.WinningTemplate
	} else if len(original.TestSettings.Variations) > 0 {
		templateID = original.TestSettings.Variations[0].TemplateID
	}
	if templateID == "" {
		return nil, fmt.Errorf("broadcast %s has no template to resend", original.ID)
	}

	segmentID, err := s.ensureNonOpenersSegment(ctx, request.WorkspaceID, original)
	if err != nil {
		return nil, err
	}

	return s.CreateBroadcast(ctx, &domain.CreateBroadcastRequest{
		WorkspaceID: request.WorkspaceID,
		Name:        domain.TruncateName("Resend to non-openers: " + original.Name),
		Audience: domain.AudienceSettings{
			List:                original.Audience.List,
			Segments:            []string{segmentID},
			ExcludeUnsubscribed: true,
		},
		TestSettings: domain.BroadcastTestSettings{
			Enabled:    false,
			Variations: []domain.BroadcastVariation{{VariationName: "A", TemplateID: templateID}},
		},
		UTMParameters: original.UTMParameters,
		DataFeed:      original.DataFeed,
		Metadata:      domain.MapOfAny{"resend_of_broadcast_id": original.ID},
	})
}

// ensureNonOpenersSegment returns the ID of the broadcast's "did not open"
// segment, creating it on first use and rebuilding it on later ones so the
// membership reflects opens recorded since.
func (s *BroadcastService) ensureNonOpenersSegment(ctx context.Context, workspaceID string, original *domain.Broadcast) (string, error) {
	segmentID := domain.NonOpenersSegmentID(original.ID)

	_, err := s.segmentService.GetSegment(ctx, &domain.GetSegmentRequest{WorkspaceID: workspaceID, ID: segmentID})
	if err == nil {
		if err := s.segmentService.RebuildSegment(ctx, workspaceID, segmentID); err != nil {
			return "", fmt.Errorf("failed to rebuild non-openers segment: %w", err)
		}
		return segmentID, nil
	}
	var notFound *domain.ErrSegmentNotFound
	if !errors.As(err, &notFound) {
		return "", fmt.Errorf("failed to look up non-openers segment: %w", err)
	}

	timezone := "UTC"
	if workspace, wsErr := s.workspaceRepo.GetByID(ctx, workspaceID); wsErr == nil && workspace.Settings.Timezone != "" {
		timezone = workspace.Settings.Timezone
	}

	_, err = s.segmentService.CreateSegment(ctx, &domain.CreateSegmentRequest{
		WorkspaceID: workspaceID,
		ID:          segmentID,
		Name:        domain.TruncateName("Did not open: " + original.Name),
		Color:       "orange",
		Tree:        domain.NonOpenersSegmentTree(original.ID),
		Timezone:    timezone,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create non-openers segment: %w", err)
	}
	return segmentID, nil
}
