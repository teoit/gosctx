# gosctx
golang service context


# use gosctx
```base
var (
	serviceName = "server-service"
	version     = "1.0.0"
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName(serviceName),
		sctx.WithComponent(fiberapp.NewFiber(configs.KeyCompFIBER)),
		sctx.WithComponent(mongodb.NewMongoDB(configs.KeyCompMongoDB, "")),
		sctx.WithComponent(jwtc.NewJWT(configs.KeyCompJWT)),
		sctx.WithComponent(sctx.NewAppLoggerDaily(configs.KeyLoggerDaily)),
		sctx.WithComponent(discord.NewDiscordClient(configs.KeyDiscordSMS)),
	)
}
```