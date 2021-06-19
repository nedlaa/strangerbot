package main

import "encoding/json"

const (
	GenderOptionType    = 1
	MatchModeOptionType = 2
	GoalOptionType      = 3

	GenderOptionMale           = 1
	GenderOptionMaleText       = "Male"
	GenderOptionMaleNoteText   = ""
	GenderOptionFemale         = 2
	GenderOptionFemaleText     = "Female"
	GenderOptionFemaleNoteText = ""

	MatchModeOptionMale             = 1
	MatchModeOptionMaleText         = "Male"
	MatchModeOptionMaleNoteText     = ""
	MatchModeOptionFemale           = 2
	MatchModeOptionFemaleText       = "Female"
	MatchModeOptionFemaleNoteText   = ""
	MatchModeOptionAnything         = 0
	MatchModeOptionAnythingText     = "Anything"
	MatchModeOptionAnythingNoteText = ""

	GoalOptionDating          = 1
	GoalOptionDatingText      = "Dating/“Tinder”"
	GoalOptionDatingNoteText  = " (Note that sending anything offensive/inappropriate can result in permanent ban; if anything illegal, such as sending explicit images, occurs, your identity will be retrieved with data handed over to the police)"
	GoalOptionFriends         = 2
	GoalOptionFriendsText     = "Friends/Wholesome Talk Only"
	GoalOptionFriendsNoteText = " (NOTE: anything rude/containing sexual, or any form of abuse results in permanent ban—your identity will be retrieved and handed to the police if anything illegal, including sending of explicit images, occurs)"
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

func (v KeyboardCallbackData) GetOptionNoteText() string {
	return GetNoteText(v.OptionType, v.OptionValue)
}

func (v KeyboardCallbackData) GetOptionText() string {
	return GetOptionText(v.OptionType, v.OptionValue)
}

func (v KeyboardCallbackData) GetOptionFullText() string {
	return GetOptionFullText(v.OptionType, v.OptionValue)
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

func GetNoteText(optionType int, optionValue int) string {
	switch optionType {
	case GenderOptionType:

		switch optionValue {
		case GenderOptionMale:
			return GenderOptionMaleNoteText
		case GenderOptionFemale:
			return GenderOptionFemaleNoteText
		}

	case MatchModeOptionType:
		switch optionValue {
		case MatchModeOptionMale:
			return GenderOptionMaleNoteText
		case MatchModeOptionFemale:
			return GenderOptionFemaleNoteText
		case MatchModeOptionAnything:
			return MatchModeOptionAnythingNoteText
		}

	case GoalOptionType:
		switch optionValue {
		case GoalOptionDating:
			return GoalOptionDatingNoteText
		case GoalOptionFriends:
			return GoalOptionFriendsNoteText
		}
	}

	return ""
}

func GetOptionFullText(optionType int, optionValue int) string {
	switch optionType {
	case GenderOptionType:

		switch optionValue {
		case GenderOptionMale:
			return GenderOptionMaleText + GenderOptionMaleNoteText
		case GenderOptionFemale:
			return GenderOptionFemaleText + GenderOptionFemaleNoteText
		}

	case MatchModeOptionType:
		switch optionValue {
		case MatchModeOptionMale:
			return GenderOptionMaleText + GenderOptionMaleNoteText
		case MatchModeOptionFemale:
			return GenderOptionFemaleText + GenderOptionFemaleNoteText
		case MatchModeOptionAnything:
			return MatchModeOptionAnythingText + MatchModeOptionAnythingNoteText
		}

	case GoalOptionType:
		switch optionValue {
		case GoalOptionDating:
			return GoalOptionDatingText + GoalOptionDatingNoteText
		case GoalOptionFriends:
			return GoalOptionFriendsText + GoalOptionFriendsNoteText
		}
	}

	return ""
}
