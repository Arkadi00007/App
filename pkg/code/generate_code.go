package code

import (
	"fmt"
	"math/rand"
	"time"
)

// генерирует случайный 6-значный код
func Generate() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}
