package routes

import (
    "gin-mongo-api/controllers"
    "github.com/gin-gonic/gin"
)

func CommentRoute(router *gin.Engine) {
    router.POST("/comments", controllers.CreateComment())
    router.GET("/comments/:commentId", controllers.GetAComment())
    router.PUT("/comments/:commentId", controllers.EditAComment())
    router.DELETE("/comments/:commentId", controllers.DeleteAComment())
    router.GET("/comments", controllers.GetAllComments())
}
