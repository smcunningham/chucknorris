package api

//API interface for creating new service
type API interface {
	NewService(rstCfg Config)
}
