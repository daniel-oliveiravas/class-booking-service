package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/daniel-oliveiravas/class-booking-service/business/classes"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AddClass(c *gin.Context) {
	var newClass classes.NewClass
	err := c.BindJSON(&newClass)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("missing body").Error()})
			return
		}
		h.cfg.Logger.Debugw("failed to bind class", "error", err.Error())
		return
	}

	ctx := c.Request.Context()

	class, err := h.cfg.ClassesUsecase.AddClass(ctx, newClass)
	if err != nil {
		if errors.Is(err, classes.ErrInvalidData) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		h.cfg.Logger.Errorw("failed to add new class", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add new class"})
		return
	}

	c.JSON(http.StatusCreated, class)
}

func (h *Handler) GetClassByID(c *gin.Context) {
	classID := c.Param("id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing path param: id"})
		return
	}

	ctx := c.Request.Context()
	class, err := h.cfg.ClassesUsecase.GetByID(ctx, classID)
	if err != nil {
		if errors.Is(err, classes.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("class with ID %s not found", classID)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get class. Retry later"})
		return
	}

	c.JSON(http.StatusOK, class)
}

func (h *Handler) UpdateClass(c *gin.Context) {
	classID := c.Param("id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing path param: id"})
		return
	}

	ctx := c.Request.Context()

	var updateClass classes.UpdateClass
	err := c.BindJSON(&updateClass)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("missing body").Error()})
			return
		}
		h.cfg.Logger.Debugw("failed to bind class", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse request body"})
		return
	}

	updatedClass, err := h.cfg.ClassesUsecase.UpdateClass(ctx, classID, updateClass)
	if err != nil {
		if errors.Is(err, classes.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("class with ID %s not found", classID)})
			return
		}
		h.cfg.Logger.Debugw("failed to update class: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update class"})
		return
	}

	c.JSON(http.StatusOK, updatedClass)
}

func (h *Handler) DeleteClass(c *gin.Context) {
	classID := c.Param("id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing path param: id"})
		return
	}

	ctx := c.Request.Context()

	err := h.cfg.ClassesUsecase.DeleteClass(ctx, classID)
	if err != nil {
		h.cfg.Logger.Debugw("failed to delete class: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete class"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListClasses(c *gin.Context) {
	pageInfo, err := h.extractClassesPageInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query params are invalid"})
		return
	}

	ctx := c.Request.Context()

	allClasses, err := h.cfg.ClassesUsecase.ListClasses(ctx, pageInfo)
	if err != nil {
		h.cfg.Logger.Debugw("failed to list classes: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list classes"})
		return
	}

	c.JSON(http.StatusOK, allClasses)
}

func (h *Handler) extractClassesPageInfo(c *gin.Context) (classes.PageInfo, error) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var page int
	var limit int
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			h.cfg.Logger.Debugw("failed to parse page param: %w", err)
			return classes.PageInfo{}, err
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.cfg.Logger.Debugw("failed to parse limit param: %w", err)
			return classes.PageInfo{}, err
		}
	}

	return classes.PageInfo{
		Limit: limit,
		Page:  page,
	}, nil
}
