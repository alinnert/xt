# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Flag `--exact`/`-e` finds exact matches only if set. Examples:   
  - `xt ./file.xsd elem` finds `elem` and `parent/elem`.
  - `xt ./file.xsd elem --exact` finds `elem` only.

### Changed

- Flag `--limit 0` now shows all found results.