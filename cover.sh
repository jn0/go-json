#!/bin/bash
set -e
rm -f cover.out
go test -coverprofile cover.out
sed -e 's,^_.*/\([^/]\+\)$,./\1,' -i cover.out
go tool cover -html=cover.out -o c.html
ls -l c.html
# EOF #
