package main

import (
	//"fmt"
  "bytes"

	"github.com/kataras/iris"
	//"github.com/kataras/iris/websocket"
  "rifle/template"
  "rifle/scoreutils"
)


var scores []scoreutils.Score

func main() {
  app := iris.New()

  app.Get("/", func(ctx iris.Context) {
    buffer := new(bytes.Buffer)
    template.WebForm(scores, buffer)
    ctx.Write(buffer.Bytes())
  })

  app.Post("/", func(ctx iris.Context){
    var score scoreutils.Score
    score.Name = ctx.PostValue("uid")
    score.Score , _ = ctx.PostValueInt64("score")
    score.Discipline , _ = ctx.PostValueInt("discipline")
    scores = append(scores, score)

    buffer := new(bytes.Buffer)
    template.WebForm(scores, buffer)
    ctx.Write(buffer.Bytes())
  })


	// x2
	// http://localhost:8080
	// http://localhost:8080
	// write something, press submit, see the result.
	app.Run(iris.Addr(":80"))
}
