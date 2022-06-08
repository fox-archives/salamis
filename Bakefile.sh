# shellcheck shell=bash

task.build() {
	go build
	cp salamis package/usr/bin
}

task.package() {
	dpkg-deb --build package salamis.deb
}

task.release () {
	namcap -i PKGBUILD
}
