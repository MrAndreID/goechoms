package middlewares

import (
	"github.com/labstack/echo/v4"
)

func (cm *CustomMiddleware) SetRouteList(routes []*echo.Route) map[string]map[string]string {
	var routeList map[string]map[string]string = make(map[string]map[string]string)

	for _, element := range routes {
		if _, ok := routeList[element.Path]; !ok {
			routeList[element.Path] = make(map[string]string)
		}

		if element.Name == "github.com/labstack/echo/v4.glob..func1" {
			continue
		}

		routeList[element.Path][element.Method] = element.Name
	}

	return routeList
}
