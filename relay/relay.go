package relay

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/svix/svix-cli/pretty"
)

// Defaults
const (
	defaultAPIHost = "api.relay.svix.com"
	apiPrefix      = "api/v1"
	defaultTimeout = 30 * time.Second
	pongWait       = 10 * time.Second
	pingPeriod     = (pongWait * 2) / 10
	writeWait      = 10 * time.Second
)

type Client struct {
	token              string
	websocketURL       string
	localURL           *url.URL
	receiveURLTemplate string
	dialer             *websocket.Dialer
	httpClient         *http.Client
	logging            bool

	conn      *websocket.Conn
	stopRead  chan struct{}
	stopWrite chan struct{}

	errChan chan error

	sendChan chan *OutgoingMessageEvent
	recChan  chan *IncomingMessage
	wg       *sync.WaitGroup
}

type ClientOptions struct {
	DisableSecurity bool
	RelayDebugUrl   string
	Logging         bool
}

func NewClient(token string, localURL *url.URL, opts *ClientOptions) *Client {
	wsProto := "wss"
	apiHost := defaultAPIHost
	logging := false
	if opts != nil {
		if opts.DisableSecurity {
			wsProto = "ws"
		}
		if opts.RelayDebugUrl != "" {
			apiHost = opts.RelayDebugUrl
		}
		if opts.Logging {
			logging = opts.Logging
			token = fmt.Sprintf("c_%s", token)
		}
	}

	return &Client{
		token:              token,
		logging:            logging,
		websocketURL:       fmt.Sprintf("%s://%s/%s/listen/", wsProto, apiHost, apiPrefix),
		localURL:           localURL,
		receiveURLTemplate: "https://play.svix.com/in/%s/",
		dialer: &websocket.Dialer{
			HandshakeTimeout: 10 * time.Second,
			Proxy:            http.ProxyFromEnvironment,
		},
		httpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: defaultTimeout,
		},
		stopRead:  make(chan struct{}, 10),
		stopWrite: make(chan struct{}, 10),

		errChan:  make(chan error, 10),
		sendChan: make(chan *OutgoingMessageEvent, 10),
		recChan:  make(chan *IncomingMessage, 10),
	}
}

type Stop = struct{}

func (c *Client) Listen(ctx context.Context) {
	if c.conn != nil {
		fmt.Printf("relay already listening\n")
		return
	}

	reconnectBackoffSteps := []time.Duration{
		0,
		100 * time.Millisecond,
		500 * time.Millisecond,
		time.Second,
		5 * time.Second,
	}
	reconnectAttempts := 0

	for {
		err := c.connect(ctx)
		if errors.Is(err, context.Canceled) {
			c.close()
			return
		}
		if err != nil {
			color.Red("Failed to connect to Webhook Relay:\n%s\n", err.Error())
			c.close()

			backoff := reconnectBackoffSteps[len(reconnectBackoffSteps)-1]
			if reconnectAttempts < len(reconnectBackoffSteps) {
				backoff = reconnectBackoffSteps[reconnectAttempts]
			}
			color.Yellow("Reattempting connection in %v\n", backoff)
			time.Sleep(backoff)

			reconnectAttempts += 1
			continue
		}

		// If the connection was successful reset reconnect counter
		reconnectAttempts = 0

		select {
		case <-ctx.Done():
			c.stopRead <- Stop{}
			c.stopWrite <- Stop{}
			return
		case <-c.errChan:
			c.stopRead <- Stop{}
			c.stopWrite <- Stop{}
			c.close()
			c.wg.Wait()
		}
	}
}
func (c *Client) close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *Client) changeConnection(conn *websocket.Conn) {
	c.conn = conn
	c.errChan = make(chan error, 10)
	c.stopRead = make(chan struct{}, 10)
	c.stopWrite = make(chan struct{}, 10)
}

func (c *Client) connect(ctx context.Context) error {
	var err error
	conn, _, err := c.dialer.Dial(c.websocketURL, nil)
	if err != nil {
		return err
	}
	c.changeConnection(conn)

	startMsgOut := &OutgoingMessageStart{
		Type:    MessageTypeStart,
		Version: version,
		Data: OutgoingMessageStartData{
			Token: c.token,
		},
	}

	err = c.conn.WriteJSON(startMsgOut)
	if err != nil {
		return err
	}

	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		if sErr, ok := err.(*websocket.CloseError); ok {
			if sErr.Code == websocket.ClosePolicyViolation {
				return fmt.Errorf("invalid token or already listening")
			}
		}
		return err
	}
	var startMsgIn IncomingMessageStart
	err = json.Unmarshal(msg, &startMsgIn)
	if err != nil {
		return err
	}
	url := fmt.Sprintf(c.receiveURLTemplate, startMsgIn.Data.Token)
	fmt.Printf(`Webhook relay is now listening at
%s

All requests on this endpoint will be forwarded to your local URL:
%s
`, pretty.MakeTerminalLink(url, url), c.localURL)
	if c.logging {
		viewUrl := fmt.Sprintf("https://play.svix.com/view/%s/", c.token)
		fmt.Printf(`
View logs and debug information at
%s
To disable logging run "svix listen --no-logging"
`,
			pretty.MakeTerminalLink(viewUrl, viewUrl),
		)
	}
	c.wg = &sync.WaitGroup{}
	c.wg.Add(2)

	go c.recLoop()
	go c.sendLoop()
	return nil
}

func (c *Client) SendMessage(msg *OutgoingMessageEvent) {
	c.sendChan <- msg
}

func (c *Client) recLoop() {
	defer c.wg.Done()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, packet, err := c.conn.ReadMessage()
		if err != nil {
			c.sendErrorMaybe(err, c.stopRead)
			return
		}
		go c.handleIncomingMessage(packet)
	}
}

func (c *Client) sendLoop() {
	defer c.wg.Done()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.sendChan:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}

			err := c.conn.WriteJSON(msg)
			if err != nil {
				// resend message when we reconnect
				c.SendMessage(msg)
				c.sendErrorMaybe(err, c.stopWrite)
				return
			}

		case <-ticker.C:
			if c.conn == nil {
				c.sendErrorMaybe(errors.New("Connection is closed"), c.stopWrite)
				return
			}

			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.sendErrorMaybe(err, c.stopWrite)
				return
			}

		case <-c.stopWrite:
			// ignore error
			return
		}
	}
}

func (c *Client) handleIncomingMessage(packet []byte) {
	var msg IncomingMessage
	if err := json.Unmarshal(packet, &msg); err != nil {
		return
	}
	switch msg.Type {
	case MessageTypeEvent:
		var msgData IncomingMessageEventData
		err := json.Unmarshal(msg.Data, &msgData)
		if err != nil {
			color.Red("Received Invalid Webhook message... skipping\n")
			return
		}
		color.Blue("<- Forwarding Message to: %s", c.localURL.String())
		res, err := c.makeLocalRequest(c.localURL, msgData)
		if err != nil {
			color.Red("Failed to make request to local server: \n%s\n", err.Error())
			return
		}

		c.processResponse(msgData.ID, res)
	default:
		return
	}
}

func formatRespHeaders(h http.Header) map[string]string {
	if h.Get("User-Agent") == "Go-http-client/1.1" {
		h.Set("User-Agent", "")
	}
	msgHeader := map[string]string{}
	for name, value := range h {
		msgHeader[name] = value[0]
	}
	return msgHeader
}

func (c *Client) makeLocalRequest(url *url.URL, msg IncomingMessageEventData) (*http.Response, error) {
	body, err := base64.StdEncoding.DecodeString(msg.Body)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: msg.Method,
		Header: http.Header{},
		URL:    url,
		Body:   io.NopCloser(bytes.NewReader(body)),
	}

	for name, value := range msg.Headers {
		if strings.ToLower(name) == "host" {
			// go requires the host to be set
			// explicitly otherwise it fails with
			// a "too many host headers" error
			req.Host = value
		} else {
			req.Header.Add(name, value)
		}
	}
	return http.DefaultClient.Do(req)
}

func (c *Client) processResponse(id string, res *http.Response) {
	buf, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	msg := &OutgoingMessageEvent{
		Type:    MessageTypeEvent,
		Version: version,
		Data: OutgoingMessageEventData{
			ID:      id,
			Status:  res.StatusCode,
			Headers: formatRespHeaders(res.Header),
			Body:    base64.StdEncoding.EncodeToString(buf),
		},
	}
	color.Green("-> Received \"%s\" response, forwarding to webhook sender\n", res.Status)
	c.SendMessage(msg)
}

func (c *Client) sendErrorMaybe(err error, stopChan chan (struct{})) {
	select {
	case <-stopChan:
	// dont send error if
	// stop was requested
	default:
		c.errChan <- err
	}
}
