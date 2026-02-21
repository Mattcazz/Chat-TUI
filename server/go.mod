module github.com/Mattcazz/Chat-TUI/server

go 1.25.1

require (
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/chi/v5 v5.2.4
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.11.1
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/cyphar/filepath-securejoin v0.5.1 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/sys/atomicwriter v0.1.0 // indirect
	github.com/moby/sys/mount v0.3.4 // indirect
	github.com/moby/sys/mountinfo v0.7.2 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/user v0.4.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/opencontainers/selinux v1.13.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.etcd.io/bbolt v1.4.3 // indirect
	golang.org/x/sys v0.38.0 // indirect
	gotest.tools/v3 v3.5.2 // indirect
)

require github.com/golang-jwt/jwt/v5 v5.3.1

require (
	github.com/Mattcazz/Chat-TUI/pkg v0.0.0
	github.com/docker/docker v28.3.3+incompatible
	golang.org/x/crypto v0.45.0
)

replace github.com/Mattcazz/Chat-TUI/pkg => ../pkg
