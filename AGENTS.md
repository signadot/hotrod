# Validating changes in a preview

Every pull request labeled `preview` gets a preview. CI builds an image from the PR's head commit,
and Argo CD application `hotrod-pr-<number>` then deploys it into a Signadot sandbox that forks
the `route` service in the shared cluster. Argo CD then runs the Playwright suite in
`playwright-tests/`, taken from your commit, against that sandbox. The result appears on your
commit as the `argocd/preview` check.

You can read Argo CD through the `argocd` MCP server. It is read-only and only shows preview
applications. You don't need cluster or Signadot credentials, and you can't sync or change
anything in Argo CD: push a commit instead.

A change is done when its preview passes, not when it compiles.

1. Open the PR with the label: `gh pr create --label preview`. If the change alters behavior a
   user can see, add or update a Playwright test for it.
2. Wait for the check on your latest commit: `gh pr checks <number> --watch --interval 30`. A
   preview takes two to three minutes, most of it the image build.
3. If `argocd/preview` fails, find out why with the MCP server:
   - `get_application` for `hotrod-pr-<number>`: `status.operationState` has the phase and message,
     and `status.operationState.syncResult.resources` lists each hook Job and its result.
   - `get_application_workload_logs` for the hook Job that failed, with `applicationNamespace`
     `argocd` and `resourceRef` `{kind: Job, group: batch, version: v1, namespace: hotrod-previews,
     name: <job>}`:
     - `hotrod-pr-<number>-sandbox`: container `wait-for-image` while it waits for CI to push the
       image (check CI with `gh pr checks`), then container `signadot` for sandbox creation.
     - `hotrod-pr-<number>-e2e`: container `e2e` for the Playwright output, or `routing-key` if
       it couldn't look up the sandbox.
   - `get_application_events` for scheduling or image pull errors.
4. Fix the cause, push, and go back to step 2. Repeat until the check passes.
5. Add the preview URL (the check's link, also the application's `Preview` info) to the PR
   description so a reviewer can try the change.
