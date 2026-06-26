package handler

import (
	"Test_App/internal/domain"
	"errors"
	"net/http"
	"strconv"

	"Test_App/internal/usecase"
	"github.com/gin-gonic/gin"
)

type SectionHandler struct {
	uc domain.SectionUseCase
}

func NewSectionHandler(uc domain.SectionUseCase) *SectionHandler {
	return &SectionHandler{uc: uc}
}

// GetSectionsBySubjectID godoc
// GET /subjects/:id/sections
func (h *SectionHandler) GetSectionsBySubjectID(c *gin.Context) {
	subjectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id предмета"})
		return
	}

	sections, err := h.uc.GetSectionsBySubjectID(c.Request.Context(), subjectID)
	if err != nil {
		if errors.Is(err, usecase.ErrSubjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения разделов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sections})
}
