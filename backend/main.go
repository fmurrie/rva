package main

import (
	"fmt"
	"log"
	"rva/application"
	"rva/helper"
	"golang.org/x/sync/errgroup"
)

func main() {
    var g errgroup.Group

    g.Go(func() error {
        return application.GetApplication("app1","./configuration").Run()
    })
    g.Go(func() error {
        return application.GetApplication("app2","./configuration").Run()
    })
    if err := g.Wait(); err != nil {
        log.Fatal(err)
    }
}

func prueba(){
    fmt.Println(helper.GetSecurityHelper().Decrypt("ca3f0a5f66881858bc8ec576c601d8330f287c8cafc44d381315f46971c8b68a47eeb65dbf90b7f983e479d47eb395d70e24dda9cb8607a796490872832bc99b4bb1f6977bc25e2a6c8d1c8099ccdc4e237e31a1645c6f711322ed4ff8d7e31b71ae798c66508bd44194f8c25c4807f6e5c0594788f5"))
}