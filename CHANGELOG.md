## 1.3.0 (2026-06-10)

### Improvements
* provider: Added new SRV service discovery provider that uses DNS SRV records for service discovery. [GH-335](https://github.com/hashicorp/go-discover/pull/335)

### Fixed
* mdns: Fixed a bug where the mDNS provider could return empty or partial results [GH-337](https://github.com/hashicorp/go-discover/pull/337)

## 1.2.0 (2026-04-24)

## 1.1.0 (2025-06-12)

### Improvements

* AWS provider: enable dual-stack support by default, allowing the use of IPv6 endpoints wherever required. [GH-271](https://github.com/hashicorp/go-discover/pull/271)

### Fixed

* Config: `String()` will always return a parseable config string. [GH-287](https://github.com/hashicorp/go-discover/pull/287)

## 1.0.0 (2025-04-22)

🚀 go-discover v1.0.0 – Official Release 🎉

We’re excited to announce the v1.0.0 release of go-discover! This marks the first feature-complete and stable version of the package.

With this release, we are officially committing to Semantic Versioning (SemVer). This means:
  - v1.0.0 signals a stable API that you can depend on.
  - Any future breaking changes will result in a major version bump.
  - Minor and patch releases will include enhancements and fixes, respectively.

Starting from this release, we’ll also be maintaining detailed changelogs to help you track what’s new, improved, or changed in every update.

Thanks to everyone who contributed, tested, or provided feedback along the way. We’re looking forward to building on this solid foundation together.
