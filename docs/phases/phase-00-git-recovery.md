# Phase 0 local Git recovery

The owner authorized a local commit of reviewed Phase 0 files only. The agent's explicit staging attempt failed creating `.git/index.lock` because the sandbox mounts original Git metadata read-only. No permission changes, alternate index for that repository, reset, clean, checkout or push were used to circumvent it.

An isolated candidate repository/commit can be created in writable temporary storage and tested independently. It is not an update of the original branch. Current report records the exact candidate SHA, evidence and any retained recovery bundle.

## Precise manifest

`phase-00-staging-manifest.txt` lists the intended source/docs/locks/CI files individually, including this guide and the manifest itself. It excludes `.env`, installed dependencies, caches, logs, generated binaries, Git metadata, local DB/storage and ignored recovery artifacts. Historical sanitized validation JSON under docs is deliberate documentary evidence, not an application runtime log.

## Authorized WSL host commands

Review any new edits first. Do not run these if unrelated changes or staged files appear. No remote push is included.

```sh
cd ~/Projects/TGP-Sys/backend
git rev-parse --show-toplevel
git branch --show-current
git status --short --untracked-files=all
test -z "$(git diff --cached --name-only)" || { echo 'Stop: existing staged work requires review'; exit 1; }
# Review the explicit paths; do not print the private environment:
cat docs/phases/phase-00-staging-manifest.txt
# Stage only this manifest, never all working-directory files:
git add --pathspec-from-file=docs/phases/phase-00-staging-manifest.txt
git diff --cached --name-only
git diff --cached --stat
# Markdown's intentional two-space hard breaks are checked by Python below.
git diff --cached --check -- . ':(glob,exclude)**/*.md'
```

Before committing, run this staged-content check without echoing private values:

```sh
python3 - <<'PY'
from pathlib import Path
import re, subprocess, sys
sys.path.insert(0, 'scripts')
from local_env import load
expected=set(Path('docs/phases/phase-00-staging-manifest.txt').read_text().splitlines())
actual=set(subprocess.check_output(['git','diff','--cached','--name-only','-z']).decode().strip('\0').split('\0'))
assert actual <= expected and '.env' not in actual, 'Unexpected staged path; stop for review'
assert subprocess.run(['git','check-ignore','-q','.env']).returncode == 0
assert not subprocess.check_output(['git','ls-files','--','.env']).strip()
settings=load()
private=[settings[k].encode() for k in ('DB_PASSWORD','POSTGRES_PASSWORD','MIGRATION_PASSWORD')]
patterns=[rb'-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----',rb'AKIA[0-9A-Z]{16}',rb'ghp_[A-Za-z0-9]{36}',rb'github_pat_[A-Za-z0-9_]{50,}']
for name in actual:
    assert not any(part in {'.git','node_modules','__pycache__','.tools'} for part in Path(name).parts)
    content=subprocess.check_output(['git','show',':'+name])
    assert not any(value in content for value in private), 'Private credential finding; stop and inspect privately'
    assert not any(re.search(pattern,content) for pattern in patterns), 'Credential-pattern finding; stop and inspect privately'
    for line in content.splitlines():
        match=re.search(rb'[ \t]+$',line)
        if match:
            assert name.endswith('.md') and match.group()==b'  ', 'Unexpected trailing whitespace; stop for review'
assert expected <= set(subprocess.check_output(['git','ls-files','-z']).decode().strip('\0').split('\0')), 'Manifest file missing from index'
print('PASS: intended files included, staged paths reviewed, private environment excluded, no credential findings')
PY
# Only after all review/security checks pass:
git commit -m 'chore(phase0): establish Go backend foundation'
git rev-parse HEAD
git ls-tree -r --name-only HEAD
git status --short
```

If any command fails, stop and report; do not reset, clean or overwrite user work. This is an owner-executed host procedure, not a claim that the agent committed the original branch. Do not blindly cherry-pick/reset to a temporary candidate over uncommitted files.

After recovery, clone the new exact revision to an isolated directory and rerun non-infrastructure onboarding/tests. The default Compose names/ports/volume belong to the existing development instance; do not run `make db-up`, `make migrate-*`, or `make verify-phase0` from the clone against that instance. Fresh live DB checks need separate isolated resources with explicit reviewed names/ports and disposable storage. Existing host evidence and browser/Bruno procedures remain separate.

The default unfiltered `git diff --check` flags existing two-space Markdown hard breaks, including the unchanged master plan. These are intentional Markdown syntax, not ignored code/test failures. The procedure checks non-Markdown files with Git and permits only exactly two spaces at Markdown line endings; tabs or other trailing whitespace fail.

The retained candidate bundle advertises `phase0-recovery-candidate`, not a default HEAD. For inspection in a new temporary directory, use `git clone --branch phase0-recovery-candidate <bundle-path> <new-empty-directory>`. Do not clone it over the working repository. The current report records its exact SHA and successful bundle integrity/file-inventory verification.
