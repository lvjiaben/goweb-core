VERSION := $(shell cat VERSION)

print-version:
	@printf "%s\n" "$(VERSION)"

release-check:
	test -f VERSION
	test -f CHANGELOG.md
	test -f docs/versioning.md
	test -f docs/compatibility.md
	test -f docs/releases/RELEASE_POLICY.md
	test -f docs/releases/RELEASE_CHECKLIST.md
	test -f docs/releases/$(VERSION).md
	grep -q "$(VERSION)" README.md
	grep -q "docs/versioning.md" README.md
	grep -q "docs/compatibility.md" README.md
	grep -q "$(VERSION)" CHANGELOG.md
