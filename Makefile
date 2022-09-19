build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./bin/go-oom-guard ./cmd/

docker:
	docker build . -t postgres_oom_guarded

docker-run:
	docker run --rm \
		--volume /sys/fs/cgroup:/sys/fs/cgroup  \
		-p 6432:5432 \
		--memory 134217728 \
		--name postgres_oom_guarded postgres_oom_guarded

run-oom-guard:
	docker exec -i postgres_oom_guarded ./go-oom-guard

test-oom-execute:
	psql -U postgres -h localhost -p 6432 -f demo/generate_oom.sql

test-oom-parse:
	python3 demo/generate_oom.py