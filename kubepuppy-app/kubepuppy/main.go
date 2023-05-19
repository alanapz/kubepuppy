package main

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"kubepuppy/app"
	openapi "kubepuppy/web"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	os.Setenv("KUBEPUPPY_KUBECONFIG", "C:\\Users\\alan\\.kube\\config")
	os.Setenv("KUBEPUPPY_SERVER_CA_FILE", "C:\\Users\\alan\\.kube\\tls.crt")
	os.Setenv("KUBEPUPPY_CLIENT_CA_FILE", "C:\\Users\\alan\\.kube\\client-ca.crt")

	cluster, err := app.InitialiseCluster(context.TODO())
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		openapi.SetClusterIntoRequest(c, cluster)
	})

	openapi.AddRoutesToRouter(router)

	router.StaticFile("/", "./assets/index.html")
	router.StaticFile("/assets/gomaster.css", "./assets/gomaster.css")

	// See: https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-without-context/server.go#L25
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Println("Preparing to start server ...")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
