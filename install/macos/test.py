#!/usr/bin/env python
import sys
import urllib
arg = sys.argv[1]
handler, fullPath = arg.split(":", 1)
path, fullArgs = fullPath.split("?", 1)
action = path.strip("/")
args = fullArgs.split("&")
params = {}
for arg in args:
    key, value = map(urllib.unquote, arg.split("=", 1))
    params[key] = value