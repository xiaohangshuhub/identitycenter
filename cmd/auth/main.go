package main

import (
	_ "github.com/xiaohangshuhub/xiaohangshu/api/auths/docs" // swagger 一定要有这行,指向你的文档地址
	"github.com/xiaohangshuhub/xiaohangshu/internal/webapi"

	"github.com/xiaohangshuhub/go-workit/pkg/webapp"
)

func main() {

	// 创建 Web 主机构建器
	builder := webapp.NewBuilder()

	// 配置依赖注入
	builder.AddServices(webapi.DependencyInjection()...)

	// 构建 web 服务
	app := builder.Build()

	// 配置swagger
	if app.Env().IsDevelopment {
		app.UseSwagger()
	}

	app.UseStaticFiles("/assets", "../../web/dist/assets")

	// 配置路由
	app.MapRoute(webapi.EndPointList...)

	// 启动应用
	app.Run()
}
