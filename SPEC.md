# Plugin manifest, v 1

A plugin brings another CLI into a terminal UI for tmux: its output as a view, and its commands on keys and on the `:` line. A plugin is one file, `plugin.json`, and nothing runs that the file does not spell out. This is the contract tmui and stui implement. A plugin folder may also hold a `language.json`, a language pack, specified in [SPEC-language.md](SPEC-language.md); tools that implement this spec ignore it. Revised 2026-11-10, rev 1.

## Where plugins live

| Directory | Read by |
|---|---|
| `$XDG_CONFIG_HOME/<tool>/plugins/<name>/plugin.json` | that tool only (`tmui`, `stui`) |
| `$XDG_CONFIG_HOME/tui-plugins/<name>/plugin.json` | every tool that implements this spec |

`$XDG_CONFIG_HOME` defaults to `~/.config`. When both have a plugin of one name, the tool's own wins and the shared one is listed as shadowed. `<tool> ext install` writes to the tool's own directory; put a plugin in the shared one by hand.

## The file

```json
{
  "v": 1,
  "name": "gh-prs",
  "description": "Pull requests of the repository the pane is in.",
  "homepage": "https://cli.github.com",
  "requires": ["gh"],
  "views": [],
  "actions": [],
  "commands": []
}
```

| Field | | |
|---|---|---|
| `v` | required | `1`. A tool refuses a higher number. |
| `name` | required | `a-z`, `0-9` and `-`, at most 32, starting with a letter or digit. The directory has the same name. |
| `description`, `homepage` | | Shown before installing and in `:ext`. |
| `requires` | | Programs that must be on `PATH`. A missing one is reported; the plugin still loads. |
| `views`, `actions`, `commands` | at least one | Below. |

Unknown fields are an error, so a misspelt `"confirm"` cannot quietly run without asking. A new field means a new `v`.

## views

A view is a table from a command's output, opened with `:<name>`.

```json
{
  "name": "prs",
  "title": "Pull requests",
  "run": ["gh", "pr", "list", "--json", "number,title,author"],
  "format": "json",
  "columns": [{ "title": "#", "field": "number" }, { "field": "title" }, { "title": "AUTHOR", "field": "author.login" }],
  "id": "number",
  "refresh": "30s",
  "timeout": "20s",
  "cwd": "pane"
}
```

| Field | | |
|---|---|---|
| `name` | required | As `name` above. It is the `:` command and the scope actions name. It must not be one of the tool's own commands. |
| `title` | | The table's title; `name` without it. |
| `run` | required | The command, as argv. Placeholders as below. |
| `format` | required | `json`, `tsv` or `lines`. |
| `columns` | required for json | `field` (required) and `title` (the field in capitals without it). Without `columns`, tsv shows its header's fields and lines its one. |
| `id` | | The field that names a row, so the cursor stays on it across refreshes. Default: the first column. |
| `refresh` | | Run again this often while the view is on screen, `1s` or more. Without it the view runs when opened and on `ctrl-r`. |
| `timeout` | | Default `10s`. A run that takes longer is stopped and shown as an error. |
| `cwd` | | `pane`: the directory of the pane the view was opened from. `home`. Nothing: the tool's own. |

Formats:

- `json`: an array of objects (`gh --json`, `kubectl -o json | jq .items`), or one object per line (`docker ps --format json`). Nested fields are joined with dots: `author.login`. A string is itself, `null` is empty, an array is its items joined by `, `, a number keeps its digits.
- `tsv`: the first line is the field names, each line after it a row, split on tabs.
- `lines`: every non-empty line is a row with one field, `line`.

A failed run, a timeout or output that does not parse shows in the view as an error, with the command's stderr. It never stops the tool.

## actions

An action is a key on the rows of some views.

```json
{ "key": "o", "title": "open in browser", "scopes": ["prs"], "run": ["gh", "pr", "view", "{number}", "--web"] }
{ "key": "S", "title": "git status", "scopes": ["panes", "windows"], "open": "git-status" }
```

| Field | | |
|---|---|---|
| `key` | required | One character, or `ctrl-a` to `ctrl-z`. |
| `title` | required | Shown with the key in the header and in help. |
| `scopes` | required | Views it works on: a plugin view's name, `sessions`, `windows`, `panes` (tmux's own), or `all`. A tool without one of tmux's views just has no such key there. |
| `run` | `run` or `open` | The command, as argv. |
| `open` | `run` or `open` | A plugin view to open, with the row's tmux placeholders (for `cwd: pane`). |
| `mode`, `confirm`, `wait`, `cwd` | | As in commands below. |

Reserved, on every view, in every tool: `j k g G h l q / : ? enter esc ctrl-b ctrl-f ctrl-r ctrl-c ctrl-e y Y` and the digits `0`-`9`. `y` copies the row's value (on a plugin view, its record's `url` when it has one, else its `id`) and `Y` the whole row; a view that wants `y` to copy a link asks its command for a `url` field. A manifest that binds one is invalid. Each tool also keeps the keys of its own views (tmui's `n a R s v p ctrl-d`, for one); a plugin key that meets one of those, or a key another plugin already took on that view, is left out and reported in `:ext`. A built-in key always wins.

## commands

A command is `:<name> [args]`.

```json
{ "name": "lazygit", "title": "lazygit here", "run": ["lazygit", "{args}"], "mode": "foreground", "cwd": "pane" }
```

| Field | | |
|---|---|---|
| `name` | required | As a view's name; no view and command may share one. |
| `title` | | Shown in help. |
| `run` | required | The command, as argv. |
| `mode` | | `background` (default): runs with no terminal and its last line of output goes to the status line. `foreground`: the tool gives the command the terminal and comes back when it exits. `popup`: a `tmux display-popup` over the tool when it runs inside tmux 3.2 or later, else as `foreground`. |
| `confirm` | | Ask y/N first, showing the command filled in. |
| `wait` | | `foreground` and `popup`: keep the output on screen until a key. |
| `cwd` | | As for views. |

After an action or command, the view on screen runs again.

## Placeholders

`run` is argv and runs without a shell. A placeholder fills part of one argument and never splits, adds or removes arguments, so a field holding `x; rm -rf ~` is one argument and runs nothing. A plugin that wants a shell says so, and passes values as positional arguments, not inside the script:

```json
["sh", "-c", "gh pr diff \"$1\" | less -R", "sh", "{number}"]
```

| Placeholder | Is |
|---|---|
| `{<field>}` | On a plugin view, a field of the row under the cursor (any field of the record, shown or not). |
| `{session}`, `{window}`, `{pane}` | The session name, window id (`@3`) and pane id (`%7`). A session stands for its active window, a window for its active pane. |
| `{pane_cwd}`, `{pane_command}` | That pane's directory and running command. |
| `{args}` | A command's words after its name. As a whole argument it is one argument per word, none when there are none; inside a longer one, the words joined by spaces. |

On tmux's views the tmux placeholders are those of the row under the cursor. On a plugin view they are those of where the view was opened from. From the `:` line with no row, they are the tool's own session.

A placeholder is `{` a letter or `_`, then letters, digits, `_`, `.` or `-`, then `}`. Any other brace is itself, so jq's `{a: .b}` and kubectl's `{.metadata.name}` need nothing; `{{` is a literal `{`, for text that would read as a placeholder (kubectl jsonpath's `{end}` is `{{end}`). An unknown placeholder is an error and nothing runs. Values are never expanded again.

On actions scoped only to `sessions`, `windows` and `panes`, only the tmux placeholders exist, and a manifest that uses another is invalid.

## Installing

`<tool> ext install <name | path | url>` reads `https://raw.githubusercontent.com/sophonz-labs/tui-plugins/main/plugins/<name>/plugin.json` for a bare name (`TUI_PLUGINS_REPO` points at a fork), checks it against this spec, prints every command it will run, and asks y/N before writing it. `<tool> ext list` and `<tool> ext remove <name>` do what they say. Inside the tool, `:ext` lists them, what they add, and what was left out. The command is `ext`, not `plugin`, because `stui plugin` already installs stui's skills for coding agents; both tools use the one name so these docs read the same for each.
