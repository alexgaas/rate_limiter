package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
)

// HTTPClientOpt is a function that sets client behavior or basic data.
// Example:
//
//	c, err := NewHTTPClient(
//	    AsService("lambda"),
//	    WithApiSubId("c06e5a28-0bdf-4d4e-9dd0-86a5a51e9a0d"),
//	    WithHTTPHost("https://localhost:8443"),
//	)

type ClientOpt func(c *Client) error

// AsService supplies service name for http client for debugging.
func AsService(name string) ClientOpt {
	return func(c *Client) error {
		if name == "" {
			return errors.New("laas: service name cannot be empty")
		}
		c.httpc.SetQueryParam("service", name)
		return nil
	}
}

func WithApiSubId(id string) ClientOpt {
	return func(c *Client) error {
		c.httpc.SetHeader("X-Limiter-Subscription-ID", id)
		return nil
	}
}

// WithHTTPHost rewrites default HTTP host in client.
func WithHTTPHost(host string) ClientOpt {
	return func(c *Client) error {
		c.httpc.SetBaseURL(host)
		return nil
	}
}

func WithTLSCert(caCertPath string) ClientOpt {
	return func(c *Client) error {
		caCert, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			return fmt.Errorf("failed to read CA cert file: %s", caCertPath)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			Renegotiation: tls.RenegotiateOnceAsClient,
			RootCAs:       caCertPool,
		}

		c.httpc.SetTLSClientConfig(tlsConfig)
		return nil
	}
}

func WithInsecureSkipVerify() ClientOpt {
	return func(c *Client) error {
		c.httpc.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		return nil
	}
}

// WithDebug enables debug resty output
func WithDebug(enable bool) ClientOpt {
	return func(c *Client) error {
		c.httpc.SetDebug(enable)
		return nil
	}
}
