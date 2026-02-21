#!/bin/sh
# Run E2E tests with ginkgo against a live Cloudflare API.
# Env: E2E_USE_EXISTING_CLUSTER - use existing cluster (default: "true")
# Env: E2E_PROCS              - parallel ginkgo procs (default: 4)
# Env: E2E_FOCUS              - ginkgo --focus regex (default: disabled)
# Env: CLUSTER_NAME            - kind cluster name (default: "abaddon", from mise.toml)
set -eu

export E2E_USE_EXISTING_CLUSTER="${E2E_USE_EXISTING_CLUSTER:-true}"

# Switch kubectl context to the dev cluster for the duration of the test run.
# Restores the original context on exit (success or failure).
target_context="kind-${CLUSTER_NAME:-abaddon}"
original_context="$(kubectl config current-context 2>/dev/null || echo "")"

restore_context() {
	if [ -n "$original_context" ] && [ "$original_context" != "$target_context" ]; then
		echo "Restoring kubectl context to ${original_context}"
		kubectl config use-context "$original_context" >/dev/null 2>&1 || true
	fi
}
trap restore_context EXIT

if [ "$(kubectl config current-context 2>/dev/null)" != "$target_context" ]; then
	echo "Switching kubectl context: $(kubectl config current-context) -> ${target_context}"
	kubectl config use-context "$target_context"
fi

echo "Running E2E tests (context: ${target_context})"

mkdir -p out

focus_flag=""
if [ -n "${E2E_FOCUS:-}" ]; then
	focus_flag="--focus=${E2E_FOCUS}"
fi

ginkgo -vv \
	--procs="${E2E_PROCS:-4}" \
	--keep-going \
	--silence-skips \
	--poll-progress-after=15s \
	$focus_flag \
	--output-dir out \
	--json-report=run.json \
	--cover --covermode atomic --coverprofile coverage.out \
	--race \
	./test/e2e
