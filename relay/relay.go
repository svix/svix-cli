package relay

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
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

	conn              *websocket.Conn
	connResetInterval time.Duration
	stopRead          chan struct{}
	stopWrite         chan struct{}

	errChan chan error

	sendChan chan *MessageOut
	recChan  chan *MessageIn
	wg       *sync.WaitGroup
}

type ClientOptions struct {
	DisableSecurity bool
	RelayDebugUrl   string
}

func NewClient(token string, localURL *url.URL, opts *ClientOptions) *Client {
	wsProto := "wss"
	httpProto := "https"
	apiHost := defaultAPIHost
	token = fmt.Sprintf("c.%s", token)
	if opts != nil {
		if opts.DisableSecurity {
			wsProto = "ws"
			httpProto = "http"
		}
		if opts.RelayDebugUrl != "" {
			apiHost = opts.RelayDebugUrl
		}
	}

	return &Client{
		token:              token,
		websocketURL:       fmt.Sprintf("%s://%s/%s/play/listen/%s/", wsProto, apiHost, apiPrefix, token),
		localURL:           localURL,
		receiveURLTemplate: fmt.Sprintf("%s://%s/%s/play/receive/%%s/", httpProto, apiHost, apiPrefix),
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
		connResetInterval: time.Minute,
		stopRead:          make(chan struct{}),
		stopWrite:         make(chan struct{}),

		errChan: make(chan error),

		// TODO should these be buffered?
		sendChan: make(chan *MessageOut),
		recChan:  make(chan *MessageIn),
	}
}

func (c *Client) Listen(ctx context.Context) {
	if c.conn != nil {
		fmt.Printf("relay already listening\n")
		return
	}

	for {
		err := c.connect(ctx)
		if err != nil {
			color.Red("Failed to connect to Webhook Relay:\n%s\n", err.Error())
			c.close()
			return
		}

		select {
		case <-ctx.Done():
			close(c.sendChan)
			close(c.recChan)

			close(c.stopRead)
			close(c.stopWrite)
			return
		case err = <-c.errChan:
			close(c.stopRead)
			close(c.stopWrite)
			c.wg.Wait()

			if sErr, ok := err.(*websocket.CloseError); ok {
				if sErr.Code == websocket.ClosePolicyViolation {
					color.Red("Unrecoverable error, please try again.")
					return
				}
			}
		case <-time.After(c.connResetInterval):
			close(c.stopRead)
			close(c.stopWrite)

			if c.conn != nil {
				c.conn.Close()
			}

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
	c.errChan = make(chan error)
	c.stopRead = make(chan struct{})
	c.stopWrite = make(chan struct{})
}

func (c *Client) connect(ctx context.Context) error {
	var err error
	conn, _, err := c.dialer.Dial(c.websocketURL, nil)
	if err != nil {
		return err
	}
	c.changeConnection(conn)

	url := fmt.Sprintf(c.receiveURLTemplate, c.token)
	playURL := fmt.Sprintf("https://play.svix.com/view/%s/", c.token)
	fmt.Printf(`Forwarding requests from
%s
to
%s

View logs and debug information on Svix Play:
%s
`, pretty.MakeTerminalLink(url, url), c.localURL, pretty.MakeTerminalLink(playURL, playURL))

	c.wg = &sync.WaitGroup{}
	c.wg.Add(2)

	go c.recLoop()
	go c.sendLoop()
	return nil
}

func (c *Client) SendMessage(msg *MessageOut) {
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
			select {
			case <-c.stopRead:
				// dont send error if
				// stop was requested
			default:
				c.errChan <- err
			}
			return
		}
		go c.handleIncommingMessage(packet)
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
				c.errChan <- err
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.errChan <- err
				return
			}
		case <-c.stopWrite:
			// ignore error
			return
		}
	}
}

func (c *Client) handleIncommingMessage(packet []byte) {
	var msg MessageIn
	if err := json.Unmarshal(packet, &msg); err != nil {
		color.Red("Recieved Invalid Message message... skipping\n")
		return
	}

	color.Blue("<- Forwarding Message to: %s", c.localURL.String())
	res, err := c.makeLocalRequest(c.localURL, msg)
	if err != nil {
		color.Red("Failed to make request to local server: \n%s\n", err.Error())
		return
	}

	c.processResponse(msg.ID, res)
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

func (c *Client) makeLocalRequest(url *url.URL, msg MessageIn) (*http.Response, error) {
	body, err := base64.StdEncoding.DecodeString(msg.Body)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
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

	msg := &MessageOut{
		ID:         id,
		StatusCode: res.StatusCode,
		Headers:    formatRespHeaders(res.Header),
		Body:       base64.StdEncoding.EncodeToString(buf),
	}
	color.Green("-> Recieved \"%s\" response, forwarding to webhook sender\n", res.Status)
	c.SendMessage(msg)
}
