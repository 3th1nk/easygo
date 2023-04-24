package logs

type Chain []Logger

func (this Chain) GetLevel() int {
	if len(this) != 0 {
		return this[0].GetLevel()
	}
	return 0
}

func (this Chain) SetLevel(level int) {
	for _, v := range this {
		v.SetLevel(level)
	}
}

func (this Chain) Error(format string, a ...interface{}) {
	for _, v := range this {
		v.Error(format, a...)
	}
}

func (this Chain) Warn(format string, a ...interface{}) {
	for _, v := range this {
		v.Warn(format, a...)
	}
}

func (this Chain) Info(format string, a ...interface{}) {
	for _, v := range this {
		v.Info(format, a...)
	}
}

func (this Chain) Debug(format string, a ...interface{}) {
	for _, v := range this {
		v.Debug(format, a...)
	}
}
