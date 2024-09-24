package configuration

type MetaConfiguration map[string]Configuration

type ConfigurationUpdater func(Configuration) Configuration

var metaConfiguration *MetaConfiguration = nil

func makeMetaConfiguration() *MetaConfiguration {
	cfg := make(MetaConfiguration)
	for _, env := range []Env{Local, Testing, Production} {
		islice := InputSlice{Env: env}
		cfg[islice.String()] = Configuration{}
	}
	return &cfg
}

func GetMetaConfiguration() *MetaConfiguration {
	if metaConfiguration == nil {
		metaConfiguration = makeMetaConfiguration()
	}
	return metaConfiguration
}

func (mcfg *MetaConfiguration) SetSlice(islice InputSlice, cfg Configuration) *MetaConfiguration {
	(*mcfg)[islice.String()] = cfg
	return mcfg
}

func (mcfg *MetaConfiguration) InheritSlice(to InputSlice, from InputSlice, updater ConfigurationUpdater) *MetaConfiguration {
	(*mcfg)[to.String()] = updater((*mcfg)[from.String()])
	return mcfg
}

func (mcfg *MetaConfiguration) GetConfiguration(islice InputSlice) Configuration {
	return (*mcfg)[islice.String()]
}
