package wrrors

type Wrror struct {
	err string
	any bool
}

var subject string = "Unknown"

func SetSubject(subj string) {
	subject = subj
}

func New(pref string) Wrror {
	return Wrror{subject + "." + pref + ": ", false}
}

func (w Wrror) Wrap(err error) error {
	if err == nil {
		if w.any {
			return w
		}
		return nil
	}
	return Wrror{w.err + err.Error(), true}
}

func (w Wrror) Error() string {
	return w.err
}

func (w Wrror) WrapS(s string) Wrror {
	return Wrror{w.err + s, true}
}

func (w Wrror) SWrap(s string) string {
	return w.err + s
}
