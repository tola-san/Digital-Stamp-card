package main

import (
    "net/http"
     "github.com/gin-gonic/gin"
)


func main() {

    //  define router to create init web 
    router := gin.Default()

    router.GET("/health",func(c*gin.Context){

        // return as JSON 
        c.JSON(http.StatusOK, gin.H{
			"message": "Digital Stamp API running",
		})
    })

    router.Run(":8080")
}