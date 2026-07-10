package tools

type Milliseconds int64

func (m *Milliseconds) AddSeconds(seconds int64) {
	*m += Milliseconds(seconds * 1000)
}

func (m *Milliseconds) AddMinutes(minutes int64) {
	*m += Milliseconds(minutes * 60 * 1000)
}

func (m *Milliseconds) AddHoursMinutes(hours int64) {
	*m += Milliseconds(hours * 60 * 60 * 1000)
}

func (m *Milliseconds) LessThanMinutes(minutes int64) bool {
	return *m < Milliseconds(minutes*60*1000)
}

func (m *Milliseconds) ToSeconds() Seconds {
	return Seconds(*m / 1000)
}

type Seconds int64

func (m *Seconds) AddSeconds(seconds int64) {
	*m += Seconds(seconds)
}

func (m *Seconds) AddMinutes(minutes int64) {
	*m += Seconds(minutes * 60)
}

func (m *Seconds) AddHoursMinutes(hours int64) {
	*m += Seconds(hours * 60 * 60)
}

func (m *Seconds) ToMilis() Milliseconds {
	return Milliseconds(*m * 1000)
}
