package main

import (
	"context"
	"qibla-backend-chat/pkg/aes"
	"qibla-backend-chat/pkg/aesfront"
	"qibla-backend-chat/pkg/env"
	"qibla-backend-chat/pkg/interfacepkg"
	"qibla-backend-chat/pkg/jwe"
	"qibla-backend-chat/pkg/jwt"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/pkg/mandrill"
	"qibla-backend-chat/pkg/mongo"
	"qibla-backend-chat/pkg/odoo"
	"qibla-backend-chat/pkg/pg"
	"qibla-backend-chat/pkg/pusher"
	"qibla-backend-chat/pkg/str"
	boot "qibla-backend-chat/server/bootstrap"
	"qibla-backend-chat/usecase"
	"time"

	"github.com/rs/cors"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v7"
	validator "gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	idTranslations "gopkg.in/go-playground/validator.v9/translations/id"
)

var (
	_, b, _, _      = runtime.Caller(0)
	basepath        = filepath.Dir(b)
	debug           = false
	host            string
	validatorDriver *validator.Validate
	uni             *ut.UniversalTranslator
	translator      ut.Translator
	envConfig       map[string]string
	corsDomainList  []string
)

// Init first time running function
func init() {
	// Load env variable from .env file
	envConfig = env.NewEnvConfig("../.env")

	// Load cors domain list
	corsDomainList = strings.Split(envConfig["APP_CORS_DOMAIN"], ",")

	host = envConfig["APP_HOST"]
	if str.StringToBool(envConfig["APP_DEBUG"]) {
		debug = true
		log.Printf("Running on Debug Mode: On at host [%v]", host)
	}
}

func main() {
	ctx := "main"

	// Connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     envConfig["REDIS_HOST"],
		Password: envConfig["REDIS_PASSWORD"],
		DB:       0,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	// Postgre DB connection
	dbInfo := pg.Connection{
		Host:    envConfig["DATABASE_HOST"],
		DB:      envConfig["DATABASE_DB"],
		User:    envConfig["DATABASE_USER"],
		Pass:    envConfig["DATABASE_PASSWORD"],
		Port:    str.StringToInt(envConfig["DATABASE_PORT"]),
		SslMode: "disable",
	}
	db, err := dbInfo.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// JWT credential
	jwtCredential := jwt.Credential{
		Secret:           envConfig["TOKEN_SECRET"],
		ExpSecret:        str.StringToInt(envConfig["TOKEN_EXP_SECRET"]),
		RefreshSecret:    envConfig["TOKEN_REFRESH_SECRET"],
		RefreshExpSecret: str.StringToInt(envConfig["TOKEN_EXP_REFRESH_SECRET"]),
	}

	// JWE credential
	jweCredential := jwe.Credential{
		KeyLocation: envConfig["APP_PRIVATE_KEY_LOCATION"],
		Passphrase:  envConfig["APP_PRIVATE_KEY_PASSPHRASE"],
	}

	// AES credential
	aesCredential := aes.Credential{
		Key: envConfig["AES_KEY"],
	}

	// AES Front credential
	aesFrontCredential := aesfront.Credential{
		Key: envConfig["AES_FRONT_KEY"],
		Iv:  envConfig["AES_FRONT_IV"],
	}

	// Mandrill credential
	mandrillCredential := mandrill.Credential{
		Key:      envConfig["MANDRILL_KEY"],
		FromMail: envConfig["MANDRILL_FROM_MAIL"],
		FromName: envConfig["MANDRILL_FROM_NAME"],
	}

	// ODOO connection
	odooInfo := odoo.Connection{
		Host: envConfig["ODOO_URL"],
		DB:   envConfig["ODOO_DATABASE"],
		User: envConfig["ODOO_ADMIN"],
		Pass: envConfig["ODOO_PASSWORD"],
	}
	odooDB, err := odooInfo.Connect()
	if err != nil {
		panic(err)
	}
	defer odooDB.Close()

	// Pusher credential
	pusherCredential := pusher.Credential{
		AppID:   envConfig["PUSHER_APP_ID"],
		Key:     envConfig["PUSHER_KEY"],
		Secret:  envConfig["PUSHER_SECRET"],
		Cluster: envConfig["PUSHER_CLUSTER"],
	}

	// MongoDB connection
	mongoConnection := mongo.Connection{
		URL:    envConfig["MONGO_URL"],
		DBName: envConfig["MONGO_DB"],
	}
	mongoDB, err := mongoConnection.Connect()
	if err != nil {
		panic(err)
	}
	ctxMongo, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer mongoDB.Disconnect(ctxMongo)

	// Validator initialize
	validatorInit()

	// Load contract struct
	contractUC := usecase.ContractUC{
		DB:          db,
		MongoDB:     mongoDB,
		MongoDBName: mongoConnection.DBName,
		Redis:       redisClient,
		EnvConfig:   envConfig,
		Jwt:         jwtCredential,
		Jwe:         jweCredential,
		Aes:         aesCredential,
		AesFront:    aesFrontCredential,
		Mandrill:    mandrillCredential,
		Odoo:        odooDB,
		Pusher:      pusherCredential,
	}

	r := chi.NewRouter()
	// Cors setup
	r.Use(cors.New(cors.Options{
		AllowedOrigins: corsDomainList,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler)

	// load application bootstrap
	bootApp := boot.Bootup{
		R:          r,
		CorsDomain: corsDomainList,
		EnvConfig:  envConfig,
		DB:         db,
		Redis:      redisClient,
		Validator:  validatorDriver,
		Translator: translator,
		ContractUC: contractUC,
		Jwt:        jwtCredential,
		Jwe:        jweCredential,
	}

	// register middleware
	bootApp.RegisterMiddleware()

	// register routes
	bootApp.RegisterRoutes()

	// Create static folder for file uploading
	filePath := envConfig["FILE_STATIC_FILE"]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, os.ModePerm)
	}

	// Register folder for a go static folder
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, filePath)
	fileServer(r, envConfig["FILE_PATH"], http.Dir(filesDir))

	// Create static folder for html picture
	filePath = envConfig["HTML_FILE_STATIC_FILE"]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, os.ModePerm)
	}
	filesDir = filepath.Join(workDir, filePath)
	fileServer(r, envConfig["HTML_FILE_PATH"], http.Dir(filesDir))

	// Log start server
	startBody := map[string]interface{}{
		"Host":     host,
		"Location": str.DefaultData(envConfig["APP_DEFAULT_LOCATION"], "Asia/Jakarta"),
	}
	logruslogger.Log(logruslogger.InfoLevel, interfacepkg.Marshall(startBody), ctx, "server_start", "")

	// Run the app
	http.ListenAndServe(host, r)
}

func validatorInit() {
	en := en.New()
	id := id.New()
	uni = ut.New(en, id)

	transEN, _ := uni.GetTranslator("en")
	transID, _ := uni.GetTranslator("id")

	validatorDriver = validator.New()

	enTranslations.RegisterDefaultTranslations(validatorDriver, transEN)
	idTranslations.RegisterDefaultTranslations(validatorDriver, transID)

	switch envConfig["APP_LOCALE"] {
	case "en":
		translator = transEN
	case "id":
		translator = transID
	}
}

// fileServer ...
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
