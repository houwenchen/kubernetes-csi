package config

type Config struct {
	DriverName    string //必须要有的，GetPluginInfo 会用到
	EndPoint      string
	NodeID        string
	VendorVersion string //必须要有的，GetPluginInfo 会用到

	VolumeDir string

	EnableLVM bool
}
