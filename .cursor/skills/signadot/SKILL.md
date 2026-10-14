---
name: signadot
description: >
  This skill should be used when the user asks to "implement and validate a feature",
  "test my changes", "verify this works", "create a sandbox", "run tests against staging",
  "closed-loop development", "validate end to end", or any task where code changes need to
  be verified against real services and dependencies. It provides the closed-loop development
  workflow: implement, deploy to an isolated sandbox, test against real dependencies, read
  results, fix issues, and iterate until the feature works. Use this skill proactively
  whenever implementing features or fixes that affect services running in Kubernetes.
---

# Signadot: Closed-Loop Development

Signadot lets coding agents verify their own work. Instead of generating code and hoping
it works, deploy changes into an isolated **sandbox** that shares real infrastructure —
databases, queues, upstream and downstream services — then test against that sandbox and
iterate until everything passes.

The core loop:

```
Implement → Deploy to Sandbox → Test → Read Results → Fix → Redeploy → Retest
```

Never consider a feature done based on code alone. Always close the loop by validating
against real dependencies.

## What is a Sandbox

A sandbox is a lightweight, isolated fork of one or more services running in a shared
Kubernetes cluster. It does not duplicate the entire environment — it only replaces the
specific services being changed, while routing test traffic through them and letting
everything else (databases, caches, other services) remain shared.

This means:
- Changes are tested against real dependencies, not mocks
- Other developers and production traffic are unaffected
- No extra clusters, no duplicate infrastructure
- Sandboxes are fast to create and update

## The Closed-Loop Workflow

### 1. Implement

Make the code changes required for the feature or fix. Focus on correctness — the sandbox
will reveal integration issues that are invisible in the code alone.

### 2. Deploy to a Sandbox

Use the `signadot` CLI to create a sandbox that includes the modified service(s).

```bash
# Create or update a sandbox from a spec
signadot sandbox apply --set name=<sandbox-name> -f sandbox.yaml --wait

# Check sandbox status
signadot sandbox get <sandbox-name>
```

A sandbox spec defines which services to fork and what image or configuration to use.
If a spec already exists in the repository (check for `.signadot/` or `sandbox.yaml`),
use it. If not, create one based on the project's deployment configuration.

When modifying multiple services, fork them all in the same sandbox so they share
routing context and can communicate with each other.

To update an existing sandbox after a code change, re-run `signadot sandbox apply` with
the same name. This updates in place — no need to delete and recreate.

### 3. Route Traffic and Test

There are several ways to test against a sandbox, depending on the project setup:

**Option A: Local proxy** — Route traffic from the local machine through the sandbox:

```bash
signadot local connect --cluster <cluster-name>
signadot local proxy --sandbox <sandbox-name> \
  --map http://<service>:<port>@localhost:<local-port>
```

Then test by hitting `localhost:<local-port>` with curl, a browser, or any test tool.

**Option B: Preview URLs** — Sandboxes can expose preview endpoints. Check sandbox
details for available URLs:

```bash
signadot sandbox get <sandbox-name>
```

**Option C: Jobs** — Run test scripts inside the cluster with sandbox routing context:

```bash
signadot job submit -f test-job.yaml --sandbox <sandbox-name> --attach --wait
```

Jobs inherit the sandbox's routing, so they automatically hit the forked services.
Use `--attach` to stream output. Artifacts (test reports, screenshots, logs) can be
downloaded after the job completes:

```bash
signadot artifact download --job <job-name> -o ./results/
```

**Option D: Smart Tests** — Run test suites with sandbox-aware routing:

```bash
signadot smart-test run --sandbox <sandbox-name> --file ./tests/ --publish
```

Choose whichever approach fits the project. If test infrastructure already exists
(job specs, smart test configs), use it. If not, use local proxy + curl for quick
validation.

### 4. Inspect Results and Iterate

When tests fail, diagnose the issue:

```bash
# View sandbox pod logs
signadot sandbox get <sandbox-name>  # find pod details

# Record and inspect traffic flowing through the sandbox
signadot traffic record --sandbox <sandbox-name> --short
signadot traffic inspect <traffic-dir>

# Check job output
signadot logs --job <job-name>
```

Read the errors. Fix the root cause in the code. Rebuild, update the sandbox, and
test again. Common failure patterns:

- **Service errors**: Read logs — usually a code bug, missing config, or dependency issue.
- **Wrong data in responses**: The implementation may not cover all code paths.
- **Timeout or connection errors**: Check service ports and inter-service communication.
- **Image issues**: Verify the build succeeded and the image is accessible.

**Do not stop after the first failure.** The value of the closed loop is iteration.
Each failed test provides signal. Fix, redeploy, retest. Continue until validation passes.

### 5. Clean Up

After the feature is validated:

```bash
signadot sandbox delete <sandbox-name>
signadot local disconnect  # if local connect was used
```

## Working with Existing Project Configuration

Before creating sandbox specs from scratch, check the repository for existing Signadot
configuration:

- `.signadot/` directory — may contain sandbox specs, job definitions, test configs
- `sandbox.yaml` or `*.sandbox.yaml` — sandbox specs with variable templates
- Job runner groups and test definitions
- CI/CD integration files (GitHub Actions, Bitbucket Pipelines, etc.)

Reuse existing specs when available. They encode project-specific knowledge about service
names, namespaces, cluster targets, image registries, and test infrastructure.

## Key Principles

- **Always close the loop.** Code alone is a guess. A sandbox test is proof.
- **Iterate on failures.** Failed tests are not a dead end — they are the feedback that
  drives the next fix. Read the error, fix the cause, redeploy, retest.
- **Reuse what exists.** Check for existing sandbox specs, job definitions, and test
  configs before creating new ones. The project may already have the infrastructure.
- **Fork only what changed.** A sandbox should only include the services that were
  modified. Everything else is shared from the baseline environment.
- **Keep the sandbox alive during iteration.** Update with `signadot sandbox apply`
  rather than deleting and recreating. Faster and preserves routing state.
- **Test the integration, not just the unit.** The sandbox exists to catch issues that
  only appear when services interact with real dependencies.

## Quick Reference

| Task | Command |
|------|---------|
| Create/update sandbox | `signadot sandbox apply --set name=<n> -f spec.yaml --wait` |
| Check sandbox status | `signadot sandbox get <name>` |
| List sandboxes | `signadot sandbox list` |
| Connect to cluster | `signadot local connect --cluster <cluster>` |
| Proxy traffic locally | `signadot local proxy --sandbox <name> --map http://svc:port@localhost:port` |
| Run test job | `signadot job submit -f job.yaml --sandbox <name> --attach` |
| Run smart tests | `signadot smart-test run --sandbox <name> --file ./tests/` |
| Record traffic | `signadot traffic record --sandbox <name> --short` |
| Download artifacts | `signadot artifact download --job <name> -o ./results/` |
| View job logs | `signadot logs --job <name>` |
| Delete sandbox | `signadot sandbox delete <name>` |
| Disconnect | `signadot local disconnect` |
