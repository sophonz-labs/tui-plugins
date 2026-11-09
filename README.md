# tui-plugins

Plugins for [tmui](https://github.com/sophonz-labs/tmui) and [stui](https://github.com/sophonz-labs/stui): a CLI's output as a view in the tool's table, and its commands on keys and on the `:` line. Each plugin is a folder holding a `plugin.json`; [SPEC.md](SPEC.md) is the format.

A plugin folder can instead, or as well, hold a `language.json`: a **language pack**, which tells [st](https://sophonz.io) how to start a language server so its editors get completion, hover, definition and diagnostics. [SPEC-language.md](SPEC-language.md) is that format. Anyone can write either kind.

## Install

```sh
tmui ext install gh-prs      # from this repository's main branch
tmui ext install ./my-plugin # a directory with a plugin.json, or a path or URL to one
tmui ext list
tmui ext remove gh-prs
```

stui takes the same commands: `stui ext install gh-prs`. Install prints every command the plugin will run and asks before writing it to `~/.config/<tool>/plugins/<name>/`; a plugin put in `~/.config/tui-plugins/<name>/` by hand is read by both. Then `:prs` opens the view, and `?` and `:ext` show what each plugin added. In stui, a key on panes works on the agents view too.

## Plugins

| Plugin | Needs | Adds |
|---|---|---|
| [gh-prs](plugins/gh-prs) | gh | `:prs` for the pane's repository (`P` on a pane or window): `o` browser, `d` diff, `c` checks, `C` check out |
| [gh-runs](plugins/gh-runs) | gh | `:gh-runs` (`A` on a pane or window): `w` watch in a popup, `o` browser, `L` failed logs, `r` rerun failed |
| [docker-ps](plugins/docker-ps) | docker | `:containers`: `L` logs, `e` shell, `i` inspect, `S` stop |
| [kubectl-pods](plugins/kubectl-pods) | kubectl | `:pods` in the current context: `L` logs, `d` describe, `e` shell, `D` delete |
| [git-status](plugins/git-status) | git | `:git-status` of the pane's directory (`S` on a pane or window): `d` diff, `a` stage, `u` unstage |
| [lazygit](plugins/lazygit) | lazygit | `L` on a session, window or pane, and `:lazygit`, in the pane's directory |

## Language packs

```sh
st ext install typescript    # from this repository's main branch
st ext install ./my-pack     # a directory with a language.json, or a path or URL to one
st ext list
st ext remove typescript
```

Install prints the pack, the server's command, every install command and where it came from, and asks before writing it to `~/.sophonz/ext/<name>/` and running the install commands. tmui and stui read only `plugin.json`, so a pack does not show up there.

| Pack | Needs | Serves |
|---|---|---|
| [java](plugins/java) | java (25 for jdtls 1.63), python3, jdtls (install it first: `brew install jdtls` on macOS) | `*.java` |
| [typescript](plugins/typescript) | node, npm; installs typescript-language-server and typescript | `*.ts *.tsx *.mts *.cts *.js *.jsx *.mjs *.cjs` |

## Writing one

Start from the closest example. Check it the way CI does:

```sh
go run ./check .
```

and try it with `tmui ext install ./plugins/<name>` (a language pack: `st ext install ./plugins/<name>`). A pull request adds the plugin's directory and its line in the table above. [CONTRIBUTING.md](CONTRIBUTING.md) has the steps and what review looks for.

## License

MIT
