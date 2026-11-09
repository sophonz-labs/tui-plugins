package language

import (
	"strings"
	"testing"
)

const good = `{
  "v": 1,
  "name": "typescript",
  "languages": [{"id": "typescript", "files": ["*.ts", "*.tsx"], "rootMarkers": ["tsconfig.json"], "lineComment": "//", "grammar": "typescript"}],
  "server": {"command": ["typescript-language-server", "--stdio"], "needs": ["node"]},
  "install": [{"npm": {"package": "typescript-language-server", "version": "6.0.3"}}, {"go": {"module": "golang.org/x/tools/gopls", "version": "v0.20.0"}}]
}`

func TestGood(t *testing.T) {
	p, err := Parse([]byte(good))
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(p.Install[0].Argv(), " "); got != "npm install -g typescript-language-server@6.0.3" {
		t.Fatal(got)
	}
	if got := strings.Join(p.Install[1].Argv(), " "); got != "go install golang.org/x/tools/gopls@v0.20.0" {
		t.Fatal(got)
	}
}

func TestRefused(t *testing.T) {
	for _, c := range []struct{ name, from, to, want string }{
		{"unknown field", `"v": 1,`, `"v": 1, "theme": null,`, "unknown field"},
		{"higher v", `"v": 1,`, `"v": 2,`, "newer version"},
		{"shell string", `["typescript-language-server", "--stdio"]`, `"typescript-language-server --stdio"`, "cannot unmarshal"},
		{"shell line in argv", `["typescript-language-server", "--stdio"]`, `["sh -c typescript-language-server"]`, "shell line"},
		{"other install kind", `{"npm": {"package": "typescript-language-server", "version": "6.0.3"}}`, `{"url": "https://x/y.tgz"}`, "unknown field"},
		{"empty install step", `{"npm": {"package": "typescript-language-server", "version": "6.0.3"}}`, `{}`, "one of npm or go"},
		{"npm range", `"6.0.3"`, `"^6.0.3"`, "not exact"},
		{"npm tag", `"6.0.3"`, `"latest"`, "not exact"},
		{"go latest", `"v0.20.0"`, `"latest"`, "not exact"},
		{"bad name", `"name": "typescript"`, `"name": "Type Script"`, "name"},
		{"no languages", `"languages": [{"id": "typescript", "files": ["*.ts", "*.tsx"], "rootMarkers": ["tsconfig.json"], "lineComment": "//", "grammar": "typescript"}]`, `"languages": []`, "no languages"},
		{"path as pattern", `"*.ts", "*.tsx"`, `"src/*.ts"`, "file-name pattern"},
	} {
		data := strings.Replace(good, c.from, c.to, 1)
		if data == good {
			t.Fatalf("%s: the replacement did not apply", c.name)
		}
		_, err := Parse([]byte(data))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: got %v, want an error with %q", c.name, err, c.want)
		}
	}
}
