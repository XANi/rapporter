package web

import (
	"fmt"
	"github.com/XANi/rapporter/db"
	"github.com/efigence/go-mon"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

type WebBackend struct {
	l   *zap.SugaredLogger
	al  *zap.SugaredLogger
	r   *gin.Engine
	cfg *Config
	db  *db.DB
}

type Config struct {
	Logger       *zap.SugaredLogger `yaml:"-"`
	AccessLogger *zap.SugaredLogger `yaml:"-"`
	DB           *db.DB             `yaml:"-"`
	ListenAddr   string             `yaml:"listen_addr"`
}

func New(cfg Config, webFS fs.FS) (backend *WebBackend, err error) {
	if cfg.Logger == nil {
		panic("missing logger")
	}
	if len(cfg.ListenAddr) == 0 {
		panic("missing listen addr")
	}
	w := WebBackend{
		l:   cfg.Logger,
		al:  cfg.AccessLogger,
		cfg: &cfg,
	}
	if cfg.AccessLogger == nil {
		w.al = w.l //.Named("accesslog")
	}
	r := gin.New()
	w.r = r
	w.db = cfg.DB
	gin.SetMode(gin.ReleaseMode)
	t, err := template.ParseFS(webFS, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("error loading templates: %s", err)
	}
	r.SetHTMLTemplate(t)
	// for zap logging
	r.Use(ginzap.GinzapWithConfig(w.al.Desugar(), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        false,
		SkipPaths:  []string{"/_status/health", "/_status/metrics"},
	}))
	//r.Use(ginzap.RecoveryWithZap(w.al.Desugar(), true))
	// basic logging to stdout
	//r.Use(gin.LoggerWithWriter(os.Stdout))
	r.Use(gin.Recovery())

	// monitoring endpoints
	r.GET("/_status/health", gin.WrapF(mon.HandleHealthcheck))
	r.HEAD("/_status/health", gin.WrapF(mon.HandleHealthcheck))
	r.GET("/_status/metrics", gin.WrapF(mon.HandleMetrics))
	defer mon.GlobalStatus.Update(mon.StatusOk, "ok")
	// healthcheckHandler, haproxyStatus := mon.HandleHealthchecksHaproxy()
	// r.GET("/_status/metrics", gin.WrapF(healthcheckHandler))

	httpFS := http.FileServer(http.FS(webFS))
	r.GET("/s/*filepath", func(c *gin.Context) {
		// content is embedded under static/ dir
		p := strings.Replace(c.Request.URL.Path, "/s/", "/static/", -1)
		c.Request.URL.Path = p
		//c.Header("Cache-Control", "public, max-age=3600, immutable")
		httpFS.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"notfound": c.Request.URL.Path,
		})
	})

	r.GET("/list", func(c *gin.Context) {
		reports, _ := w.db.GetLatestReports()
		cfg.Logger.Infof("reports: %d", len(reports))
		c.HTML(http.StatusOK, "list.tmpl", gin.H{
			"reports": reports,
		})
	})
	r.GET("/report/:device_id/:component_id", func(c *gin.Context) {
		deviceId := c.Param("device_id")
		componentId := c.Param("component_id")
		report, err := w.db.GetReport(deviceId, componentId)
		if err != nil {
			c.Header("Refresh", "1;/list")
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
				"notfound": c.Request.URL.Path,
			})
			return
		}
		content := w.markdownParse(report.Content)

		c.HTML(http.StatusOK, "report.tmpl", gin.H{
			"title":   report.Title,
			"report":  report,
			"content": template.HTML(content),
		})
	})
	rav1 := r.Group("/api/v1")
	rav1.POST("/report/:device_id/:component_id", w.V1PostReport)
	rav1.DELETE("/report/:device_id/:component_id", w.V1DeleteReport)
	rav1.POST("/report/:device_id/:component_id/:status", w.V1PostReport)
	return &w, nil
}

func (b *WebBackend) Run() error {
	b.l.Infof("listening on %s", b.cfg.ListenAddr)
	return b.r.Run(b.cfg.ListenAddr)
}
