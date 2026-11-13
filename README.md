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
wrap [-c wrap.config.yaml] [-v | -d] <program> [arguments...]
```

Examples:

```sh
wrap docker compose up
wrap -c ./production.yaml kubectl get pods
wrap -d curl
```

Available flags:

- `-c, --config <path>`: Sets configuration file path.
- `-v, --verbose`: Outputs more verbose logs.
- `-d, --debug`: Prints resolved command without executing.
- `-h, --help`: Shows usage information.

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
### Flag names

Keys that start with a custom prefix like `-` are kept as written. Other keys receive a `--` prefix by default. Values are separated from flags by a space by default.

To use different separators like `=` instead, include the complete flag and value in the key and set it to `true`:

```yaml
verbose: true                 # --verbose
-v: true                      # -v
output: wide                  # --output wide
-H: value                     # -H value
"--group-by=paths": true      # --group-by=paths
```

### Flag values

The YAML value controls what follows a flag:

- `true` emits the flag by itself, such as `--verbose`; `false` emits nothing.
- A non-empty string or number emits the flag and value as separate arguments, such as `--output wide`.
- An empty string emits the flag followed by an empty argument.
- An array repeats the flag for each item. For example, `[true, 2, three]` becomes `--items --items 2 --items three`.
- `null`, an empty array, or an empty mapping emits nothing.

### Subcommands and arguments

Nested mappings define subcommands. Flags are added at the level where they are configured: program-level flags come before the subcommand, and subcommand-level flags come after it.

`wrap` matches the leading arguments against configured subcommands. At the first unmatched argument, it stops matching; that argument and everything after it are passed to the program unchanged.

When unmatched arguments are present, a `$args` mapping adds its flags after them. Explicit arguments are not parsed or deduplicated, so repeated flags are passed through and handled by the target program.

## Contribution

`wrap` is developed in Go. To build the binary, execute `just go-build`, which will generate the binary at `./build/wrap`. To publish a new version, use `just release <patch|minor|major>`. This command will automatically increment the version according to semantic versioning.
