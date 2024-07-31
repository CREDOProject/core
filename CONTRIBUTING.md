# Contribution Guidelines

This document provides a set of guidelines that aim to streamline the
collaborative aspect of this project.

Thank you for your interest in contributing to the Credo project.
We welcome contributions from everyone.
To make the process smooth and efficient, please follow these guidelines.

The keywords "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL
NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in
[RFC 2119](https://www.rfc-editor.org/rfc/rfc2119.txt).

## Prerequisites

Ensure you have the following tools installed:
- git

### For contributors inside this organization

#### Pull

To begin working on Credo you MUST pull the repository from the designated
current upstream source.

```bash
git pull git@github.com:CREDOProject/core.git
```

When the pull is complete, you can begin working with the repository.

```bash
cd core
```

We use a `dev` branch as the base for all feature development.
Contributions MUST be made through feature branches and MUST be
integrated via pull requests.

#### Workflow

##### Create a Feature Branch

Ensure you are on the `dev` branch and the local branch is up-to-date with
the upstream.

```bash
git checkout dev
git fetch origin
git rebase origin/dev
```

Create a new feature branch.

```bash
git checkout -b <username>/<feature-name>
```

You MUST replace `<username>` with your GitHub username and `<feature-name>`
with a concise description of the feature or fix.

##### Making commits

You MUST follow the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/). Use the following
format for your commit messages.

```
<type>[(<optional scope>)]: <description>

[optional body]

[optional footer(s)]
```

Common types include:

- `fix`: A bug fix
- `feat`: A new feature
- `chore`: Changes to the build process or auxiliary tools.
- `docs`: Changes affecting documentation.

It is RECOMMENDED to refer to [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/).

Examples:

```
feat(apt): add apt functionality
fix: change logic in argument parsing
```

##### Keeping your branch up-to-date

You MUST rebase your feature branch to incorporate changes from the `dev`
branch.

```bash
git fetch origin
git rebase origin/dev
```

##### Opening an issue

It is RECOMMENDED to open an issue before starting to work on a feature or an
issue to avoid duplicate efforts.

##### Submitting a pull request

You MUST push your feature branch before opening a pull request.

```bash
git push origin <username>/<feature-name>
```

Then you MUST go to Credo repository, and you should see a prompt to open a
pull request.
You MUST ensure the PR targets the dev branch of the original repository.
It is RECOMMENDED to provide a clear and descriptive title and description for
your PR, linking to the relevant issue(s) if applicable.

The maintainers will review your PR and may request changes or improvements.
Once your PR is approved, a maintainer will merge it into the `dev` branch.
