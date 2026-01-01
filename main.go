package main

import (
	"context"
	"embed"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/pkg/dispatcher"
	"teaching_manage/backend/pkg/logger"
	"teaching_manage/backend/repository"
	"teaching_manage/backend/service"
	"teaching_manage/backend/wirex"

	"github.com/wailsapp/wails/v2"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	// setup logger
	zaplog := wirex.InitLogger()
	logger.SetGlobalLogger(zaplog)

	// Setup database
	db, err := wirex.NewGormDB()
	if err != nil {
		logger.Error("failed to connect database", logger.ErrorType(err))
		panic(err)
	}

	// Setup teacher manager
	teacherDao := dao.NewTeacherDao(db)
	teacherRepository := repository.NewTeacherRepository(teacherDao)
	teacherManager := service.NewTeacherManager(teacherRepository)

	// Setup student manager
	studentDao := dao.NewStudentDao(db)
	studentRepository := repository.NewStudentRepository(studentDao)
	studentManager := service.NewStudentManager(studentRepository, teacherRepository)

	// Setup order manager
	orderDao := dao.NewOrderDao(db)
	orderRepository := repository.NewOrderRepository(orderDao)
	orderManager := service.NewOrderManager(orderRepository, studentRepository)

	// Setup record manager
	recordDao := dao.NewRecordDao(db)
	recordRepository := repository.NewRecordRepository(recordDao)
	recordManager := service.NewRecordManager(recordRepository, studentRepository)

	// Setup Dashboard manager
	dashboardManager := service.NewDashboardManager()

	// Setup dispatcher
	dispatcher := dispatcher.New()

	// Create an instance of the app structure
	app := NewApp(dispatcher)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "teaching_manage",
		Width:  1224,
		Height: 868,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop:     true,
			DisableWebViewDrop: true,
			CSSDropProperty:    "--wails-drop-target",
			CSSDropValue:       "drop",
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			teacherManager.Ctx = ctx
			studentManager.Ctx = ctx
			orderManager.Ctx = ctx
			recordManager.Ctx = ctx
			dashboardManager.Ctx = ctx

			// Register routes
			studentManager.RegisterRoute(dispatcher)
			teacherManager.RegisterRoute(dispatcher)
			orderManager.RegisterRoute(dispatcher)
			recordManager.RegisterRoute(dispatcher)
			dashboardManager.RegisterRoute(dispatcher)
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
