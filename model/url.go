package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/rbo13/chibyurl/words"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// URL represents a url to be
// shortened.
type URL struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Alias     string             `bson:"alias" json:"alias"`
	URL       string             `bson:"url" json:"url"`
	Click     int32              `bson:"click" json:"click"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Urls []URL

// Generate generates a new url alias
// if the user does not provide the alias.
func (u *URL) Generate() string {
	seeder := int64(time.Now().UnixNano() * int64(len(words.Adjective)*len(words.Animal)*len(words.Verb)))
	rand.Seed(seeder) // create seed value first

	randomAdjective := removeSpace(strings.Title(words.Adjective[rand.Intn(len(words.Adjective))]))
	randomAnimal := removeSpace(strings.Title(words.Animal[rand.Intn(len(words.Animal))]))
	randomVerb := removeSpace(strings.Title(words.Verb[rand.Intn(len(words.Verb))]))

	return fmt.Sprintf("%s%s%s", randomAdjective, randomVerb, randomAnimal)
}

func removeSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
