package main

import (
	"context"
	"embed"
	"log"
	"teaching_manage/dao"
	"teaching_manage/pkg/dispatcher"
	"teaching_manage/repository"
	"teaching_manage/service"
	"teaching_manage/wirex"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Setup database
	db, err := wirex.NewGormDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Setup student manager
	studentDao := dao.NewStudentDao(db)
	studentRepository := repository.NewStudentRepository(studentDao)
	studentManager := service.NewStudentManager(studentRepository)

	// Setup teacher manager
	teacherDao := dao.NewTeacherDao(db)
	teacherRepository := repository.NewTeacherRepository(teacherDao)
	teacherManager := service.NewTeacherManager(teacherRepository)

	dispatcher := dispatcher.New()

	// Create an instance of the app structure
	app := NewApp(dispatcher)

	// Create application with options
	err = wails.Run(&options.App{
		Title:    "teaching_manage",
		Width:    1024,
		Height:   768,
		LogLevel: logger.DEBUG,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			teacherManager.Ctx = ctx
			studentManager.Ctx = ctx
			// Register routes
			studentManager.RegisterRoute(dispatcher)
			teacherManager.RegisterRoute(dispatcher)
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
