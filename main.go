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

	// setup DAOs
	teacherDao := dao.NewTeacherDao(db)
	studentDao := dao.NewStudentDao(db)
	courseDao := dao.NewStudentCourseDao(db)
	subjectDao := dao.NewSubjectDao(db)
	orderDao := dao.NewRechargeOrderDao(db)
	recordDao := dao.NewRecordDao(db)
	// setup repositories

	courseRepository := repository.NewCourseRepository(courseDao)
	teacherRepository := repository.NewTeacherRepository(teacherDao)
	studentRepository := repository.NewStudentRepository(studentDao)
	subjectRepository := repository.NewSubjectRepository(subjectDao)
	recordRepository := repository.NewRecordRepository(recordDao)
	orderRepository := repository.NewOrderRepository(orderDao)

	// Setup teacher manager

	teacherManager := service.NewTeacherManager(teacherRepository, courseRepository)

	// Setup student manager

	studentManager := service.NewStudentManager(studentRepository, courseRepository)

	// Set up subject manager

	subjectManager := service.NewSubjectManager(subjectRepository, courseRepository)

	// Set up course manager

	courseManager := service.NewCourseManager(courseRepository, studentRepository)

	// Setup order manager

	orderManager := service.NewOrderManager(orderRepository, studentRepository)

	// Setup record manager

	recordManager := service.NewRecordManager(recordRepository, courseRepository, subjectRepository, studentRepository)

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
			subjectManager.Ctx = ctx
			courseManager.Ctx = ctx
			orderManager.Ctx = ctx
			recordManager.Ctx = ctx
			dashboardManager.Ctx = ctx

			// Register routes
			studentManager.RegisterRoute(dispatcher)
			teacherManager.RegisterRoute(dispatcher)
			subjectManager.RegisterRoute(dispatcher)
			courseManager.RegisterRoute(dispatcher)
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
