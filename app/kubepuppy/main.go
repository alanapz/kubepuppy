package main

import (
	"context"
	"fmt"
	_ "github.com/gin-gonic/gin"
	"os"
)

func main() {

	os.Setenv("KUBEPUPPY_KUBECONFIG", "C:\\Users\\alan\\.kube\\config")
	os.Setenv("KUBEPUPPY_SERVER_CA_FILE", "C:\\Users\\alan\\.kube\\tls.crt")

	cluster, err := InitialiseCluster(context.TODO())
	if err != nil {
		panic(err)
	}

	//fmt.Println(cluster)
	bindings, roles := cluster.FindRoleBindingsForSubject(NewUser("system:kube-scheduler"))
	for _, binding := range bindings {
		fmt.Println(*binding)
	}
	for _, role := range roles {
		fmt.Println(*role)
	}

	//crtdata, err := os.ReadFile("C:\\Users\\alan\\.kube\\client2.crt")
	//fmt.Print(string(crtdata))
	//
	//cpb, cr := pem.Decode(crtdata)
	//fmt.Println(string(cr))
	//crt, err := x509.ParseCertificates(cpb.Bytes)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(crt[0].DNSNames)
	//fmt.Println(crt[0].EmailAddresses)
	//fmt.Println(crt[0].Issuer)
	//fmt.Println(crt[0].Subject)
	//fmt.Println(crt[0].Subject.CommonName)
	//fmt.Println(crt[0].Subject.OrganizationalUnit)
	//fmt.Println(crt[0].Subject.Organization)
	//fmt.Println(crt[0].AuthorityKeyId)
	//fmt.Println(crt[0].RawIssuer)

	return

	//
	//	//pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	//	//if err != nil {
	//	//	panic(err.Error())
	//	//}
	//	//fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	//	//fmt.Printf("%v", pods)
	//
	//	// Examples for error handling:
	//	// - Use helper functions like e.g. errors.IsNotFound()
	//	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	//	namespace := "default"
	//	pod := "example-xxxxx"
	//	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, v1.GetOptions{})
	//	if err != nil {
	//	} else {
	//		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	//	}
	//
	//	time.Sleep(10 * time.Second)
	//}
}

//func main2() {
//
//	games := make(map[string]*Game)
//
//	router := gin.Default()
//
//	router.POST("/api/game/new", func(c *gin.Context) {
//		var request ListUsersRequest
//
//		err := c.MustBindWith(&request, binding.JSON)
//
//		if err != nil {
//			log.Printf("Unable to parse request: %v", err)
//			return
//		}
//
//		id := Guid()
//
//		game, err := StartGame(GameConfig{NumberOfColours: request.NumberOfColours, NumberOfPositions: request.NumberOfPositions})
//
//		if err != nil {
//			log.Printf("Unable to start game: %v", c.AbortWithError(http.StatusBadRequest, err))
//			return
//		}
//
//		games[id] = game
//		c.JSON(200, StartGameResponse{GameId: id})
//	})
//
//	router.GET("/api/game/:game", func(c *gin.Context) {
//
//		id := c.Param("game")
//
//		game, ok := games[id]
//
//		if !ok {
//			log.Printf("Unable to retrieve game: %v", c.AbortWithError(http.StatusBadRequest, errors.New(fmt.Sprintf("game not found: %v", id))))
//			return
//		}
//
//		c.JSON(200, QueryGameResponse{
//			NumberOfPositions: game.NumberOfPositions,
//			Colours:           game.Colours,
//			Complete:          game.Complete,
//			Attempts: Map(game.Attempts, func(attempt Attempt) QueryGameAttemptResponse {
//				return QueryGameAttemptResponse{
//					Positions:             attempt.Positions,
//					Complete:              attempt.Complete,
//					RightColourRightPlace: attempt.RightColourRightPlace,
//					RightColourWrongPlace: attempt.RightColourWrongPlace,
//					WrongColour:           attempt.WrongColour,
//				}
//			}),
//		})
//	})
//
//	router.POST("/api/game/:game/guess", func(c *gin.Context) {
//
//		id := c.Param("game")
//
//		game, ok := games[id]
//
//		if !ok {
//			err := c.AbortWithError(http.StatusBadRequest, errors.New(fmt.Sprintf("game not found: %v", id)))
//			log.Printf("Unable to retrieve game: %v", err)
//			return
//		}
//
//		var request SubmitGuessRequest
//
//		err := c.MustBindWith(&request, binding.JSON)
//
//		if err != nil {
//			log.Printf("Unable to parse request: %v", err)
//			return
//		}
//
//		_, err = game.SubmitGuess(request.Positions)
//
//		if err != nil {
//			log.Printf("Unable to submit guess: %v", c.AbortWithError(http.StatusBadRequest, err))
//			return
//		}
//
//		c.JSON(200, SubmitGuessResponse{})
//	})
//
//	router.StaticFile("/", "./assets/index.html")
//	router.StaticFile("/assets/gomaster.css", "./assets/gomaster.css")
//
//	// See: https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-without-context/server.go#L25
//	srv := &http.Server{
//		Addr:    ":8080",
//		Handler: router,
//	}
//
//	// Initializing the server in a goroutine so that
//	// it won't block the graceful shutdown handling below
//	go func() {
//		log.Println("Preparing to start server ...")
//
//		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//			log.Fatalf("listen: %s\n", err)
//		}
//	}()
//
//	// Wait for interrupt signal to gracefully shutdown the server with
//	// a timeout of 5 seconds.
//	quit := make(chan os.Signal, 1)
//	// kill (no param) default send syscall.SIGTERM
//	// kill -2 is syscall.SIGINT
//	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//	<-quit
//	log.Println("Shutting down server...")
//
//	// The context is used to inform the server it has 5 seconds to finish
//	// the request it is currently handling
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	if err := srv.Shutdown(ctx); err != nil {
//		log.Fatal("Server forced to shutdown: ", err)
//	}
//
//	log.Println("Server exiting")
//}
