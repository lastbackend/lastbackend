#!/bin/bash

mode="local"
src="/opt/src/github.com/lastbackend"
dst="/go/src/github.com/lastbackend"

function host() {
	if [[ "$mode" == "local" ]]; then
		eval "$(docker-machine env $1)"
	else
		export DOCKER_HOST="$1.lstbknd.net:2376"
	fi
}

function start() {
		case $1 in
		"api")
			host "genesis"
			app=lastbackend
			docker run -d -it --restart=always \
				-v "${src}/${app}":"${dst}"/lastbackend \
				-v "${src}/${app}/contrib/config.yml":/etc/lastbackend/config.yml \
				-v /etc/nginx/ssl:/ssl:rw \
				-w "${dst}"/lastbackend \
				--name=lastbackend-api \
				--net=host \
				lastbackend/lastbackend go run ./cmd/kit/kit.go api -c /etc/lastbackend/config.yml
			;;

		"ctl")
			host "genesis"
			app=lastbackend
			docker run -d -it --restart=always \
				-v "${src}/${app}":"${dst}"/lastbackend \
				-v "${src}/${app}"/contrib/config.yml:/etc/lastbackend/config.yml \
				-w "${dst}"/lastbackend \
				--name=lastbackend-ctl --net=host \
				lastbackend/lastbackend go run ./cmd/kit/kit.go ctl -c /etc/lastbackend/config.yml
			;;

		"sdl")
			host "genesis"
			app=lastbackend
			docker run -d -it --restart=always \
				-v "${src}/${app}":"${dst}"/lastbackend \
				-v "${src}/${app}"/contrib/config.yml:/etc/lastbackend/config.yml \
				-w "${dst}"/lastbackend \
				--name=lastbackend-sdl --net=host \
				lastbackend/lastbackend go run ./cmd/kit/kit.go sdl -c /etc/lastbackend/config.yml
			;;

		"node-00")
			host "node-00"
			app=lastbackend
			docker run -d -it --restart=always \
				-v "${src}/${app}":"${dst}"/lastbackend \
				-v "${src}/${app}"/contrib/node.yml:/etc/lastbackend/config.yml \
				-v /var/run/docker.sock:/var/run/docker.sock \
				-v /var/lib/lastbackend:/var/lib/lastbackend \
				-w "${dst}"/lastbackend \
				--net=host --privileged --name="lastbackend-node-00" \
				lastbackend/lastbackend go run ./cmd/node/node.go -c /etc/lastbackend/config.yml
			;;
		"node-01")
			host "node-01"
			app=lastbackend
			docker run -d -it --restart=always \
				-v "${src}/${app}":"${dst}"/lastbackend \
				-v "${src}/${app}"/contrib/node.yml:/etc/lastbackend/config.yml \
				-v /var/run/docker.sock:/var/run/docker.sock \
				-v /var/lib/lastbackend:/var/lib/lastbackend \
				-w "${dst}"/lastbackend \
				--net=host --privileged --name="lastbackend-node-01" \
				lastbackend/lastbackend go run ./cmd/node/node.go -c /etc/lastbackend/config.yml
			;;
		esac
		;;
	esac
}

function stop() {

	case $2 in
	"node-00")
		host "node-00"
		docker rm -vf lastbackend-node-00
		;;
	"node-01")
		host "node-01"
		docker rm -vf lastbackend-node-01
		;;
	*)
		host "genesis"
		docker rm -vf "$1-$2"
		;;
	esac
}

function logs() {

	case $2 in
	"node-00")
		host "node-00"
		docker logs "${@:3}" lastbackend-node-00
		;;
	"node-01")
		host "node-01"
		docker logs "${@:3}" lastbackend-node-01
		;;
	*)
		host "genesis"
		docker logs "${@:3}" "$1-$2"
		;;
	esac
}

function restart() {
	stop ${@}
	start ${@}
}

function genesis() {
	case "$1" in
	"start")
		start "genesis" "${@:2}"
		logs "genesis" "${@:2}" -f
		;;
	"stop") stop "genesis" "${@:2}" ;;
	"logs") logs "genesis" "${@:2}" ;;
	"restart")
		stop "genesis" "${@:2}"
		start "genesis" "${@:2}"
		logs "genesis" "${@:2}" -f
		;;
	*)
		echo "unknown command"
		exit 1
		;;
	esac
}

function cluster() {
	case "$1" in
	"start")
		start "lastbackend" "${@:2}"
		logs "lastbackend" "${@:2}" -f
		;;
	"stop") stop "lastbackend" "${@:2}" ;;
	"logs") logs "lastbackend" "${@:2}" ;;
	"restart")
		stop "lastbackend" "${@:2}"
		start "lastbackend" "${@:2}"
		logs "lastbackend" "${@:2}" -f
		;;
	*)
		echo "unknown command"
		exit 1
		;;
	esac
}

case "$1" in
"genesis") genesis "${@:2}" ;;
"cluster") cluster "${@:2}" ;;
*)
	echo "unknown partial"
	exit 1
	;;
esac
