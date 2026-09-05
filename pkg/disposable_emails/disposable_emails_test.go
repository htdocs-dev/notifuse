package disposable_emails

import (
	"slices"
	"testing"
)

func TestIsDisposableEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{email: "test@example.com", want: false},
		{email: "example.com", want: false},
		{email: "", want: false},
		{email: "@", want: false},
		{email: "0-180.com", want: true},
		{email: "someone@0-180.com", want: true},
		{email: "Someone@0-180.COM", want: true},
		{email: " someone@zzzz1717.com ", want: true},
	}

	for _, test := range tests {
		if got := IsDisposableEmail(test.email); got != test.want {
			t.Errorf("IsDisposableEmail(%q) = %v, want %v", test.email, got, test.want)
		}
	}
}

// The lookup is a binary search, so the list must stay sorted in byte order
// when it is refreshed from upstream.
func TestDisposableListIsSorted(t *testing.T) {
	if !slices.IsSorted(disposableEmails) {
		t.Fatal("disposableEmails must be sorted in byte order (LC_ALL=C sort)")
	}
}
