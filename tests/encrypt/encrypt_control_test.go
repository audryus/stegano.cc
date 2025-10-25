package encrypt

import (
	"testing"

	"github.com/audryus/steganocc/config"
	"github.com/audryus/steganocc/encrypt/control"
	"github.com/audryus/steganocc/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestEncrypt(t *testing.T) {
	t.Run("encrypt message", func(t *testing.T) {
		fxtest.New(t, fx.Provide(
			logger.New,
			config.New),
			fx.Invoke(func() {
				encrypted, err := control.Encrypt("password123", "staff update garment toy tornado bamboo wrap sign skin spike oak toddler harvest void slow champion push spread exercise giraffe citizen check reject wagon")
				if err != nil {
					t.Errorf("Encrypt returned error: %v", err)
				}

				if len(encrypted) == 0 {
					t.Fatal("Empty encrypted string")
				}
			}))
	})
}
