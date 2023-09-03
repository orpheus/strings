package threadsvc

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
	"sort"
)

type StringMap struct {
	Default       []*core.String
	Deleted       []*core.String
	Archived      []*core.String
	sourceStrings map[uuid.UUID]*core.String
}

func NewStringMap(sourceStrings []*core.String) *StringMap {
	stringMap := &StringMap{
		sourceStrings: make(map[uuid.UUID]*core.String),
	}
	for _, str := range sourceStrings {
		stringMap.sourceStrings[str.StringId] = str

		// in order of precedence
		if str.Deleted {
			stringMap.Deleted = append(stringMap.Deleted, str)
		} else if str.Archived {
			stringMap.Archived = append(stringMap.Archived, str)
		} else {
			stringMap.Default = append(stringMap.Default, str)
		}
	}

	return stringMap
}

func (s *StringMap) UpdateFrom(strings []*core.String) (map[uuid.UUID]*core.String, error) {
	updatedStrings := make(map[uuid.UUID]*core.String)

	for _, clientString := range strings {
		if _, exists := s.sourceStrings[clientString.StringId]; exists {
			serverString := s.sourceStrings[clientString.StringId]

			if _, exists := updatedStrings[serverString.StringId]; exists {
				return nil, fmt.Errorf("duplicate client string provided for string id %s", serverString.StringId)
			}

			if serverString.Diff(clientString) {
				if serverString.Locked() {
					return nil, ErrStringCannotBeUpdated
				}
				// remember, because we're updating a pointer, the string updates in both the updatedStrings map
				// and the Default array
				serverString.UpdateFromClient(clientString)
				updatedStrings[serverString.StringId] = serverString
			}
		}
	}

	return updatedStrings, nil
}

func (s *StringMap) GetNewStrings(strings []*core.String) []*core.String {
	var newStrings []*core.String

	for _, clientString := range strings {
		if _, exists := s.sourceStrings[clientString.StringId]; !exists {
			newStrings = append(newStrings, clientString)
		}
	}

	return newStrings
}

func (s *StringMap) OrderStrings() error {
	maxRange := len(s.Default)

	var orderedStrings []*core.String
	var unorderedStrings []*core.String

	for _, stringItem := range s.Default {
		if stringItem.Order < 0 {
			return fmt.Errorf("string (%s) order is less than 0", stringItem.StringId)
		}

		if stringItem.Order == 0 {
			unorderedStrings = append(unorderedStrings, stringItem)
			continue
		}

		if stringItem.Order > maxRange {
			return fmt.Errorf("string (%s) order is greater than max range %d", stringItem.StringId, maxRange)
		}

		orderedStrings = append(orderedStrings, stringItem)
	}

	sort.Slice(orderedStrings, func(i, j int) bool {
		return orderedStrings[i].Order < orderedStrings[j].Order
	})

	for index, stringItem := range orderedStrings {
		if stringItem.Order != index+1 {
			return fmt.Errorf("invalid string order")
		}
	}

	for _, newString := range unorderedStrings {
		newString.Order = len(orderedStrings) + 1
		orderedStrings = append(orderedStrings, newString)
	}

	s.Default = orderedStrings

	return nil
}

func (s *StringMap) Include(strings []*core.String) {
	for _, str := range strings {
		if str.Deleted {
			s.Default = append(s.Default, str)
		} else if str.Archived {
			s.Archived = append(s.Archived, str)
		} else {
			s.Default = append(s.Default, str)
		}
	}
}
