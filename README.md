# Costroid demo

This repository shows [Costroid](https://github.com/Costroid/costroid), the cost check for pull requests, on a small Go project. Every pull request gets a `costroid` check and a comment that state what the coding agent's review cost, with the figure's basis.

## What a pull request gets

| Pull request | Workflow | What it spends |
| --- | --- | --- |
| From a branch of this repository, by a member, or with the label `costroid-agent` a maintainer adds | [`costroid.yml`](.github/workflows/costroid.yml): `claude-code-action` reviews the change under a cap of USD 5 a run, and Costroid reports what the run cost | The run's tokens, on the key the environment `demo` holds |
| From a fork | [`costroid-recorded.yml`](.github/workflows/costroid-recorded.yml): Costroid reports the rows of a recorded run | Nothing. The comment's first line says the figures are workflow run 42368537950's, recorded on 2026-11-07, not the pull request's |

No agent runs on a fork's pull request. GitHub holds a first-time contributor's workflow runs until a maintainer approves them, so on such a pull request the recorded comment arrives after that approval. [`recipes/costroid-run.yml`](recipes/costroid-run.yml) is the same check with Costroid running Claude Code, which stops the run at the budget. It is not in `.github/workflows`, so each pull request pays for one run.

The project is [`split`](split/split.go), which divides a bill among people in whole cents. `go test ./...` runs its tests.

## Try Costroid on your own repository

[Install Costroid](https://github.com/Costroid/costroid/blob/main/docs/install.md), then print the comment for the same recorded run, with no key:

```sh
costroid demo
```

In your repository, write the workflow and `costroid.yml`:

```sh
costroid init github --agent claude-code-action
```

`init` prints the commands that are left: the key's secret, a branch and the pull request.

## Licence

Apache-2.0.
