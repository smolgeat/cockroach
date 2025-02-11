#!/usr/bin/env bash

set -xeuo pipefail

dir="$(dirname $(dirname $(dirname $(dirname "${0}"))))"
source "$dir/teamcity-support.sh"
source "$dir/teamcity-bazel-support.sh"

tc_start_block "Run compose tests"

bazel build //pkg/cmd/bazci --config=ci
BAZEL_BIN=$(bazel info bazel-bin --config=ci)
BAZCI=$BAZEL_BIN/pkg/cmd/bazci/bazci_/bazci

bazel build //pkg/cmd/cockroach //pkg/compose/compare/compare:compare_test --config=ci --config=crosslinux --config=test --config=with_ui
CROSSBIN=$(bazel info bazel-bin --config=ci --config=crosslinux --config=test --config=with_ui)
COCKROACH=$CROSSBIN/pkg/cmd/cockroach/cockroach_/cockroach
COMPAREBIN=$CROSSBIN/pkg/compose/compare/compare/compare_test_/compare_test
ARTIFACTS_DIR=$PWD/artifacts
mkdir -p $ARTIFACTS_DIR

exit_status=0
$BAZCI --process_test_failures --artifacts_dir=$ARTIFACTS_DIR -- \
       test --config=ci //pkg/compose:compose_test \
       "--sandbox_writable_path=$ARTIFACTS_DIR" \
       "--test_tmpdir=$ARTIFACTS_DIR" \
       --test_env=GO_TEST_WRAP_TESTV=1 \
       --test_arg -cockroach --test_arg $COCKROACH \
       --test_arg -compare --test_arg $COMPAREBIN \
       --test_timeout=1800 || exit_status=$?

tc_end_block "Run compose tests"
exit $exit_status
