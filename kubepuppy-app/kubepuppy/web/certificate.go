package openapi

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidateCertificate(c *gin.Context) {
	subjectId := c.GetString("subjectId")

	if len(subjectId) == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("subjectId required"))
		return
	}

	cluster := GetClusterFromRequest(c)

	cluster.LockForRead(func() {

		subject, ok := cluster.Subjects[subjectId]

		if !ok {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("subject not found: %v", subjectId))
			return
		}

		c.JSON(http.StatusOK, SubjectDetails{
			//SubjectId: subject.SubjectId,
			Kind:      subject.Kind,
			Name:      subject.Name,
			Namespace: subject.Namespace,
		})
	})
}
