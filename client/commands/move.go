package commands

func Move(in, out string) {
	Copy(in, out)
	Remove(in)
}
