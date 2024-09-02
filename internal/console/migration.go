package console

import (
	"database/sql"
	"log"

	"github.com/tubagusmf/payment-service-gb1/internal/config"
	"github.com/tubagusmf/payment-service-gb1/internal/helper"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var (
	direction string
	step      int = 1
)

func init() {
	rootCMd.AddCommand(migrationCmd)

	migrationCmd.Flags().StringVarP(&direction, "direction", "d", "up", "Migration direction")
	migrationCmd.Flags().IntVarP(&step, "step", "s", 1, "Migration step")
}

var migrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Run:   migrateDB,
}

func migrateDB(cmd *cobra.Command, args []string) {
	config.LoadWithViper()

	connDB, err := sql.Open("mysql", helper.GetConnectionString())
	if err != nil {
		log.Panicf("Gagal koneksi ke database: %s", err.Error())
	}
	defer connDB.Close()

	migrations := &migrate.FileMigrationSource{Dir: "./db/migrations"}

	var n int
	if direction == "down" {
		n, err = migrate.ExecMax(connDB, "mysql", migrations, migrate.Down, step)
	} else {
		n, err = migrate.ExecMax(connDB, "mysql", migrations, migrate.Up, step)
	}
	if err != nil {
		log.Panicf("Gagal migrasi database: %s", err.Error())
	}

	log.Printf("Sukses melakukan migrations %d", n)
}
