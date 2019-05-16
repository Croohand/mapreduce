package commands

func Copy(in, out string) {
	Merge([]string{in}, out)
}
