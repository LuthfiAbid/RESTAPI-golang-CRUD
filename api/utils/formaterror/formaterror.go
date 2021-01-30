package formaterror

import(
	"errors"
	"strings"
)

func FormatError(err string) error{
	if strings.Contains(err, "username"){
		return errors.New("Username exists!")
	}
	if strings.Contains(err, "hashedPassword"){
		return errors.New("Incorrect Password")
	}

	return errors.New("Incorrect Details")
}