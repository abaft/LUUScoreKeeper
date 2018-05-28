package main

import (
	//"fmt"
	"github.com/kataras/iris"
	//"github.com/kataras/iris/websocket"
	RT "github.com/abaft/LUUScoreKeeper/routes"
	"github.com/boltdb/bolt"
	"log"
)

//var scores []scoreutils.Score

func main() {

	// SETUP DATABASE
	db, err := bolt.Open("users.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("userauth"))
		tx.CreateBucketIfNotExists([]byte("scores"))
		return err
	})

	db.Close()

	app := iris.New()

	app.Get("/", RT.LoginForm)
	app.Get("/view", RT.View)

	app.Post("/makeuser", RT.MakeUser)
	app.Post("/add_score", RT.SubmitScore)
	app.Post("/login", RT.LoginUser)

	app.Any("/logout", RT.LogoutUser)

	app.Run(iris.Addr(":25565"))
}
