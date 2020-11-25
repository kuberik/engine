# Releases

This page describes the release process and the currently planned schedule for upcoming releases.

## Release schedule

Release cadence of first pre-releases being cut is 6 weeks.

## Release coordinator responsibilities

The release coordinator is responsible for the entire release series of a minor release, meaning all pre- and patch releases of a minor release. The process formally starts with the initial pre-release, but some preparations should be done a few days in advance.

* We aim to keep the master branch in a working state at all times. In principle, it should be possible to cut a release from master at any time. In practice, things might not work out as nicely. A few days before the pre-release is scheduled, the coordinator should check the state of master. Following their best judgement, the coordinator should try to expedite bug fixes that are still in progress but should make it into the release. On the other hand, the coordinator may hold back merging last-minute invasive and risky changes that are better suited for the next minor release.
* On the date listed in the table above, the release coordinator cuts the first pre-release (using the suffix `-rc.0`) and creates a new branch called  `release-<major>.<minor>` starting at the commit tagged for the pre-release. In general, a pre-release is considered a release candidate (that's what `rc` stands for) and should therefore not contain any known bugs that are planned to be fixed in the final release.
* With the pre-release, the release coordinator is responsible for making sure there's no major issues with the release  within next 3 days, after which, if everything is alright, the pre-release is promoted to a stable release.
* If regressions or critical bugs are detected, they need to get fixed before cutting a new pre-release (called `-rc.1`, `-rc.2`, etc.).

See the next section for details on cutting an individual release.

## How to cut an individual release

These instructions are currently valid for all [kuberik projects](https://github.com/kuberik).

### Branch management and versioning strategy

We use [Semantic Versioning](https://semver.org/).

We maintain a separate branch for each minor release, named `release-<major>.<minor>`, e.g. `release-0.3`, `release-1.0`.

Note that branch protection kicks in automatically for any branches whose name starts with `release-`. Never use names starting with `release-` for branches that are not release branches.

The usual flow is to merge new features and changes into the master branch and to merge bug fixes into the latest release branch. Bug fixes are then merged into master from the latest release branch. The master branch should always contain all commits from the latest release branch. As long as master hasn't deviated from the release branch, new commits can also go to master, followed by merging master back into the release branch.

If a bug fix got accidentally merged into master after non-bug-fix changes in master, the bug-fix commits have to be cherry-picked into the release branch, which then have to be merged back into master. Try to avoid that situation.

Maintaining the release branches for older minor releases happens on a best effort basis.

### 0. Updating dependencies

A few days before a major or minor release, consider updating the dependencies:

```
export GO111MODULE=on
go get -u ./...
go mod tidy
git add go.mod go.sum
```

Then create a pull request against the master branch.

Note that after a dependency update, you should look out for any weirdness that
might have happened. Such weirdnesses include but are not limited to: flaky
tests, differences in resource usage, panic.

In case of doubt or issues that can't be solved in a reasonable amount of time,
you can skip the dependency update or only update select dependencies. In such a
case, you have to create an issue or pull request in the GitHub project for
later follow-up.

### 1. Prepare your release

At the start of a new major or minor release cycle create the corresponding release branch based on the master branch. For example if we're releasing `0.4.0` and the previous stable release is `0.3.0` we need to create a `release-0.4` branch. Note that all releases are handled in protected release branches, see the above `Branch management and versioning` section. Release candidates and patch releases for any given major or minor release happen in the same `release-<major>.<minor>` branch. Do not create `release-<version>` for patch or release candidate releases.

Changes for a patch release or release candidate should be merged into the previously mentioned release branch via pull request.

Bump the version in the `VERSION` file and update `CHANGELOG.md`. Do this in a proper PR pointing to the release branch as this gives others the opportunity to chime in on the release in general and on the addition to the changelog in particular. For a release candidate, append something like `-rc.0` to the version (with the corresponding changes to the tag name, the release name etc.).

Note that `CHANGELOG.md` should only document changes relevant to users of kuberik, including external API changes, performance improvements, and new features. Do not document changes of internal interfaces, code refactorings and clean-ups, changes to the build process, etc. People interested in these are asked to refer to the git history.

For release candidates still update `CHANGELOG.md`, but when you cut the final release later, merge all the changes from the pre-releases into the one final update.

Entries in the `CHANGELOG.md` are meant to be in this order:

* `[CHANGE]`
* `[FEATURE]`
* `[ENHANCEMENT]`
* `[BUGFIX]`

### 2. Draft the new release

Tag the new release via the following commands:

```bash
tag=$(< VERSION)
git tag -s "v${tag}" -m "v${tag}"
git push origin "v${tag}"
```

Signing a tag with a GPG key is appreciated, but in case you can't add a GPG key to your Github account using the following [procedure](https://help.github.com/articles/generating-a-gpg-key/), you can replace the `-s` flag by `-a` flag of the `git tag` command to only annotate the tag without signing.

TODO: Process from bellow paragraphs is still manual

Once a tag is created, the release process will be triggered for this tag and new draft the GitHub release using the `halbot` account.

Finally, wait for the build step for the tag to finish. The point here is to wait for tarballs to be uploaded to the Github release and the container images to be pushed to the Docker Hub. Once that has happened, click _Publish release_, which will make the release publicly visible and create a GitHub notification.
** Note: ** for a release candidate version ensure the _This is a pre-release_ box is checked when drafting the release in the Github UI. The CI job should take care of this but it's a good idea to double check before clicking _Publish release_.`