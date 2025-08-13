package fiberapp

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/teoit/gosctx"
)

const (
	defaultPort = 3000
	defaultMode = "debug"
)

type FiberComponent interface {
	GetPort() int
	GetApp() *fiber.App
	SetEngineConfig(engine *html.Engine)
}

type Config struct {
	port      int
	fiberMode string
}

type fiberApp struct {
	*Config
	name   string
	id     string
	logger gosctx.Logger
	app    *fiber.App
	engine *html.Engine
}

func NewFiber(id string) *fiberApp {
	return &fiberApp{
		Config: new(Config),
		id:     id,
	}
}

func (fs *fiberApp) ID() string {
	return fs.id
}

func (fs *fiberApp) Activate(sv gosctx.ServiceContext) error {
	fs.logger = sv.Logger(fs.id)
	fs.name = sv.GetName()

	fs.app = fiber.New(fiber.Config{
		CaseSensitive: true, // Set to false if you want case-insensitive routes
		BodyLimit:     100 * 1024 * 1024,
		StrictRouting: true,
	})
	return nil
}

func (fs *fiberApp) SetConfig(config *fiber.Config) {
	fs.app = fiber.New(*config)
}

func (fs *fiberApp) SetEngineConfig(engine *html.Engine) {
	fs.engine = engine
	cfg := fiber.Config{
		Views:       engine,
		BodyLimit:   100 * 1024 * 1024,
		ViewsLayout: "layouts/master",
	}
	fs.app = fiber.New(cfg)
}

func (fs *fiberApp) Stop() error {
	return nil
}

func (fs *fiberApp) InitFlags() {
	flag.IntVar(&fs.Config.port, "fiber-port", defaultPort, "Fiber server port. Default 3000")
	flag.StringVar(&fs.Config.fiberMode, "fiber-mode", defaultMode, "Fiber mode (debug | release). Default debug")
}

func (fs *fiberApp) GetPort() int {
	return fs.port
}

func (fs *fiberApp) GetApp() *fiber.App {
	return fs.app
}
