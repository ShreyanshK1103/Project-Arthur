package siteurl

import "testing"

func TestForSlugDefaults(t *testing.T) {
	t.Setenv("SITE_SCHEME", "")
	t.Setenv("SITE_BASE_DOMAIN", "")

	got := ForSlug("my-site-a1b2")
	want := "http://my-site-a1b2.localhost:8080"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestForSlugProduction(t *testing.T) {
	t.Setenv("SITE_SCHEME", "https")
	t.Setenv("SITE_BASE_DOMAIN", "arthurapps.dev")

	got := ForSlug("my-site-a1b2")
	want := "https://my-site-a1b2.arthurapps.dev"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}