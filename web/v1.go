package web

import (
	"bufio"
	"fmt"
	"github.com/XANi/rapporter/db"
	"github.com/efigence/go-mon"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (w *WebBackend) V1PostReport(c *gin.Context) {
	r := db.Report{}
	deviceId := c.Param("device_id")
	componentId := c.Param("component_id")
	status := strings.TrimSpace(strings.ToLower(c.Param("status")))
	raw := strings.ToLower(c.Query("raw"))
	rawContent := false
	if raw == "true" || raw == "1" {
		rawContent = true
	}

	if strings.Contains(strings.ToLower(c.ContentType()), "json") {
		err := c.BindJSON(&r)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		content, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("%s", err))
			return
		}
		if rawContent {
			scanner := bufio.NewScanner(strings.NewReader(string(content)))
			scanner.Split(bufio.ScanWords)
			for scanner.Scan() {
				r.Content = r.Content + "    " + scanner.Text() + "\n"
			}
		} else {
			r.Content = string(content)
		}
		r.Title = c.GetHeader("x-report-title")
		state, err := strconv.Atoi(c.GetHeader("x-report-state"))
		if len(status) > 0 {
			switch status {
			case "ok", "okay":
				state = int(mon.StateOk)
			case "warn", "warning":
				state = int(mon.StateWarning)
			case "crit", "critical":
				state = int(mon.StateCritical)
			case "unk", "unknown":
				state = int(mon.StateUnknown)
			default:
				c.String(http.StatusBadRequest,
					"only ok/warn/crit/unk are recognized status codes, not %q",
					status,
				)
			}
		} else if err != nil || state > 4 || state < 1 {
			c.String(http.StatusBadRequest,
				"x-report-state must be set to 1 (OK),2(warning),3(critical) or 4(unknown), not [%s] (%s)",
				c.GetHeader("x-report-state"),
			)

			return
		}
		r.State = mon.State(state)
		ttl, err := strconv.Atoi(c.GetHeader("x-report-ttl"))
		if len(c.GetHeader("x-report-ttl")) > 0 {
			if err != nil {
				c.String(http.StatusBadRequest,
					"x-report-ttl error: %s", err,
				)
				return
			} else {
				if ttl < 1 {
					c.String(http.StatusBadRequest,
						"x-report-ttl needs to be larger than 1")
					return
				}
				r.TTL = uint(ttl)
			}
		}
		if len(c.GetHeader("x-report-title")) > 0 {
			r.Title = c.GetHeader("x-report-title")
		} else {
			r.Title = componentId
		}
		if len(c.GetHeader("x-report-expire-in")) > 0 {
			ttl, _ := strconv.Atoi(c.GetHeader("x-report-expire-in"))
			r.ExpireIn = uint(ttl)
		}
		if len(c.GetHeader("x-report-category")) > 0 {
			r.Category = c.GetHeader("x-report-category")
		}
	}
	r.DeviceID = deviceId
	r.ComponentID = componentId
	err := r.Validate()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	_, err = w.db.AddReport(r)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, c.ContentType())
}
func (w *WebBackend) V1DeleteReport(c *gin.Context) {
	deviceId := c.Param("device_id")
	componentId := c.Param("component_id")
	err := w.db.DeleteReport(deviceId, componentId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.String(http.StatusOK, "ok")
	}
}
