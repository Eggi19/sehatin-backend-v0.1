package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/database"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/handlers"
	"github.com/tsanaativa/sehatin-backend-v0.1/middlewares"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/tsanaativa/sehatin-backend-v0.1/ws"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouterOpts struct {
	User                *handlers.UserHandler
	Auth                *handlers.AuthHandler
	Pharmacy            *handlers.PharmacyHandler
	Doctor              *handlers.DoctorHandler
	Specialist          *handlers.SpecialistHandler
	PharmacyManager     *handlers.PharmacyManagerHandler
	UserAddress         *handlers.UserAddressHandler
	PharmacyProduct     *handlers.PharmacyProductHandler
	Category            *handlers.CategoryHandler
	WebSocket           *ws.WebSocketHandler
	Product             *handlers.ProductHandler
	Consultation        *handlers.ConsultationHandler
	ProductField        *handlers.ProductFieldHandler
	Cart                *handlers.CartHandler
	Location            *handlers.LocationHandler
	StockHistory        *handlers.StockHistoryHandler
	ShppingMethod       *handlers.ShippingMethodHandler
	ResetPassword       *handlers.ResetPasswordHandler
	Order               *handlers.OrderHandler
	Admin               *handlers.AdminHandler
	MutationStatus      *handlers.MutationStatusHandler
	StockTransfer       *handlers.StockTransferHandler
	StockHistoryReport  *handlers.StockHistoryReportHandler
	SalesReport         *handlers.SalesReportHandler
	SalesReportCategory *handlers.SalesReportCategoryHandler
	MostBoughtUser      *handlers.MostBoughtUserHandler
}

func createRouter(config utils.Config, hub *ws.Hub) *gin.Engine {
	db, err := database.ConnectDB(config)
	if err != nil {
		log.Fatalf("error connecting to DB: %s", err.Error())
	}

	userRepo := repositories.NewUserRepositoryPostgres(&repositories.UserRepoOpt{Db: db})
	doctorRepo := repositories.NewDoctorRepositoryPostgres(&repositories.DoctorRepoOpt{Db: db})
	pharmacyManagerRepo := repositories.NewPharmacyManagerRepositoryPostgres(&repositories.PharmacyManagerRepoOpt{Db: db})
	adminRepo := repositories.NewAdminRepositoryPostgers(&repositories.AdminRepoOpt{Db: db})
	pharmacyRepo := repositories.NewPharmacyRepositoryPostgres(&repositories.PharmacyRepoOpts{Db: db})
	specialistsRepo := repositories.NewSpecialistRepositoryPostgres(&repositories.SpecialistRepoOpts{Db: db})
	userAddressRepo := repositories.NewUserAddressRepositoryPostgres(&repositories.UserAddressRepoOpts{Db: db})
	pharmacyProductRepo := repositories.NewPharmacyProductRepositoryPostgres(&repositories.PharmacyProductRepoOpt{Db: db})
	genderRepo := repositories.NewGenderRepositoryPostgres(&repositories.GenderRepoOpts{Db: db})
	categoryRepo := repositories.NewCategoryRepositoryPostgres(&repositories.CategoryRepoOpts{Db: db})
	shippingMethodRepo := repositories.NewShippingMethodRepositoryPostgres(&repositories.ShippingMethodRepoOpt{Db: db})
	productRepo := repositories.NewProductRepositoryPostgres(&repositories.ProductRepoOpts{Db: db})
	productCategoryRepo := repositories.NewProductCategoryRepositoryPostgres(&repositories.ProductCategoryRepoOpts{Db: db})
	consultationRepo := repositories.NewConsultationRepositoryPostgres(&repositories.ConsultationRepoOpts{Db: db})
	productFieldRepo := repositories.NewProductFieldRepositoryPostgres(&repositories.ProductCategoryRepoOpts{Db: db})
	cartRepo := repositories.NewCartRepositoryPostgres(&repositories.CartRepoOpts{Db: db})
	chatRepo := repositories.NewChatRepositoryPostgres(&repositories.ChatRepoOpts{Db: db})
	locationRepo := repositories.NewLocationRepositoryPostgres(&repositories.LocationRepoOpt{Db: db})
	pharmacyAddressesRepo := repositories.NewPharmacyAddressRepositoryPostgres(&repositories.PharmacyAddressRepoOpts{Db: db})
	stockHistoryRepo := repositories.NewStockHistoryRepositoryPostgres(&repositories.StockHistoryRepoOpts{Db: db})
	pharmacyAddressRepo := repositories.NewPharmacyAddressRepositoryPostgres(&repositories.PharmacyAddressRepoOpts{Db: db})
	userResetPasswordRepo := repositories.NewUserResetPasswordRepositoryPostgres(&repositories.UserResetPasswordRepoOpts{Db: db})
	doctorResetPasswordRepo := repositories.NewDoctorResetPasswordRepositoryPostgres(&repositories.DoctorResetPasswordRepoOpts{Db: db})
	orderRepo := repositories.NewOrderRepositoryPostgres(&repositories.OrderRepoOpts{Db: db})
	mutationStatusRepo := repositories.NewMutationStatusRepositoryPostgres(&repositories.MutationStatusRepoOpts{Db: db})
	stockTransferRepo := repositories.NewStockTransferRepositoryPostgres(&repositories.StockTransferRepoOpts{Db: db})
	stockHistoryReportRepo := repositories.NewStockHistoryReportRepositoryPostgres(&repositories.StockHistoryReportOpts{Db: db})
	salesReportRepo := repositories.NewSalesReportRepositoryPostgres(&repositories.SalesReportRepoOpts{Db: db})
	salesReportCategoryRepo := repositories.NewSalesReportCategoryRepositoryPostgres(&repositories.SalesReportCatgoryRepoOpts{Db: db})
	mostBoughtUserRepo := repositories.NewMostBoughtUserRepositoryPostgres(&repositories.MostBoughtUserRepoOpts{Db: db})

	userUsecase := usecases.NewUserUsecaseImpl(&usecases.UserUsecaseOpts{
		UserRepo:        userRepo,
		UserAddressRepo: userAddressRepo,
		GenderRepo:      genderRepo,
	})
	loginUsecase := usecases.NewLoginUsecaseImpl(&usecases.LoginUsecaseOpts{
		UserRepo:            userRepo,
		UserAddressRepo:     userAddressRepo,
		DoctorRepo:          doctorRepo,
		PharmacyManagerRepo: pharmacyManagerRepo,
		AdminRepo:           adminRepo,
		HashAlgorithm:       utils.NewBCryptHasher(),
		AuthTokenProvider:   utils.NewJwtProvider(config),
	})
	oauthUsecase := usecases.NewOAuthUsecaseImpl(&usecases.OAuthUsecaseOpts{
		UserRepo:          userRepo,
		UserAddressRepo:   userAddressRepo,
		AuthTokenProvider: utils.NewJwtProvider(config),
		GoogleSigner:      utils.NewGoogleSigner(config),
	})
	pharmacyUsecase := usecases.NewPharmacyUsecaseImpl((&usecases.PharmacyUsecaseOpts{
		PharmacyRepo:        pharmacyRepo,
		PharmacyAddressRepo: pharmacyAddressesRepo,
		ShippingMethodRepo:  shippingMethodRepo,
		PharmacyProductRepo: pharmacyProductRepo,
		StockHistoryRepo:    stockHistoryRepo,
	}))
	doctorUsecase := usecases.NewDoctorUsecaseImpl(&usecases.DoctorUsecaseOpts{
		DoctorRepo: doctorRepo,
	})
	specialistUsecase := usecases.NewSpecialistUsecaseImpl(&usecases.SpecialistUsecaseOpts{SpecialistRepo: specialistsRepo})
	pharmacyManagerUsecase := usecases.NewPharmacyManagerUsecaseImpl(&usecases.PharmacyManagerOpts{PharmacyManagerRepo: pharmacyManagerRepo})
	userAddresUsecase := usecases.NewUserAddressUsecaseImpl(&usecases.UserAddressUsecaseOpts{UserAddressRepo: userAddressRepo, UserRepo: userRepo})
	categoryUsecase := usecases.NewCategoryUsecaseImpl(&usecases.CategoryUsecaseOpts{CategoryRepo: categoryRepo})

	registerUsecase := usecases.NewRegisterUsecaseImpl(&usecases.RegisterUsecaseOpts{
		HashAlgorithm:       utils.NewBCryptHasher(),
		EmailSender:         utils.NewGoogleEmailSender(),
		Transactor:          repositories.NewTransactor(db),
		AuthTokenProvider:   utils.NewJwtProvider(config),
		UploadFile:          utils.NewCloudinaryUploadFile(),
		DoctorRepo:          doctorRepo,
		UserRepo:            userRepo,
		PharmacyManagerRepo: pharmacyManagerRepo,
	})
	verifyUsecase := usecases.NewVerifyUsecaseImpl(&usecases.VerifyUsecaseOpts{
		UserRepo:          userRepo,
		DoctorRepo:        doctorRepo,
		HashAlgorithm:     utils.NewBCryptHasher(),
		Transactor:        repositories.NewTransactor(db),
		AuthTokenProvider: utils.NewJwtProvider(config),
		EmailSender:       utils.NewGoogleEmailSender(),
	})
	refreshTokenUsecase := usecases.NewRefreshTokenImpl((&usecases.RefreshTokenOpts{
		AuthTokenProvider: utils.NewJwtProvider(config),
	}))
	shippingMethodUsecase := usecases.NewShippingMethodUsecaseImpl(&usecases.ShippingMethodOpts{
		ShippingMethodRepository:  shippingMethodRepo,
		UserAddressRepository:     userAddressRepo,
		PharmacyAddressRepository: pharmacyAddressRepo,
	})
	pharmacyProductUsecase := usecases.NewPharmacyProductUsecaseImpl(&usecases.PharmacyProductUsecaseOpts{
		PharmacyProductRepository: pharmacyProductRepo,
		PharmacyRepository:        pharmacyRepo,
		CategoryRepository:        categoryRepo,
		ShippingMethodUsecase:     shippingMethodUsecase,
		Transactor:                repositories.NewTransactor(db),
		ProductRepository:         productRepo,
		StockHistoryRepository:    stockHistoryRepo,
	})
	productUsecase := usecases.NewProductUsecaseImpl(&usecases.ProductUsecaseOpts{
		ProductRepo:         productRepo,
		CategoryRepos:       categoryRepo,
		ProductCategoryRepo: productCategoryRepo,
	})
	productCategoryUsecase := usecases.NewProductCategoryUsecaseImpl(&usecases.ProductCategoryUsecaseOpts{
		ProductCategoryRepo: productCategoryRepo,
	})
	consultationUsecase := usecases.NewConsultationUsecaseImpl(&usecases.ConsultationUsecaseOpts{
		ConsultationRepo: consultationRepo,
		DoctorRepo:       doctorRepo,
		ChatRepo:         chatRepo,
		ProductRepo:      productRepo,
		PharmacyRepo:     pharmacyRepo,
		UserAddressRepo:  userAddressRepo,
		CartRepo:         cartRepo,
		UploadFile:       utils.NewCloudinaryUploadFile(),
	})
	productFieldUsecase := usecases.NewProductFieldUsecaseImpl(&usecases.ProductFieldUsecaseOpts{ProductFieldRepo: productFieldRepo})
	cartUsecase := usecases.NewCartUsecaseImpl(&usecases.CartUsecaseOpts{
		CartRepo:                  cartRepo,
		PharmacyProductRepository: pharmacyProductRepo,
		ShippingMethodUsecase:     shippingMethodUsecase,
	})
	locationUsecase := usecases.NewLocationUsecaseImpl(&usecases.LocationUsecaseOpts{LocationRepo: locationRepo})
	stockHistoryUsecase := usecases.NewStockHistoryUsecaseImpl(&usecases.StockHistoryUsecaseOpts{StockHistoryRepo: stockHistoryRepo})
	userResetPasswordUsecase := usecases.NewUserResetPasswordUsecaseImpl(&usecases.UserResetPasswordUsecaseOpts{
		UserResetPasswordRepo: userResetPasswordRepo,
		UserRepo:              userRepo,
		HashAlgorithm:         utils.NewBCryptHasher(),
		EmailSender:           utils.NewGoogleEmailSender(),
		AuthTokenProvider:     utils.NewJwtProvider(config),
		Transactor:            repositories.NewTransactor(db),
	})
	doctorResetPasswordUsecase := usecases.NewDoctorResetPasswordUsecaseImpl(&usecases.DoctorResetPasswordUsecaseOpts{
		DoctorResetPasswordRepo: doctorResetPasswordRepo,
		DoctorRepo:              doctorRepo,
		HashAlgorithm:           utils.NewBCryptHasher(),
		EmailSender:             utils.NewGoogleEmailSender(),
		AuthTokenProvider:       utils.NewJwtProvider(config),
		Transactor:              repositories.NewTransactor(db),
	})

	orderUsecase := usecases.NewOrderUsecaseImpl(&usecases.OrderUsecaseOpts{
		OrderRepository:           orderRepo,
		CartRepository:            cartRepo,
		PharmacyProductRepository: pharmacyProductRepo,
		StockHistoryRepository:    stockHistoryRepo,
		Transactor:                repositories.NewTransactor(db),
		UploadFile:                utils.NewCloudinaryUploadFile(),
	})
	adminUsecase := usecases.NewAdminUsecaseImpl(&usecases.AdminUsecaseOpts{
		AdminRepository: adminRepo,
		Transactor:      repositories.NewTransactor(db),
		HashAlgorithm:   utils.NewBCryptHasher(),
	})
	mutationStatusUsecase := usecases.NewMutationSatusUsecaseImpl(&usecases.MutationSatusUsecaseOpts{MutationStatusRepo: mutationStatusRepo})
	stockTransferUsecase := usecases.NewStockTransferUsecaseImpl(&usecases.StockTransferUsecaseOpts{
		StockTransferRepo:   stockTransferRepo,
		PharmacyRepo:        pharmacyRepo,
		PharmacyProductRepo: pharmacyProductRepo,
		Transactor:          repositories.NewTransactor(db),
		StockHistoryRepo:    stockHistoryRepo,
	})
	stockHistoryReportUsecase := usecases.NewStockHistoryReportUsecaseImpl(&usecases.StockHistoryReportUsecaseOpts{StockHistoryReportRepo: stockHistoryReportRepo})
	salesReportUsecase := usecases.NewSalesReportUsecaseImpl(&usecases.SalesReportUsecaseOpts{
		SalesReportRepo: salesReportRepo,
		CategoryRepo:    categoryRepo,
	})
	salesReporctCategoryUsecase := usecases.NewSalesReportCategoryUsecaseImpl(&usecases.SalesReportCategoryUsecaseOpts{SalesReportCategoryRepo: salesReportCategoryRepo})
	mostBoughtUserUsecase := usecases.NewMostBoughtUserUsecaseImpl(&usecases.MostBoughtUserUsecaseOpts{
		MostBoughtUserRepo:  mostBoughtUserRepo,
		CategoryRepo:        categoryRepo,
		PharmacyProductRepo: pharmacyProductRepo,
		PharmacyRepo:        pharmacyRepo,
	})

	userHandler := handlers.NewUserHandler(&handlers.UserHandlerOpts{
		UserUsecase: userUsecase,
		UploadFile:  utils.NewCloudinaryUploadFile(),
	})
	authHandler := handlers.NewAuthHandler(&handlers.AuthHandlerOpts{
		LoginUsecase:        loginUsecase,
		RegisterUsecase:     registerUsecase,
		VerifyUsecase:       verifyUsecase,
		RefreshTokenUsecase: refreshTokenUsecase,
		OAuthUsecase:        oauthUsecase,
		AuthTokenProvider:   utils.NewJwtProvider(config),
	})
	pharmacyHandler := handlers.NewPharmacyHandler(&handlers.PharmacyHandlerOpts{
		PharmacyUsecase: pharmacyUsecase,
	})
	doctorHandler := handlers.NewDoctorHandler(&handlers.DoctorHandlerOpts{
		DoctorUsecase: doctorUsecase,
		UploadFile:    utils.NewCloudinaryUploadFile(),
	})
	specialistHandler := handlers.NewSpecialistHandler(&handlers.SpecialistHandlerOpts{SpecialistUsecase: specialistUsecase})
	pharmacyManagerHandler := handlers.NewPharmacyManagerHandler(&handlers.PharmacyManagerHandlerOpts{PharmacyManagerUsecase: pharmacyManagerUsecase})
	userAddressHandler := handlers.NewUserAddressHandler(&handlers.UserAddressHandlerOpts{UserAddressUsecase: userAddresUsecase})
	pharmacyProductHandler := handlers.NewPharmacyProductHandler(&handlers.PharmacyProductHandlerOpts{
		PharmacyProductUsecase: pharmacyProductUsecase,
	})
	categoryHandler := handlers.NewCategoryHandler(&handlers.CategoryHandlerOpts{CategoryUsecase: categoryUsecase})
	productHandler := handlers.NewProductHandler(&handlers.ProductHandlerOpts{
		ProductUsecase:         productUsecase,
		ProductCategoryUsecase: productCategoryUsecase,
		UploadFile:             utils.NewCloudinaryUploadFile(),
	})
	consultationHandler := handlers.NewConsultationHandler(&handlers.ConsultationHandlerOpts{
		ConsultationUsecase: consultationUsecase,
	})

	wsHandler := ws.NewWebSocketHandler(&ws.WebSocketHandlerOpts{
		Hub:                 hub,
		ConsultationUsecase: consultationUsecase,
	})

	productFieldHandler := handlers.NewProductFieldHandler(&handlers.ProductFieldHandlerOpts{ProductFieldUsecase: productFieldUsecase})
	cartHandler := handlers.NewCartHandler(&handlers.CartHandlerOpts{CartUsecase: cartUsecase})
	locationHandler := handlers.NewLocationHandler(&handlers.LocationHandlerOpts{LocationUsecase: locationUsecase})
	stockHistoryHandler := handlers.NewStockHistoryHandler(&handlers.StockHistoryHandlerOpts{StockHistoryUsecase: stockHistoryUsecase})
	shippingMethodHandler := handlers.NewShippingMethodHandler(&handlers.ShippingMethodHandlerOpts{
		ShippingMethodUsecase: shippingMethodUsecase,
	})
	resetPasswordHandler := handlers.NewUserResetPasswordHandler(&handlers.ResetPasswordHandler{
		UserResetPasswordUsecase:   userResetPasswordUsecase,
		DoctorResetPasswordUsecase: doctorResetPasswordUsecase,
		AuthTokenProvider:          utils.NewJwtProvider(config),
	})
	orderHandler := handlers.NewOrderHandler(&handlers.OrderHandlerOpts{OrderUsecase: orderUsecase})
	adminHandler := handlers.NewAdminHandler(&handlers.AdminHandlerOpts{AdminUsecase: adminUsecase})
	mutationStatusHandler := handlers.NewMutationStatusHandler(&handlers.MutationSatusHandlerOpts{MutationSatusUsecase: mutationStatusUsecase})
	stockTransferHandler := handlers.NewStockTransferHandler(&handlers.StockTransferHandler{StockTransferUsecase: stockTransferUsecase})
	stochHistoryReportHandler := handlers.NewStockHistoryReportHandler(&handlers.StockHistoryReportHandlerOpts{StockHistoryReportUsecase: stockHistoryReportUsecase})
	salesReportHandler := handlers.NewSalesReportHandler(&handlers.SalesReportHandlerOpts{SalesReportUsecase: salesReportUsecase})
	salesReportCategoryHandler := handlers.NewSalesReportCategoryHandler(&handlers.SalesReportCategoryHandlerOpts{SalesReportCategoryUsecase: salesReporctCategoryUsecase})
	mostBoughtUserHandler := handlers.NewMostBoughtUserHandler(&handlers.MostBoughtUserHandlerOpts{MostBoughtUserUsecase: mostBoughtUserUsecase})

	return NewRouter(config, &RouterOpts{
		User:                userHandler,
		Auth:                authHandler,
		Pharmacy:            pharmacyHandler,
		Doctor:              doctorHandler,
		Specialist:          specialistHandler,
		PharmacyManager:     pharmacyManagerHandler,
		UserAddress:         userAddressHandler,
		PharmacyProduct:     pharmacyProductHandler,
		Category:            categoryHandler,
		WebSocket:           wsHandler,
		Product:             productHandler,
		Consultation:        consultationHandler,
		ProductField:        productFieldHandler,
		Cart:                cartHandler,
		Location:            locationHandler,
		StockHistory:        stockHistoryHandler,
		ShppingMethod:       shippingMethodHandler,
		ResetPassword:       resetPasswordHandler,
		Order:               orderHandler,
		Admin:               adminHandler,
		MutationStatus:      mutationStatusHandler,
		StockTransfer:       stockTransferHandler,
		StockHistoryReport:  stochHistoryReportHandler,
		SalesReport:         salesReportHandler,
		SalesReportCategory: salesReportCategoryHandler,
		MostBoughtUser:      mostBoughtUserHandler,
	})
}

func Init() {
	config, err := utils.ConfigInit()
	if err != nil {
		log.Fatalf("error getting env: %s", err.Error())
	}

	hub := ws.NewHub()
	go hub.Run()

	router := createRouter(config, hub)

	srv := http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", config.Port),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown: ", err)
		}
	}()

	<-ctx.Done()
	log.Println("Server exiting")

}

func NewRouter(config utils.Config, handlers *RouterOpts) *gin.Engine {
	router := gin.Default()

	router.ContextWithFallback = true

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	router.Use(middlewares.CORS, middlewares.RequestId, middlewares.Logger(log), middlewares.ErrorHandling)

	publicRouter := router.Group("/")
	publicRouter.Use(middlewares.SetPublic())
	{
		authRouter := publicRouter.Group("/auth")
		{
			authRouter.POST("/login", handlers.Auth.Login)
			authRouter.POST("/register/user", handlers.Auth.RegisterUser)
			authRouter.POST("/refresh-token", handlers.Auth.RefreshToken)
			authRouter.POST("/verify", handlers.Auth.Verification)
			authRouter.POST("/verify/resend", handlers.Auth.ResendVerification)
			authRouter.POST("/register/doctor", handlers.Auth.RegisterDoctor)
			authRouter.POST("/oauth/google", handlers.Auth.GoogleOauth)
			authRouter.POST("/forgot-password", handlers.ResetPassword.ForgotPassword)
			authRouter.POST("/reset-password", handlers.ResetPassword.ResetPassword)

			privateAuthRouter := authRouter.Group("/")
			privateAuthRouter.Use(middlewares.JwtAdminAuthMiddleware(config))
			privateAuthRouter.POST("/register/pharmacy-manager", handlers.Auth.RegisterPharmacyManager)
		}

		pharmacyRouter := publicRouter.Group("/pharmacies")
		{
			pharmacyRouter.GET("/:id", handlers.Pharmacy.GetPharmacyById)
		}

		doctorRouter := publicRouter.Group("/doctors")
		{
			doctorRouter.GET("/verified", handlers.Doctor.GetAllDoctor)
			doctorRouter.GET("/:id", handlers.Doctor.GetDoctorById)
			doctorRouter.GET("/subscribe", handlers.WebSocket.JoinDoctorRoom)
			doctorRouter.GET("/:id/consultations/:consultationId/rooms", handlers.WebSocket.JoinRoomAsDoctor)
		}

		userRouter := publicRouter.Group("/users")
		{
			userRouter.GET("/:id/consultations/:consultationId/rooms", handlers.WebSocket.JoinRoomAsUser)
		}

		specialistRouter := publicRouter.Group("/specialists")
		{
			specialistRouter.GET("", handlers.Specialist.GetAllSpecialist)
		}

		productRouter := publicRouter.Group("/products")
		{
			productRouter.GET("/:id", handlers.Product.GetProductById)
			productRouter.GET("/", handlers.Product.GetAllProduct)
			productRouter.GET("/nearest", handlers.PharmacyProduct.GetNearestProducts)
			productRouter.GET("/nearest/search", handlers.PharmacyProduct.GetAllNearestPharmacyProducts)
			productRouter.GET("/detail", handlers.PharmacyProduct.ProductDetail)
		}

		categoryRouter := publicRouter.Group("/categories")
		{
			categoryRouter.GET("", handlers.Category.GetAllCategory)
			categoryRouter.GET("/:id", handlers.Category.GetCategoryById)
		}

		locationRouter := publicRouter.Group("/loc")
		{
			locationRouter.GET("/provinces", handlers.Location.GetAllProvinces)
			locationRouter.GET("/cities/:id", handlers.Location.GetCitiesByProvinceId)
			locationRouter.GET("/districts/:id", handlers.Location.GetDistrictsByCityId)
			locationRouter.GET("/sub-districts/:id", handlers.Location.GetSubDistrictsByDistrictId)
			locationRouter.GET("/reverse", handlers.Location.ReverseCoordinate)
		}

		mostBoughtUserRouter := publicRouter.Group("/most-boughts/search")
		{
			mostBoughtUserRouter.GET("/", handlers.MostBoughtUser.GetMostBought)
		}
	}

	privateRouter := router.Group("/")
	{
		privateRouter.Use(middlewares.JwtAuthMiddleware(config))

		adminPrivate := privateRouter.Group("/admins")
		{
			adminPrivate.Use(middlewares.JwtAdminAuthMiddleware(config))
			adminPrivate.POST("", handlers.Admin.CreateAdmin)
			adminPrivate.DELETE("/:id", handlers.Admin.DeleteAdmin)
			adminPrivate.GET("/:id", handlers.Admin.GetAdminById)
			adminPrivate.GET("", handlers.Admin.GetAllAdmin)

			adminOrderRouter := adminPrivate.Group("/orders")
			adminOrderRouter.GET("/", handlers.Order.GetAllOrderByAdmin)
			adminOrderRouter.PATCH("/:id/approve", handlers.Order.UpdateOrderStatusToProcessing)
			adminOrderRouter.PATCH("/:id/cancel", handlers.Order.CancelOrderByAdmin)
		}

		authPrivateRouter := privateRouter.Group("/auth")
		{
			authPrivateRouter.Use(middlewares.JwtMultiRoleMiddleware(config, []string{constants.UserRole, constants.DoctorRole}))
			authPrivateRouter.POST("/change-password", handlers.ResetPassword.ChangePassword)
		}

		privateUserRouter := privateRouter.Group("/users")
		{
			adminUserRouter := privateUserRouter.Group("/")
			{
				adminUserRouter.Use(middlewares.JwtAdminAuthMiddleware(config))
				adminUserRouter.GET("", handlers.User.GetAllUser)
				adminUserRouter.GET("/:id", handlers.User.GetUserById)
				adminUserRouter.DELETE("/:id", handlers.User.DeleteUser)
				adminUserRouter.PUT("/:id", handlers.User.UpdateUser)
				adminUserRouter.POST("/:id/addresses", handlers.UserAddress.CreateUserAddress)
				adminUserRouter.GET("/:id/addresses/:addressId", handlers.UserAddress.GetAddressById)
				adminUserRouter.PUT("/:id/addresses/:addressId", handlers.UserAddress.UpdateUserAddress)
				adminUserRouter.DELETE("/:id/addresses/:addressId", handlers.UserAddress.DeleteUserAddress)
			}
			userRouter := privateUserRouter.Group("/")
			{
				userRouter.Use(middlewares.JwtUserAuthMiddleware(config))
				userRouter.GET("/profile", handlers.User.GetUserById)
				userRouter.PUT("/profile", handlers.User.UpdateUser)
				userRouter.POST("/profile/addresses", handlers.UserAddress.CreateUserAddress)
				userRouter.GET("/profile/addresses/:addressId", handlers.UserAddress.GetAddressById)
				userRouter.PUT("/profile/addresses/:addressId", handlers.UserAddress.UpdateUserAddress)
				userRouter.DELETE("/profile/addresses/:addressId", handlers.UserAddress.DeleteUserAddress)

				userConsultRouter := userRouter.Group("/consultations")
				userConsultRouter.GET("", handlers.Consultation.GetAllConsultationByUser)
				userConsultRouter.POST("", handlers.Consultation.CreateConsultation)
				userConsultRouter.GET("/:id", handlers.Consultation.GetConsultationById)
				userConsultRouter.POST("/:id/chats", handlers.Consultation.CreateChat)
				userConsultRouter.POST("/:id/chats/file", handlers.Consultation.CreateChatFile)
				userConsultRouter.POST("/:id/end", handlers.Consultation.EndConsultation)
				userConsultRouter.POST("/:id/prescription/add", handlers.Consultation.AddPrescriptionToCart)
				userConsultRouter.POST("/rooms", handlers.WebSocket.CreateRoom)

				userOrderRouter := userRouter.Group("/orders")
				userOrderRouter.POST("/", handlers.Order.CreateOrder)
				userOrderRouter.GET("/", handlers.Order.GetAllOrderByUser)
				userOrderRouter.PATCH("/:orderId/complete", handlers.Order.UpdateOrderStatusToCompleted)
				userOrderRouter.POST("/payment-proof", handlers.Order.UploadPaymentProof)
				userOrderRouter.PATCH("/:orderId/cancel", handlers.Order.UpdateOrderStatusToCanceled)
			}
		}

		privatePharmacyRouter := privateRouter.Group("/pharmacies")
		{
			adminPrivatePharmacyRouter := privatePharmacyRouter.Group("")
			{
				adminPrivatePharmacyRouter.Use(middlewares.JwtMultiRoleMiddleware(config, []string{constants.AdminRole, constants.PharmacyManagerRole}))
				adminPrivatePharmacyRouter.DELETE("/:id", handlers.Pharmacy.DeletePharmacyById)
			}

			pmPrivatePharmacyRouter := privatePharmacyRouter.Group("")
			{
				pmPrivatePharmacyRouter.Use(middlewares.JwtPharmacyManagerMiddleware(config))
				pmPrivatePharmacyRouter.PUT("/:id", handlers.Pharmacy.UpdatePharmacy)
				pmPrivatePharmacyRouter.GET("", handlers.Pharmacy.GetAllPharmacyByPharmacyManager)
				pmPrivatePharmacyRouter.GET("/:id/products", handlers.PharmacyProduct.GetPharmacyProductsByPharmacyId)
				pmPrivatePharmacyRouter.GET("/products/:id", handlers.PharmacyProduct.GetPharmacyProductById)
			}
		}

		privateManagerRouter := privateRouter.Group("/pharmacy-managers")
		{
			stockTransferRouter := privateManagerRouter.Group("/stock-mutations")
			{
				stockTransferRouter.Use(middlewares.JwtPharmacyManagerMiddleware(config))
				stockTransferRouter.GET("/", handlers.StockTransfer.GetAllStockTransfer)
				stockTransferRouter.POST("/", handlers.StockTransfer.CreateStockTransfer)
				stockTransferRouter.PUT("/:id", handlers.StockTransfer.UpdateMutationStatus)
			}

			admiOnlyPrivateManagerRouter := privateManagerRouter.Group("")
			{
				admiOnlyPrivateManagerRouter.Use(middlewares.JwtAdminAuthMiddleware(config))
				admiOnlyPrivateManagerRouter.POST("/:id/pharmacies", handlers.Pharmacy.CreatePharmacy)
			}

			adminPrivateManagerRouter := privateManagerRouter.Group("")
			{
				adminPrivateManagerRouter.Use(middlewares.JwtMultiRoleMiddleware(config, []string{constants.AdminRole, constants.PharmacyManagerRole}))
				adminPrivateManagerRouter.GET("", handlers.PharmacyManager.GetAllPharmacyManager)
				adminPrivateManagerRouter.GET("/:id", handlers.PharmacyManager.GetPharmacyManagerById)
				adminPrivateManagerRouter.DELETE("/:id", handlers.PharmacyManager.DeletePharmacyMaagerById)
				adminPrivateManagerRouter.PUT("/:id", handlers.PharmacyManager.UpdatePharmacyManager)
				adminPrivateManagerRouter.GET("/:id/pharmacies", handlers.Pharmacy.GetAllPharmacyByPharmacyManager)
				adminPrivateManagerRouter.POST("/pharmacies", handlers.Pharmacy.CreatePharmacy)
			}

			pharmacyManagerRouter := privateManagerRouter.Group("/")
			{
				pharmacyManagerRouter.Use(middlewares.JwtPharmacyManagerMiddleware(config))

				pharmacyManagerOrderRouter := pharmacyManagerRouter.Group("/orders")
				pharmacyManagerOrderRouter.GET("/", handlers.Order.GetAllOrderByPharmacyManager)
				pharmacyManagerOrderRouter.PATCH("/:orderId/ship", handlers.Order.UpdateOrderStatusToShipped)
				pharmacyManagerOrderRouter.PATCH("/:orderId/cancel", handlers.Order.CancelOrderByPharmacyManager)
			}
		}

		privateDoctorRouter := privateRouter.Group("/doctors")
		{
			adminPrivateDoctorRouter := privateDoctorRouter.Group("")
			{
				adminPrivateDoctorRouter.Use(middlewares.JwtAdminAuthMiddleware(config))
				adminPrivateDoctorRouter.GET("", handlers.Doctor.GetAllDoctor)
				adminPrivateDoctorRouter.PUT("/:id", handlers.Doctor.UpdateDoctor)
				adminPrivateDoctorRouter.DELETE("/:id", handlers.Doctor.DeleteDoctor)

			}
			doctorRouter := privateDoctorRouter.Group("/")
			{
				doctorRouter.Use(middlewares.JwtDoctorMiddleware(config))
				doctorRouter.GET("/profile", handlers.Doctor.GetDoctorById)
				doctorRouter.PUT("/profile", handlers.Doctor.UpdateDoctor)
				doctorRouter.POST("/toggle-is-online", handlers.Doctor.ToggleDoctorIsOnline)

				doctorConsultRouter := doctorRouter.Group("/consultations")
				doctorConsultRouter.GET("/:id", handlers.Consultation.GetConsultationById)
				doctorConsultRouter.POST("/:id/certificate", handlers.Consultation.CreateCertificate)
				doctorConsultRouter.POST("/:id/prescription", handlers.Consultation.CreatePrescription)
				doctorConsultRouter.GET("", handlers.Consultation.GetAllConsultationByDoctor)
				doctorConsultRouter.POST("/:id/chats", handlers.Consultation.CreateChat)
				doctorConsultRouter.POST("/:id/end", handlers.Consultation.EndConsultation)
			}
		}

		privateCategoryRouter := privateRouter.Group("/categories")
		{
			privateCategoryRouter.Use(middlewares.JwtAdminAuthMiddleware(config))
			privateCategoryRouter.POST("", handlers.Category.CreateCategory)
			privateCategoryRouter.PUT("/:id", handlers.Category.UpdateCategory)
			privateCategoryRouter.DELETE("/:id", handlers.Category.DeleteCategory)
		}

		privateProductRouter := privateRouter.Group("/products")
		{
			adminPrivateProductRouter := privateProductRouter.Group("")
			{
				adminPrivateProductRouter.Use(middlewares.JwtAdminAuthMiddleware(config))
				adminPrivateProductRouter.POST("/", handlers.Product.CreateProduct)
				adminPrivateProductRouter.PUT("/:id", handlers.Product.UpdateProduct)
				adminPrivateProductRouter.GET("/forms", handlers.ProductField.GetAllForm)
				adminPrivateProductRouter.GET("/classifications", handlers.ProductField.GetAllClassification)
				adminPrivateProductRouter.GET("/manufactures", handlers.ProductField.GetAllManufacture)
				adminPrivateProductRouter.DELETE("/:id", handlers.Product.DeleteProduct)
			}
		}

		privateCartRouter := privateRouter.Group("/carts")
		{
			privateCartRouter.Use(middlewares.JwtUserAuthMiddleware(config))
			privateCartRouter.POST("", handlers.Cart.CreateCartItem)
			privateCartRouter.PUT("/:id/increase", handlers.Cart.IncreaseCartItem)
			privateCartRouter.PUT("/:id/decrease", handlers.Cart.DecreaseCartItem)
			privateCartRouter.DELETE("/:id", handlers.Cart.DeleteCartItem)
			privateCartRouter.GET("", handlers.Cart.GetAllCartItem)
		}
		privatePharmacyProductRouter := privateRouter.Group("/pharmacy-products")
		{
			privatePharmacyProductRouter.Use(middlewares.JwtPharmacyManagerMiddleware(config))
			privatePharmacyProductRouter.POST("/", handlers.PharmacyProduct.CreatePharmacyProduct)
			privatePharmacyProductRouter.PUT("/:id", handlers.PharmacyProduct.UpdatePharmacyProduct)
			privatePharmacyProductRouter.DELETE("/:id", handlers.PharmacyProduct.DeletePharmacyProduct)
		}

		privateStockHistoryRouter := privateRouter.Group("/stock-histories")
		{
			privateStockHistoryRouter.Use(middlewares.JwtPharmacyManagerMiddleware(config))
			privateStockHistoryRouter.GET("/:pharmacyId", handlers.StockHistory.GetStockHistoriesByPharmacyId)
		}

		privateShippingCostRouter := privateRouter.Group("/shipping-costs")
		{
			privateShippingCostRouter.Use(middlewares.JwtUserAuthMiddleware(config))
			privateShippingCostRouter.POST("/official", handlers.ShppingMethod.GetOfficialShippingCost)
			privateShippingCostRouter.POST("/non-official", handlers.ShppingMethod.GetNonOfficialShippingCost)
		}

		privateStockHistoryReport := privateRouter.Group("/stock-history-reports")
		{
			privateStockHistoryReport.Use(middlewares.JwtMultiRoleMiddleware(config, []string{constants.PharmacyManagerRole, constants.AdminRole}))
			privateStockHistoryReport.GET("/", handlers.StockHistoryReport.GetStockHistoryReports)
		}

		privateSalesReport := privateRouter.Group("/sales-reports")
		{
			privateSalesReport.Use(middlewares.JwtMultiRoleMiddleware(config, []string{constants.AdminRole, constants.PharmacyManagerRole}))
			privateSalesReport.GET("/", handlers.SalesReport.GetSalesReports)
		}

		privateSalesReportCategory := privateRouter.Group("/sales-report-categories")
		{
			privateSalesReportCategory.Use(middlewares.JwtMultiRoleMiddleware(config, []string{constants.AdminRole, constants.PharmacyManagerRole}))
			privateSalesReportCategory.GET("/", handlers.SalesReportCategory.GetSalesReportCategories)
		}

		privateOrder := privateRouter.Group("/orders")
		{
			privateOrder.Use(middlewares.JwtMultiRoleMiddleware(config, []string{constants.UserRole, constants.AdminRole, constants.PharmacyManagerRole}))
			privateOrder.GET("/:id", handlers.Order.GetOrderDetail)
		}
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, dtos.ErrResponse{Message: constants.EndpointNotFoundErrMsg})
	})

	return router
}
