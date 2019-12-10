package router

import (
	"hr-server/handler/api/record"
	"hr-server/handler/api/statistics"
	"hr-server/pkg/auth"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "hr-server/docs" // docs is generated by Swag CLI, you have to import it.
	"hr-server/handler/api/audit"
	"hr-server/handler/api/group"
	transfer "hr-server/handler/api/grouptransfer"
	"hr-server/handler/api/permission"
	"hr-server/handler/api/profile"
	"hr-server/handler/api/role"
	"hr-server/handler/api/salary"
	"hr-server/handler/api/sd"
	"hr-server/handler/api/tag"
	"hr-server/handler/api/templateaccount"
	"hr-server/handler/api/tool"
	"hr-server/handler/api/user"
	"hr-server/handler/api/usergroup"
	"hr-server/router/middleware"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	//g.LoadHTMLGlob("templates/*")

	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	// g.Use(middleware.Options)
	// g.Use(middleware.Secure)
	g.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	g.StaticFS("/api/download", http.Dir("./export"))
	g.StaticFS("/api/upload", http.Dir("./upload"))
	g.StaticFS("/api/backup", http.Dir("./backup"))
	g.StaticFS("/api/static", http.Dir("./assets"))
	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//g.GET("/api", handler.Home)
	// pprof router
	//pprof.Register(g)

	// api for authentication functionalities
	g.POST("/api/login", user.Login)

	e := auth.GetEnforcer("./conf/authz_model.conf", "./conf/authz_policy.csv")
	// The user handlers, requiring authentication
	u := g.Group("/api/v1/user")
	u.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		u.GET("", user.List)
		u.GET("/:username", user.Get)
		u.POST("", user.Create)
		u.POST("/freeze", user.Freeze)
		u.POST("/active", user.Active)
		u.POST("/password/reset", user.ResetPassword)
		u.POST("/password/change", user.ChangePassword)
		u.PUT("/:id", user.Update)
		u.DELETE("/:id", user.Delete)
	}

	ugPublic := g.Group("/api/v1/usergroup")
	ugPublic.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		ugPublic.GET("", usergroup.List)
		ugPublic.GET("/:id/users", usergroup.UserList)
	}

	ugPrivate := g.Group("/api/v1/usergroup")
	ugPrivate.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		ugPrivate.POST("", usergroup.Create)
		ugPrivate.POST("/:id/users", usergroup.RelateUsers)
		ugPrivate.POST("/:id/users/remove", usergroup.RemoveRelateUsers)
		ugPrivate.PUT("/:id", usergroup.Update)
		ugPrivate.DELETE("/:id", usergroup.Delete)
	}

	rolePublic := g.Group("/api/v1/role")
	rolePublic.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		rolePublic.GET("", role.List)
		rolePublic.GET("/:id/users", role.UserList)
	}

	rolePrivate := g.Group("/api/v1/role")
	rolePrivate.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		rolePrivate.POST("", role.Create)
		rolePrivate.POST("/:id/users", role.RelateUsers)
		rolePrivate.POST("/:id/users/remove", role.RemoveRelateUsers)
		rolePrivate.PUT("/:id", role.Update)
		rolePrivate.DELETE("/:id", role.Delete)
	}
	p := g.Group("/api/v1/profile")
	p.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	// u.Use(middleware.AuthMiddleware())
	{
		p.GET("", profile.List)
		p.GET("/:id", profile.Get)
		p.GET("/:id/detail", profile.GetDetail)
		p.GET("/:id/transfer", profile.GetTransfer)
		p.GET("/:id/tags", profile.ListProfileTags)
		p.PUT("/:id/tags", profile.RelateTags)
		p.POST("", profile.Create)
		p.POST("/transfer", transfer.Create)
		//p.POST("/delete", profile.Delete)
		p.POST("/unfreeze", profile.UnFreeze)
		p.POST("/freeze", profile.Freeze)
		p.PUT("/:id", profile.Update)
	}

	in := g.Group("/api/v1/import")
	in.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		in.POST("/profile", profile.Import)
		in.POST("/tag", tag.Import)
		in.POST("/group", group.Import)
		in.POST("/group/tags", group.ImportTags)
		in.POST("/salary", salary.Import)
	}

	gPublic := g.Group("/api/v1/group")
	gPublic.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		gPublic.GET("", group.List)
		gPublic.GET("/:id", group.Get)
		gPublic.GET("/:id/users", group.UserList)
		gPublic.GET("/:id/profiles", group.ProfileList)
	}

	gPrivate := g.Group("/api/v1/group")
	gPrivate.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		gPrivate.GET("/:id/lock", group.Lock)
		gPrivate.GET("/:id/unlock", group.UnLock)
		gPrivate.GET("/:id/invalid", group.Invalid)
		gPrivate.GET("/:id/valid", group.Valid)
		gPrivate.POST("", group.Create)
		gPrivate.POST("/:id/profiles", group.RelateProfiles)
		gPrivate.POST("/:id/profiles/remove", group.RemoveRelateProfiles)
		gPrivate.PUT("/:id", group.Update)
		gPrivate.PUT("/:id/move", group.Move)
		gPrivate.DELETE("/:id", group.Delete)
		gPrivate.PUT("/:id/tags", group.RelateTags)
	}

	tagPublic := g.Group("/api/v1/tag")
	tagPublic.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		tagPublic.GET("", tag.List)
		tagPublic.GET("/:id/child", tag.GetChild)
		tagPublic.GET("/:id/users", tag.RelatedUserList)
		tagPublic.GET("/:id/profiles", tag.ProfileList)
	}

	tagPrivate := g.Group("/api/v1/tag")
	tagPrivate.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		tagPrivate.POST("", tag.Create)
		tagPrivate.PUT("/:id", tag.Update)
		tagPrivate.POST("/:id/profiles", tag.RelateProfiles)
		tagPrivate.POST("/:id/profiles/remove", tag.RemoveRelateProfiles)
		// tagPrivate.POST("/delete", tag.DeleteList)
		tagPrivate.DELETE("/:id", tag.Delete)
	}

	salaryPublic := g.Group("/api/v1/salary")
	salaryPublic.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		salaryPublic.POST("/calculate", salary.Calculate)
		salaryPublic.POST("/template/config", salary.TemplateConfig)
		salaryPublic.POST("/template/order", salary.TemplateOrder)
		salaryPublic.POST("/tax/config", salary.TaxSetting)
		salaryPublic.GET("/tax/config", salary.GetTaxSetting)
		salaryPublic.POST("/pre_deduction_rate/config", salary.PreDeductionRateSetting)
		salaryPublic.GET("/pre_deduction_rate/config", salary.GetPreDeductionRateSetting)
		salaryPublic.POST("/base/config", salary.BaseConfig)
		salaryPublic.POST("/profile/config", salary.SalaryProfileConfig)
		salaryPublic.GET("/profile/config", salary.GetSalaryProfileConfigList)
		salaryPublic.DELETE("/profile/config/:id", salary.DeleteSalaryProfileConfig)
		salaryPublic.GET("/base", salary.GetBaseSalary)
		salaryPublic.POST("/upload", salary.Upload)
		salaryPublic.GET("/template", salary.ListTemplate)
		salaryPublic.GET("/template/:id", salary.GetTemplate)
		salaryPublic.GET("/template/:id/audit", salary.GetAuditTemplate)
		salaryPublic.DELETE("/template/:id", salary.DeleteTemplate)
		salaryPublic.GET("/account", templateaccount.List)
		salaryPublic.GET("/account/:id", templateaccount.Get)
		salaryPublic.GET("/year/:year/month/:month/profile/:id", salary.GetProfileMonthSalary)
		salaryPublic.GET("/account/:id/template", templateaccount.GetAccountTemplate)
		salaryPublic.GET("/account/:id/year/:year/fields", templateaccount.GetAccountFields)
		salaryPublic.DELETE("/account/:id", templateaccount.Delete)
		salaryPublic.POST("/account", templateaccount.Create)
		salaryPublic.GET("/export", salary.Export)
		salaryPublic.GET("/payroll", salary.GetPayroll)
	}

	tc := g.Group("/api/v1/templateaccount")
	tc.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		tc.GET("/:id/templates", templateaccount.ListTemplateWithFields)
	}

	r := g.Group("/api/v1/record")
	tc.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		r.GET("/list", record.List)
		r.GET("/transfer", transfer.List)
		r.GET("/operation", record.OperationList)
	}

	per := g.Group("/api/v1/permission")
	per.Use(middleware.AuthMiddleware())
	{
		per.GET("/role/:id", permission.Get)
		per.POST("", permission.Create)
	}

	pkgPublic := g.Group("/api/v1/tool")
	//pkgPublic.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		pkgPublic.GET("/func", tool.ListFunc)
		pkgPublic.GET("/backup/files", tool.ListBackupFiles)
		pkgPublic.GET("/backup", tool.Backup)
	}

	statisticsPublic := g.Group("/api/v1/statistics")
	salaryPublic.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		statisticsPublic.GET("/annual_income/employee", statistics.EmployeeAnnualIncome)
		statisticsPublic.GET("/annual_income/department", statistics.DepartmentAnnualIncome)
		statisticsPublic.POST("/query/detail", statistics.DetailQuery)
		statisticsPublic.POST("/query/department/income", statistics.DepartmentIncomeQuery)
	}

	auditPrivate := g.Group("/api/v1/audit")
	auditPrivate.Use(middleware.AuthMiddleware(), middleware.Authority(e))
	{
		auditPrivate.GET("/list", audit.List)
		auditPrivate.POST("/:id", audit.UpdateState)
	}
	// The health check handlers
	svcd := g.Group("/api/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}