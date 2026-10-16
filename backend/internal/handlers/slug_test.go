package handlers

import (
	"regexp"
	"strings"
	"testing"
)

func TestSlugBase(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"My Portfolio", "my-portfolio"},
		{"  Hello__World!! ", "hello-world"},
		{"café", "caf"},
		{"!!!", "project"},
		{"", "project"},
		// capped at 30 chars, and never ends in a dash after the cut
		{strings.Repeat("a", 50), strings.Repeat("a", 30)},
		{strings.Repeat("a", 29) + " bbb", strings.Repeat("a", 29)},
	}

	for _, c := range cases {
		if got := slugBase(c.in); got != c.want {
			t.Errorf("slugBase(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

var projectSlugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*-[0-9a-f]{4}$`)

func TestNewSlugShape(t *testing.T) {
	for _, name := range []string{"portfolio", "My Cool App!!", "", strings.Repeat("x", 100)} {
		slug, err := newSlug(name)
		if err != nil {
			t.Fatalf("newSlug(%q) error: %v", name, err)
		}
		if !projectSlugRe.MatchString(slug) {
			t.Errorf("newSlug(%q) = %q, bad shape", name, slug)
		}
		// must fit in one DNS label
		if len(slug) > 63 {
			t.Errorf("newSlug(%q) = %q is %d chars, max is 63", name, slug, len(slug))
		}
	}
}

// The server tries the project slug first, then falls back to a deployment
// prefix. These must never be confusable.
func TestSlugNeverLooksLikeDeploymentPrefix(t *testing.T) {
	for i := 0; i < 500; i++ {
		slug, err := newSlug("abcdef12")
		if err != nil {
			t.Fatal(err)
		}
		if deploymentPrefixRe.MatchString(slug) {
			t.Fatalf("project slug %q matches the deployment prefix pattern", slug)
		}
	}
}

func TestDeploymentPrefixRe(t *testing.T) {
	good := []string{"a1b2c3d4-e5f"}
	bad := []string{"a1b2c3d4e5f6", "_", "%", "a1b2c3d4-e5", "A1B2C3D4-E5F", "a1b2c3d4-e5f6", ""}

	for _, s := range good {
		if !deploymentPrefixRe.MatchString(s) {
			t.Errorf("%q should match", s)
		}
	}
	for _, s := range bad {
		if deploymentPrefixRe.MatchString(s) {
			t.Errorf("%q should NOT match", s)
		}
	}
}

func TestValidRepoURL(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"https://github.com/user/repo", true},
		{"https://github.com/user/repo.git", true},
		{"http://github.com/user/repo", false},
		{"file:///etc", false},
		{"ssh://git@github.com/x/y", false},
		{"git@github.com:x/y.git", false},
		{"ext::sh -c id", false},
		{"--upload-pack=touch /tmp/x", false},
		{"https://user:pass@github.com/x/y", false},
		{"https://", false},
		{"", false},
	}

	for _, c := range cases {
		if got := validRepoURL(c.in); got != c.want {
			t.Errorf("validRepoURL(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}