
pkgname=ncp
pkgver=1.0.0.internal.0
pkgrel=1

arch=armhf

mkdir -p debian

cat > debian/rules << EOF
#!/usr/bin/make -f
%:
	dh \$@

EOF

# libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev gstreamer1.0-plugins-good
cat > debian/control << EOF
Source: ncp
Maintainer: a-wing <1@233.email>
Build-Depends: debhelper (>= 8.0.0), golang (>= 1.11)
Standards-Version: 3.9.3
Section: utils

Package: ncp
Priority: extra
Depends: gstreamer1.0-plugins-good (>= 1.14)
Architecture: ${arch}
Description: Node control protocol
EOF

cat > debian/ncp.install << EOF
ncp usr/lib/${pkgname}/
scripts usr/lib/${pkgname}/
conf/ncp.service lib/systemd/system/
conf/ncp@.service lib/systemd/system/
conf/config-dist.yml etc/ncp/
EOF

cat > debian/ncp.links << EOF
usr/lib/${pkgname}/ncp /usr/bin/ncp
EOF

cat > debian/changelog << EOF
${pkgname} (${pkgver}-${pkgrel}) unstable; urgency=low

  * Initial release.

 -- a-wing <1@233.email>  Sun, 07 Apr 2019 22:07:56 +0800
EOF

echo 9 > debian/compat

# repo golang 1.11.x   Need use golang 1.13.x
debuild --prepend-path=`which go | sed s'#/go$##'` -us -uc

rm -r debian

