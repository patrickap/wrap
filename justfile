PROJECT_NAME := "patrickap/wrap"
PROJECT_VERSION := "VERSION"

[private]
get-version:
  @cat {{PROJECT_VERSION}}

[private]
set-version version:
  @echo {{version}} > {{PROJECT_VERSION}}

[private]
backup-version:
  @cp {{PROJECT_VERSION}} {{PROJECT_VERSION}}.bak

[private]
restore-version:
  @cp {{PROJECT_VERSION}}.bak {{PROJECT_VERSION}}

[private]
go-build:
  @go mod download
  @GOOS=darwin GOARCH=amd64 go build -o ./build/darwin-amd64/wrap
  @GOOS=darwin GOARCH=arm64 go build -o ./build/darwin-arm64/wrap
  @GOOS=linux GOARCH=amd64 go build -o ./build/linux-amd64/wrap
  @GOOS=linux GOARCH=arm64 go build -o ./build/linux-arm64/wrap

[private]
git-publish:
  @git add . -- ':!{{PROJECT_VERSION}}.bak'
  @git commit -m "chore(release): v$(just get-version)"
  @git push
  @git tag -a "v$(just get-version)" -m "Release v$(just get-version)"
  @git push --tags origin

[private]
clean-up:
  @rm {{PROJECT_VERSION}}.bak

[private]
release-patch:
  @just backup-version
  @just set-version $(just get-version | awk -F. -v OFS=. '{$3++; print}')
  @just go-build && just git-publish || just restore-version
  @just clean-up

[private]
release-minor:
  @just backup-version
  @just set-version $(just get-version | awk -F. -v OFS=. '{$2++; $3=0; print}')
  @just go-build && just git-publish || just restore-version
  @just clean-up

[private]
release-major:
  @just backup-version
  @just set-version $(just get-version | awk -F. -v OFS=. '{$1++; $2=0; $3=0; print}')
  @just go-build && just git-publish || just restore-version
  @just clean-up

release type:
  @just release-{{type}}
