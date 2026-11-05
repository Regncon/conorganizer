#!/usr/bin/env bash
set -eu
bash "$(dirname "$0")/files/make-fixture-repo.sh" fixture-repo uncommitted-retry
