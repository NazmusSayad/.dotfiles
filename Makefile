.PHONY: compile git macos shell

compile:
	go run ./src/compile-scripts/main.go

git:
	bash ./etc/git-config.sh

shell:
	bash -c 'source ./etc/shell-config.sh'

macos: git shell
	bash ./src/macos/setup.sh
