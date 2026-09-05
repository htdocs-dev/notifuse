package service

// Fork patch 7: contact tags. Thin authorisation layer over the repository.

import (
	"context"
	"fmt"

	"github.com/Notifuse/notifuse/internal/domain"
)

func (s *ContactService) AddContactTags(ctx context.Context, workspaceID string, emails, tags []string) (int, error) {
	ctx, err := s.authorizeContactsWrite(ctx, workspaceID)
	if err != nil {
		return 0, err
	}
	return s.repo.AddContactTags(ctx, workspaceID, emails, tags)
}

func (s *ContactService) RemoveContactTags(ctx context.Context, workspaceID string, emails, tags []string) (int, error) {
	ctx, err := s.authorizeContactsWrite(ctx, workspaceID)
	if err != nil {
		return 0, err
	}
	return s.repo.RemoveContactTags(ctx, workspaceID, emails, tags)
}

func (s *ContactService) authorizeContactsWrite(ctx context.Context, workspaceID string) (context.Context, error) {
	ctx, _, userWorkspace, err := s.authService.AuthenticateUserForWorkspace(ctx, workspaceID)
	if err != nil {
		return ctx, fmt.Errorf("failed to authenticate user: %w", err)
	}
	if !userWorkspace.HasPermission(domain.PermissionResourceContacts, domain.PermissionTypeWrite) {
		return ctx, domain.NewPermissionError(
			domain.PermissionResourceContacts,
			domain.PermissionTypeWrite,
			"Insufficient permissions: write access to contacts required",
		)
	}
	return ctx, nil
}
