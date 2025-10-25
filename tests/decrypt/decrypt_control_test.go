package encrypt

import (
	"os"
	"testing"

	"github.com/audryus/steganocc/config"
	"github.com/audryus/steganocc/decrypt/control"
	"github.com/audryus/steganocc/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestPNGDecrypt(t *testing.T) {
	t.Run("decrypt png wallpaper", func(t *testing.T) {
		fxtest.New(t, fx.Provide(
			logger.New,
			config.New),
			fx.Invoke(func() {
				t.Log("open file")
				file, err := os.Open("wallpaper_embed.png")
				if err != nil {
					t.Fatalf("failed to open file: %v", err)
				}
				defer file.Close()
				message, fileText, err := control.Decrypt("password123", "ImpGrFNZik8V/8ikw8ZWjGUNEbLYtgVqbpKycZ4DDbA2CRxAkdMF/1AhwSmTyP3br1epJNoIof/45Lqj1DEaxHNIA2dSWMEEPeCUV/URhJ6pTp6wFVnjUGE+atGDpX0dua+wpYSitFXMaKpQOFU05BVHDVblI8KzCYP3FV7XnSXpy1vdmD1kip2YqhZ6YrJ4EbzThvIids3FImv/iymUTEvh5U/WDQhd8/kMPUhCk5Z6NlL6o5hxudIUseHQhd4pZGA0he5upw==", file)
				if err != nil {
					t.Errorf("Encrypt returned error: %v", err)
				}
				expected := "staff update garment toy tornado bamboo wrap sign skin spike oak toddler harvest void slow champion push spread exercise giraffe citizen check reject wagon"
				if message != expected {
					t.Fatalf("unexpected message: got=%q want=%q", message, expected)
				}

				if fileText != expected {
					t.Fatalf("unexpected file text: got=%q want=%q", fileText, expected)
				}
			}))
	})
}

func TestJPEGDecrypt(t *testing.T) {
	t.Run("decrypt former jpeg wallpaper", func(t *testing.T) {
		fxtest.New(t, fx.Provide(
			logger.New,
			config.New),
			fx.Invoke(func() {
				t.Log("open file")
				file, err := os.Open("wallpaper_embed.png")
				if err != nil {
					t.Fatalf("failed to open file: %v", err)
				}
				defer file.Close()
				message, fileText, err := control.Decrypt("password123", "ImpGrFNZik8V/8ikw8ZWjGUNEbLYtgVqbpKycZ4DDbA2CRxAkdMF/1AhwSmTyP3br1epJNoIof/45Lqj1DEaxHNIA2dSWMEEPeCUV/URhJ6pTp6wFVnjUGE+atGDpX0dua+wpYSitFXMaKpQOFU05BVHDVblI8KzCYP3FV7XnSXpy1vdmD1kip2YqhZ6YrJ4EbzThvIids3FImv/iymUTEvh5U/WDQhd8/kMPUhCk5Z6NlL6o5hxudIUseHQhd4pZGA0he5upw==", file)
				if err != nil {
					t.Errorf("Encrypt returned error: %v", err)
				}
				expected := "staff update garment toy tornado bamboo wrap sign skin spike oak toddler harvest void slow champion push spread exercise giraffe citizen check reject wagon"
				if message != expected {
					t.Fatalf("unexpected message: got=%q want=%q", message, expected)
				}

				if fileText != expected {
					t.Fatalf("unexpected file text: got=%q want=%q", fileText, expected)
				}
			}))
	})
}
