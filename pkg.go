package mnemosyne

import "time"

type Initializer[V any] func() (V, time.Duration)
