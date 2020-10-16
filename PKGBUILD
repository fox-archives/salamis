# Maintainer: Edwin Kofler <>

pkgname=salamis
pkgver=0.2.3
pkgrel=1
pkgdesc="Manage vscode extensions"
arch=('x86_64')
url="https://github.com/eankeen/salamis"
license=('')
depends=()
makedepends=()
source=(https://github.com/eankeen/salamis/releases/download/v$pkgver/salamis_${pkgver}_Linux_x86_64.tar.gz)
sha256sums=('4e4364e614f5d966d7fe3c9b3c271c68f69882fb861a5236c995999427ceba1a')

build() {
	cd "$srcdir/$pkgname-$pkgver"

  	just --no-dotenv build
}

package() {
	cd "$srcdir/$pkgname-$pkgver"

	install -Dm0755 "${pkgdir}/usr/bin/salamis"
}
