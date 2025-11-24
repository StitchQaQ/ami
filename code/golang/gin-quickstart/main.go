package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StructA struct {
	FieldA string `form:"field_a"`
}

type StructB struct {
	NestedStruct StructA
	FieldB       string `form:"field_b"`
}

type StructC struct {
	NestedStructPointer *StructA
	FieldC              string `form:"field_c"`
}

type StructD struct {
	NestedAnonyStruct struct {
		FieldX string `form:"field_x"`
		FieldY string `form:"field_y"`
	}
	FieldD string `form:"field_d"`
}

func GetDataB(c *gin.Context) {
	var b StructB
	c.Bind(&b)
	c.JSON(http.StatusOK, gin.H{
		"a": b.NestedStruct,
		"b": gin.H{"FieldB": b.FieldB},
	})

}

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang":    "GO语言",
			"tag":     "<br>",
			"comment": "GO语言真的很不错！",
		}

		c.AsciiJSON(http.StatusOK, data)
	})

	// 对比：使用普通 JSON（会保留中文字符）
	router.GET("/normalJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang":    "GO语言",
			"tag":     "<br>",
			"comment": "GO语言真的很不错！",
		}

		c.JSON(http.StatusOK, data)
	})

	// 对比：使用 PureJSON（不转义 HTML 标签）
	router.GET("/pureJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang":    "GO语言",
			"tag":     "<br>",
			"comment": "GO语言真的很不错！",
		}

		c.PureJSON(http.StatusOK, data)
	})

	router.GET("/getb", GetDataB)

	router.Run(":8080")
}
