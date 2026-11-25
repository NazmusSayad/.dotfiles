.PHONY: skills compile git macos linux shell init

skills:
	go run ./src/build-skills/main.go

compile:
	go run ./src/compile-scripts/main.go

git:
	bash ./etc/git-config.sh

shell:
	bash -c 'source ./etc/shell-config.sh'

macos: git shell
	bash ./src/macos/setup.sh

linux: git shell
	bash ./src/linux/setup.sh
	bash ./.dotsys/linux/setup.sh

init:
	MSYS=winsymlinks:nativestrict ln -sf ../../.githooks/pre-push ../../.githooks/reference-transaction .git/hooks/
	MSYS=winsymlinks:nativestrict ln -sf ../../../.githooks/pre-push ../../../.githooks/reference-transaction .dotsys/.git/hooks/
