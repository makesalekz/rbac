package data

import (
	"context"
	"os"

	"gitlab.calendaria.team/services/rbac/ent" //nolint:typecheck
	"gitlab.calendaria.team/services/rbac/internal/conf"
	u_config "gitlab.calendaria.team/services/utils/v1/config"
	u_jwt "gitlab.calendaria.team/services/utils/v1/jwt"
	u_tracing "gitlab.calendaria.team/services/utils/v2/tracing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	_ "gitlab.calendaria.team/services/rbac/ent/runtime"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	u_config.NewConfig,
	u_jwt.NewJwtProcessor,
	u_tracing.NewTracer,
	NewNatsClient,
	NewRoleRepo,
	NewTeamsRepo,
	NewPermissionRepo,
	NewAssignedRolesRepo,
)

// Data .
type Data struct {
	log *log.Helper
	db  *ent.Client
}

// NewData .
func NewData(bc *conf.Bootstrap, c *u_config.Config, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(logger)

	dbDsn := bc.Db // read from local config
	if dbDsn == "" {
		// read from vault
		secret, err := c.ReadSecretsFor(context.Background(), "db-dsn")
		if err != nil {
			l.Fatalf("db dsn not found: %v", err)
			return nil, nil, err
		}
		dbDsn = secret["data"].(string)
	}

	autoMigrate := os.Getenv("AUTOMIGRATE")
	entLogging := os.Getenv("ENT_LOGGING")
	var options []ent.Option
	if entLogging != "" {
		options = append(options, ent.Debug(), ent.Log(l.Debug))
	}

	client, err := ent.Open("postgres", dbDsn, options...)
	if err != nil {
		l.Fatalf("failed opening connection to postgres: %v", err)
		return nil, nil, err
	}

	if autoMigrate != "" {
		if err := client.Schema.Create(context.Background()); err != nil {
			l.Errorf("failed creating schema resources: %v", err)
			return nil, nil, err
		}
	}

	l.Info("Connected to postgres")

	cleanup := func() {
		if err := client.Close(); err != nil {
			l.Error(err)
		}
	}

	return &Data{
		log: log.NewHelper(logger),
		db:  client,
	}, cleanup, nil
}

// extractSlice extracts a slice of R from a slice of E using a function that returns
// a value and a bool. If the bool is true, the value is appended to the result slice.
//
// Example:
//
//	type Person struct {
//		Name string
//		Age  int
//	}
//
//	persons := []Person{
//		{Name: "Alice", Age: 30},
//		{Name: "Bob", Age: 35},
//	}
//
//	names := extractSlice(persons, func(p Person) (string, bool) {
//		return p.Name, true
//	})
//
//	fmt.Println(names) // Output: [Alice Bob]
//
// The function is generic and works with any slice of any type.
func ExtractSlice[S ~[]E, E, R any](slice S, extract func(E) (R, bool)) []R {
	result := make([]R, 0, len(slice))
	for _, item := range slice {
		if value, ok := extract(item); ok {
			result = append(result, value)
		}
	}
	return result
}
