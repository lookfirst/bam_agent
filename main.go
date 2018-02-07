package main

import (
	"flag"
	"fmt"
	"os"
	//"log"
	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/blockassets/bam_agent/controller"
	"net/http"
	"log"
)

var (
	// Makefile build
	version = ""
)

func main() {
	port := flag.String("port", "1111", "The address to listen on.")
	flag.Parse()

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())

	e.Use(middleware.Recover())

	e.GET("/favicon.ico", echo.WrapHandler(http.FileServer(rice.MustFindBox("static").HTTPBox())))

	controller.Init(e)

	log.Printf("bam_agent %s %s", os.Args[0], version)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
