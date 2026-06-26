package handler

import (
	"Test_App/internal/domain"
	"errors"
	"net/http"
	"strconv"

	"Test_App/internal/usecase"
	"github.com/gin-gonic/gin"
)

type SubjectHandler struct {
	uc domain.SubjectUseCase
}

func NewSubjectHandler(uc domain.SubjectUseCase) *SubjectHandler {
	return &SubjectHandler{uc: uc}
}

// GetAllSubjects godoc
// GET /subjects
func (h *SubjectHandler) GetAllSubjects(c *gin.Context) {
	subjects, err := h.uc.GetAllSubjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения предметов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subjects})
}

// GetSubjectByID godoc
// GET /subjects/:id
func (h *SubjectHandler) GetSubjectByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id предмета"})
		return
	}

	subject, err := h.uc.GetSubjectByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrSubjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения предмета"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subject})
}
