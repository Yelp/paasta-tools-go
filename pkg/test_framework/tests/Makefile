fail:
	$(error $*)

fail-close-cluster-stop:
	$(error $*)

test-sleep05-operator-start:
	@echo "test-sleep05-operator-start $${RND} $${TEST_OPERATOR_NS}"
	@sleep 0.5s

%-operator-start:
	@echo "$*-operator-start $${RND} $${TEST_OPERATOR_NS}"

%-operator-stop:
	@echo "$*-operator-stop $${RND} $${TEST_OPERATOR_NS}"

%-cleanup:
	@echo "$*-cleanup $${RND} $${TEST_OPERATOR_NS}"

%-env:
	@echo "export RND=$$(xxd -l8 -p /dev/urandom)"

%:
	@echo "$* $${RND}"
