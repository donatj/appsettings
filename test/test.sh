#!/bin/bash

set -e

BASEDIR=$(dirname "$0")
BIN="$BASEDIR/appsettings"
CFG="$BASEDIR/cfg.json"

go build -o "$BIN" "$BASEDIR"/../cmd/appsettings

if $BIN -file "$CFG" get "soupx" 2> /dev/null ; then
  echo "Error: Undefined Key Should Fail"
  exit 1
else
  echo "OK"
fi

$BIN -file "$CFG" set "foo" "bar"
if [ "$($BIN -file "$CFG" get "foo")" != "bar" ]; then
  echo "Error looking up root key"
  exit 1
else
    echo "OK"
fi

$BIN -file "$CFG" delete "foo"
if $BIN -file "$CFG" get "foo" 2> /dev/null ; then
  echo "Error: Foo should be deleted"
  exit 1
else
  echo "OK"
fi

$BIN -file "$CFG" set "foo.bar.baz.buzz.bigly" "down under"

if [ "$($BIN -file "$CFG" get "foo.bar.baz.buzz.bigly")" != "down under" ]; then
  echo "Error looking up nested key"
  exit 1
else
    echo "OK"
fi
