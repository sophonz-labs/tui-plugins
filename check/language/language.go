// Copied to github.com/sophonz-labs/sophonz-tunnel-go internal/ext/language/language.go; keep the two in step.

// Package language reads language.json: a language pack, the plugin that
// tells a tool how to start a language server and which files it serves.
// The format is the tui-plugins spec (SPEC-language.md, v 1). Standard
// library only, so any tool can copy it.
package language

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Version is the language.json version this package reads.
const Version = 1

// Pack is one language.json.
type Pack struct {
	V           int        `json:"v"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Homepage    string     `json:"homepage,omitempty"`
	Requires    []string   `json:"requires,omitempty"` // programs that must be on PATH
	Languages   []Language `json:"languages"`
	Server      Server     `json:"server"`
	Install     []Install  `json:"install,omitempty"`
}

// Language is one language the server serves.
type Language struct {
	ID          string   `json:"id"`
	Files       []string `json:"files"`                 // base-name patterns: *.java
	RootMarkers []string `json:"rootMarkers,omitempty"` // nearest ancestor holding one is the root
	LineComment string   `json:"lineComment,omitempty"`
	Grammar     string   `json:"grammar,omitempty"` // a grammar the editor already has
}

// Server starts the language server.
type Server struct {
	Command []string `json:"command"`           // argv, never a shell
	Version []string `json:"version,omitempty"` // arguments that ask its version; a failure counts as missing
	Needs   []string `json:"needs,omitempty"`   // programs it needs besides itself
}

// Install is one step that puts the server on the machine. Exactly one kind.
type Install struct {
	Npm *Package `json:"npm,omitempty"`
	Go  *Module  `json:"go,omitempty"`
}

type Package struct {
	Package string `json:"package"`
	Version string `json:"version"`
}

type Module struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

// Argv is the command the step runs.
func (i Install) Argv() []string {
	switch {
	case i.Npm != nil:
		return []string{"npm", "install", "-g", i.Npm.Package + "@" + i.Npm.Version}
	case i.Go != nil:
		return []string{"go", "install", i.Go.Module + "@" + i.Go.Version}
	}
	return nil
}

var (
	nameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,31}$`)
	idRe   = regexp.MustCompile(`^[a-z0-9][a-z0-9_+-]{0,31}$`)
	// Exact versions only: a range or a tag installs something nobody read.
	npmVersionRe = regexp.MustCompile(`^\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$`)
	goVersionRe  = regexp.MustCompile(`^v\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$`)
	npmPackageRe = regexp.MustCompile(`^(@[a-z0-9][a-z0-9._-]*/)?[a-z0-9][a-z0-9._-]*$`)
	goModuleRe   = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._~-]*(/[A-Za-z0-9._~-]+)+$`)
)

// Parse reads and checks a pack. Unknown fields are an error, as in
// plugin.json: a misspelt field must not quietly change what runs.
func Parse(data []byte) (*Pack, error) {
	var probe struct {
		V int `json:"v"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}
	if probe.V > Version {
		return nil, fmt.Errorf("language.json v %d needs a newer version (this reads v %d)", probe.V, Version)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var p Pack
	if err := dec.Decode(&p); err != nil {
		return nil, err
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

// Validate is every rule of the spec that one pack can check alone.
func (p *Pack) Validate() error {
	var errs []error
	bad := func(f string, a ...any) { errs = append(errs, fmt.Errorf(f, a...)) }
	if p.V != Version {
		bad(`"v" must be %d`, Version)
	}
	if !nameRe.MatchString(p.Name) {
		bad("name %q: lower-case letters, digits and -, at most 32", p.Name)
	}
	if len(p.Languages) == 0 {
		bad("no languages")
	}
	ids := map[string]bool{}
	for i, l := range p.Languages {
		at := fmt.Sprintf("languages[%d] %q", i, l.ID)
		if !idRe.MatchString(l.ID) {
			bad("%s: id: lower-case letters, digits, _, + and -, at most 32", at)
		}
		if ids[l.ID] {
			bad("%s: id appears twice", at)
		}
		ids[l.ID] = true
		if len(l.Files) == 0 {
			bad("%s: no files", at)
		}
		for _, f := range l.Files {
			if f == "" || strings.ContainsAny(f, `/\`) {
				bad("%s: files: %q is not a file-name pattern", at, f)
			}
		}
		for _, m := range l.RootMarkers {
			if m == "" || strings.ContainsAny(m, `/\`) || m == "." || m == ".." {
				bad("%s: rootMarkers: %q is not a file name", at, m)
			}
		}
	}
	argv := func(field string, a []string, required bool) {
		if len(a) == 0 {
			if required {
				bad("server.%s: required, as argv", field)
			}
			return
		}
		for _, s := range a {
			if s == "" {
				bad("server.%s: an empty argument", field)
			}
		}
		// A shell line in one string: argv only, so "a b" as the program is refused.
		if field == "command" && strings.ContainsAny(a[0], " \t\n;|&$`") {
			bad("server.command: %q is a shell line; give argv, one argument per string", a[0])
		}
	}
	argv("command", p.Server.Command, true)
	argv("version", p.Server.Version, false)
	for _, n := range append(append([]string{}, p.Requires...), p.Server.Needs...) {
		if n == "" || strings.ContainsAny(n, " \t/\\") {
			bad("%q is not a program name", n)
		}
	}
	for i, in := range p.Install {
		at := fmt.Sprintf("install[%d]", i)
		switch {
		case in.Npm != nil && in.Go != nil:
			bad("%s: one of npm or go, not both", at)
		case in.Npm != nil:
			if !npmPackageRe.MatchString(in.Npm.Package) {
				bad("%s: npm package %q", at, in.Npm.Package)
			}
			if !npmVersionRe.MatchString(in.Npm.Version) {
				bad("%s: npm version %q is not exact (1.2.3)", at, in.Npm.Version)
			}
		case in.Go != nil:
			if !goModuleRe.MatchString(in.Go.Module) {
				bad("%s: go module %q", at, in.Go.Module)
			}
			if !goVersionRe.MatchString(in.Go.Version) {
				bad("%s: go version %q is not exact (v1.2.3)", at, in.Go.Version)
			}
		default:
			bad("%s: one of npm or go", at)
		}
	}
	return errors.Join(errs...)
}
