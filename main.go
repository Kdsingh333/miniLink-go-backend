package main

import (
	"log"

	"github.com/Kdsingh333/miniLink-go-backend/database"
	"github.com/Kdsingh333/miniLink-go-backend/routes"
)

func init(){
   database.Setup();
}
func main() {
	engine := routes.Routers();
	err := engine.Run(":8080");
	if err != nil{
		log.Fatal(err);
	}
}