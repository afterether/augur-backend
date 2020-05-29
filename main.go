package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)
var (
	augur_srv *AugurServer
	// thes variables should be removed on the next code reorg task
	market_order_id int64 = 0
	fill_order_id int64 = 0
)
func main() {
	port_plain := os.Getenv("AUGUR_HTTP_PORT")
	port_secure := os.Getenv("AUGUR_HTTPS_PORT")

	if port_plain == "" {
		port_plain = "9090"
		log.Printf("Defaulting plain HTTP to port %s", port_plain)
	}
	if port_secure== "" {
		port_secure= "9443"
		log.Printf("Defaulting secure protocol to port %s", port_secure)
	}

	augur_srv = create_augur_server()

	r := gin.New()
	r.LoadHTMLGlob("html/templates/*")
	r.Use(gin.Logger())

	// Main HTML
	r.GET("/", main_page)
	r.GET("/index.html", main_page)
	r.GET("/index.htm", main_page)
	r.GET("/INDEX.HTM", main_page)

	r.GET("/markets.html",markets)
	r.GET("/statistics.html",statistics)
	r.GET("/categories.html",categories)
	r.GET("/explorer.html",explorer)
	r.GET("/market/:market",  market_trades)
	r.GET("/mdepth/:market/:outcome", market_depth)
	r.GET("/mphist/:market/:outcome", market_price_history)
	r.GET("/search", search)

	r.Static("/imgs", "./html/imgs")
	r.Static("/res", "./html/res")			// resources (static)
	r.StaticFile("/favicon.ico", "./html/res/favicon.ico")

	// Listen and serve on defined port
	log.Printf("Listening on port %s", port_plain)
	r.Run(":" + port_plain)
}
