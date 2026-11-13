package builder

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func Build(program string, args []string, config map[string]any) []string {
	conf, ok := config[program].(map[string]any)

	if !ok {
		return args
	}

	return build(args, conf)
}

func build(args []string, config map[string]any) []string {
	flags := resolveFlags(config)

	if len(args) > 0 {
		conf, ok := config[args[0]].(map[string]any)

		if ok {
			return append(
				append(flags, args[0]),
				build(args[1:], conf)...,
			)
		}
	}

	conf, ok := config["$args"].(map[string]any)

	if ok {
		if len(args) == 0 {
			return flags
		}

		return append(
			append(flags, args...),
			resolveFlags(conf)...,
		)
	}

	return append(flags, args...)
}

func resolveFlags(config map[string]any) []string {
	flags := make([]string, 0, len(config))

	for _, key := range slices.Sorted(maps.Keys(config)) {
		flags = append(flags, formatFlags(key, config[key])...)
	}

	return flags
}

func formatFlags(key string, value any) []string {
	flag := key

	if !strings.HasPrefix(flag, "-") {
		flag = "--" + flag
	}

	switch value := value.(type) {
	case bool:
		if value {
			return []string{flag}
		}
	case []any:
		flags := make([]string, 0, len(value))

		for _, item := range value {
			flags = append(flags, formatFlags(key, item)...)
		}

		return flags
	case map[string]any, nil:
		return nil
	default:
		return []string{flag, fmt.Sprint(value)}
	}

	return nil
}
