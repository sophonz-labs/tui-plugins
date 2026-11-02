// check validates every plugins/<name>/plugin.json against SPEC.md, and
// the rules a single manifest cannot check: its directory, the views its
// scopes and opens name, no two plugins taking one : name, and a line for
// it in README.md.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sophonz-labs/tui-plugins/check/plugin"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	files, _ := filepath.Glob(filepath.Join(root, "plugins", "*", "plugin.json"))
	sort.Strings(files)
	readme, _ := os.ReadFile(filepath.Join(root, "README.md"))
	var errs []string
	bad := func(f string, a ...any) { errs = append(errs, fmt.Sprintf(f, a...)) }
	ms := map[string]*plugin.Manifest{}
	names := map[string]string{} // : name -> plugin
	views := map[string]bool{}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			bad("%s: %v", f, err)
			continue
		}
		m, err := plugin.Parse(data)
		if err != nil {
			bad("%s: %s", f, strings.ReplaceAll(err.Error(), "\n", "\n  "))
			continue
		}
		if dir := filepath.Base(filepath.Dir(f)); dir != m.Name {
			bad("%s: name %q, directory %q", f, m.Name, dir)
		}
		if !strings.Contains(string(readme), "plugins/"+m.Name+")") {
			bad("%s: README.md does not list it as [..](plugins/%s)", f, m.Name)
		}
		claim := func(n string) {
			if o := names[n]; o != "" {
				bad("%s: :%s is taken by %s too", f, n, o)
			}
			names[n] = m.Name
		}
		for _, v := range m.Views {
			claim(v.Name)
			views[v.Name] = true
		}
		for _, c := range m.Commands {
			claim(c.Name)
		}
		ms[f] = m
	}
	for _, f := range files {
		m := ms[f]
		if m == nil {
			continue
		}
		for _, a := range m.Actions {
			for _, s := range a.Scopes {
				if s != "all" && !views[s] && !contains(plugin.TmuxScopes, s) {
					bad("%s: key %s: no view %q here", f, a.Key, s)
				}
			}
			if a.Open != "" && !views[a.Open] {
				bad("%s: key %s opens :%s, which no plugin here has", f, a.Key, a.Open)
			}
		}
	}
	if len(files) == 0 {
		bad("no plugins/*/plugin.json under %s", root)
	}
	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e)
	}
	if len(errs) > 0 {
		os.Exit(1)
	}
	fmt.Printf("%d plugins ok\n", len(files))
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
