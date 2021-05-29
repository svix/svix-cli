package pretty

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	svix "github.com/svixhq/svix-libs/go"
	prettyJson "github.com/tidwall/pretty"
)

type PrintOptions struct {
	JSON bool
}

func Print(v interface{}, opts *PrintOptions) {
	if opts != nil && opts.JSON {
		tryPrintJson(v)
		return
	}

	switch t := v.(type) {
	case *svix.ListResponseApplicationOut:
		printApplicationList(t)
	case *svix.ApplicationOut:
		printApplicationOut(t)
	case *svix.EventTypeInOut:
		printEventTypeInOut(t)
	case *svix.EndpointOut:
		printEndpointOut(t)
	case *svix.ListResponseEndpointOut:
		printListResponseEndpointOut(t)
	case *svix.EndpointSecret:
		printEndpointSecret(t)
	case *svix.ListResponseMessageOut:
		printListResponseMessageOut(t)
	case *svix.MessageOut:
		printMessageOut(t)
	case *svix.ListResponseMessageAttemptOut:
		printListResponseMessageAttemptOut(t)
	case *svix.ListResponseMessageEndpointOut:
		printListResponseMessageEndpointOut(t)
	case *svix.ListResponseMessageAttemptEndpointOut:
		printListResponseMessageAttemptEndpointOut(t)
	case *svix.MessageAttemptOut:
		printMessageAttemptOut(t)
	case *svix.DashboardAccessOut:
		printDashboardAccessOut(t)
	default:
		// if all else fails try to print as json
		tryPrintJson(t)
	}
}

func getTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.Debug)
}

func makeTerminalHyperlink(name, url string) string {
	return fmt.Sprintf("\u001B]8;;%s\a%s\u001B]8;;\a", url, name)
}

func tryPrintJson(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("%+v\n", v)
	}
	fmt.Println(string(prettyJson.Pretty(b)))
}

func fmtStringPtr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return string(*s)
}

func printApplicationList(l *svix.ListResponseApplicationOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "Name\tID\tUID\tCreated At")
	for _, app := range l.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", app.Name, app.Id, fmtStringPtr(app.Uid), app.CreatedAt)
	}
	w.Flush()
}

func printApplicationOut(a *svix.ApplicationOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "Name\tID\tUID\tCreated At")
	uid := ""
	if a.Uid != nil {
		uid = *a.Uid
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", a.Name, a.Id, uid, a.CreatedAt)
	w.Flush()
}

func printEventTypeInOut(et *svix.EventTypeInOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "Name\tDescription")
	fmt.Fprintf(w, "%s\t%s\n", et.Name, et.Description)
	w.Flush()
}

func printEndpointOut(ep *svix.EndpointOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tURL\tDescription\tFilter Types")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ep.Id, ep.Url, fmtStringPtr(ep.Description), ep.FilterTypes)
	w.Flush()
}

func printListResponseEndpointOut(l *svix.ListResponseEndpointOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tURL\tDescription\tFilter Types")
	for _, ep := range l.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ep.Id, ep.Url, fmtStringPtr(ep.Description), ep.FilterTypes)
	}
	w.Flush()
}

func printEndpointSecret(secret *svix.EndpointSecret) {
	w := getTabWriter()
	fmt.Fprintln(w, "Endpoint Secret Key")
	fmt.Fprintf(w, "%s\n", secret.Key)
	w.Flush()
}

func printListResponseMessageOut(l *svix.ListResponseMessageOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tTimestamp\tEvent ID\tPayload")
	for _, msg := range l.Data {
		jsonPayload, err := json.Marshal(msg.Data)
		if err != nil {
			jsonPayload = []byte("Error Displaying Payload")
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", msg.Id, msg.Timestamp.Format(time.RFC3339), fmtStringPtr(msg.EventId), string(jsonPayload))
	}
	w.Flush()
}

func printMessageOut(msg *svix.MessageOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tTimestamp\tEvent ID\tPayload")
	jsonPayload, err := json.Marshal(msg.Data)
	if err != nil {
		jsonPayload = []byte("Error Displaying Payload")
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", msg.Id, msg.Timestamp.Format(time.RFC3339), fmtStringPtr(msg.EventId), string(jsonPayload))
	w.Flush()
}

func printListResponseMessageAttemptOut(l *svix.ListResponseMessageAttemptOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tTimestamp\tEvent ID\tResponse")
	for _, attempt := range l.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", attempt.Id, attempt.Timestamp, attempt.EndpointId, attempt.Response)
	}
	w.Flush()
}

func printListResponseMessageEndpointOut(l *svix.ListResponseMessageEndpointOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tTimestamp\tEvent ID\tResponse")
	for _, ep := range l.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\n", ep.Id, ep.Url, fmtStringPtr(ep.Description))
	}
	w.Flush()
}

func printListResponseMessageAttemptEndpointOut(l *svix.ListResponseMessageAttemptEndpointOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tTimestamp\tStatus\tResponse")
	for _, ep := range l.Data {
		fmt.Fprintf(w, "%s\t%s\t%dt%s\n", ep.Id, ep.Timestamp, ep.Status, fmtStringPtr(&ep.Response))
	}
	w.Flush()
}

func printMessageAttemptOut(ma *svix.MessageAttemptOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tTimestamp\tEndpoint ID\tStatus\tResponse")
	fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", ma.Id, ma.Timestamp, ma.EndpointId, ma.Status, ma.Response)
	w.Flush()
}

func printDashboardAccessOut(da *svix.DashboardAccessOut) {
	fmt.Printf(`You can access the Dashboard at the following URL:
%s
`, makeTerminalHyperlink(da.Url, da.Url))
}
