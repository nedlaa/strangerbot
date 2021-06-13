package main

import "encoding/json"

const (
	GenderOptionType    = 1
	MatchModeOptionType = 2
	GoalOptionType      = 3

	GenderOptionMale       = 1
	GenderOptionMaleText   = "Male"
	GenderOptionFemale     = 2
	GenderOptionFemaleText = "Female"

	MatchModeOptionMale         = 1
	MatchModeOptionMaleText     = "Male"
	MatchModeOptionFemale       = 2
	MatchModeOptionFemaleText   = "Female"
	MatchModeOptionAnything     = 0
	MatchModeOptionAnythingText = "Anything"

	GoalOptionDating      = 1
	GoalOptionDatingText  = "Dating"
	GoalOptionFriends     = 2
	GoalOptionFriendsText = "Friends"
)

var (
	GenderMale = KeyboardCallbackData{
		OptionType:  GenderOptionType,
		OptionValue: GenderOptionMale,
	}.toString()
	GenderFemale = KeyboardCallbackData{
		OptionType:  GenderOptionType,
		OptionValue: GenderOptionFemale,
	}.toString()
	MatchModeMale = KeyboardCallbackData{
		OptionType:  MatchModeOptionType,
		OptionValue: MatchModeOptionMale,
	}.toString()
	MatchModeFemale = KeyboardCallbackData{
		OptionType:  MatchModeOptionType,
		OptionValue: MatchModeOptionFemale,
	}.toString()
	MatchModeAnything = KeyboardCallbackData{
		OptionType:  MatchModeOptionType,
		OptionValue: MatchModeOptionAnything,
	}.toString()
	GoalDating = KeyboardCallbackData{
		OptionType:  GoalOptionType,
		OptionValue: GoalOptionDating,
	}.toString()
	GoalFriends = KeyboardCallbackData{
		OptionType:  GoalOptionType,
		OptionValue: GoalOptionFriends,
	}.toString()
)

type KeyboardCallbackData struct {
	OptionType  int `json:"ot"`
	OptionValue int `json:"ov"`
}

func (v KeyboardCallbackData) toString() string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

func (v KeyboardCallbackData) GetOptionText() string {
	return GetOptionText(v.OptionType, v.OptionValue)
}

func GetOptionText(optionType int, optionValue int) string {
	switch optionType {
	case GenderOptionType:

		switch optionValue {
		case GenderOptionMale:
			return GenderOptionMaleText
		case GenderOptionFemale:
			return GenderOptionFemaleText
		}

	case MatchModeOptionType:
		switch optionValue {
		case MatchModeOptionMale:
			return GenderOptionMaleText
		case MatchModeOptionFemale:
			return GenderOptionFemaleText
		case MatchModeOptionAnything:
			return MatchModeOptionAnythingText
		}

	case GoalOptionType:
		switch optionValue {
		case GoalOptionDating:
			return GoalOptionDatingText
		case GoalOptionFriends:
			return GoalOptionFriendsText
		}
	}

	return ""
}
