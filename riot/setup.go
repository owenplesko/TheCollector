package riot

import (
	"github.com/go-playground/validator/v10"
)

var key string
var validate *validator.Validate
var matchesAfter int64

func Setup(riotKey string, after int64) {
	key = riotKey
	matchesAfter = after
	validate = validator.New()
}
