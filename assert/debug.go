package assert

var Debug = &hide{}

//func (he *hide) Lan(v []string) {
//	he.lan = v
//}
//
//func (he *hide) Edition(v string) {
//	he.edition = v
//}
//
//func (he *hide) Hostname(v string) {
//	he.hostname = v
//}
//
//func (he *hide) VIP(v []string) {
//	he.vip = v
//}

type hide struct {
	Lan      []string
	Vip      []string
	Edition  string
	Hostname string
}
