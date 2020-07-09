package migrate

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-admin/database"
	orm "go-admin/global"
	"go-admin/models"
	"go-admin/models/gorm"
	"go-admin/tools"
	"go-admin/tools/config"

	"github.com/spf13/cobra"
)

var (
	configYml string
	mode      string
	StartCmd  = &cobra.Command{
		Use:   "init",
		Short: "initialize the database",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func run() {
	usage := `start init`
	fmt.Println(usage)
	//1. 读取配置
	config.ConfigSetup(configYml)
	//2. 设置日志
	tools.InitLogger()
	//3. 初始化数据库链接
	database.Setup(config.DatabaseConfig.Driver)

	//4. 数据库迁移
	_ = migrateModel()
	log.Println("数据库结构初始化成功！")
	//5. 数据初始化完成
	if err := models.InitDb(); err != nil {
		log.Fatal("数据库基础数据初始化失败！")
	}
	usage = `数据库基础数据初始化成功`
	fmt.Println(usage)
}

func migrateModel() error {
	if config.DatabaseConfig.Driver == "mysql" {
		orm.Eloquent = orm.Eloquent.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	return gorm.AutoMigrate(orm.Eloquent)
}