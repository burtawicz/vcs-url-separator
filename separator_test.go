package vcs_url_separator

import (
	"reflect"
	"testing"
)

func Test_stripHttpPrefix(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			"valid link - http prefix",
			"http://github.com/burtawicz/vcs-url-separator",
			"github.com/burtawicz/vcs-url-separator",
		},
		{
			"valid link - https prefix",
			"https://github.com/burtawicz/vcs-url-separator",
			"github.com/burtawicz/vcs-url-separator",
		},
		{
			"valid link - no prefix",
			"github.com/burtawicz/vcs-url-separator",
			"github.com/burtawicz/vcs-url-separator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripHttpPrefix(tt.url); got != tt.want {
				t.Errorf("stripHttpPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stripTopLevelDomain(t *testing.T) {
	tests := []struct {
		name        string
		providerUrl string
		want        string
	}{
		{
			"dot com",
			"github.com",
			"github",
		},
		{
			"dot org",
			"bitbucket.org",
			"bitbucket",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripTopLevelDomain(tt.providerUrl); got != tt.want {
				t.Errorf("stripTopLevelDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchProvider(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		want         string
	}{
		{
			"github",
			"github",
			"GitHub",
		},
		{
			"bitbucket",
			"bitbucket",
			"BitBucket",
		},
		{
			"gitlab",
			"gitlab",
			"GitLab",
		},
		{
			"unknown",
			"unknown",
			"unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchProvider(tt.providerName); got != tt.want {
				t.Errorf("matchProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeparateVcsUrl(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    UrlParts
		wantErr bool
	}{
		{
			"happy path - no nested directories - no file path",
			"https://github.com/burtawicz/vcs-url-separator",
			UrlParts{
				Provider:       "GitHub",
				Organization:   "burtawicz",
				Project:        "vcs-url-separator",
				SubDirectories: make([]string, 0),
				FilePath:       "",
			},
			false,
		},
		{
			"happy path - with nested directories - no file path",
			"https://github.com/burtawicz/vcs-url-separator/some-nested-dir/another-nested-dir",
			UrlParts{
				Provider:       "GitHub",
				Organization:   "burtawicz",
				Project:        "vcs-url-separator",
				SubDirectories: []string{"some-nested-dir", "another-nested-dir"},
				FilePath:       "",
			},
			false,
		},
		{
			"happy path - no nested directories - with file path",
			"https://github.com/burtawicz/vcs-url-separator/somefile.py",
			UrlParts{
				Provider:       "GitHub",
				Organization:   "burtawicz",
				Project:        "vcs-url-separator",
				SubDirectories: make([]string, 0),
				FilePath:       "somefile.py",
			},
			false,
		},
		{
			"happy path - with nested directories - with file path",
			"https://github.com/burtawicz/vcs-url-separator/some-nested-dir/another-nested-dir/somefile.py",
			UrlParts{
				Provider:       "GitHub",
				Organization:   "burtawicz",
				Project:        "vcs-url-separator",
				SubDirectories: []string{"some-nested-dir", "another-nested-dir"},
				FilePath:       "somefile.py",
			},
			false,
		},
		{
			"unhappy path - multiple file paths",
			"https://github.com/burtawicz/vcs-url-separator/some-file.py/some-other-file.js",
			UrlParts{},
			true,
		},
		{
			"unhappy path - invalid characters",
			"https://github.com/burtawicz/vcs-url-separator|[]{}\\",
			UrlParts{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := SeparateVcsUrl(tt.url)
			if (gotErr != nil) != tt.wantErr {
				t.Fatalf("SeparateVcsUrl() gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SeparateVcsUrl() got: %v, want: %v", got, tt.want)
			}
		})
	}
}
