# This is a filter that converts the output of autotune
# to an argument list for the setbasals command.
# Usage:
#   setbasals $(jq -r -f autotune.jq < autotune.json)
# Don't forget the -r (--raw-output) flag!

.[] | (.start | rtrimstr(":00")), .rate
