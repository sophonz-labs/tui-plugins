# tui-plugins

Plugins for [tmui](https://github.com/sophonz-labs/tmui) and [stui](https://github.com/sophonz-labs/stui): a CLI's output as a view in the tool's table, and its commands on keys and on the `:` line. Each plugin is one `plugin.json`; [SPEC.md](SPEC.md) is the format.

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

## Writing one

Start from the closest example. Check it the way CI does:

```sh
go run ./check .
```

and try it with `tmui ext install ./plugins/<name>`. A pull request adds the plugin's directory and its line in the table above. [CONTRIBUTING.md](CONTRIBUTING.md) has the steps and what review looks for.

## License

MIT
