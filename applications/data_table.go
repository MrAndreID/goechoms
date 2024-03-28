package applications

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func (app *Application) DataTable(ctx context.Context, queryBuilder *gorm.DB, searchField []string, orderField string, sortField string, defaultOrder string, defaultSort string, page int, length *int, search string, filter bool) (filterStatus bool) {
	var (
		db    *gorm.DB = app.Database.WithContext(ctx)
		limit int      = 10
	)

	if search != "" {
		filter = true

		if len(searchField) == 1 {
			statement := fmt.Sprintf("lower(%s) like ?", searchField[0])

			db = db.Where(statement, fmt.Sprintf("%%%s%%", strings.TrimLeft(strings.TrimRight(strings.ToLower(string(search)), " "), " ")))
		} else {
			for _, value := range searchField {
				statement := fmt.Sprintf("lower(%s) like ?", value)

				db = db.Or(statement, fmt.Sprintf("%%%s%%", strings.TrimLeft(strings.TrimRight(strings.ToLower(string(search)), " "), " ")))
			}
		}
	}

	if page <= 0 {
		page = 1
	}

	if length == nil || cast.ToInt(length) <= 0 {
		*length = limit
	} else {
		limit = *length
	}

	pagination := (page * limit) - limit

	order := orderField

	if order == "" {
		order = defaultOrder
	}

	sort := sortField

	if sort == "" {
		sort = defaultSort
	}

	orderByQuery := fmt.Sprintf("%s %s", order, sort)

	queryBuilder.Where(db).Order(orderByQuery).Limit(limit).Offset(pagination)

	return filter
}
