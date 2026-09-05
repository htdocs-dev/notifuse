package domain

import "testing"

func TestTagContactsRequestValidate(t *testing.T) {
	req := &TagContactsRequest{WorkspaceID: "ws", Emails: []string{" A@B.com "}, Tags: []string{" vip "}}
	if err := req.Validate(); err != nil {
		t.Fatal(err)
	}
	if req.Emails[0] != "a@b.com" || req.Tags[0] != "vip" {
		t.Fatalf("not normalised: %v %v", req.Emails, req.Tags)
	}
	for _, bad := range []*TagContactsRequest{
		{Emails: []string{"a@b.com"}, Tags: []string{"x"}},
		{WorkspaceID: "ws", Tags: []string{"x"}},
		{WorkspaceID: "ws", Emails: []string{"a@b.com"}},
		{WorkspaceID: "ws", Emails: []string{"a@b.com"}, Tags: []string{" "}},
	} {
		if err := bad.Validate(); err == nil {
			t.Fatalf("expected error for %+v", bad)
		}
	}
}
