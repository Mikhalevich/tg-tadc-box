package imagepayload

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type PayloadType string

const (
	PayloadTypeChest  PayloadType = "chest"
	PayloadTypeReward PayloadType = "reward"
)

type ImagePayload struct {
	Type    PayloadType
	Reward  reward.Reward
	BoxType box.Type
}

func (p ImagePayload) GOBEncode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, fmt.Errorf("gob enbode: %w", err)
	}

	return buf.Bytes(), nil
}

func GOBDecode(b []byte) (ImagePayload, error) {
	var ip ImagePayload
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&ip); err != nil {
		return ImagePayload{}, fmt.Errorf("gob decode: %w", err)
	}

	return ip, nil
}
