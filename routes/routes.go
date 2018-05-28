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
	scores := SU.GetScores()

	db, _ := bolt.Open("users.db", 0600, nil)
	password := make([]byte, 32)
	var (
		averages  []float32
		usernames []string
	)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("userauth"))
		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			averages = append(averages, SU.Average(SU.FindLast(scores, string(k), 0, 10)))
			usernames = append(usernames, string(k))
		}
	})

	db.Close()

	TP.LoginForm(usernames, averages, buffer)
	ctx.Write(buffer.Bytes())
	buffer.Reset()
}

func LoginUser(ctx iris.Context) {
	db, _ := bolt.Open("users.db", 0600, nil)
	password := make([]byte, 32)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("userauth"))
		copy(password, b.Get([]byte(ctx.PostValue("username"))))
		return nil
	})

	db.Close()

	rawInputHash := sha256.New()
	rawInputHash.Write([]byte(ctx.PostValue("password")))

	if password == nil || !bytes.Equal(rawInputHash.Sum(nil)[:], password) {
		ctx.Redirect("/", 302)
		return
	}

	session := sess.Start(ctx)
	session.Set("username", ctx.PostValue("username"))
	session.Set("auth", true)

	ctx.Redirect("/view", 302)
}

func View(ctx iris.Context) {
	session := sess.Start(ctx)
	if !session.GetBooleanDefault("auth", false) {
		ctx.Redirect("/", 302)
		return
	}

	buffer := new(bytes.Buffer)

	scores := SU.GetScores()
	scores = SU.FindLast(scores, session.GetString("username"), 0, 10)
	TP.View(session.GetString("username"), scores, buffer)
	ctx.Write(buffer.Bytes())
	buffer.Reset()
}

func SubmitScore(ctx iris.Context) {
	session := sess.Start(ctx)
	if !session.GetBooleanDefault("auth", false) {
		ctx.Redirect("/", 302)
		return
	}

	score, _ := ctx.PostValueInt("score")
	discipline, _ := ctx.PostValueInt("discipline")
	SU.AddScore(score, discipline, session.GetString("username"))
	ctx.Redirect("/view", 302)
}

func MakeUser(ctx iris.Context) {

	if ctx.PostValue("password") == "" {
		ctx.Redirect("/", 302)
		return
	}

	db, err := bolt.Open("users.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("userauth"))
		rawInputHash := sha256.New()
		rawInputHash.Write([]byte(ctx.PostValue("password")))

		if b.Get([]byte(ctx.PostValue("username"))) == nil {
			err := b.Put([]byte(ctx.PostValue("username")), rawInputHash.Sum(nil)[:])
			session := sess.Start(ctx)
			session.Set("username", ctx.PostValue("username"))
			session.Set("auth", true)

			ctx.Redirect("/view", 302)
			return err
		} else {
			ctx.Redirect("/", 302)
			return nil
		}
	})

	defer db.Close()
}

func LogoutUser(ctx iris.Context) {
	sess.Destroy(ctx)
	ctx.Redirect("/", 302)
}
