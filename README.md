# wrap

`wrap` runs existing command-line programs with defaults defined in YAML. It is a lightweight argument wrapper, not a task runner or alias manager. It starts programs directly without invoking a shell.

## Installation

Use the Go version declared in [`go.mod`](go.mod). Build the executable from the repository root:

```sh
go build -o wrap .
```

## Usage

```text
wrap [-c config.yaml] [-d] <program> [command path] [arguments...]
```

Examples:

```sh
wrap docker compose up
wrap -c ./production.yaml kubectl get pods
wrap git log --oneline
wrap -d docker compose up
```

Wrapper options must appear before the program name. `-c` or `--config` selects a configuration file; without it, `./config.yaml` is used. A missing or unreadable configuration file is an error. `-d` or `--debug` prints the resolved command without running it.

Arguments after the program are passed to the target as separate arguments. `wrap` does not use a shell, so spaces and special characters are not interpreted by a shell.

## Configuration

See [`examples/config.yaml`](examples/config.yaml) for a complete example:

```yaml
docker:
  debug: true
  compose:
    up:
      detach: true
      build: true
      force-recreate: true
        
kubectl:
  get:
    pods:
      all-namespaces: true
      output: wide

cp:
  $args:
    archive: true
    verbose: true
```

Objects represent command levels. Values become default flags:

- `true` emits a boolean flag; `false` emits nothing.
- Strings and numbers emit `--name=value`.
- Arrays emit one flag per item using the same value rules: strings and numbers emit `--name=value`, `true` emits `--name`, and `false` emits nothing.

Flag keys control their prefix. A key without a prefix receives `--`; an existing prefix is preserved. For example, `verbose: true` emits `--verbose`, `"-v": true` emits `-v`, and `"--output": wide` emits `--output=wide`.

Use the reserved `$args` key to pass through all remaining arguments before adding its nested defaults:

```yaml
restic:
  backup:
    $args:
      verbose: true
```

With `wrap restic backup data`, the wrapper emits `restic --repo=... backup data --verbose`. `$args` preserves every remaining argument, including flags and positionals, before adding its nested defaults.

## Command Resolution

The command path must immediately follow the program name. `wrap` follows leading arguments only while each one matches a configured command object. At each matched level, it inserts that level's defaults before the next command or remaining arguments. Once an argument does not match a configured command, the rest are passed through unchanged and are not inspected as commands.

For example, `wrap docker compose up --detach` resolves the configured `docker`, `compose`, and `up` levels, then passes `--detach` through to Docker. Explicit flags are not parsed or deduplicated, so configured defaults and user-provided arguments may contain identical flags. The target program determines how repeated flags are handled.

## Development

Run the test suite from the repository root:

```sh
go test ./...
```