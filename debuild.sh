
pkgname=ncp
pkgver=0.1.0
arch=arm
#arch=amd64

srcdir=${pkgname}_${pkgver}_linux_${arch}

GOOS=linux GOARCH=${arch} go build


bundler install --path=vendor

arch=armhf

mkdir -p debian

cat > debian/rules << EOF
#!/usr/bin/make -f
%:
	dh \$@

EOF

cat > debian/control << EOF
Source: ncp
Maintainer: a-wing <1@233.email>
Build-Depends: debhelper (>= 8.0.0), golang (>= 1.11)
Standards-Version: 3.9.3
Section: utils

Package: ncpgo
Priority: extra
Architecture: ${arch}
Description: Go Node control protocol

Package: ncpcmd
Priority: extra
Architecture: ${arch}
Depends: ncpgo, ruby, ruby-bundler
Description: Node control protocol
EOF

cat > debian/ncpgo.install << EOF
ncpgo usr/bin/
ncpgo.service lib/systemd/system/
config-dist.yml etc/ncp/
EOF

cat > debian/ncpcmd.install << EOF
ncp.rb usr/lib/ncp/
Gemfile* usr/lib/ncp/
lib usr/lib/ncp/
scripts usr/lib/ncp/
log usr/lib/ncp/
vendor usr/lib/ncp/vendor
.bundle usr/lib/ncp/

ncp.service lib/systemd/system/
EOF

cat > debian/changelog << EOF
${pkgname} (${pkgver}-0) unstable; urgency=low

  * Initial release.

 -- a-wing <1@233.email>  Sun, 07 Apr 2019 22:07:56 +0800
EOF

echo 9 > debian/compat

debuild -us -uc

rm -r debian

