package pretty

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	svix "github.com/svixhq/svix-libs/go"
)

func getTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.Debug)
}

func fmtStringPtr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return string(*s)
}

func PrintApplicationList(l *svix.ListResponseApplicationOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "Name\tID\tUID\tCreated At")
	for _, app := range l.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", app.Name, app.Id, fmtStringPtr(app.Uid), app.CreatedAt)
	}
	w.Flush()
}

func PrintApplicationOut(a *svix.ApplicationOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "Name\tID\tUID\tCreated At")
	uid := ""
	if a.Uid != nil {
		uid = *a.Uid
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", a.Name, a.Id, uid, a.CreatedAt)
	w.Flush()
}

func PrintEventTypeInOut(et *svix.EventTypeInOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "Name\tDescription")
	fmt.Fprintf(w, "%s\t%s\n", et.Name, et.Description)
	w.Flush()
}

func PrintEndpointOut(ep *svix.EndpointOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tURL\tDescription\tFilter Types")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ep.Id, ep.Url, fmtStringPtr(ep.Description), ep.FilterTypes)
	w.Flush()
}

func PrintEndpointSecret(endpointID string, secret *svix.EndpointSecret) {
	w := getTabWriter()
	fmt.Fprintln(w, "Endpoint ID\tSecret")
	fmt.Fprintf(w, "%s\t%s\n", endpointID, secret.Key)
	w.Flush()
}

func PrintListResponseMessageOut(l *svix.ListResponseMessageOut) {
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

func PrintMessageOut(msg *svix.MessageOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tTimestamp\tEvent ID\tPayload")
	jsonPayload, err := json.Marshal(msg.Data)
	if err != nil {
		jsonPayload = []byte("Error Displaying Payload")
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", msg.Id, msg.Timestamp.Format(time.RFC3339), fmtStringPtr(msg.EventId), string(jsonPayload))
	w.Flush()
}

func PrintListResponseMessageAttemptOut(l *svix.ListResponseMessageAttemptOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tTimestamp\tEvent ID\tResponse")
	for _, attempt := range l.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", attempt.Id, attempt.Timestamp, attempt.EndpointId, attempt.Response)
	}
	w.Flush()
}

func PrintListResponseMessageEndpointOut(l *svix.ListResponseMessageEndpointOut) {
	w := getTabWriter()
	fmt.Fprintln(w, "ID\tTimestamp\tEvent ID\tResponse")
	for _, ep := range l.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\n", ep.Id, ep.Url, fmtStringPtr(ep.Description))
	}
	w.Flush()
}

func PrintDashboardURL(appID, url string) {
	fmt.Printf(`You can access the Dashboard for %s at the following URL:
%s
`, appID, makeTerminalHyperlink(url, url))
}

func makeTerminalHyperlink(name, url string) string {
	return fmt.Sprintf("\u001B]8;;%s\a%s\u001B]8;;\a", url, name)
}
