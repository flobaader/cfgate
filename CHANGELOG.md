# Changelog

All notable changes to cfgate are documented in this file.


## [0.1.0-alpha.13] - 2026-02-21

### Bug Fixes

- Sort ingress rules by path specificity, harden SDK and controllers
- **(controller)** Add isGatewayParentRef guard in findGatewaysForHTTPRoute
- DNS cleanup path, ownership matching, status constants, CEL rules
- **(crd)** Tighten DNS hostname validation per RFC 1035

### Documentation

- Rewrite README, extract service-mesh guide, add cross-references
- Apply prose style conventions across documentation

### Maintenance

- Generate changelog
- Bind container image
- Add CI workflows, ignore E2E output directory

### Other

- Merge pull request #1 from cfgate/v0.1.0-alpha.13

v0.1.0-alpha.13

## [0.1.0-alpha.12] - 2026-02-19

### Features

- Add cosign keyless signing to container releases

### Bug Fixes

- Unmark GH release as pre-release

### Maintenance

- **(actions)** Add timesouts, pin versions

## [0.1.0-alpha.11] - 2026-02-19

### Features

- **(site)** Add "images"
- New image release tags for AH

### CI/CD

- Add Artifact Hub metadata and annotations

### Maintenance

- **(site)** Update packages
- **(site)** Update images
- Rename Cloudflare worker

### Other

- Inherent-design/cfgate -> cfgate/cfgate

## [0.1.0-alpha.10] - 2026-02-17

### Features

- Scaffold Astro alongside Hono worker
- **(site)** Wire static assets and OG tags into Astro layout
- **(site)** Integrate brand design system, i18n, and Starwind UI
- **(site)** Add Hindi translation, fix zh locale label

### Bug Fixes

- **(test)** Add fallback credentials to deletion invariant DNS resource
- Patch task scripts for empty-arg bug, fragile cd, contract docs
- **(site)** Update English subtitle to match translation structure

### Testing

- Fix invariant assertion, add conflict retry to bare Get/Update sites

### Documentation

- Convert ASCII diagrams to Mermaid

### Refactoring

- Extract shared task scripts for mise/CI invariance
- **(site)** Switch to published @inherent.design/brand package
- **(site)** Use brand components, theme button globally

### Maintenance

- Update pnpm-lock
- **(site)** Update packages
- **(site)** Remove stale scripts, rename deploy:cf -> deploy
- **(site)** Fix scripts
- Update brand package
- Update README
- Sync chart + app version

## [0.1.0-alpha.9] - 2026-02-09

### Testing

- Add E2E invariant tests for structural property verification

### Documentation

- Alpha.9 documentation overhaul, purge origin-no-tls-verify, fix examples

## [0.1.0-alpha.8] - 2026-02-09

### Bug Fixes

- **(controller)** Register v1beta1 scheme, sync AccessPolicy CRD, demote noisy log

### Documentation

- Fix deployment names, remove dead dns-sync annotations from examples

### Maintenance

- Update changelog

## [0.1.0-alpha.7] - 2026-02-09

### Features

- **(controller)** Alpha.7 reconciler stabilization and HTTPRoute credential inheritance

### Documentation

- Add shields.io badges to README
- Fix CRD table, add credential resolution and troubleshooting

### Maintenance

- Update changelog for alpha.6 and fix badge layout
- **(chart)** Bump to v1.0.3 / appVersion 0.1.0-alpha.7

## [0.1.0-alpha.6] - 2026-02-08

### Features

- Alpha.6 comprehensive stabilization (unreleased)

### Bug Fixes

- **(controller)** Alpha.6 reconcile, deletion, and API stabilization
- **(controller)** Logging guard and em-dash removal

### Testing

- **(e2e)** Alpha.6 coverage expansion and 94/94 stabilization

### Documentation

- Add git-cliff changelog and update project docs

### CI/CD

- Use git-cliff for release notes generation

### Maintenance

- Local dev fixes (docker cache, mise tasks)
- Reset kustomization.yaml after local deploy

## [0.1.0-alpha.5] - 2026-02-06

### Features

- Helm chart v1.0.1

### Bug Fixes

- Alpha.5 controller stabilization

### Infrastructure

- Cfgate.io v0.1.2 custom_domain for auto DNS

### Maintenance

- Local dev tasks + docs

## [0.1.0-alpha.4] - 2026-02-06

### Features

- Alpha.3 implementation
- CloudflareDNS CRD (composable architecture)
- Add helm chart v1.0.0

### Bug Fixes

- SA1019 events API migration + reconciliation bugs

### Testing

- Alpha.3 E2E suite (85/85 passing)

### Documentation

- Alpha.3 samples and examples
- Godoc comments and logging audit

### Infrastructure

- Initialize cfgate.io as wrangler
- Add version injection to builds
- Use kubectl kustomize
- Cfgate.io v0.1.1 with route fix
- Alpha.4 CI/CD improvements

### Maintenance

- Pin doc2go version
- Organize mise.toml
- Fix cfgate.io bootstrap

## [0.1.0-alpha.2] - 2026-02-02

### Bug Fixes

- Use release version tag in install.yaml

### Documentation

- Update README and examples for v0.1.0-alpha.1

## [0.1.0-alpha.1] - 2026-02-02

### Features

- Initial commit
- Add Dockerfile for container builds
- Add docs, ci, mise tooling

### Bug Fixes

- Remove deprecated Connections field (SA1019)
- Add kustomize directory structure

### Documentation

- Update Gateway API version and consolidate test tasks

### CI/CD

- Separate e2e, drop mise
- Bump golangci-lint to v2.8.0
- Update workflows; remove e2e
- Add path filter to pull_request trigger
- Remove workflow_dispatch from release

### Maintenance

- Clean CI workflows and dead code
