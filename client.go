package pipehub

import (
	"context"

	"github.com/pkg/errors"
)

// ClientConfig holds the client configuration.
type ClientConfig struct {
	Server          ClientConfigServer
	HTTP            []ClientConfigHTTP
	Pipe            []ClientConfigPipe
	AsyncErrHandler func(error)
}

func (c *ClientConfig) setDefaultValues() {
	c.Server.setDefaultValues()
	if c.AsyncErrHandler == nil {
		c.AsyncErrHandler = func(error) {}
	}
}

// ClientConfigServer holds the server configuration.
type ClientConfigServer struct {
	HTTP   ClientConfigServerHTTP
	Action ClientConfigServerAction
}

func (c *ClientConfigServer) setDefaultValues() {
	c.HTTP.setDefaultValues()
}

// ClientConfigServerAction holds the action server configuration.
type ClientConfigServerAction struct {
	NotFound string
	Panic    string
}

// ClientConfigServerHTTP holds the http server configuration.
type ClientConfigServerHTTP struct {
	Port int
}

func (c *ClientConfigServerHTTP) setDefaultValues() {
	if c.Port == 0 {
		c.Port = 80
	}
}

// ClientConfigHTTP holds the configuration to direct the request through pipes.
type ClientConfigHTTP struct {
	Endpoint string
	Handler  string
}

// ClientConfigPipe holds the pipe configuration.
type ClientConfigPipe struct {
	Path    string
	Module  string
	Version string
	Config  map[string]interface{}
}

// Client is pipehub entrypoint.
type Client struct {
	cfg         ClientConfig
	server      *server
	pipeManager *pipeManager
}

// Start pipehub.
func (c *Client) Start() error {
	if err := c.server.start(); err != nil {
		return errors.Wrap(err, "server start error")
	}
	return nil
}

// Stop the pipehub.
func (c *Client) Stop(ctx context.Context) error {
	if err := c.server.stop(ctx); err != nil {
		return errors.Wrap(err, "server stop error")
	}

	if err := c.pipeManager.close(ctx); err != nil {
		return errors.Wrap(err, "pipe manager close error")
	}
	return nil
}

// nolint: gocritic
func (c *Client) init(cfg ClientConfig) error {
	cfg.setDefaultValues()
	c.cfg = cfg
	if c.cfg.AsyncErrHandler == nil {
		c.cfg.AsyncErrHandler = func(error) {}
	}

	var err error
	c.pipeManager, err = newPipeManager(c)
	if err != nil {
		return errors.Wrap(err, "pipe manager new instance error")
	}
	c.server = newServer(c)
	return nil
}

// NewClient return a configured pipehub client.
// nolint: gocritic
func NewClient(cfg ClientConfig) (Client, error) {
	var c Client
	if err := c.init(cfg); err != nil {
		return c, errors.Wrap(err, "client new instance error")
	}
	return c, nil
}
