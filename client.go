package rebuilderd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	defaultEndpoint = "http:127.0.0.1:8484"

	AuthCookieHeader   = "X-Auth-Cookie"
	WorkerKeyHeader    = "X-Worker-Key"
	SignupSecretHeader = "X-Signup-Secret"
)

type Client struct {
	// baseURL
	url    url.URL
	client *http.Client
	tr     http.RoundTripper

	AuthCookie   string
	WorkerKey    string
	SignupSecret string
}

func NewClient(endpoint string, hc *http.Client) (*Client, error) {
	var c Client

	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	c.url = *u

	c.client = hc
	if c.client == nil {
		c.client = http.DefaultClient
	}

	c.tr = c.client.Transport
	if c.tr == nil {
		c.tr = http.DefaultTransport
	}
	c.client.Transport = &c

	return &c, nil
}

func (c *Client) RoundTrip(r *http.Request) (*http.Response, error) {
	if c.AuthCookie != "" {
		r.Header.Set(AuthCookieHeader, c.AuthCookie)
	}
	if c.SignupSecret != "" {
		r.Header.Set(SignupSecretHeader, c.SignupSecret)
	}
	if c.WorkerKey != "" {
		r.Header.Set(WorkerKeyHeader, c.WorkerKey)
	}
	return c.tr.RoundTrip(r)
}

// func (c *Client) BuildPing(ticket *QueueItem) error {
// 	err := c.postJSON("/api/v0/build/ping", ticket, nil)
// 	return err
// }
//
// func (c *Client) BuildReport(ticket *BuildReport) error {
// 	err := c.postJSON("/api/v0/build/report", ticket, nil)
// 	return err
// }
//
// // sync_suite?
// func (c *Client) PkgsSync(im *SuiteImport) error {
// 	err := c.postJSON("/api/v0/pkgs/sync", im, nil)
// 	return err
// }

// PkgsList lists the current packages and their status,
// no auth required
func (c *Client) PkgsList(list *ListPkgs) ([]PkgRelease, error) {
	var pkgs []PkgRelease
	err := c.getJSON("/api/v0/pkgs/list", list.Values(), &pkgs)
	return pkgs, err
}

// func (c *Client) QueueList(list *ListQueue) ([]QueueList, error) {
// 	var ql []QueueList
// 	err := c.postJSON("/api/v0/queue/list", list, &ql)
// 	return ql, err
// }
//
// func (c *Client) QueuePush(push *PushQueue) error {
// 	err := c.postJSON("/api/v0/queue/push", push, nil)
// 	return err
// }
//
// func (c *Client) QueuePop(query *WorkQuery) (*JobAssignment, error) {
// 	var ja JobAssignment
// 	err := c.postJSON("/api/v0/queue/pop", query, &ja)
// 	return &ja, err
// }
//
// func (c *Client) QueueDrop(query *DropQueueItem) error {
// 	err := c.postJSON("/api/v0/queue/drop", query, nil)
// 	return err
// }

// Workers lists the current workers
// requires AuthCookie
func (c *Client) Workers() ([]Worker, error) {
	var workers []Worker
	err := c.getJSON("/api/v0/workers", nil, &workers)
	return workers, err
}

func (c *Client) postJSON(path string, in, out interface{}) error {
	var r io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}
	u := c.url
	u.Path = path
	req, err := http.NewRequest(http.MethodPost, u.String(), r)
	if err != nil {
		return err
	}
	return c.doJSON(req, out)
}

func (c *Client) getJSON(path string, val url.Values, out interface{}) error {
	u := c.url
	u.Path = path
	u.RawQuery = val.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	return c.doJSON(req, out)
}

func (c *Client) doJSON(req *http.Request, out interface{}) error {
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf(res.Status)
	}
	r, err := httputil.DumpResponse(res, true)
	if err == nil {
		fmt.Println(string(r))
	}

	err = json.NewDecoder(res.Body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}
