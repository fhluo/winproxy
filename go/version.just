set shell := ["nu", "-c"]
set script-interpreter := ["nu"]
set indentation := "  "
set default-list := true

[script]
bump version:
  let ver = ("{{version}}" | into semver --loose | into record | reject prefix | into semver | to text)

  let status = (git status --porcelain --untracked-files=no)
  if ($status | is-not-empty) {
    print $"(ansi red)working tree is not clean, please commit or stash changes first(ansi reset)"
    exit 1
  }

  let tag = $"go/v($ver)"

  if (git tag -l $tag | is-not-empty) {
    print $"(ansi red)tag ($tag) already exists(ansi reset)"
    exit 1
  }

  print $"(ansi light_gray)Tagging ($tag)...(ansi reset)"
  git tag $tag -m $"($tag)"

  print $"(ansi green)✓ Bumped to ($tag)(ansi reset)"
