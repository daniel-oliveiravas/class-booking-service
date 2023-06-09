package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Readiness(c *gin.Context) {
	ctx := c.Request.Context()

	status := "ok"
	statusCode := http.StatusOK
	if err := h.cfg.PgProbe.Check(ctx); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, gin.H{"status": status})
}

func (h *Handler) Liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}
