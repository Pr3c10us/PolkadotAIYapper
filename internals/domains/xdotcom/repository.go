package xdotcom

type Repository interface {
	Tweet(tweet Tweet) (string, error)
}
