fail:
	$(error $*)

fail-close-cluster-stop:
	$(error $*)

test-sleep05-operator-start:
	@echo "test-sleep05-operator-start $${RND}"
	@sleep 0.5s

%-env:
	@echo "export RND=$$(xxd -l8 -p /dev/random)"

%:
	@echo "$* $${RND}"
