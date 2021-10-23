package api_service

import (
	"errors"
	"net/url"
	"strings"

	"gopkg.in/yaml.v3"
)

type InMemoryAppStore struct {
	db map[string]*AppMetaData
}

func NewInMemoryAppStore() *InMemoryAppStore {
	return &InMemoryAppStore{
		db: make(map[string]*AppMetaData),
	}
}

func (s *InMemoryAppStore) Create(yamlPayload []byte) ([]byte, error) {
	entry, err := LoadFrom(yamlPayload)
	if err != nil {
		return nil, err
	}
	s.db[entry.Id] = entry
	encodedStr, err := yaml.Marshal(entry)
	if err != nil {
		return nil, err
	}
	return encodedStr, nil
}

var ErrNotFound = errors.New("metadata not found")

func (s *InMemoryAppStore) Get(id string) ([]byte, error) {
	if entry, ok := s.db[id]; !ok {
		return nil, ErrNotFound
	} else {
		encodedStr, err := yaml.Marshal(entry)
		if err != nil {
			return nil, err
		}
		return encodedStr, nil
	}
}

//Please note this search is inefficient, as there is no index to assist the search, The Time complexity
//is O(N), N is the number of entries in the DB

func matchCondition(m *AppMetaData, params url.Values) bool {
	if len(params) == 0 {
		return true
	}

	mp := StructToMap(m)
	allMatch := true
	for k, v := range params {
		if v2, ok := mp[k]; !ok {
			allMatch = false
			break
		} else {
			found := false
			s := strings.ToLower(v2.(string))
			substr := strings.ToLower(v[0])
			if strings.HasPrefix(substr, "in:") {
				substr = substr[len("in:"):]
				found = strings.Contains(s, substr)
			} else {
				found = s == substr
			}
			if !found {
				allMatch = false
				break
			}
		}
	}

	return allMatch
}

func (s *InMemoryAppStore) GetAll2Yaml(params url.Values) ([]byte, error) {
	if len(s.db) == 0 {
		return []byte(""), nil
	}
	res := make([]*AppMetaData, 0)
	for _, v := range s.db {
		if matchCondition(v, params) {
			res = append(res, v)
		}
	}
	encodedStr, err := yaml.Marshal(res)
	if err != nil {
		return nil, err
	}
	return encodedStr, nil
}

func (s *InMemoryAppStore) Update(id string, yamlPayload []byte) ([]byte, error) {
	if _, ok := s.db[id]; !ok {
		return nil, ErrNotFound
	}
	entry, err := LoadFrom(yamlPayload)
	if err != nil {
		return nil, err
	}
	s.db[id] = entry
	encodedStr, err := yaml.Marshal(entry)
	if err != nil {
		return nil, err
	}
	return encodedStr, nil
}

//partially update

func (s *InMemoryAppStore) Patch(id string, yamlPayload []byte) ([]byte, error) {
	entry, ok := s.db[id]
	if !ok {
		return nil, ErrNotFound
	}
	newMt := AppMetaData{}
	err := yaml.Unmarshal(yamlPayload, &newMt)
	if err != nil {
		return nil, err
	}
	if newMt.Title != "" {
		entry.Title = newMt.Title
	}
	if newMt.Version != "" {
		entry.Version = newMt.Version
	}
	if len(newMt.Maintainers) > 0 {
		entry.Maintainers = newMt.Maintainers
	}
	if newMt.Company != "" {
		entry.Company = newMt.Company
	}
	if newMt.Website != "" {
		entry.Website = newMt.Website
	}
	if newMt.Source != "" {
		entry.Source = newMt.Source
	}
	if newMt.License != "" {
		entry.License = newMt.License
	}
	if newMt.Description != "" {
		entry.Description = newMt.Description
	}
	s.db[id] = entry
	encodedStr, err := yaml.Marshal(entry)
	if err != nil {
		return nil, err
	}
	return encodedStr, nil
}
