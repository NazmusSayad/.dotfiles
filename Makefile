.PHONY: skills compile git macos shell init

skills:
	bash ./etc/bundle-skills.sh

compile:
	go run ./src/compile-scripts/main.go

git:
	bash ./etc/git-config.sh

shell:
	bash -c 'source ./etc/shell-config.sh'

macos: git shell
	bash ./src/macos/setup.sh

init:
	ln -sf ../../.githooks/pre-push ../../.githooks/reference-transaction .git/hooks/
	ln -sf ../../../.githooks/pre-push ../../../.githooks/reference-transaction .dotsys/.git/hooks/
