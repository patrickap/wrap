# wrap

`wrap` is a configuration-driven wrapper for command-line programs. It reads its configuration from a `wrap.config.yaml` file and adds the defined flags when you run a program. `wrap` is a lightweight wrapper, not a task runner.

## Installation

Use the Go version declared in [`go.mod`](go.mod) to build from the repository root:

```sh
just go-build
```

Or download the binary for your OS and architecture from `Releases` and install it in a directory on your `PATH`:

```sh
sudo install -m 755 wrap /usr/local/bin/wrap
```

## Usage

```sh
wrap [-c wrap.config.yaml] [-d] <program> [command path] [arguments...]
```

Examples:

```sh
wrap docker compose up
wrap -c ./production.yaml kubectl get pods
wrap -d curl
```

Available flags:

- `-c, --config <path>`: Read defaults from this YAML file. Defaults to `./wrap.config.yaml` in the current directory.
- `-d, --debug`: Print the resolved command without running it.
- `-h, --help`: Show usage information.

Everything after the program name is passed to the target as separate arguments.

## Configuration

Configuration mappings represent command levels. For example:

```yaml
# docker --debug compose up --detach --build --force-recreate
docker:
  debug: true
  compose:
    up:
      detach: true
      build: true
      force-recreate: true
   
# kubectl get pods --all-namespaces --output=wide
kubectl:
  get:
    pods:
      all-namespaces: true
      output: wide

# curl -H "Accept: application/json" -H "Authorization: Bearer TOKEN"
curl:
  -H:
    - "Accept: application/json"
    - "Authorization: Bearer TOKEN"

# cp /from /to --archive --verbose
cp:
  $args:
    archive: true
    verbose: true
```

Values are rendered as flags:

- `true` emits a flag without a value, such as `--verbose`; `false` emits nothing.
- Strings and numbers emit `--name=value`. An empty string emits `--name=` with an explicitly empty value.
- Arrays emit one flag per item. Each item uses the same rules, so `[true, 2, three]` becomes `--items --items=2 --items=three`.
- `null`, empty arrays, and empty mappings emit no flag.

Keys without a prefix receive `--`; an existing prefix is preserved. For example, `verbose: true` emits `--verbose`, `"-v": true` emits `-v`, `"--output": wide` emits `--output=wide`, and `"-H": value` emits `-H=value`.

Nested mappings name subcommands. Flags are added at the level where they are configured: program-level flags appear after the program name and before its subcommand, while subcommand flags appear after that subcommand.

The first arguments after the program are matched against configured subcommands. Once an argument does not match, it and all remaining arguments pass through unchanged and are not treated as subcommands.

A configured `$args` mapping adds its flags after those remaining arguments:

Explicit arguments are not parsed or deduplicated. If a configured default and a user argument repeat a flag, both are passed to the target; the target determines how to handle them.

## Contribution

`wrap` is developed in Go. To build the binary, execute `just go-build`, which will generate the binary at `./build/wrap`. To publish a new version, use `just release <patch|minor|major>`. This command will automatically increment the version according to semantic versioning.
