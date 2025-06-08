package djictl

type SubsystemConfigurer Device

func (d *Device) Configurer() *SubsystemConfigurer {
	return (*SubsystemConfigurer)(d)
}

func (s *SubsystemConfigurer) SubsystemID() SubsystemID {
	return SubsystemIDConfigurer
}

func (s *SubsystemConfigurer) Device() *Device {
	return (*Device)(s)
}
