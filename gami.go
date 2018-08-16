// Copyright 2014 Jovany Leandro G.C <bit4bit@riseup.net>. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

package gami

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	errNotEvent    = errors.New("not a AMI Event")
	errNotResponse = errors.New("not a AMI Response")
)

// Params for the actions
type Params map[string]string

// AMIClient a connection to AMI server
type AMIClient struct {
	bufReader        *bufio.Reader
	connRaw          net.Conn
	mutexAsyncAction *sync.RWMutex

	address     string
	amiUser     string
	amiPass     string
	useTLS      bool
	unsecureTLS bool
	loggedIn    bool

	// TLSConfig for secure connections
	tlsConfig *tls.Config

	// network wait for a new connection
	waitNewConnection chan struct{}

	response map[string]chan *AMIResponse

	// Events for client parse
	Events chan *AMIEvent

	// Error Raise on logic
	Error chan error

	//NetError a network error
	NetError chan error
}

// AMIResponse from action
type AMIResponse struct {
	ID     string
	Status string
	Params map[string]string
}

// AMIEvent it's a representation of Event readed
type AMIEvent struct {
	//Identification of event Event: xxxx
	ID string

	Privilege []string

	// Params  of arguments received
	Params map[string]string
}

func UseTLS(c *AMIClient) {
	c.useTLS = true
}

func UseTLSConfig(config *tls.Config) func(*AMIClient) {
	return func(c *AMIClient) {
		c.tlsConfig = config
		c.useTLS = true
	}
}

func UnsecureTLS(c *AMIClient) {
	c.unsecureTLS = true
}

// Login authenticate to AMI
func (client *AMIClient) Login(username, password string) error {
	response, err := client.Action("Login", Params{"Username": username, "Secret": password}, time.Second*5)
	if err != nil {
		return err
	}

	if (*response).Status == "Error" {
		return errors.New((*response).Params["Message"])
	}
	client.loggedIn = true
	client.amiUser = username
	client.amiPass = password
	return nil
}

// Reconnect the session, autologin if a new network error it put on client.NetError
func (client *AMIClient) Reconnect() error {
	client.connRaw.Close()
	client.loggedIn = false
	client.bufReader = nil
	err := client.NewConn()

	if err != nil {
		client.NetError <- err
		return err
	}

	client.waitNewConnection <- struct{}{}

	if err := client.Login(client.amiUser, client.amiPass); err != nil {
		return err
	}

	return nil
}

// AsyncAction return chan for wait response of action with parameter *ActionID* this can be helpful for
// massive actions,
func (client *AMIClient) AsyncAction(action string, params Params) (<-chan *AMIResponse, error) {
	var outBuf bytes.Buffer

	client.mutexAsyncAction.Lock()
	defer client.mutexAsyncAction.Unlock()

	if params == nil {
		params = Params{}
	}
	if _, ok := params["ActionID"]; !ok {
		params["ActionID"] = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	if _, ok := client.response[params["ActionID"]]; !ok {
		client.response[params["ActionID"]] = make(chan *AMIResponse, 1)
	}
	outBuf.WriteString(fmt.Sprintf("Action: %s\n", strings.TrimSpace(action)))
	for k, v := range params {
		outBuf.WriteString(fmt.Sprintf("%s: %s\n", k, strings.TrimSpace(v)))
	}
	outBuf.WriteString("\n")
	if _, err := client.connRaw.Write(outBuf.Bytes()); err != nil {
		return nil, err
	}

	return client.response[params["ActionID"]], nil
}

func (client *AMIClient) AsyncActionNoResponse(action string, params Params) error {
	_, err := client.AsyncAction(action, params)
	return err
}

// Action send with params, wait for response and return the response
func (client *AMIClient) Action(action string, params Params, actionTimeout time.Duration) (*AMIResponse, error) {
	var ActionID string

	if _, ok := params["ActionID"]; !ok {
		ActionID = fmt.Sprintf("%d", time.Now().UnixNano())
		params["ActionID"] = ActionID
	} else {
		ActionID = params["ActionID"]
	}
	resp, err := client.AsyncAction(action, params)
	if err != nil {
		return nil, err
	}
	time.AfterFunc(actionTimeout, func() {
		client.mutexAsyncAction.Lock()
		if c, exists := client.response[ActionID]; exists {
			delete(client.response, ActionID)
			client.mutexAsyncAction.Unlock()
			response := &AMIResponse{ID: ActionID, Status: "Error", Params: make(map[string]string)}
			response.Params["Error"] = "Timeout"
			c <- response
			close(c)
			return
		}
		client.mutexAsyncAction.Unlock()
	})
	response := <-resp
	return response, nil
}

func readMessage(r *bufio.Reader) (map[string]string, error) {
	var (
		responseFollows bool
		responseBuffer  []string
		err             error
		kv              []byte
	)
	m := make(map[string]string)
	for {
		kv, _, err = r.ReadLine()
		//fmt.Printf("ReadLine (%s): '%s'\n", err, string(kv))
		if len(kv) == 0 {
			break
		}

		var key string
		i := bytes.IndexByte(kv, ':')
		if i >= 0 {
			endKey := i
			for endKey > 0 && kv[endKey-1] == ' ' {
				endKey--
			}
			key = string(kv[:endKey])
		}
		if key == "" && !responseFollows {
			if err != nil {
				break
			}
			continue
		}

		if responseFollows && (key != "Privilege" && key != "ActionID") {
			if string(kv) != "--END COMMAND--" {
				responseBuffer = append(responseBuffer, string(kv))
			}
			if err != nil {
				break
			}
			continue
		}

		i++
		for i < len(kv) && (kv[i] == ' ' || kv[i] == '\t') {
			i++
		}
		value := string(kv[i:])
		if key == "Response" && value == "Follows" {
			responseFollows = true
			responseBuffer = make([]string, 0, 64)
		}

		m[key] = value

		if err != nil {
			break
		}
	}
	if responseFollows {
		m["CommandResponse"] = strings.Join(responseBuffer, "\n")
	}
	return m, err
}

// Run process socket waiting events and responses
func (client *AMIClient) Run() {
	go func() {
		for {
			data, err := readMessage(client.bufReader)
			if err != nil {
				switch err {
				case syscall.ECONNABORTED:
					fallthrough
				case syscall.ECONNRESET:
					fallthrough
				case syscall.ECONNREFUSED:
					fallthrough
				case io.EOF:
					client.NetError <- err
					<-client.waitNewConnection
				default:
					client.Error <- err
				}
				continue
			}
			if ev, err := newEvent(data); err != nil {
				if err != errNotEvent {
					client.Error <- err
				}
			} else {
				client.Events <- ev
			}

			if response, err := newResponse(data); err == nil {
				client.notifyResponse(response)
			}

		}
	}()
}

// Close the connection to AMI
func (client *AMIClient) Close() {
	if client.loggedIn {
		client.Action("Logoff", nil, time.Second*3)
	}
	client.connRaw.Close()
}

func (client *AMIClient) notifyResponse(response *AMIResponse) {
	go func() {
		client.mutexAsyncAction.RLock()
		_, exists := client.response[response.ID]
		client.mutexAsyncAction.RUnlock()
		if exists {
			client.mutexAsyncAction.Lock()
			if c, exists := client.response[response.ID]; exists {
				delete(client.response, response.ID)
				client.mutexAsyncAction.Unlock()
				c <- response
				close(c)
				return
			}
			client.mutexAsyncAction.Unlock()
		}
	}()
}

//newResponse build a response for action
func newResponse(data map[string]string) (*AMIResponse, error) {
	r, found := data["Response"]
	if !found {
		return nil, errNotResponse
	}
	response := &AMIResponse{
		ID:     data["ActionID"],
		Status: r,
		Params: make(map[string]string),
	}
	for k, v := range data {
		if k == "Response" {
			continue
		}
		response.Params[k] = v
	}
	return response, nil
}

//newEvent build event
func newEvent(data map[string]string) (*AMIEvent, error) {
	e, found := data["Event"]
	if !found {
		return nil, errNotEvent
	}
	ev := &AMIEvent{
		ID:        e,
		Privilege: strings.Split(data["Privilege"], ","),
		Params:    make(map[string]string),
	}
	for k, v := range data {
		if k == "Event" || k == "Privilege" {
			continue
		}
		ev.Params[k] = v
	}
	return ev, nil
}

// Dial create a new connection to AMI
func Dial(address string, options ...func(*AMIClient)) (*AMIClient, error) {
	client := &AMIClient{
		address:           address,
		amiUser:           "",
		amiPass:           "",
		mutexAsyncAction:  new(sync.RWMutex),
		waitNewConnection: make(chan struct{}),
		response:          make(map[string]chan *AMIResponse),
		Events:            make(chan *AMIEvent, 100),
		Error:             make(chan error, 1),
		NetError:          make(chan error, 1),
		useTLS:            false,
		unsecureTLS:       false,
		tlsConfig:         new(tls.Config),
	}
	for _, op := range options {
		op(client)
	}
	err := client.NewConn()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewConn create a new connection to AMI
func (client *AMIClient) NewConn() (err error) {
	if client.useTLS {
		client.tlsConfig.InsecureSkipVerify = client.unsecureTLS
		nd := new(net.Dialer)
		nd.Timeout = time.Second * 10
		client.connRaw, err = tls.DialWithDialer(nd, "tcp", client.address, client.tlsConfig)
	} else {
		client.connRaw, err = net.DialTimeout("tcp", client.address, time.Second*10)
	}

	if err != nil {
		return err
	}

	client.bufReader = bufio.NewReader(client.connRaw)
	return nil
}
