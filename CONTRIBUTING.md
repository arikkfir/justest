# Contributing

Hi there! We're thrilled that you'd like to contribute to this project. Your help is essential for keeping it great.

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## Setup

Only `Go` is required. Make sure the version you have installed matches or supports the version listed in the `go.mod` file.

## Testing

Use `go test` to run the tests locally:

```bash
$ go test ./...
```

## Commit messages

While interim commits in non-main branches do not have to conform to anything, before you can merge a PR, you will have to make sure each one of them both:

- Conforms to the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) standard
- Only uses allowed types & scopes as listed below
- Properly documents breaking changes (see below)
- Message body is required (unlike in the conventional commits standard)
- PR reference is required as a footer (see example below)

### Allowed types & scopes

- `build`: for changes to how the project is built (e.g. changes to Go version), with one of the following scopes:
  - `go`
- `ci`: for changes to CI files like GitHub workflows, with one of the following scopes:
  - `github`
- `docs`: for changes to documentation without corresponding functionality changes, scopes can be any subject matter
- `feat`: for new functionality & features, scopes can be any subject matter
- `fix`: for bug fixes, scopes can be any subject matter
- `perf`: for changes that improve performance, scopes can be any subject matter
- `refactor`: for code structure changes without changing functionality, no scope allowed
- `style`: for fixing code formatting or styling, no scope allowed
- `test`: for adding tests, scopes can be any subject matter

### Breaking changes

Commits that introduce breaking changes **must**:

- Use `!` following their type & scope declaration
- Prefix the commit body with the `BREAKING CHANGE: ` prefix

### Examples

Here's a simple example introducing a new feature:

```
feat(matcher:not): add a new matcher called "Not"

This change introduces a new `Matcher` implementation: `Not`

Refs: #187
```

Here's another example that introduces a breaking change:

```
fix(api)!: remove accidentally exposed method

BREAKING CHANGE: this change removes the `DoSomething()` method which was accidentally
added to the public API of this package.

Refs: #43
```

## Issues and PRs

If you have suggestions for how this project could be improved, or want to report a bug, open an issue! We'd love all and any contributions. If you have questions, too, we'd love to hear them.

We'd also love PRs. If you're thinking of a large PR, we advise opening up an issue first to talk about it, though! Look at the links below if you're not sure how to open a PR.

### Submitting a pull request

1. Make sure you set up your local development environment as described above.
2. [Fork](https://github.com/arikkfir/justest/fork) and clone the repository.
3. Create a new branch: `git checkout -b my-branch-name`.
4. Make your change and make sure the entire tests suite passes (see above)
5. Push to your fork, and submit a pull request
6. Pat your self on the back and wait for your pull request to be reviewed and merged.

Here are a few things you can do that will increase the likelihood of your pull request being accepted:

- Ensure you write and/or update tests to match the new/modified functionality
- Keep your PR as focused as possible, avoid PRs that add multiple unrelated features or fix multiple bugs at once
- Keep your PR as small as possible, to make reviewing easier and faster
- Provide any and all necessary information in the PR description to help reviewers understand the context and impact of your change.
- Work in Progress pull requests are also welcome to get feedback early on, or if there is something blocked you
  - Make sure you use GitHub's "Draft PR" feature for these

### Merging strategy

The `main` branch must retain linear history, and that means the only allowed merge strategy in this repository is the "Rebase" strategy. It also means that reviews must be done on the tip of the PR, and once you rebase or force-push to restructure commits, another review will be required.

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)
