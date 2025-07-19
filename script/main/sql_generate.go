package main

import (
	"github.com/RussellLuo/kun/pkg/caseconv"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := "root:pi=3.1415@tcp(127.0.0.1:3306)/poetry?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn))
	config := gen.Config{
		OutPath:          "../../src/internal/infrastructure/repository", // 查询代码输出目录
		ModelPkgPath:     "../../src/internal/domain/entity",             // 模型代码输出目录
		Mode:             gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:    false, // 允许空字段生成指针类型
		FieldWithTypeTag: true,
		WithUnitTest:     true,
	}
	config.WithOpts(gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
		return caseconv.ToLowerCamelCase(columnName)
	},
	))

	// 初始化生成器
	g := gen.NewGenerator(config)
	g.UseDB(db) // 绑定数据库连接

	// 生成所有表
	g.GenerateAllTable()

	// 执行生成
	g.Execute()
}
