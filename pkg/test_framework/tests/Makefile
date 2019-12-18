fail:
	$(error $*)

fail-close-cluster-stop:
	$(error $*)

test-sleep25-operator-start:
	@echo "test-sleep25-operator-start $${RND}"
	@sleep 2.5s

%-env:
	@echo "export RND=$$(xxd -l8 -p /dev/random)"

%:
	@echo "$* $${RND}"
