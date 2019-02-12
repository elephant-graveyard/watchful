#!/usr/bin/env bash

set -euo pipefail # The hook should fail if the make targets are not executed correctly
make analysis test # Test the actual code using analysis and test targets
