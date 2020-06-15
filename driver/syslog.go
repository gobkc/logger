package driver

import "log/syslog"

const (
	SYS_LOC = iota
	SYS_UDP
	SYS_TCP
)

//syslog driver
type Syslog struct {
	Server   string
	Protocol int
	Tag      string
}

//系统日志驱动
func (s *Syslog) Write(p []byte) (n int, err error) {
	var sysLog *syslog.Writer
	sysLog, err = syslog.Dial("", "", syslog.LOG_INFO, s.Tag)
	if err != nil {
		return n, err
	}
	sysLog.Emerg(string(p))
	n = len(string(p))
	return n, err
}
