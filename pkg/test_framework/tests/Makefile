export KUBECONFIG
export CLUSTER
export CIDFILE
export K3DIMAGE

itest-env:
	echo "export KUBECONFIG=$(shell mktemp /tmp/kube_XXXX.yaml)"
	echo "export CLUSTER=itest-$(shell xxd -p -l4 /dev/urandom)"
	echo "export CIDFILE=$(shell mktemp /tmp/kube_XXXX.txt)"
	echo "export K3DIMAGE=itest-$(shell xxd -p -l4 /dev/urandom)"

itest-cluster-start:
	docker build -t $${K3DIMAGE} .
	rm -f $${CIDFILE}
	docker container run -d --cidfile $${CIDFILE} -e CLUSTER=$${CLUSTER} -v /var/run/docker.sock:/var/run/docker.sock $${K3DIMAGE} /bin/sleep infinity
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'python3 port_reserve.py > port.txt && rm port_reserve.py'
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'k3d delete --name $${CLUSTER} &>/dev/null || echo -n'
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'k3d create --name $${CLUSTER} --api-port $$(cat port.txt) && rm port.txt'
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'k3d start --name $${CLUSTER}'
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'until k3d get-kubeconfig --name $${CLUSTER} &>/dev/null; do sleep 0.5s; done'
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'cat $$(k3d get-kubeconfig --name $${CLUSTER})' > $${KUBECONFIG}

itest-cluster-stop:
	docker container exec $$(cat $${CIDFILE}) /bin/bash -c 'k3d delete --all' || echo -n
	docker container rm -f $$(cat $${CIDFILE}) &>/dev/null || echo -n

itest-operator-start:
	echo "operator-start $${KUBECONFIG} $${TEST_OPERATOR_NS}"

itest-operator-stop:
	echo "operator-stop $${KUBECONFIG} $${TEST_OPERATOR_NS}"

itest-cleanup:
	echo "cleanup $${KUBECONFIG} $${TEST_OPERATOR_NS}"

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
