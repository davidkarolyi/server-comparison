package docker

// IDocker handles docker containers related to a single benchmark.
type IDocker interface {
	BuildServerImage() error
	StartServerContainer() (chan *Stats, error)
	Clean() error
}
