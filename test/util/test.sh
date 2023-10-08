#!/bin/bash

function curlRc {
  curl --connect-timeout 3 --fail "$1"
  return "$?"
}
