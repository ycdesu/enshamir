package enshamir

import "github.com/AlecAivazis/survey/v2"

func AskPassword() ([]byte, error) {
	prompt := survey.Password{
		Message: "Encryption password:",
	}
	var pwd string
	if err := survey.AskOne(&prompt, &pwd); err != nil {
		return nil, err
	}
	return []byte(pwd), nil
}
