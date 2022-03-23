package face

import (
	"errors"
	"strings"
)

// ErrInvalidTopicEmptyString is the error returned when a topic string
// is passed in that is 0 length
// var ErrInvalidTopicEmptyString = errors.New("invalid Topic; empty string")
// var ErrInvalidTopicEmptyLevel = errors.New("invalid Topic; empty level")
// var ErrInvalidTopicMatchError = errors.New("invalid Topic; not match")
// var ErrInvalidTopicMultilevel = errors.New("invalid Topic; multi-level wildcard must be last level")

var ErrInvalidPatterunTopic = errors.New("invalid Topic; publish topic should not be pattern")
var ErrInvalidTopicFormat = errors.New("Invalid topic")

// // ErrInvalidTopicEmptyString is the error returned when a topic string
// // is passed in that is 0 length
// var ErrInvalidTopicEmptyString = errors.New("invalid Topic; empty string")

// // ErrInvalidTopicMultilevel is the error returned when a topic string
// // is passed in that has the multi level wildcard in any position but
// // the last
// var ErrInvalidTopicMultilevel = errors.New("invalid Topic; multi-level wildcard must be last level")

const _topicLevelExp = "^[0-9a-zA-Z_.:-]+$"

func ValidatePattern(pattern string) error {
	if len(pattern) == 0 {
		return ErrInvalidTopicFormat
	}
	levels := strings.Split(pattern, "/")
	for i, level := range levels {
		if level == "#" && i != len(levels)-1 {
			return ErrInvalidTopicFormat
		}
	}
	return nil
	// Skip char check
	// levels := strings.Split(pattern, "/")
	// for i, level := range levels {
	// 	if level == "" {
	// 		return ErrInvalidTopicFormat
	// 	}
	// 	if i == 0 && len(levels) > 1 {
	// 		if level == "$SYS" || level == "$USR" || level == "$CID" {
	// 			continue
	// 		}
	// 	}
	// 	if level == "+" {
	// 		continue
	// 	}
	// 	if level == "#" {
	// 		if i != len(levels)-1 {
	// 			return ErrInvalidTopicFormat
	// 		}
	// 		continue
	// 	}
	// 	match, err := regexp.MatchString(_topicLevelExp, level)
	// 	if err != nil {
	// 		return err
	// 	} else if !match {
	// 		return ErrInvalidTopicFormat
	// 	}
	// }
	// return nil
}

func ValidateTopic(topic string) error {
	if strings.Contains(topic, "#") || strings.Contains(topic, "+") {
		return ErrInvalidPatterunTopic
	}
	return ValidatePattern(topic)
	// if strings.HasPrefix(topic, "/") || strings.HasSuffix(topic, "/") {
	// 	return ErrInvalidTopicFormat
	// }
	// levels := strings.Split(topic, "/")
	// for i, level := range levels {
	// 	if level == "" {
	// 		return ErrInvalidTopicFormat
	// 	}
	// 	if i == 0 && len(levels) > 1 {
	// 		if level == "$SYS" || level == "$USR" || level == "$CID" {
	// 			continue
	// 		}
	// 	}
	// 	match, err := regexp.MatchString(_topicLevelExp, level)
	// 	if err != nil {
	// 		return err
	// 	} else if !match {
	// 		return ErrInvalidTopicFormat
	// 	}
	// }
	// return nil
}

func MatchTopic(pattern, topic string) error {
	if err := ValidatePattern(pattern); err != nil {
		return err
	}
	if err := ValidateTopic(topic); err != nil {
		return err
	}
	patterns := strings.Split(pattern, "/")
	topics := strings.Split(topic, "/")
	for i := 0; i < len(patterns); i++ {
		if patterns[i] == "+" {
			continue
		}
		if patterns[i] == "#" {
			return nil
		}
		if patterns[i] != topics[i] {
			return errors.New("not match")
		}
	}
	// regStr := "^" + pattern
	// if !strings.HasSuffix(pattern, "#") {
	// 	regStr += "$"
	// }
	// regStr = strings.ReplaceAll(regStr, "#", "")
	// regStr = strings.ReplaceAll(regStr, "+", "[0-9a-zA-Z_.:-]+")
	// if matched, _ := regexp.MatchString(regStr, topic); matched {
	// 	return true
	// }
	return nil
}
