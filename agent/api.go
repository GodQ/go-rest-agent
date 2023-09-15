package agent

import (
	"github.com/gin-gonic/gin"
)

func SetApi(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")
	controller := NewController()

	apiV1.GET("/tasks", controller.ListTasks)
	apiV1.GET("/tasks/:task_id", controller.GetTask)
	apiV1.POST("/tasks", controller.CreateTask)
	apiV1.GET("/file", controller.GetFile)
	apiV1.POST("/file", controller.UploadFile)
}
