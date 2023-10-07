#!/bin/bash

function route_ex {
  if [[ "$1" == '4' ]]
  then
    rts="$(ip route show)"
  else
    rts="$(ip -6 route show)"
  fi
  if grep -q "$2" <<< "$rts"
  then
    return 1
  else
    return 0
  fi
}

function route_add {
    net="$1"
    rt="route add $net via $PROXY_HOST"

    if grep -q ':' <<< "$net" && grep -q ':' <<< "$PROXY_HOST"
    then
      if ! route_ex '6' "$rt"
      then
        ip -6 $rt
      fi
    else
      if ! route_ex '4' "$rt"
      then
        ip $rt
      fi
    fi
}

function route_rm {
    net="$1"
    rt="route del $net via $PROXY_HOST"

    if grep -q ':' <<< "$net" && grep -q ':' <<< "$PROXY_HOST"
    then
      if route_ex '6' "$rt"
      then
        ip -6 $rt
      fi
    else
      if route_ex '4' "$rt"
      then
        ip $rt
      fi
    fi
}

