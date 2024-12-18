package data

import (
	"context"
	"os"
	"slices"

	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/conf"
	u_config "gitlab.calendaria.team/services/utils/v1/config"
	u_jwt "gitlab.calendaria.team/services/utils/v2/jwt"
	u_nats "gitlab.calendaria.team/services/utils/v2/nats"
	u_tracing "gitlab.calendaria.team/services/utils/v2/tracing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	_ "gitlab.calendaria.team/services/rbac/ent/runtime"
)

// Hidden role ids.
const (
	AdminRoleID = 1
	BasicRoleID = 2
)

// ProviderSet is data providers.
//
//nolint:gochecknoglobals // this global variable is required for wire
var ProviderSet = wire.NewSet(
	NewData,
	NewNatsClient,
	u_config.NewConfig,
	u_jwt.NewJwtProcessor,
	u_nats.NewQueueManager,
	NewRoleRepo,
	NewTeamsRepo,
	NewPermissionRepo,
	NewAssignedRolesRepo,
	u_tracing.NewTracer,
)

// Data .
type Data struct {
	log *log.Helper
	db  *ent.Client
}

// NewData .
func NewData(bc *conf.Bootstrap, c *u_config.Config, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(logger)

	dbDsn := bc.GetDb() // read from local config
	if dbDsn == "" {
		// read from vault
		secret, err := c.ReadSecretsFor(context.Background(), "db-dsn")
		if err != nil {
			l.Fatalf("db dsn not found: %v", err)
			return nil, nil, err
		}
		dsn, ok := secret["data"].(string)
		if !ok {
			l.Fatalf("db dsn not found")
			return nil, nil, err
		}
		dbDsn = dsn
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

// ExtractUnique extracts a slice of R from a slice of E using a function that returns
// a value and a bool. If the bool is true, the value is appended to the result slice.
// The result slice contains only unique values.
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
//		{Name: "Bob", Age: 40},
//	}
//
//	names := extractSlice(persons, func(p Person) (string, bool) {
//		return p.Name, true
//	})
//
//	fmt.Println(names) // Output: [Alice Bob]
//
// The function is generic and works with any slice of any type.
func ExtractUnique[S ~[]E, E, R comparable](slice S, extract func(E) (R, bool)) []R {
	result := make([]R, 0, len(slice))
	uniques := make(map[R]struct{}, len(slice))
	for _, item := range slice {
		if value, ok := extract(item); ok {
			if _, ok = uniques[value]; !ok {
				uniques[value] = struct{}{}
				result = append(result, value)
			}
		}
	}
	return result
}

func Diff[T comparable, S ~[]T](slice S, other S) []T {
	var result []T
	for _, item := range slice {
		if !slices.Contains(other, item) {
			result = append(result, item)
		}
	}
	return result
}

func Filter[S ~[]E, E any](slice S, filter func(E) bool) S {
	var result []E
	for _, item := range slice {
		if filter(item) {
			result = append(result, item)
		}
	}

	return result
}
