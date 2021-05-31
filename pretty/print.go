package pretty

import (
	"encoding/json"
	"fmt"

	prettyJson "github.com/tidwall/pretty"
)

type PrintOptions struct {
	Color bool
}

func Print(v interface{}, opts *PrintOptions) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("%+v\n", v)
	}
	b = prettyJson.Pretty(b)
	if opts != nil && opts.Color {
		b = prettyJson.Color(b, nil)
	}
	fmt.Println(string(b))
}

// func makeTerminalHyperlink(name, url string) string {
// 	return fmt.Sprintf("\u001B]8;;%s\a%s\u001B]8;;\a", url, name)
// }
