package fromenv

// os.Getenv mocking concept from alexellis
// https://gist.github.com/alexellis/adc67eb022b7fdca31afc0de6529e5ea
type fakeEnv struct {
	values map[string]string
}

func (env fakeEnv) Getenv(k string) string {
	return env.values[k]
}

func (env fakeEnv) Setenv(k string, v string) {
	env.values[k] = v
}

func (env fakeEnv) Unsetenv(k string) {
	delete(env.values, k)
}

func NewFakeEnv() fakeEnv {
	e := fakeEnv{}
	e.values = make(map[string]string)
	return e
}
