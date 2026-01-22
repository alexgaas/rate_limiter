package internal

import (
	"github.com/alexgaas/bucket_quoter"

	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
)

type Api struct {
	g *CmdGlobal

	core *Core

	limiterMap map[string]*bucket_quoter.BucketQuoter
}

func CreateApi(g *CmdGlobal) (*Api, error) {
	var api Api

	api.g = g

	api.limiterMap = make(map[string]*bucket_quoter.BucketQuoter)
	for key, l := range api.g.Opts.Buckets.Buckets {
		api.limiterMap[key] = bucket_quoter.NewBucketQuoter(int64(l.Inflow), int64(l.Capacity), false, nil)
	}

	return &api, nil
}

const (
	API_UNIXSOCKET = 0
	API_HTTPS      = 1
)

// RunHTTPServer provide run http in tcp or unix socket mode
func (a *Api) RunHTTPServer(mode int) error {
	id := "(http) (server)"
	var err error

	r := a.routerEngine()
	if mode == API_UNIXSOCKET {

		socket := a.g.Opts.RateLimiter.Socket

		if Exists(socket) {
			err = os.Remove(socket)
			if err != nil {
				return err
			}
		}

		a.g.Log.Debug(fmt.Sprintf("%s run http server on socket:'%s'", id, socket))

		// http web server needs to write permissions
		// on unix socket to write
		unix.Umask(0000)

		err = r.RunUnix(socket)

		return err
	}
	if mode == API_HTTPS {
		certfile := a.g.Opts.RateLimiter.Certfile
		keyfile := a.g.Opts.RateLimiter.Keyfile
		port := a.g.Opts.RateLimiter.Port

		if port > 0 {
			a.g.Log.Debug(fmt.Sprintf("%s run http server on port:'%d'", id, port))
			err = r.RunTLS(fmt.Sprintf(":%d", port), certfile, keyfile)
			if err != nil {
				a.g.Log.Error(fmt.Sprintf("%s run failed on the port:'%d', err %s", id, port, err.Error()))
			}
		}
	}

	return err
}

// Main entry point for http requests
func (a *Api) routerEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(a.apiRequestLogger)

	r.GET("/ping", a.getPing)
	// register limiter API
	r.GET("/limiter", a.isAPIAvailableWithLimiter)

	return r
}

func (a *Api) Apiloop(wg *sync.WaitGroup, mode int) {
	defer wg.Done()

	a.RunHTTPServer(mode)
	return
}

func (a *Api) apiRequestLogger(c *gin.Context) {
	path, ip := a.apiRequestString(c)
	a.g.Log.Info(fmt.Sprintf("%s '%s %s' %d %s", ip, c.Request.Method, path,
		c.Writer.Status(), c.Request.UserAgent()))
	c.Next()
}

func (a *Api) apiRequestString(c *gin.Context) (string, string) {
	ip := c.ClientIP()
	if len(ip) == 0 {
		ip = "socket"
	}

	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}

	return path, ip
}

// ResponseInfo contains a code and message returned by the API as errors or
// informational messages inside the response.
type ResponseInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response is a template.  There will also be a result struct.  There will be a
// unique response type for each response, which will include this type.
type Response struct {
	Success  bool           `json:"success"`
	Errors   []ResponseInfo `json:"errors,omitempty"`
	Messages []ResponseInfo `json:"messages,omitempty"`
}

func (a *Api) apiSendError(c *gin.Context, code int, errStr string) {
	var r Response

	r.Success = false
	var ri ResponseInfo
	ri.Code = code
	ri.Message = errStr
	r.Errors = append(r.Errors, ri)

	response, _ := json.Marshal(r)

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(code, string(response))
}

func (a *Api) apiSendOK(c *gin.Context, code int, msgStr string) {
	var r Response

	r.Success = true
	var ri ResponseInfo
	ri.Code = code
	ri.Message = msgStr
	r.Messages = append(r.Messages, ri)

	responseBody, _ := json.Marshal(r)

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(code, string(responseBody))
}
