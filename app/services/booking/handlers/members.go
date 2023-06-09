package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/daniel-oliveiravas/class-booking-service/business/members"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AddMember(c *gin.Context) {
	var newMember members.NewMember
	err := c.BindJSON(&newMember)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("missing body").Error()})
			return
		}
		h.cfg.Logger.Debugw("failed to bind member", "error", err.Error())
		return
	}

	ctx := c.Request.Context()

	member, err := h.cfg.MembersUsecase.AddMember(ctx, newMember)
	if err != nil {
		if errors.Is(err, members.ErrInvalidData) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		h.cfg.Logger.Errorw("failed to add new member", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add new member"})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func (h *Handler) GetMemberByID(c *gin.Context) {
	memberID := c.Param("id")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing path param: id"})
		return
	}

	ctx := c.Request.Context()
	member, err := h.cfg.MembersUsecase.GetByID(ctx, memberID)
	if err != nil {
		if errors.Is(err, members.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("member with ID %s not found", memberID)})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get member. Retry later"})
		return
	}

	c.JSON(http.StatusOK, member)
}

func (h *Handler) UpdateMember(c *gin.Context) {
	memberID := c.Param("id")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing path param: id"})
		return
	}

	ctx := c.Request.Context()

	var updateMember members.UpdateMember
	err := c.BindJSON(&updateMember)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("missing body").Error()})
			return
		}
		h.cfg.Logger.Debugw("failed to bind member", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse request body"})
		return
	}

	updatedMember, err := h.cfg.MembersUsecase.UpdateMember(ctx, memberID, updateMember)
	if err != nil {
		if errors.Is(err, members.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("member with ID %s not found", memberID)})
			return
		}
		h.cfg.Logger.Debugw("failed to update member: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update member"})
		return
	}

	c.JSON(http.StatusOK, updatedMember)
}

func (h *Handler) DeleteMember(c *gin.Context) {
	memberID := c.Param("id")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing path param: id"})
		return
	}

	ctx := c.Request.Context()

	err := h.cfg.MembersUsecase.DeleteMember(ctx, memberID)
	if err != nil {
		h.cfg.Logger.Debugw("failed to delete member: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete member"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ListMembers(c *gin.Context) {
	pageInfo, err := h.extractPageInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query params are invalid"})
		return
	}

	ctx := c.Request.Context()

	allMembers, err := h.cfg.MembersUsecase.ListMembers(ctx, pageInfo)
	if err != nil {
		h.cfg.Logger.Debugw("failed to list members: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list members"})
		return
	}

	c.JSON(http.StatusOK, allMembers)
}

func (h *Handler) extractPageInfo(c *gin.Context) (members.PageInfo, error) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	var page int
	var limit int
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			h.cfg.Logger.Debugw("failed to parse page param: %w", err)
			return members.PageInfo{}, err
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.cfg.Logger.Debugw("failed to parse limit param: %w", err)
			return members.PageInfo{}, err
		}
	}

	return members.PageInfo{
		Limit: limit,
		Page:  page,
	}, nil
}
