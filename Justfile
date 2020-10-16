check:
	namcap -i PKGBUILD

build:
	go build
	cp salamis package/usr/bin

package:
	dpkg-deb --build package salamis.deb
