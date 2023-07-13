package main

import (
	"fmt"
	"github.com/mattellis91/libretasks-server/db"
)

func main ()  {
	db.InitDb()	

 	if db.Db != nil {
		fmt.Println("Db connection already exists");
		db.CloseConnection()
	}
}