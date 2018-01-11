package main

import (
	"flag"
)

const CONFFILE = "./conf.json"

func main() {

	cfile := flag.String("cfile", CONFFILE, "specify path to the configuration file")
	lenabled := flag.Bool("l", false, "enable app logging")
	clear := flag.Bool("c", false, "clear the database beforehand")
	flag.Parse()

	app := Initiate(*cfile, *lenabled, *clear)

	defer app.Exit()

	app.Run()
}

// post1 := Post{
// 	Name: "Update",
// 	Body: "sudo ionice -c3 sudo swapoff -a&& sudo swapon -a",
// 	Tags: []Tag{
// 		Tag{Name: "mysql"},
// 		Tag{Name: "ntp"}},
// }
//
// post2 := Post{
// 	Name: "convert crt to pem",
// 	Body: "openssl x509 -in softcall.me.crt -out domain_cert_tls.softcall.me.pem -outform PEM",
// 	Tags: []Tag{
// 		Tag{Name: "certs"},
// 		Tag{Name: "ssl"}},
// }
//
// app.DB.Save(&post1)
// app.DB.Save(&post2)
