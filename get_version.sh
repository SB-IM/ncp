
HASH=`git log | head -1 | cut -d " " -f 2`
TAG=$1
VERSION=${TAG:-${HASH}}

cat > version.go << EOF
package main
var version = "${VERSION}"
EOF
