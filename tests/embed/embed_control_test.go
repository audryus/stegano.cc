package encrypt

import (
	"os"
	"testing"

	"github.com/audryus/steganocc/config"
	"github.com/audryus/steganocc/embed/control"
	"github.com/audryus/steganocc/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestEmbedPNG(t *testing.T) {
	t.Run("embed message into png", func(t *testing.T) {
		fxtest.New(t, fx.Provide(
			logger.New,
			config.New),
			fx.Invoke(func() {
				file, err := os.Open("wallpaper.png")
				if err != nil {
					t.Fatalf("failed to open file: %v", err)
				}
				defer file.Close()
				buf, err := control.Embed("ImpGrFNZik8V/8ikw8ZWjGUNEbLYtgVqbpKycZ4DDbA2CRxAkdMF/1AhwSmTyP3br1epJNoIof/45Lqj1DEaxHNIA2dSWMEEPeCUV/URhJ6pTp6wFVnjUGE+atGDpX0dua+wpYSitFXMaKpQOFU05BVHDVblI8KzCYP3FV7XnSXpy1vdmD1kip2YqhZ6YrJ4EbzThvIids3FImv/iymUTEvh5U/WDQhd8/kMPUhCk5Z6NlL6o5hxudIUseHQhd4pZGA0he5upw==", "wallpaper.png", file)
				if err != nil {
					t.Errorf("Encrypt returned error: %v", err)
				}

				if buf == nil || buf.Len() == 0 {
					t.Fatal("Empty buffer")
				}
				f, err := os.Create("wallpaper_embed.png")
				if err != nil {
					t.Errorf("Erro creating file: %v", err)
				}
				buf.WriteTo(f)
				f.Close()
			}))
	})
}

func TestEmbedJPEG(t *testing.T) {
	t.Run("embed message into jpeg", func(t *testing.T) {
		fxtest.New(t, fx.Provide(
			logger.New,
			config.New),
			fx.Invoke(func() {
				file, err := os.Open("wallpaper.jpg")
				if err != nil {
					t.Fatalf("failed to open file: %v", err)
				}
				defer file.Close()
				buf, err := control.Embed("ImpGrFNZik8V/8ikw8ZWjGUNEbLYtgVqbpKycZ4DDbA2CRxAkdMF/1AhwSmTyP3br1epJNoIof/45Lqj1DEaxHNIA2dSWMEEPeCUV/URhJ6pTp6wFVnjUGE+atGDpX0dua+wpYSitFXMaKpQOFU05BVHDVblI8KzCYP3FV7XnSXpy1vdmD1kip2YqhZ6YrJ4EbzThvIids3FImv/iymUTEvh5U/WDQhd8/kMPUhCk5Z6NlL6o5hxudIUseHQhd4pZGA0he5upw==", "wallpaper.jpg", file)
				if err != nil {
					t.Errorf("Encrypt returned error: %v", err)
				}

				if buf == nil || buf.Len() == 0 {
					t.Fatal("Empty buffer")
				}

				f, err := os.Create("wallpaper_embed.jpg")
				if err != nil {
					t.Errorf("Erro creating file: %v", err)
				}
				buf.WriteTo(f)
				f.Close()
			}))
	})
}
