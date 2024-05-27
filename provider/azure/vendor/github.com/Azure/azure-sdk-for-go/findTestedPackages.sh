# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

dirname $(find | grep _test.go | grep -v vendor) | sort -u