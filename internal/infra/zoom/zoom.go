package zoom

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
)

const linkIDLength = 15

type Zoom struct{}

func NewZoom() *Zoom {
	return &Zoom{}
}

func (z *Zoom) Create(_ context.Context) (string, error) {
	linkID, err := z.generateRandomString(linkIDLength)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://fakelink.ru/%s", linkID), nil
}

func (z *Zoom) Cancel(_ context.Context, _ string) error {
	return nil
}

func (z *Zoom) generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, length)
	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		bytes[i] = charset[num.Int64()]
	}

	return string(bytes), nil
}
