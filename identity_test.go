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
