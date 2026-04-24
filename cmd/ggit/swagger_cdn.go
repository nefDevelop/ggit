//go:build swagger_cdn

package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/nef734/ggit/docs"
)

func init() {
	registerSwagger = func(r *gin.Engine) {
		r.GET("/swagger", func(c *gin.Context) {
			c.Redirect(301, "/swagger/index.html")
		})

		r.GET("/swagger/doc.json", func(c *gin.Context) {
			c.File("docs/swagger.json")
		})

		r.GET("/swagger/*any", func(c *gin.Context) {
			c.Header("Content-Type", "text/html")
			html := `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>Swagger UI (CDN)</title>
				<link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css" >
				<style>html { box-sizing: border-box; overflow-y: scroll; } *, *:before, *:after { box-sizing: inherit; } body { margin:0; background: #fafafa; }</style>
			</head>
			<body>
				<div id="swagger-ui"></div>
				<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"> </script>
				<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-standalone-preset.js"> </script>
				<script>
				window.onload = function() {
					window.ui = SwaggerUIBundle({
						url: "/swagger/doc.json",
						dom_id: '#swagger-ui',
						deepLinking: true,
						presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
						plugins: [SwaggerUIBundle.plugins.DownloadUrl],
						layout: "StandaloneLayout"
					});
				};
				</script>
			</body>
			</html>`
			c.String(200, html)
		})
	}
}
