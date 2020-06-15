package driver

//mysql driver
type Mysql struct {
	Server   string
	Port     int
	User     string
	Password string
	Table    string
}

//mysql日志驱动
func (m *Mysql) Write(p []byte) (n int, err error) {
	return n, err
}
