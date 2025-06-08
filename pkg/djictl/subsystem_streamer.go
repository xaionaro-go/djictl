package djictl

type SubsystemStreamer Device

func (d *Device) Streamer() *SubsystemStreamer {
	return (*SubsystemStreamer)(d)
}

func (s *SubsystemStreamer) Device() *Device {
	return (*Device)(s)
}

func (s *SubsystemStreamer) SubsystemID() SubsystemID {
	return SubsystemIDStreamer
}
