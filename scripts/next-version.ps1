# Compute the next semantic version from conventional commits.
#
#   feat: / feat(scope):            -> minor
#   fix:  / fix(scope):             -> patch
#   "BREAKING CHANGE" / "type!:"    -> major
#   anything else                   -> patch
#
# No release at all when nothing changed in the shipped library since the
# last tag (docs-only, CI/workflow, or script-only pushes are skipped — they
# batch up into the next library release instead).
#
# Writes skip/next/bump/changelog to $env:GITHUB_OUTPUT (GitHub Actions),
# or prints them when run locally.
#
# Usage: ./scripts/next-version.ps1

$ErrorActionPreference = 'Continue'

git fetch --tags --force 2>$null

$out = @{
  skip      = 'false'
  next      = ''
  bump      = ''
  changelog = ''
}

$headMsg = git log -1 --format=%s
if ($headMsg -match '\[skip release\]') {
  $out.skip = 'true'
  Write-Host 'HEAD asks to skip the release.'
} else {
  $latest = git tag --list 'v*' --sort=-v:refname | Select-Object -First 1
  if ([string]::IsNullOrEmpty($latest)) { $range = 'HEAD'; $base = '0.0.0' }
  else { $range = "$latest..HEAD"; $base = $latest.TrimStart('v') }

  if ((git rev-list --count $range) -eq '0') {
    $out.skip = 'true'
    Write-Host "No new commits since $latest - nothing to release."
  } else {
    # Files changed since the last release (first release: everything tracked).
    # Only these paths affect the shipped module — anything else (docs, CI,
    # release scripts) batches up silently into the next library release.
    if ([string]::IsNullOrEmpty($latest)) { $files = @(git ls-files) }
    else { $files = @(git diff --name-only $latest HEAD) }
    $appChanged = @($files | Where-Object {
      (($_ -match '\.go$') -and ($_ -notmatch '^scripts/')) `
        -or ($_ -match '^(go\.mod|go\.sum)$') `
        -or ($_ -match '^integration/(go\.mod|go\.sum)$')
    })
    if (-not $appChanged) {
      $out.skip = 'true'
      Write-Host "No library changes since $(if ($latest) { $latest } else { 'the beginning' }) - docs/meta only, skipping release."
    } else {
    $subjects = @(git log $range --format=%s)
    $bodies = git log $range --format=%b | Out-String

    $bump = 'patch'
    if (($subjects -match '^[a-zA-Z]+(\([^)]*\))?!:') -or ($bodies -match 'BREAKING[ -]CHANGE')) { $bump = 'major' }
    elseif ($subjects -match '^feat(\([^)]*\))?:') { $bump = 'minor' }
    elseif ($subjects -match '^fix(\([^)]*\))?:') { $bump = 'patch' }

    $p = $base.Split('.'); $ma = [int]$p[0]; $mi = [int]$p[1]; $pa = [int]$p[2]
    switch ($bump) {
      'major' { $ma++; $mi = 0; $pa = 0 }
      'minor' { $mi++; $pa = 0 }
      default { $pa++ }
    }
    $next = "v$ma.$mi.$pa"
    if (git rev-parse --verify --quiet "refs/tags/$next") {
      $out.skip = 'true'
      Write-Host "Tag $next already exists - skipping."
    } else {
      $breaking = @($subjects | Where-Object { $_ -match '^[a-zA-Z]+(\([^)]*\))?!:' })
      $feats = @($subjects | Where-Object { ($_ -match '^feat(\([^)]*\))?:') -and ($_ -notmatch '!:' ) })
      $fixes = @($subjects | Where-Object { ($_ -match '^fix(\([^)]*\))?:') -and ($_ -notmatch '!:' ) })
      $other = @($subjects | Where-Object { ($_ -notmatch '^[a-zA-Z]+(\([^)]*\))?!:') -and ($_ -notmatch '^feat(\([^)]*\))?:') -and ($_ -notmatch '^fix(\([^)]*\))?:') })

      $cl = @()
      if ($breaking.Count) { $cl += '### Breaking Changes'; $breaking | ForEach-Object { $cl += "- $_" }; $cl += '' }
      if ($feats.Count) { $cl += '### Features'; $feats | ForEach-Object { $cl += "- $_" }; $cl += '' }
      if ($fixes.Count) { $cl += '### Fixes'; $fixes | ForEach-Object { $cl += "- $_" }; $cl += '' }
      if ($other.Count) { $cl += '### Other'; $other | ForEach-Object { $cl += "- $_" }; $cl += '' }
      if ($latest) { $cl += "_Bump: $bump (previous: $latest)_" }
      else { $cl += "_Bump: $bump (first release)_" }

      $out.next = $next
      $out.bump = $bump
      $out.changelog = ($cl -join "`n")
      Write-Host "Releasing $next ($bump)"
    }
    }
  }
}

if ($env:GITHUB_OUTPUT) {
  "skip=$($out.skip)" | Out-File -FilePath $env:GITHUB_OUTPUT -Append -Encoding utf8
  "next=$($out.next)" | Out-File -FilePath $env:GITHUB_OUTPUT -Append -Encoding utf8
  "bump=$($out.bump)" | Out-File -FilePath $env:GITHUB_OUTPUT -Append -Encoding utf8
  'changelog<<EOF' | Out-File -FilePath $env:GITHUB_OUTPUT -Append -Encoding utf8
  $out.changelog | Out-File -FilePath $env:GITHUB_OUTPUT -Append -Encoding utf8
  'EOF' | Out-File -FilePath $env:GITHUB_OUTPUT -Append -Encoding utf8
} else {
  $out
}

# The tag-existence check above (`git rev-parse --verify --quiet`) exits 1
# when the tag does NOT exist — the normal, expected case. GitHub's pwsh
# wrapper ends every step with `exit $LastExitCode`, so a stale 1 here would
# fail the step even though the script succeeded. Reset unconditionally.
$global:LastExitCode = 0
