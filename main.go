package main

import (
	"flag"
)

func main() {

	cfile := flag.String("cfile", "", "specify path to the configuration file")
	lfok := flag.Bool("f", false, "enable http logging into file")
	ltok := flag.Bool("t", false, "enable http logging into terminal")
	dbok := flag.Bool("d", false, "enable logging db queries")
	clear := flag.Bool("c", false, "clear the database beforehand")
	flag.Parse()

	app := Initiate(*cfile, *clear)

	app.Run(*lfok, *ltok, *dbok)
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
