package utils

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/xid"
)

// FormData is the global form fill data instance
var FormData FormFillData

func init() {
	FormData = DefaultFormFillData
}

// FormFillData contains suggestions for form filling
type FormFillData struct {
	Email       string `yaml:"email"`
	Color       string `yaml:"color"`
	Password    string `yaml:"password"`
	PhoneNumber string `yaml:"phone"`
	Placeholder string `yaml:"placeholder"`
}

var DefaultFormFillData = FormFillData{
	Email:       fmt.Sprintf("%s@katanacrawler.io", xid.New().String()),
	Color:       "#e66465",
	Password:    "katanaP@assw0rd1",
	PhoneNumber: "2124567890",
	Placeholder: "katana",
}

// FormParam is an input for a form field
type FormParam struct {
	Type       string
	Name       string
	Value      string
	Attributes map[string]string
}

// FormParamFillSuggestions returns a list of form filling suggestions
// for inputs returning the specified recommended values.
func FormParamFillSuggestions(inputs []FormParam) map[string]string {
	data := make(map[string]string)

	// Fill checkboxes and radioboxes first or default values first
	for _, input := range inputs {
		switch input.Type {
		case "radio":
			// Use a single radio name per value
			if _, ok := data[input.Name]; !ok {
				data[input.Name] = input.Value
			}
		case "checkbox":
			data[input.Name] = input.Value

		default:
			// If there is a value, use it for the input. Else
			// infer the values based on input types.
			if input.Value != "" {
				data[input.Name] = input.Value
			}
		}
	}

	// Fill rest of the inputs based on their types or name and ids
	for _, input := range inputs {
		if input.Value != "" {
			continue
		}

		switch input.Type {
		case "email":
			data[input.Name] = FormData.Email
		case "color":
			data[input.Name] = FormData.Color
		case "number", "range":
			var err error
			var max, min, step, val int

			if min, err = strconv.Atoi(input.Attributes["min"]); err != nil {
				min = 1
			}
			if max, err = strconv.Atoi(input.Attributes["max"]); err != nil {
				max = 10
			}
			if step, err = strconv.Atoi(input.Attributes["step"]); err != nil {
				step = 1
			}
			val = min + step
			if val > max {
				val = max - step
			}
			data[input.Name] = strconv.Itoa(val)
		case "password":
			data[input.Name] = FormData.Password
		case "tel":
			data[input.Name] = FormData.Password
		default:
			data[input.Name] = FormData.Placeholder
		}
	}
	return data
}

// ConvertGoquerySelectionToFormParam converts goquery selection to form input
func ConvertGoquerySelectionToFormParam(item *goquery.Selection) FormParam {
	param := FormParam{Attributes: make(map[string]string)}

	param.Name, _ = item.Attr("name")

	if item.Is("input") || item.Is("button") {
		param.Type = item.AttrOr("type", "")
		if item.AttrOr("type", "") == "checkbox" {
			param.Value = item.AttrOr("value", "on")
		} else {
			param.Value = item.AttrOr("value", "")
		}
	} else if item.Is("select") {
		param.Type = "radio"
		param.Value = item.Children().Last().AttrOr("value", "")
	} else if item.Is("texarea") {
		param.Type = "text"
		param.Value = item.Text()
	}

	for _, attribute := range item.Nodes[0].Attr {
		param.Attributes[attribute.Key] = attribute.Val
	}

	return param
}
