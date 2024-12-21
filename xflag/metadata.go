package xflag

import (
	"fmt"
	"strings"
)

const (
	_defaultPadding  int    = 7
	_defaultSynopsis string = "[global options]"
)

// Option represents a single command-line option.
type Option struct {
	// Name is the name of the option.
	Name string

	// Shorthand is a single-character shorthand for the option.
	Shorthand string

	// Type is the type of value the option accepts (e.g., "FILE", "DIRECTORY").
	Type string

	// Description is a short description of the option.
	Description string

	// Default is the default value for the option.
	Default string
}

// Metadata represents metadata for a command-line application.
type Metadata struct {
	// Name is the name of the application.
	Name string

	// Description is a short description of the application.
	Description string

	// Ver is the version of the application.
	Ver string

	// Synopsis is the command usage pattern shown in the usage section. Defaults to "[global options]" if empty.
	Synopsis string

	// Options is a list of options for the application.
	Options []Option
}

// Usage returns the formatted usage information for the application.
//
// The returned string contains the following sections:
//
//  1. NAME: Application name and brief description.
//  2. USAGE: Command-line usage pattern with global options.
//  3. VERSION: Current version of the application.
//  4. GLOBAL OPTIONS: Available command-line options with their details and default values.
func (m *Metadata) Usage() string {
	var builder strings.Builder

	synopsis := m.Synopsis
	if synopsis == "" {
		synopsis = _defaultSynopsis
	}

	fmt.Fprintf(&builder, "NAME:\n   %s - %s\n\n", m.Name, m.Description)
	fmt.Fprintf(&builder, "USAGE:\n   %s %s\n\n", m.Name, synopsis)
	fmt.Fprintf(&builder, "VERSION:\n   %s\n\n", m.Ver)

	if len(m.Options) > 0 {
		fmt.Fprintf(&builder, "GLOBAL OPTIONS:\n")

		fmt.Fprint(&builder, m.formatOptions())
	}

	return builder.String()
}

// Version returns the version string of the application.
func (m *Metadata) Version() string {
	return m.Name + " version " + m.Ver
}

// calculateMaxLength returns the maximum length needed for option formatting.
func (m *Metadata) calculateMaxLength() int {
	var maxLen int

	for _, opt := range m.Options {
		length := len(opt.Name) + len(opt.Shorthand) + 5
		if opt.Type != "" {
			length += len(opt.Type)*2 + 4
		}

		if length > maxLen {
			maxLen = length
		}
	}

	return maxLen + _defaultPadding
}

// formatOptions returns the formatted options section.
func (m *Metadata) formatOptions() string {
	var (
		builder strings.Builder
		maxLen  = m.calculateMaxLength()
	)

	for _, opt := range m.Options {
		var optPart string

		if opt.Type != "" {
			optPart = fmt.Sprintf("   --%s <%s>, -%s <%s>",
				opt.Name, opt.Type,
				opt.Shorthand, opt.Type)
		} else {
			optPart = fmt.Sprintf("   --%s, -%s",
				opt.Name, opt.Shorthand)
		}

		padding := strings.Repeat(" ", maxLen-len(optPart))

		if opt.Default != "" {
			fmt.Fprintf(&builder, "%s%s%s (default: %s)\n",
				optPart,
				padding,
				opt.Description,
				opt.Default,
			)
		} else {
			fmt.Fprintf(&builder, "%s%s%s\n",
				optPart,
				padding,
				opt.Description,
			)
		}
	}

	return builder.String()
}
