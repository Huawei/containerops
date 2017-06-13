/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	log "github.com/Sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/Huawei/containerops/dockyard/model"
	"github.com/Huawei/containerops/dockyard/setting"
)

// databasecmd is sub command which migrate/backup/restore Dockyard's database.
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database sub command migrate/backup/restore Dockyard's database.",
	Long:  ``,
}

// migrateDatabaseCmd is sub command migrate Dockyard's database.
var migrateDatabaseCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate sub command migrate Dockyard's database.",
	Long:  ``,
	Run:   migrateDatabase,
}

// backupDatabaseCmd is sub command backup Dockyard's database.
var backupDatabaseCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup sub command backup Dockyard's database.",
	Long:  ``,
	Run:   backupDatabase,
}

// restoreDatabaseCmd is sub command restore Dockyard's database.
var restoreDatabaseCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore sub command restore Dockyard's database.",
	Long:  ``,
	Run:   restoreDatabase,
}

// init()
func init() {
	RootCmd.AddCommand(databaseCmd)

	databaseCmd.AddCommand(migrateDatabaseCmd)
	databaseCmd.AddCommand(backupDatabaseCmd)
	databaseCmd.AddCommand(restoreDatabaseCmd)

	migrateDatabaseCmd.Flags().StringVarP(&configFilePath, "config", "c", "./conf/runtime.conf", "path of the config file.")
}

// migrateDatabase is auto-migrate database of Dockyard.
func migrateDatabase(cmd *cobra.Command, args []string) {
	if err := setting.SetConfig(configFilePath); err != nil {
		log.Fatalf("Failed to init settings: %s", err.Error())
		return
	}

	model.OpenDatabase(&setting.DBConfig)
	model.Migrate()
}

// backupDatabase is
func backupDatabase(cmd *cobra.Command, args []string) {

}

// restoreDatabase is
func restoreDatabase(cmd *cobra.Command, args []string) {

}
