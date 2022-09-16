itest-env:
	echo "export KUBECONFIG=$(shell mktemp /tmp/kube_XXXX.yaml)"
	echo "export CLUSTER=itest-$(shell xxd -p -l4 /dev/urandom)"
	echo "export CIDFILE=$(shell mktemp /tmp/kube_XXXX.txt)"
	echo "export K3DIMAGE=itest-$(shell xxd -p -l4 /dev/urandom)"

itest-cluster-start:
	docker build -t $${K3DIMAGE} .
	rm -f $${CIDFILE}
	docker container run -d --init --cidfile $${CIDFILE} -e CLUSTER=$${CLUSTER} -e DOCKER_HOST -v /var/run/docker.sock:/var/run/docker.sock $${K3DIMAGE} /bin/sleep infinity
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'python3 port_reserve.py > port.txt && rm port_reserve.py'
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'k3d --trace cluster delete $${CLUSTER} &>/dev/null || echo -n'
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'k3d --trace cluster create $${CLUSTER} --api-port $$(cat port.txt) --wait --timeout 30s && rm port.txt'
	sleep 5
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'k3d --trace cluster start $${CLUSTER}'
	sleep 5
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'k3d kubeconfig get $${CLUSTER}' > $${KUBECONFIG}
	echo "Using KUBECONFIG=$${KUBECONFIG}, API $$(grep -Eo 'server: http.*$$' $${KUBECONFIG})"

itest-cluster-stop:
	if [ -n "$${CIDFILE}" ] && [ -s "$${CIDFILE}" ]; then \
	  docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'k3d cluster delete --all'; \
	  docker container stop $$(cat $${CIDFILE}) >/dev/null 2>&1; \
	  docker container rm $$(cat $${CIDFILE}) >/dev/null 2>&1; \
	  rm $${CIDFILE}; \
	fi

test-sleep05-foo:
	@echo "test-sleep05-foo $${RND} $${TEST_OPERATOR_NS} $${TEST_COUNT}$${FOO}"

test-sleep05-bar:
	@echo "test-sleep05-bar $${RND} $${TEST_OPERATOR_NS} $${TEST_COUNT}, expected to raise an error"
	false

fail:
	$(error $*)

fail-close-cluster-stop:
	$(error $*)

test-sleep05-operator-start:
	@echo "test-sleep05-operator-start $${RND} $${TEST_OPERATOR_NS} $${TEST_COUNT}"
	@sleep 0.5s

%-operator-start:
	@echo "$*-operator-start $${RND} $${TEST_OPERATOR_NS} $${TEST_COUNT}"

%-operator-stop:
	@echo "$*-operator-stop $${RND} $${TEST_OPERATOR_NS} $${TEST_COUNT}"

%-cleanup:
	@echo "$*-cleanup $${RND} $${TEST_OPERATOR_NS} $${TEST_COUNT}"

%-env:
	@echo "export RND=$$(xxd -l8 -p /dev/urandom)"

%:
	@echo "$* $${RND}"
