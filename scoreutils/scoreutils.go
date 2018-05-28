package scoreutils

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type Score struct {
	Score      int
	Discipline int
	Name       string
	Date       time.Time
}

func AddScore(score int, discipline int, username string) {
	s := Score{
		Score:      score,
		Name:       username,
		Discipline: discipline,
		Date:       time.Now(),
	}

	db, err := bolt.Open("users.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("scores"))
		id, _ := b.NextSequence()
		rawScore, _ := json.Marshal(s)

		err := b.Put([]byte(string(id)), rawScore)
		return err
	})
}

func GetScores() []Score {
	db, err := bolt.Open("users.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var scores []Score
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("scores"))
		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var score Score
			json.Unmarshal(v, &score)

			scores = append(scores, score)
		}

		return nil
	})
	return scores
}

func FindLast(scores []Score, username string, discipline int, count uint) []Score {
	tmp := make([]Score, count)

	for _, s := range scores {
		if s.Name == username && s.Discipline == discipline {
			if s.Date.After(tmp[0].Date) {
				tmp[0] = s
			}
		}
	}

	for i := uint(1); i < count; i++ {
		for _, s := range scores {
			if s.Name == username && s.Discipline == discipline {
				if tmp[i-1].Date.Before(tmp[i].Date) && s.Date.After(tmp[i].Date) {
					tmp[i] = s
				}
			}
		}
	}
	return tmp
}

func Average(scores []Score, username string, discipline int) Score {
	sum := 0
	for _, s := range scores {
		sum += s.Score
	}
	rtn := scores[0]
	rtn.Score = sum / len(scores)
	return rtn
}

func (s Score) GetDiscipline() string {
	switch s.Discipline {
	case 0:
		return ".22lr Prone"
	case 1:
		return ".22lr Kneeling"
	case 2:
		return ".22lr Offhand Carbine"
	}
	return "misc"
}
