---
name: Release check list
about: Check list to release new version
title: release [v1.x.x] check list
labels: ''
assignees: icemint0828

---

Please review the following tasks.

- [ ] Ensure milestones are at 100%.
- [ ] Cut out the working branch.  
`git checkout -b release-[v1.x.x]`
- [ ] Update the version of cmd/main.go.
- [ ] Update CHANGELOG.md.
- [ ] Push changes related to the release.  
`git push origin`
- [ ] Create new version of tag.  
`git tag [v1.x.x]`  
`git push origin [v1.x.x]`  
- [ ] Create release note according to CHANGELOG.md.
- [ ] Attach binaries to release note.  
`make build`
- [ ] Update [homebrew-tap](https://github.com/icemint0828/homebrew-tap) file.
