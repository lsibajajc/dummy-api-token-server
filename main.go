package main

import (
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var token string = "23gq3DtjdWxRnTsYsyLyuX4a007c73Md"

const listenPort = "LISTEN_PORT"
const defaultPort = "8080"
const authHeader = "authorization"
const bearerRe = `(?i)^bearer (.+)$`

func main() {
	r := gin.Default()

	r.GET("/token", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"token": token})
	})

	r.GET("/auth", authorized)

	r.Run(getListenPort())
}

func getListenPort() string {
	port := os.Getenv(listenPort)
	if len(port) == 0 {
		port = defaultPort
	}
	intPort, err := strconv.Atoi(port)
	if err != nil || intPort < 1 || 65535 < intPort {
		port = defaultPort
	}

	return ":" + port
}

func authorized(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get(authHeader)
	tokenRe := regexp.MustCompile(bearerRe)
	if len(authHeader) == 0 {
		authHeaderMissing(ctx)
	} else {
		matches := tokenRe.FindAllStringSubmatch(authHeader, -1)
		if len(matches) == 0 || !(token == matches[0][1]) {
			tokenMissmatch(ctx)
		} else {
			statusOK(ctx)
		}
	}
}

func authHeaderMissing(ctx *gin.Context) {
	ctx.Writer.Header().Set("WWW-Authenticate", "Bearer realm=\"token_required\"")
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"authorized": false,
		"error":      "missing Header: " + authHeader,
	})
}

func tokenMissmatch(context *gin.Context) {
	context.Writer.Header().Set("WWW-Authenticate", "Bearer realm=\"token_required\" error=\"invalid_token\"")
	context.JSON(http.StatusUnauthorized, gin.H{
		"authorized": false,
		"error":      "token mismatch",
	})
}

func statusOK(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"authorized": true,
	})
}
