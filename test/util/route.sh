#!/bin/bash

function route_ex {
  if [[ "$1" == '4' ]]
  then
    rts="$(ip route show | grep "via ${PROXY_HOST}")"
  else
    rts="$(ip -6 route show | grep "via ${PROXY_HOST}")"
  fi
  if grep -q "$2" <<< "$rts"
  then
    return 0
  else
    return 1
  fi
}

function route_add {
    if [[ -z "$1" ]] || [[ -z "$PROXY_HOST" ]]
    then
      return
    fi

    net="$1"
    rt="route add ${net} via ${PROXY_HOST}"

    if grep -q ':' <<< "$net" && grep -q ':' <<< "$PROXY_HOST"
    then
      if ! route_ex '6' "$net"
      then
        echo "Adding route: ${rt}"
        sudo ip -6 $rt
      fi
    else
      if ! route_ex '4' "$net"
      then
        echo "Adding route: ${rt}"
        sudo ip $rt
      fi
    fi
}

function route_rm {
    if [[ -z "$1" ]] || [[ -z "$PROXY_HOST" ]]
    then
      return
    fi

    net="$1"
    rt="route del ${net} via ${PROXY_HOST}"

    if grep -q ':' <<< "$net" && grep -q ':' <<< "$PROXY_HOST"
    then
      if route_ex '6' "$net"
      then
        echo "Removing route: ${rt}"
        sudo ip -6 $rt
      fi
    else
      if route_ex '4' "$net"
      then
        echo "Removing route: ${rt}"
        sudo ip $rt
      fi
    fi
}

