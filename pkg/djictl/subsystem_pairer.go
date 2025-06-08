package djictl

type SubsystemPairer Device

func (d *Device) Pairer() *SubsystemPairer {
	return (*SubsystemPairer)(d)
}

func (s *SubsystemPairer) SubsystemID() SubsystemID {
	return SubsystemIDPairer
}

func (s *SubsystemPairer) Device() *Device {
	return (*Device)(s)
}
