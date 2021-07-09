package handlers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handle struct {
	Logger    *logrus.Logger
	RemoteURL string
}

func New(remoteURL, logSeverity string) (*Handle, error) {

	if remoteURL == "" {
		return nil, fmt.Errorf("remoteURL parameter is empty")
	}

	logger := &logrus.Logger{
		Out:          nil,
		Hooks:        nil,
		Formatter:    nil,
		ReportCaller: false,
		Level:        0,
		ExitFunc:     nil,
	}

	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		ForceQuote:                true,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	switch logSeverity {
	case "Info":
		logger.SetLevel(logrus.InfoLevel)
	case "Debug":
		logger.SetLevel(logrus.DebugLevel)
	case "Warn":
		logger.SetLevel(logrus.WarnLevel)
	default:
		return nil, fmt.Errorf("wrong logger level: %v", logSeverity)
	}

	return &Handle{
		Logger:    logger,
		RemoteURL: remoteURL,
	}, nil
}

// HandleAnyRoute is handler for / hosted at http://localhost:8080
func (h *Handle) HandleAnyRoute(c *gin.Context) {
	path := c.Param("path")
	parsedUrl := fmt.Sprint(h.RemoteURL, path)
	remote, err := url.Parse(parsedUrl)
	if err != nil {
		h.Logger.Errorf("Unable to parse remote address: %v", err)
		c.JSON(http.StatusNotFound, "Unable to parse request")
	}

	clientMethod := c.Request.Method

	h.Logger.Debug("Preparing PROXY...")

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Header.Set("X-Forwarded-For", req.Host)
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")
	}

	h.Logger.Debug("Configuring underlying request for HTTP method", clientMethod)
	h.Logger.Debug("Making underlying request...")
	proxy.ServeHTTP(c.Writer, c.Request)
	h.Logger.Debug("Request was ended without issues")
	h.Logger.Debug("Underlying HTTP status code: ", c.Writer.Status())
}
