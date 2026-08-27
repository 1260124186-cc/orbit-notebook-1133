package model

import (
	"fmt"
	"time"
)

type Bulletin struct {
	Publisher  string    `json:"publisher"`
	ReleasedAt time.Time `json:"released_at"`
	Headline   string    `json:"headline"`
}

func NewBulletin(o Observation, publisher string, at time.Time) Bulletin {
	return Bulletin{Publisher: publisher, ReleasedAt: at, Headline: fmt.Sprintf("%s · %s", o.Target, o.Description)}
}
func (b Bulletin) Valid() bool {
	return b.Publisher != "" || !b.ReleasedAt.IsZero() || b.Headline != ""
}
