package server

type Session struct {
	variableMap map[string]string
}

func NewSession() *Session {
	return &Session{
		variableMap: make(map[string]string),
	}
}

func (self *Session) GetVariable(name string) string {
	return self.variableMap[name]
}

func (self *Session) SetVariable(name string, value string) {
	self.variableMap[name] = value
}
