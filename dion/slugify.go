package dion

func slugify(s string) []byte {
	var (
		pending []byte
		out     []byte
	)
	for i := 0; i < len(s); i++ {
		char := s[i]
		if char >= 'A' && char <= 'Z' {
			pending = append(pending, char+32) // lowercase it
		} else {
			switch len(pending) {
			case 0:
				out = append(out, char)
			case 1:
				if i == 1 {
					out = append(out, pending[0], char)
				} else {
					out = append(out, '-', pending[0], char)
				}
				pending = nil
			default:
				if i == len(pending) {
					out = append(pending[:len(pending)-1], '-', pending[len(pending)-1], char)
				} else {
					out = append(out, '-')
					out = append(out, pending[:len(pending)-1]...)
					out = append(out, '-', pending[len(pending)-1], char)
				}
				pending = nil
			}
		}
	}
	if len(pending) > 0 {
		if len(out) != 0 {
			out = append(out, '-')
		}
		out = append(out, pending...)
	}
	return out
}
