package main

import (
	"github.com/docker/distribution/context"
	flags "github.com/jessevdk/go-flags"
	"go.uber.org/zap"

	"github.com/pyr-sh/keybase-notarybot/bot/api"
	"github.com/pyr-sh/keybase-notarybot/bot/database"
	"github.com/pyr-sh/keybase-notarybot/bot/keybase"
	"github.com/pyr-sh/keybase-notarybot/bot/storage"
)

/*
var (
	storageKind      = flag.String("storage_kind", "local", "where to keep the static files")
	storageParams    = flag.String("storage_params", "", )
	storageContainer = flag.String("storage_container", "", "name of the storage container to use")

	inputPath       = flag.String("i", "./template.pdf", "input path")
	outputPath      = flag.String("o", "./output.pdf", "output path")
	signaturePath   = flag.String("s", "./signature.png", "signature path")
	signatureX      = flag.Float64("x", 0, "sig x")
	signatureY      = flag.Float64("y", 0, "sig y")
	signatureWidth  = flag.Float64("w", 1, "sig w")
	signatureHeight = flag.Float64("h", 1, "sig h")
)
*/

func main() {
	var opts struct {
		Debug bool `short:"d" long:"debug" env:"DEBUG" description:"Show verbose debug information"`
		Keys  struct {
			HMAC string `long:"hmac" env:"HMAC" default:"helloworld" description:"HMAC used to sign the URLs"`
		} `env-namespace:"KEYS" namespace:"keys" group:"Secret keys"`
		HTTP struct {
			URL  string `long:"url" env:"URL" default:"http://localhost:4000" description:"Base URL of the frontend"`
			Addr string `short:"a" long:"address" env:"ADDRESS" default:":4001" description:"Address to bind the HTTP server to"`
		} `env-namespace:"HTTP" namespace:"http" group:"HTTP server"`
		Database struct {
			Driver string `long:"driver" env:"DRIVER" default:"postgres" choice:"postgres" description:"Database driver to use"`
			DSN    string `long:"dsn" env:"DSN" description:"DSN used to connect"`
		} `env-namespace:"DB" namespace:"db" group:"Database connectivity"`
		Keybase struct {
			BinaryPath string `long:"binary_path" env:"BINARY_PATH" description:"Path to the binary path"`
			HomeDir    string `long:"home_dir" env:"HOMEDIR" description:"Path to the home dir"`
			Username   string `long:"username" env:"USERNAME" description:"If provided, the bot gets provisioned using oneshot"`
			PaperKey   string `long:"paperkey" env:"PAPERKEY" description:"If provided, the bot gets provisioned using oneshot"`
			LogPath    string `long:"log_path" env:"LOG_PATH" description:"If not set, logs are printed out to stdout/stderr"`
		} `env-namespace:"KB" namespace:"kb" group:"Keybase bot settings"`
		Storage struct {
			Kind      string            `long:"kind" env:"KIND" choice:"google" choice:"local" choice:"s3" choice:"swift" default:"swift" description:"What storage driver to use"`
			Params    map[string]string `long:"param" env:"PARAMS" description:"A map of storage driver configuration"`
			Container string            `long:"container" env:"CONTAINER" description:"Name of the storage container to use"`
		} `env-namespace:"STORAGE" namespace:"storage" group:"File storage"`
	}
	if _, err := flags.Parse(&opts); err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			return
		}
		panic(err)
	}

	var (
		logger *zap.Logger
		err    error
	)
	if opts.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	storage, err := storage.New(storage.Config{
		Kind:      opts.Storage.Kind,
		Params:    opts.Storage.Params,
		Container: opts.Storage.Container,

		Log: logger,
	})
	if err != nil {
		panic(err)
	}

	database, err := database.New(database.Config{
		Context: ctx,
		Driver:  opts.Database.Driver,
		DSN:     opts.Database.DSN,

		Log: logger,
	})
	if err != nil {
		panic(err)
	}

	bot, err := keybase.New(keybase.Config{
		BinaryPath: opts.Keybase.BinaryPath,
		HomeDir:    opts.Keybase.HomeDir,
		Username:   opts.Keybase.Username,
		PaperKey:   opts.Keybase.PaperKey,
		LogPath:    opts.Keybase.LogPath,

		HTTPURL: opts.HTTP.URL,
		HMACKey: []byte(opts.Keys.HMAC),

		Context: ctx,
		Log:     logger,
	})
	if err != nil {
		panic(err)
	}

	api, err := api.New(api.Config{
		Addr:    opts.HTTP.Addr,
		Debug:   opts.Debug,
		HMACKey: []byte(opts.Keys.HMAC),

		Log:      logger,
		Storage:  storage,
		Database: database,
	})
	if err != nil {
		panic(err)
	}

	if err := bot.Start(ctx); err != nil {
		panic(err)
	}

	if err := api.Start(ctx); err != nil {
		panic(err)
	}

	/*
		(func() {
			c := creator.New()

			signatureImage, err := c.NewImageFromFile(*signaturePath)
			if err != nil {
				panic(err)
			}

			f, err := os.OpenFile(*inputPath, 0644, 0)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			pdf, err := model.NewPdfReader(f)
			if err != nil {
				panic(err)
			}

			pagesCount, err := pdf.GetNumPages()
			if err != nil {
				panic(err)
			}

			for i := 1; i <= pagesCount; i++ {
				page, err := pdf.GetPage(i)
				if err != nil {
					panic(err)
				}

				c.AddPage(page)

				// x,y,w,h are proportional
				xa := c.Context().PageWidth * *signatureX
				ya := c.Context().PageHeight * *signatureY
				wa := c.Context().PageWidth * *signatureWidth
				ha := c.Context().PageHeight * *signatureHeight

				// We want to take up all of the space provided by the (xa, ya), (xa + wa, ya + ha) rectangle

				sigRatio := signatureImage.Width() / signatureImage.Height()
				boxRatio := wa / ha

				if boxRatio >= sigRatio {
					signatureImage.ScaleToHeight(ha)

					// We want to place the image in the middle of the horizontal field.
					signatureImage.SetPos(xa+wa/2-signatureImage.Width()/2, ya)
					if err := c.Draw(signatureImage); err != nil {
						panic(err)
					}
				} else {
					signatureImage.ScaleToWidth(wa)

					// We want to place the image in the middle of the vertical field.
					signatureImage.SetPos(xa, ya+ha/2-signatureImage.Height()/2)
				}

				if err := c.Draw(signatureImage); err != nil {
					panic(err)
				}
			}

			c.SetOutlineTree(pdf.GetOutlineTree())
			c.SetForms(pdf.AcroForm)

			if err := c.WriteToFile(*outputPath); err != nil {
				panic(err)
			}

			fmt.Println("Done!")
		})()
	*/
}
