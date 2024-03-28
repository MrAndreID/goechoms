package applications

import (
	"encoding/json"

	jsonIterator "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type CustomJSON struct{}

func (app *Application) NewCustomJSON() *CustomJSON {
	return &CustomJSON{}
}

func (cjson *CustomJSON) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())

	if indent != "" {
		enc.SetIndent("", indent)
	}

	return enc.Encode(i)
}

func (cjson *CustomJSON) Deserialize(c echo.Context, i interface{}) error {
	var (
		ConfigCompatibleWithStandardLibrary = jsonIterator.Config{CaseSensitive: true}.Froze()
		customJSON                          = ConfigCompatibleWithStandardLibrary
	)

	return customJSON.NewDecoder(c.Request().Body).Decode(i)
}
