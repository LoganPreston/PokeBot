package bot

import (
	"encoding/json"
	"io/ioutil"
)

var QuestionAry []fact
var (
	question string
	answer   string
)

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
		return err
	}

	//unmarshall the file into the main struct. handle err.
	if err = json.Unmarshal(file, &info); err != nil {
		return err
	}

	QuestionAry = info.Facts

	return nil
}

func replyToTriviaQuestionMessage() string {
	randIdx := getRandomIntBetween(0, len(QuestionAry))
	question = QuestionAry[randIdx].Question
	answer = QuestionAry[randIdx].Answer
	return question
}

func replyToTriviaAnswerMessage() string {
	if answer == "" {
		answer = "Try prompting a question first!"
	}
	return answer
}
