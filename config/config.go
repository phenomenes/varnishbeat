// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Varnishbeat VarnishbeatConfig
}

type VarnishbeatConfig struct {
	Period string `config:"period"`
}
