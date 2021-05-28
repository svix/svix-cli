package pretty

import (
	"fmt"

	svix "github.com/svixhq/svix-libs/go"
)

func PrintApplicationList(l *svix.ListResponseApplicationOut) {
	for _, app := range l.Data {
		fmt.Printf("%+v\n", app)
	}
}

func PrintApplicationOut(a *svix.ApplicationOut) {
	fmt.Printf("%+v", a)
}

func PrintEventTypeInOut(et *svix.EventTypeInOut) {
	fmt.Printf("%+v", et)
}

func PrintEndpointOut(ep *svix.EndpointOut) {
	fmt.Printf("%+v", ep)
}

func PrintEndpointSecret(endpointID string, secret *svix.EndpointSecret) {
	fmt.Printf(`Endpoint: %s
Secret: %s
`, endpointID, secret.Key)
}

func PrintListResponseMessageOut(l *svix.ListResponseMessageOut) {
	fmt.Printf("%+v", l)
}

func PrintMessageOut(m *svix.MessageOut) {
	fmt.Printf("%+v", m)
}

func PrintListResponseMessageAttemptOut(l *svix.ListResponseMessageAttemptOut) {
	fmt.Printf("%+v", l)
}

func PrintListResponseMessageEndpointOut(l *svix.ListResponseMessageEndpointOut) {
	fmt.Printf("%+v", l)
}

func PrintDashboardURL(appID, url string) {
	fmt.Printf(`You can access the Dashboard for %s at the following URL:
%s
`, appID, makeTerminalHyperlink(url, url))
}

func makeTerminalHyperlink(name, url string) string {
	return fmt.Sprintf("\u001B]8;;%s\a%s\u001B]8;;\a", url, name)
}
