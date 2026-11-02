# Contributing a plugin

A plugin is one `plugins/<name>/plugin.json`. [SPEC.md](SPEC.md) is the format; this page is how to get one into this repository.

## Add it

1. Fork the repository and make a branch.
2. Copy the closest example from `plugins/` to `plugins/<name>/plugin.json`. A table from a command's output is a view (`gh-prs`, `docker-ps`); a program that takes the terminal is a command or a `foreground` action (`lazygit`).
3. Fill in the manifest:
   - `"v": 1` and `name`. The name is `a-z`, `0-9` and `-`, at most 32, and starts with a letter or digit. The directory has the same name. Name it after the CLI and what it shows (`kubectl-pods`, `gh-runs`).
   - `description`, `homepage` (the CLI's page) and `requires` (every program `run` calls, so a missing one is reported).
   - At least one of `views`, `actions`, `commands`. Unknown fields are an error.
   - Keys: one character or `ctrl-a` to `ctrl-z`. Never `j k g G h l q / : ? enter esc ctrl-b ctrl-f ctrl-r ctrl-c ctrl-e y Y` or a digit; a manifest that binds one is invalid. Avoid the tools' own keys on their views (tmui's `n a R s v p ctrl-d` on tmux's views): a built-in key wins and yours is left out.
   - `:` names (views and commands) must not be a tool's own command or another plugin's name here.
   - `run` is argv with no shell. If you need a pipe, use `["sh", "-c", "... \"$1\" ...", "sh", "{field}"]` so values arrive as arguments, never inside the script.
   - `confirm: true` on anything that changes or deletes something.
4. Add one line for it to the Plugins table in [README.md](README.md), linked as `[<name>](plugins/<name>)`. The check fails without it.

## Test it before the PR

Run what CI runs, from the repository root:

```sh
test -z "$(gofmt -l .)"
go vet ./...
go run ./check .
```

`go run ./check .` validates every manifest against the spec, and the rules across plugins: the directory matches the name, every scope and `open` names a view that exists, no two plugins take one `:` name, and README.md lists the plugin. It prints `N plugins ok` and exits 0 when it passes.

Then install it in tmui or stui and use every view, key and command once:

```sh
tmui ext install ./plugins/<name>      # from your checkout
stui ext install ./plugins/<name>
```

After you push, check that a bare name installs from your branch the way users will get it from `main`:

```sh
TUI_PLUGINS_REPO=https://raw.githubusercontent.com/<you>/tui-plugins/<branch>/plugins tmui ext install <name>
```

Install prints every command the plugin will run. Read that list: it is what reviewers read too. `:ext` inside the tool shows anything that was left out, such as a key a built-in already uses. `tmui ext remove <name>` takes it out again.

## The pull request

- One plugin per pull request.
- In the description: the external CLI it needs, a link to it and its license, and the commands the plugin runs (the list `ext install` printed).
- Say which tool you tried it in (tmui, stui or both) and on what OS.

## What reviewers check

- It runs only the CLI it names in `requires`, plus `sh` and common tools such as `less` where the example plugins use them. Nothing is downloaded or installed.
- It makes no network calls of its own. Network access comes only from the CLI doing its usual job (`gh` talking to GitHub, `kubectl` to the cluster).
- No secrets, tokens, hostnames or account names in the manifest. Credentials stay in the CLI's own config.
- Anything that changes or deletes state has `confirm: true`.
- Field values reach commands only as placeholders in argv or as positional arguments to `sh -c`, never spliced into a shell string.
- The external CLI is something users install themselves; the plugin does not bundle or vendor it.

## Update or remove a plugin

To change a plugin, send a pull request that edits its `plugin.json`, with the same checks. Users get the change by installing it again: `tmui ext install <name>` shows the new commands and asks before replacing the old file. Keep `name` and existing `:` names stable, since users type them; a rename is a new plugin plus the removal of the old one.

To remove a plugin, delete its directory and its README line in one pull request, and say why (the CLI is gone, a better plugin replaces it). Copies already installed keep working until users run `tmui ext remove <name>`.

Changes to the format itself go to SPEC.md and the tools first; open an issue before a pull request.
