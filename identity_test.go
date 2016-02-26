package mango

import "testing"

func TestBasicIdentityUserIDReturnsUsername(t *testing.T) {
	want := "Jeff"
	i := BasicIdentity{Username: "Jeff"}
	got := i.UserID()
	if got != want {
		t.Errorf("UserID = %q, want %q", got, want)
	}
}

func TestBasicIdentityEmailReturnsEmptyString(t *testing.T) {
	want := ""
	i := BasicIdentity{Username: "Jeff"}
	got := i.Email()
	if got != want {
		t.Errorf("Email = %q, want %q", got, want)
	}
}

func TestBasicIdentityFullnameReturnsEmptyString(t *testing.T) {
	want := ""
	i := BasicIdentity{Username: "Jeff"}
	got := i.Fullname()
	if got != want {
		t.Errorf("Fullname = %q, want %q", got, want)
	}
}

func TestBasicIdentityOrganizationReturnsEmptyString(t *testing.T) {
	want := ""
	i := BasicIdentity{Username: "Jeff"}
	got := i.Organization()
	if got != want {
		t.Errorf("Organization = %q, want %q", got, want)
	}
}
