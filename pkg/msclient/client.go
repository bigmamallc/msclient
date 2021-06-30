package msclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	cfg *Cfg
	log zerolog.Logger
	mx  *metrics

	httpClient *http.Client
}

var RedirectedError = errors.New("got a redirect")

func New(cfg *Cfg, log zerolog.Logger, mxReg *prometheus.Registry, mxSubsystem string) (c *Client, err error) {
	c = &Client{
		cfg: cfg,
		log: log,
		mx:  newMetrics(mxReg, mxSubsystem),
	}

	c.httpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return RedirectedError
		},
		Timeout: cfg.RequestTimeout,
	}

	return
}

type UnexpectedStatusCodeErr struct {
	Expected []int
	Got      int
}
func (e *UnexpectedStatusCodeErr) Error() string {
	return fmt.Sprintf("expected status codes %+v got %d", e.Expected, e.Got)
}

type ResponseReadError struct {
	Cause error
}
func (e *ResponseReadError) Error() string {
	return fmt.Sprintf("response read error: %v", e.Cause)
}
func (e *ResponseReadError) Unwrap() error {
	return e.Cause
}

type ResponseParseError struct {
	Cause error
}
func (e *ResponseParseError) Error() string {
	return fmt.Sprintf("response parse error: %v", e.Cause)
}
func (e *ResponseParseError) Unwrap() error {
	return e.Cause
}

func (c *Client) Get(uri string, resp interface{}, args... interface{}) error {
	mxLabels := map[string]string{labelURI: uri, labelMethod: "GET"}
	c.mx.requests.With(mxLabels)

	timer := prometheus.NewTimer(c.mx.successfulResponseTimes.With(mxLabels))

	urlEncodedArgs := make([]interface{}, len(args))
	for i, v := range args {
		urlEncodedArgs[i] = url.QueryEscape(fmt.Sprintf("%v", v))
	}

	reqUrl := c.cfg.BaseURL + fmt.Sprintf(uri, urlEncodedArgs...)
	httpResp, err := c.httpClient.Get(reqUrl)
	if err != nil {
		c.mx.requestErrors.With(mxLabels).Inc()
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		c.mx.statusCodeErrors.With(mxLabels).Inc()
		return &url.Error{
			Op:  "GET",
			URL: reqUrl,
			Err: &UnexpectedStatusCodeErr{
				Expected: []int{http.StatusOK},
				Got:      httpResp.StatusCode,
			},
		}
	}

	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		c.mx.responseReadErrors.With(mxLabels).Inc()
		return &url.Error{
			Op:  "GET",
			URL: reqUrl,
			Err: &ResponseReadError{Cause: err},
		}
	}

	if err := json.Unmarshal(data, resp); err != nil {
		c.mx.responseParseErrors.With(mxLabels).Inc()
		return &url.Error{
			Op: "GET",
			URL: reqUrl,
			Err: &ResponseParseError{Cause: err},
		}
	}

	timer.ObserveDuration()
	c.mx.successfulResponses.With(mxLabels).Inc()
	return nil
}
