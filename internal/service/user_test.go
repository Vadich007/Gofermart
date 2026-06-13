package service

import (
	"testing"
	"time"
)

func TestBuildAndParseToken(t *testing.T) {
	userID := 42

	token, err := buildToken(userID)
	if err != nil {
		t.Fatalf("buildToken: %v", err)
	}

	got, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}

	if got != userID {
		t.Errorf("ParseToken returned user_id=%d, want %d", got, userID)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	cases := []struct {
		name  string
		token string
	}{
		{"empty", ""},
		{"garbage", "not.a.token"},
		{"wrong signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.badSig"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseToken(tc.token)
			if err == nil {
				t.Errorf("ParseToken(%q): expected error, got nil", tc.token)
			}
		})
	}
}

func TestBuildToken_DifferentUsers(t *testing.T) {
	t1, _ := buildToken(1)
	t2, _ := buildToken(2)

	if t1 == t2 {
		t.Error("tokens for different users must differ")
	}

	id1, _ := ParseToken(t1)
	id2, _ := ParseToken(t2)

	if id1 != 1 {
		t.Errorf("want user_id=1, got %d", id1)
	}
	if id2 != 2 {
		t.Errorf("want user_id=2, got %d", id2)
	}
}

func TestBuildToken_TTL(t *testing.T) {
	t1, _ := buildToken(1)
	time.Sleep(time.Second)
	t2, _ := buildToken(1)

	if _, err := ParseToken(t1); err != nil {
		t.Errorf("first token invalid: %v", err)
	}
	if _, err := ParseToken(t2); err != nil {
		t.Errorf("second token invalid: %v", err)
	}

	if t1 == t2 {
		t.Error("expected different tokens due to different timestamps")
	}
}
