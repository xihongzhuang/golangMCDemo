package api_service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"

	"github.com/lithammer/shortuuid/v3"

	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyVersion = errors.New("empty version")
	ErrInvalidEmail = errors.New("invalid maintainer email")
)

type PersonInfo struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type AppMetaData struct {
	Id          string       `yaml:"id"`
	Title       string       `yaml:"title"`
	Version     string       `yaml:"version"`
	Maintainers []PersonInfo `yaml:"maintainers"`
	Company     string       `yaml:"company"`
	Website     string       `yaml:"website"`
	Source      string       `yaml:"source"`
	License     string       `yaml:"license"`
	Description string       `yaml:"description"`
}

func (r *AppMetaData) Print() {
	fmt.Printf("\tId:%s\n\tTitle: %s\n\tVersion: %s\n\tMaintainers: %v\n\tCompany: %s\n\tWebSite: %s\n\tSource: %s\n\tLicense: %s\n\tDescription: %s\n",
		r.Id,
		r.Title,
		r.Version,
		r.Maintainers,
		r.Company,
		r.Website,
		r.Source,
		r.License,
		r.Description,
	)
}

var regexEmail = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)

func IsValidEmail(email string) bool {
	return regexEmail.MatchString(email)
}

func LoadFrom(data []byte) (*AppMetaData, error) {
	mt := AppMetaData{}
	err := yaml.Unmarshal(data, &mt)
	if err != nil {
		return nil, err
	}
	if mt.Version == "" {
		return nil, ErrEmptyVersion
	}
	for _, p := range mt.Maintainers {
		if !IsValidEmail(p.Email) {
			return nil, ErrInvalidEmail
		}
	}
	if mt.Id == "" {
		mt.Id = shortuuid.New()
	}
	return &mt, nil
}

func LoadAppMeta(fileName string) (*AppMetaData, error) {
	yfile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return LoadFrom(yfile)
}

func StructToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("yaml")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
