package routes

import (
	"bytes"
	"crypto/sha256"
	SU "github.com/abaft/LUUScoreKeeper/scoreutils"
	TP "github.com/abaft/LUUScoreKeeper/template"
	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"log"
)

var (
	cookieNameForSessionID = "LUUScoreKeeper"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)

func LoginForm(ctx iris.Context) {
	buffer := new(bytes.Buffer)
	TP.LoginForm(SU.TopScores, buffer)
	ctx.Write(buffer.Bytes())
	buffer.Reset()
}

func LoginUser(ctx iris.Context) {
	db, _ := bolt.Open("users.db", 0600, nil)
	var password []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("userauth"))
		password = b.Get([]byte(ctx.PostValue("username")))
		return nil
	})

	db.Close()

	rawInputHash := sha256.New()
	rawInputHash.Write([]byte(ctx.PostValue("password")))

	log.Println(password)
	log.Println(rawInputHash.Sum(nil))

	//if password == nil {
	//	ctx.Redirect("/", 302)
	//	return
	//}

	session := sess.Start(ctx)
	session.Set("username", ctx.PostValue("username"))
	session.Set("auth", true)

	ctx.Redirect("/view", 302)
}

func View(ctx iris.Context) {
	session := sess.Start(ctx)
	if !session.GetBooleanDefault("auth", false) {
		ctx.Redirect("/view", 302)
		return
	}

	buffer := new(bytes.Buffer)

	TP.View(session.GetString("username"), SU.TopScores, buffer)
	ctx.Write(buffer.Bytes())
	buffer.Reset()
}

func MakeUser(ctx iris.Context) {
	db, err := bolt.Open("users.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("userauth"))
		rawInputHash := sha256.New()
		rawInputHash.Write([]byte(ctx.PostValue("password")))

		err := b.Put([]byte(ctx.PostValue("username")), []byte(rawInputHash.Sum(nil)[:32]))
		return err
	})

	session := sess.Start(ctx)
	session.Set("username", ctx.PostValue("username"))
	session.Set("auth", true)

	ctx.Redirect("/view", 302)

	defer db.Close()
}

func LogoutUser(ctx iris.Context) {
	sess.Destroy(ctx)
	ctx.Redirect("/", 302)
}
