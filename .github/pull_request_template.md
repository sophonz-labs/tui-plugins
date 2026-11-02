<!-- One plugin per pull request. See CONTRIBUTING.md. -->

Plugin: `plugins/<name>`

External CLI: <name>, <link>, <license>

Commands it runs (from `ext install`):

```
```

Tried in: tmui / stui, on <OS>

- [ ] `test -z "$(gofmt -l .)"`, `go vet ./...` and `go run ./check .` pass
- [ ] Installed with `ext install ./plugins/<name>` and used every view, key and command
- [ ] A line in the README table, linked as `[<name>](plugins/<name>)`
- [ ] `requires` lists every program `run` calls
- [ ] No network calls of its own, no secrets, tokens or hostnames
- [ ] `confirm: true` on anything that changes or deletes state
- [ ] Values reach commands as placeholders in argv or `sh -c` positional arguments only
