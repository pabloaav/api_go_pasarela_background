package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/api/background"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/api/middlewares"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/api/middlewares/middlewareinterno"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/api/routes"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/database"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/internal/storage"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/administracion"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/apilink"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/auditoria"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/banco"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/cierrelote"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/commons"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/email"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/reportes"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/util"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela_background/pkg/domains/webhook"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
)

func InicializarApp(clienteHttp *http.Client, clienteSql *database.MySQLClient, clienteFile *os.File) *fiber.App {
	//Servicios comunes
	middlewares := middlewares.MiddlewareManager{HTTPClient: clienteHttp}
	fileRepository := commons.NewFileRepository(clienteFile)
	commonsService := commons.NewCommons(fileRepository)

	utilRepository := util.NewUtilRepository(clienteSql)
	utilService := util.NewUtilService(utilRepository)

	runEndpoint := util.NewRunEndpoint(clienteHttp, utilService)

	//Valida si existe un correo para solicitud de nuevas cuentas si no existe lo crea.
	utilService.FirstOrCreateConfiguracionService("EMAIL_SOLICITUD_CUENTA", "Email que recibirá la solicitud de apertura de cuenta", "developmenttelco@gmail.com")

	//ApiLink
	apiLinkRemoteRepository := apilink.NewRemote(clienteHttp, utilService)
	apiLinkRepository := apilink.NewRepository(clienteSql, utilService)
	apiLinkService := apilink.NewService(apiLinkRemoteRepository, apiLinkRepository)

	// webhooks
	webhooksRepository := webhook.NewRemote(clienteHttp)

	//Store Service
	storeService := storage.NewS3Session()
	storeServiceEst := cierrelote.NewStore(storeService)

	auditoriaRespository := auditoria.NewAuditoriaRepository(clienteSql)
	auditoriaService := auditoria.AuditoriaService(auditoriaRespository)

	administracionRepository := administracion.NewRepository(clienteSql, auditoriaService, utilService)
	administracionService := administracion.NewService(administracionRepository, apiLinkService, commonsService, utilService, webhooksRepository, storeServiceEst)
	middlewaresPasarela := middlewareinterno.MiddlewareManagerPasarela{Service: administracionService}
	/* MOVIMIENTOS BANCO: servicio para consultar y validar movimientos de pagos acreditados en la cuenta de telco*/
	movimientosBancoRemoteRepository := banco.NewRemote(clienteHttp)
	movimientosBancoService := banco.NewService(movimientosBancoRemoteRepository, utilService, administracionService)
	/* Envio de correos */
	emailService := email.NewService(commonsService)
	/* REPORTES CLIENTES */
	reportesRepository := reportes.NewRepository(clienteSql, auditoriaService, utilService)
	reportesService := reportes.NewService(reportesRepository, administracionService, utilService, commonsService, storeServiceEst, emailService)
	cierreloteRepository := cierrelote.NewRepository(clienteSql, utilRepository)
	storage := storage.NewS3Session()
	reafileStore := cierrelote.NewStore(storage)
	cierreloteService := cierrelote.NewService(cierreloteRepository, commonsService, utilService, reafileStore, administracionService, movimientosBancoService)

	// descomentar esto en servidor
	engine := html.New(filepath.Join(filepath.Base(config.DIR_BASE), "api", "views"), ".html")
	//descomentar esto en local
	// engine := html.New("views", ".html")
	engine.Delims("${", "}")

	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var msg string
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				msg = e.Message
			}

			if msg == "" {
				msg = "No se pudo procesar el llamado a la api: " + err.Error()
			}

			_ = ctx.Status(code).JSON(internalError{
				Message: msg,
			})

			return nil
		},
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.ALLOW_ORIGIN + ", " + config.AUTH,
		AllowHeaders: "",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))
	app.Get("/", func(ctx *fiber.Ctx) error {
		msj := fmt.Sprintf("Corrientes Telecomunicaciones Api %v", config.API_NOMBRE)
		return ctx.Send([]byte(msj))
	})

	//Procesos en segundo plano
	background.BackgroudServices(administracionService, cierreloteService, utilService, movimientosBancoService, reportesService, runEndpoint)

	/* Routes */
	// cierre de lote
	cierreloteRouter := app.Group("/cierrelote")
	routes.CierreLoteRoutes(cierreloteRouter, cierreloteService, administracionService, utilService, middlewares)

	// reportes
	reportesRouter := app.Group("/reporte")
	routes.ReporteRoutes(reportesRouter, middlewares, middlewaresPasarela, reportesService, emailService, runEndpoint)

	// administracion
	administracionRouter := app.Group("/administracion")
	routes.AdministracionRoutes(administracionRouter, middlewares, administracionService, utilService, movimientosBancoService)

	// banco
	banco := app.Group("/banco")
	routes.BancoRoutes(banco, movimientosBancoService, middlewares)

	return app
}

func main() {
	var HTTPTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false, // <- this is my adjustment
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	var HTTPClient = &http.Client{
		Transport: HTTPTransport,
	}

	//HTTPClient.Timeout = time.Second * 120 //Todo validar si este tiempo está bien
	clienteSQL := database.NewMySQLClient()
	osFile := os.File{}

	app := InicializarApp(HTTPClient, clienteSQL, &osFile)

	_ = app.Listen(":3400")
}

type internalError struct {
	Message string `json:"message"`
}
