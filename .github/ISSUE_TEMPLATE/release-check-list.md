---
name: release check list
about: check list to release new version
title: release [v1.x.x] check list
labels: ''
assignees: icemint0828

---

Please review the following tasks.

- [ ] Ensure milestones are at 100%.
- [ ] Update the version of cmd/main.go.
- [ ] Update CHANGELOG.md.
- [ ] Create new version of tag.
git tag [v1.x.x]
git push origin [v1.x.x]
- [ ] Create release note according to CHANGELOG.md.
- [ ] Attach binaries to release note.
make build
