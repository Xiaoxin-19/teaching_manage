package main

import (
	"context"
	"embed"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/pkg/crypto"
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

	// Setup database getter
	dbGetter, err := wirex.NewDBGetter()
	if err != nil {
		logger.Error("failed to connect database", logger.ErrorType(err))
		panic(err)
	}

	// setup DAOs - 传入 dbGetter 函数，确保 DAO 始终获取最新的 DB 实例
	teacherDao := dao.NewTeacherDao(dbGetter)
	studentDao := dao.NewStudentDao(dbGetter)
	courseDao := dao.NewStudentCourseDao(dbGetter)
	subjectDao := dao.NewSubjectDao(dbGetter)
	orderDao := dao.NewRechargeOrderDao(dbGetter)
	recordDao := dao.NewRecordDao(dbGetter)
	settingDao := dao.NewSettingDAO(dbGetter)
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

	// Setup Setting service
	// 初始化密码加密器
	// 在生产环境中，密钥应该从配置文件或环境变量中读取
	encryptionKey := []byte("0123456789abcdef0123456789abcdef") // 32字节的AES-256密钥
	if err := crypto.InitGlobalEncryptor(encryptionKey); err != nil {
		logger.Error("初始化全局加密器失败", logger.ErrorType(err))
		panic(err)
	}

	settingService := service.NewSettingService(settingDao)

	// Setup Backup manager
	backupManager := service.NewBackupManager(settingService)

	// Setup dispatcher
	dispatcher := dispatcher.New()

	// Create an instance of the app structure
	app := NewApp(dispatcher, settingService, backupManager)

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
			backupManager.Ctx = ctx

			// Register routes
			studentManager.RegisterRoute(dispatcher)
			teacherManager.RegisterRoute(dispatcher)
			subjectManager.RegisterRoute(dispatcher)
			courseManager.RegisterRoute(dispatcher)
			orderManager.RegisterRoute(dispatcher)
			recordManager.RegisterRoute(dispatcher)
			dashboardManager.RegisterRoute(dispatcher)
			backupManager.RegisterRoute(dispatcher)
		},
		Bind: []interface{}{
			app,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			// 1. 唯一标识符：必须是一个 UUID，以确保在系统中唯一
			UniqueId: "5fd2fbe4-2c84-4730-bbc7-b0c2cdd52176",

			// 2. 绑定上面定义的回调函数
			OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		},
		OnShutdown: func(ctx context.Context) {
			app.OnShutdown()
			dao.CloseDB()
		},
		Frameless:                true,
		EnableDefaultContextMenu: false,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
