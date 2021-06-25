package pretty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	svix "github.com/svix/svix-libs/go"
	prettyJson "github.com/tidwall/pretty"
)

type PrinterOptions struct {
	Color bool
}

type Printer struct {
	opts *PrinterOptions
}

func NewPrinter(opts *PrinterOptions) *Printer {
	return &Printer{
		opts: opts,
	}
}

func (p *Printer) Print(v interface{}) {
	var b []byte
	switch msg := v.(type) {
	case []byte:
		b = msg
	default:
		var err error
		var buf bytes.Buffer

		// disable html escaping
		// otherwise & gets transformed to \u0026
		// so url can become invalid
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)

		err = enc.Encode(v)
		if err != nil {
			fmt.Printf("%+v\n", v)
			return
		}
		b = buf.Bytes()
	}

	b = prettyJson.Pretty(b)
	if p.opts != nil && p.opts.Color {
		b = prettyJson.Color(b, nil)
	}
	fmt.Println(string(b))
}

func (p *Printer) CheckErr(msg interface{}) {
	if msg != nil {
		if err, ok := msg.(*svix.Error); ok {
			p.Print(err.Body())
		}
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}

func MakeTerminalLink(name, url string) string {
	return fmt.Sprintf("\u001B]8;;%s\a%s\u001B]8;;\a", url, name)
}
