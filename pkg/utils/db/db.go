package utils

import (
	"database/sql"
	"fmt"
	conf "github.com/JMURv/effectiveMobile/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"go.uber.org/zap"
	"path/filepath"
	"strconv"
	"strings"
)

func ApplyMigrations(db *sql.DB, conf *conf.DBConfig) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	path, _ := filepath.Abs("db/migration")
	path = filepath.ToSlash(path)
	m, err := migrate.NewWithDatabaseInstance("file://"+path, conf.Database, driver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	zap.L().Info("Migrations applied successfully")
	return nil
}

func BuildFilterQuery(filters map[string]any) (string, []any) {
	conds := make([]string, 0, len(filters))
	args := make([]any, 0, len(filters))

	for key, value := range filters {
		newArg := strconv.Itoa(len(args) + 1)
		switch key {
		case "group":
			conds = append(conds, "group_name ILIKE $"+newArg)
			args = append(args, "%"+value.(string)+"%")
		case "song":
			conds = append(conds, "song_name ILIKE $"+newArg)
			args = append(args, "%"+value.(string)+"%")
		case "min_release_date":
			conds = append(conds, "release_date >= $"+newArg)
			args = append(args, value)
		case "max_release_date":
			conds = append(conds, "release_date <= $"+newArg)
			args = append(args, value)
		case "release_date":
			conds = append(conds, "release_date = $"+newArg)
			args = append(args, value)
		case "link":
			conds = append(conds, "link ILIKE $"+newArg)
			args = append(args, "%"+value.(string)+"%")
		}
	}

	var q strings.Builder
	if len(conds) > 0 {
		q.WriteString(" WHERE ")
		q.WriteString(strings.Join(conds, " AND "))
	}

	return q.String(), args
}
