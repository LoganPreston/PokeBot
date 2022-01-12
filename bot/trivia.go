package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var QuestionAry []fact

type trivia struct {
	Facts []fact `json:"trivia"`
}

type fact struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func ReadTrivia() error {
	var (
		info *trivia
		file []byte
		err  error
	)
	//Read the configuration file, token and prefix. Handle the err
	if file, err = ioutil.ReadFile("./trivia.json"); err != nil {
		fmt.Println(err.Error())
		return err
	}

	//unmarshall the file into the main struct. handle err.
	if err = json.Unmarshal(file, &info); err != nil {
		fmt.Println(err.Error())
		return err
	}

	QuestionAry = info.Facts

	return nil
}

func replyToTriviaQuestionMessage() string {
	return ""
}

func replyToTriviaAnswerMessage() string {
	return ""
}
