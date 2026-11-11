# Language pack, v 1

A language pack tells a tool how to start a language server for a language, so its editor gets completion, hover, definition and diagnostics. A pack is one file, `language.json`, in a plugin folder, beside a `plugin.json` ([SPEC.md](SPEC.md)) or alone. Nothing runs that the file does not spell out. st implements this spec (`st ext`). Draft 2026-11-10.

It keeps plugin.json's rules: `v`, and a tool refuses a higher number; the same `name` rules; `description`, `homepage` and `requires`; argv only, never a shell; unknown fields are an error, so a new field means a new `v`.

## The file

```json
{
  "v": 1,
  "name": "typescript",
  "description": "TypeScript and JavaScript through typescript-language-server.",
  "homepage": "https://github.com/typescript-language-server/typescript-language-server",
  "requires": ["node", "npm"],
  "languages": [
    {
      "id": "typescript",
      "files": ["*.ts", "*.tsx", "*.js"],
      "rootMarkers": ["tsconfig.json", "package.json"],
      "lineComment": "//",
      "grammar": "typescript"
    }
  ],
  "server": {
    "command": ["typescript-language-server", "--stdio"],
    "needs": ["node"]
  },
  "install": [
    { "npm": { "package": "typescript-language-server", "version": "6.0.3" } },
    { "npm": { "package": "typescript", "version": "6.0.3" } }
  ]
}
```

| Field | | |
|---|---|---|
| `v` | required | `1`. A tool refuses a higher number. |
| `name` | required | `a-z`, `0-9` and `-`, at most 32, starting with a letter or digit. The folder has the same name. |
| `description`, `homepage` | | Shown before installing. A server the pack cannot install says here how to get it. |
| `requires` | | Programs that must be on `PATH`, shown before installing. |
| `languages` | required | At least one. Below. |
| `server` | required | Below. |
| `install` | | Steps that put the server on the machine. Below. |

## languages

| Field | | |
|---|---|---|
| `id` | required | The language id the editor and the server use (`typescript`, `java`): `a-z`, `0-9`, `_`, `+` and `-`, at most 32. One pack may not name an id twice. An installed pack whose id is one a tool has built in replaces the built-in. |
| `files` | required | Base-name patterns of the files it serves (`*.java`). No `/`. |
| `rootMarkers` | | File or folder names. The nearest folder above a file holding one is the project root the server starts in; with none, the file's own folder. |
| `lineComment` | | The line comment, for the editor's toggle (`//`, `#`). |
| `grammar` | | The name of a syntax grammar the editor already has (`java`, `typescript`). A pack cannot bring a grammar. |

## server

| Field | | |
|---|---|---|
| `command` | required | argv. The first element is looked up on the login shell's `PATH`, and is never a shell line: `["jdtls"]`, not `"jdtls"` and not `["sh -c jdtls"]`. It runs in the project root, with LSP on stdin and stdout. |
| `version` | | Arguments that ask the server's version (`["--version"]`). When given, a failing answer counts as missing. |
| `needs` | | Programs the server needs besides itself (`node`, `java`). A missing one makes the server missing. |

## install

A list of steps, run in order after the owner says yes, in the login shell's environment. Each step is exactly one of:

| Kind | Runs |
|---|---|
| `{"npm": {"package": "<name>", "version": "<x.y.z>"}}` | `npm install -g <name>@<x.y.z>` |
| `{"go": {"module": "<path>", "version": "<vX.Y.Z>"}}` | `go install <path>@<vX.Y.Z>` (the server is then looked for in `~/go/bin` too) |

The step names what to install, never where: the host decides. st installs under its own home, one folder per server (`<home>/tools/lsp/<server>/`, with `npm install -g --prefix <that folder>` and `GOBIN=<that folder>/bin`), looks there before `PATH`, and leaves a global install alone; removing the folder falls back to it.

Versions are exact: `6.0.3`, `v0.20.0`, with an optional pre-release suffix. A range, a tag such as `latest`, or any other kind of step is an error. v 1 has no downloads by URL; a server that ships as an archive (jdtls) has no `install` and says in `description` how to get it.

`install` is a list because a server can need a second package: typescript-language-server needs `typescript` beside it when the project has none.

## Installing

`st ext install <name | path | url>` reads `https://raw.githubusercontent.com/sophonz-labs/tui-plugins/main/plugins/<name>/language.json` for a bare name (`TUI_PLUGINS_REPO` points at a fork), checks it against this spec, prints the pack, the server's command, every install command and where it came from, and asks y/N (`--yes` answers for a script; without a terminal it is required). It then writes `~/.sophonz/ext/<name>/language.json` and runs the steps, showing their output. `st ext list` and `st ext remove <name>` do what they say; remove does not uninstall what the steps installed.
