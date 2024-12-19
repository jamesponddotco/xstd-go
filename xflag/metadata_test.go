package xflag_test

import (
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xflag"
)

func TestMetadata_Usage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give xflag.Metadata
		want string
	}{
		{
			name: "Empty metadata",
			give: xflag.Metadata{},
			want: "NAME:\n    - \n\n" +
				"USAGE:\n    [global options]\n\n" +
				"VERSION:\n   \n\n",
		},
		{
			name: "Basic metadata without options",
			give: xflag.Metadata{
				Name:        "myapp",
				Description: "does cool things",
				Ver:         "1.0.0",
			},
			want: "NAME:\n   myapp - does cool things\n\n" +
				"USAGE:\n   myapp [global options]\n\n" +
				"VERSION:\n   1.0.0\n\n",
		},
		{
			name: "Custom synopsis",
			give: xflag.Metadata{
				Name:        "myapp",
				Description: "does cool things",
				Ver:         "1.0.0",
				Synopsis:    "command [arguments]",
			},
			want: "NAME:\n   myapp - does cool things\n\n" +
				"USAGE:\n   myapp command [arguments]\n\n" +
				"VERSION:\n   1.0.0\n\n",
		},
		{
			name: "Single option",
			give: xflag.Metadata{
				Name:        "myapp",
				Description: "does cool things",
				Ver:         "1.0.0",
				Options: []xflag.Option{
					{
						Name:        "config",
						Shorthand:   "c",
						Type:        "FILE",
						Description: "config file path",
						Default:     "/etc/myapp.conf",
					},
				},
			},
			want: "NAME:\n   myapp - does cool things\n\n" +
				"USAGE:\n   myapp [global options]\n\n" +
				"VERSION:\n   1.0.0\n\n" +
				"GLOBAL OPTIONS:\n\n" +
				"   --config <FILE>, -c <FILE>  config file path (default: /etc/myapp.conf)\n",
		},
		{
			name: "Multiple options with various formats",
			give: xflag.Metadata{
				Name:        "myapp",
				Description: "does cool things",
				Ver:         "1.0.0",
				Options: []xflag.Option{
					{
						Name:        "config",
						Shorthand:   "c",
						Type:        "FILE",
						Description: "config file path",
						Default:     "/etc/myapp.conf",
					},
					{
						Name:        "verbose",
						Shorthand:   "v",
						Description: "enable verbose logging",
					},
					{
						Name:        "port",
						Shorthand:   "p",
						Type:        "NUMBER",
						Description: "server port",
					},
				},
			},
			want: "NAME:\n   myapp - does cool things\n\n" +
				"USAGE:\n   myapp [global options]\n\n" +
				"VERSION:\n   1.0.0\n\n" +
				"GLOBAL OPTIONS:\n\n" +
				"   --config <FILE>, -c <FILE>    config file path (default: /etc/myapp.conf)\n" +
				"   --verbose, -v                 enable verbose logging\n" +
				"   --port <NUMBER>, -p <NUMBER>  server port\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.give.Usage()
			got = strings.ReplaceAll(got, "\r\n", "\n")

			want := strings.ReplaceAll(tt.want, "\r\n", "\n")

			if got != want {
				t.Errorf("\ngot:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

func TestMetadata_Version(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give xflag.Metadata
		want string
	}{
		{
			name: "Standard version",
			give: xflag.Metadata{
				Name: "myapp",
				Ver:  "1.0.0",
			},
			want: "myapp version 1.0.0",
		},
		{
			name: "Empty version",
			give: xflag.Metadata{
				Name: "myapp",
			},
			want: "myapp version ",
		},
		{
			name: "Empty name",
			give: xflag.Metadata{
				Ver: "1.0.0",
			},
			want: " version 1.0.0",
		},
		{
			name: "Empty both",
			give: xflag.Metadata{},
			want: " version ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.give.Version()
			if got != tt.want {
				t.Errorf("Version() = %q, want %q", got, tt.want)
			}
		})
	}
}
