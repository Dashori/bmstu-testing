package registry

import (
	config "backend/config"
	hasherImplementation "backend/internal/pkg/hasher/implementation"
	"backend/internal/repository"
	postgres_repo "backend/internal/repository/postgres_repo"
	"backend/internal/services"
	servicesImplementation "backend/internal/services/implementation"
	"github.com/charmbracelet/log"
	"os"
	"github.com/opentracing/opentracing-go"
    "github.com/opentracing/opentracing-go/ext"
    "github.com/opentracing/opentracing-go/log"
)

type AppServiceFields struct {
	ClientService services.ClientService
	DoctorService services.DoctorService
	PetService    services.PetService
	RecordService services.RecordService
}

type App struct {
	Config       config.Config
	Repositories *AppRepositoryFields
	Services     *AppServiceFields
	Logger       *log.Logger
}

type AppRepositoryFields struct {
	ClientRepository repository.ClientRepository
	DoctorRepository repository.DoctorRepository
	PetRepository    repository.PetRepository
	RecordRepository repository.RecordRepository
}

func (a *App) initRepositories(fields *postgres_repo.PostgresRepositoryFields) *AppRepositoryFields {
	f := &AppRepositoryFields{
		ClientRepository: postgres_repo.CreateClientPostgresRepository(fields),
		DoctorRepository: postgres_repo.CreateDoctorPostgresRepository(fields),
		PetRepository:    postgres_repo.CreatePetPostgresRepository(fields),
		RecordRepository: postgres_repo.CreateRecordPostgresRepository(fields),
	}

	a.Logger.Info("Success initialization of repositories")

	return f
}

func (a *App) initServices(r *AppRepositoryFields) *AppServiceFields {
	passwordHasher := hasherImplementation.NewBcryptHasher()

	u := &AppServiceFields{
		ClientService: servicesImplementation.NewClientServiceImplementation(r.ClientRepository, passwordHasher, a.Logger),
		DoctorService: servicesImplementation.NewDoctorServiceImplementation(r.DoctorRepository, passwordHasher, a.Logger),
		PetService:    servicesImplementation.NewPetServiceImplementation(r.PetRepository, r.ClientRepository, a.Logger),
		RecordService: servicesImplementation.NewRecordServiceImplementation(r.RecordRepository, r.DoctorRepository,
			r.ClientRepository, r.PetRepository, a.Logger),
	}

	a.Logger.Info("Success initialization of services")
	return u
}

func (a *App) Init() error {
	a.initLogger()
	
	a.InitJaeger("backend")

	fields, err := postgres_repo.CreatePostgresRepositoryFields(a.Config.Postgres, a.Logger)
	if err != nil {
		a.Logger.Fatal("Error create postgres repository fields", "err", err)
		return err
	}

	a.Repositories = a.initRepositories(fields)
	a.Services = a.initServices(a.Repositories)

	return nil
}

func (a *App) initLogger() {
	f, err := os.OpenFile(a.Config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Logger := log.New(f)

	log.SetFormatter(log.LogfmtFormatter)
	Logger.SetReportTimestamp(true)
	Logger.SetReportCaller(true)

	if a.Config.LogLevel == "debug" {
		Logger.SetLevel(log.DebugLevel)
	} else if a.Config.LogLevel == "info" {
		Logger.SetLevel(log.InfoLevel)
	} else {
		log.Fatal("Error log level")
	}

	Logger.Print("\n")
	Logger.Info("Success initialization of new Logger!")

	a.Logger = Logger
}

func (a *App) InitJaeger(serviceName string) (opentracing.Tracer, io.Closer) {
    cfg := config.Configuration{
        ServiceName: serviceName,
        Sampler: &config.SamplerConfig{
            Type:  "const",
            Param: 1,
        },
        Reporter: &config.ReporterConfig{
            LogSpans: true,
        },
    }

    tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
    if err != nil {
        a.Logger.Error("Could not initialize jaeger tracer:", err)
    }
    return tracer, closer
}

func (a *App) Run() error {
	err := a.Init()

	if err != nil {
		a.Logger.Error("Error init app", "err", err)
		return err
	}

	return nil
}
