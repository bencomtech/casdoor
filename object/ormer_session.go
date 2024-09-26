// Copyright 2023 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"fmt"
	"strings"

	"github.com/casdoor/casdoor/conf"
	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/xorm"
)

func GetSession(owner string, offset, limit int, field, value, sortField, sortOrder string) *xorm.Session {
	session := ormer.Engine.Prepare()
	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}
	if owner != "" {
		session = session.And("owner=?", owner)
	}
	if field != "" && value != "" {
		if util.FilterField(field) {
			session = session.And(fmt.Sprintf("%s like ?", util.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
		}
	}
	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}
	if sortOrder == "ascend" {
		session = session.Asc(util.SnakeString(sortField))
	} else {
		session = session.Desc(util.SnakeString(sortField))
	}
	return session
}

func GetSessionFilterCriteria(owner string, offset, limit int, filterCriteria []*FilterCriteria, sortField, sortOrder string) *xorm.Session {
	session := ormer.Engine.Prepare()
	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}
	if owner != "" {
		session = session.And("owner=?", owner)
	}
	for _, criteria := range filterCriteria {
		if util.FilterField(criteria.Field) {
			session = session.And(fmt.Sprintf("%s like ?", util.SnakeString(criteria.Field)), fmt.Sprintf("%%%s%%", criteria.Value))
		}
	}
	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}
	if sortOrder == "ascend" {
		session = session.Asc(util.SnakeString(sortField))
	} else {
		session = session.Desc(util.SnakeString(sortField))
	}
	return session
}

func GetSessionForUser(owner string, offset, limit int, field, value, sortField, sortOrder string) *xorm.Session {
	session := ormer.Engine.Prepare()
	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}
	if owner != "" {
		if offset == -1 {
			session = session.And("owner=?", owner)
		} else {
			session = session.And("a.owner=?", owner)
		}
	}
	if field != "" && value != "" {
		if util.FilterField(field) {
			if offset != -1 {
				field = fmt.Sprintf("a.%s", field)
			}
			session = session.And(fmt.Sprintf("%s like ?", util.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
		}
	}
	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}

	tableNamePrefix := conf.GetConfigString("tableNamePrefix")
	tableName := tableNamePrefix + "user"
	if offset == -1 {
		if sortOrder == "ascend" {
			session = session.Asc(util.SnakeString(sortField))
		} else {
			session = session.Desc(util.SnakeString(sortField))
		}
	} else {
		if sortOrder == "ascend" {
			session = session.Alias("a").
				Join("INNER", []string{tableName, "b"}, "a.owner = b.owner and a.name = b.name").
				Select("b.*").
				Asc("a." + util.SnakeString(sortField))
		} else {
			session = session.Alias("a").
				Join("INNER", []string{tableName, "b"}, "a.owner = b.owner and a.name = b.name").
				Select("b.*").
				Desc("a." + util.SnakeString(sortField))
		}
	}

	return session
}

func GetSessionForUserFilterCriteria(owner string, offset, limit int, filterCriteria []*FilterCriteria, sortField, sortOrder string) *xorm.Session {
	session := ormer.Engine.Prepare()
	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}
	if owner != "" {
		if offset == -1 {
			session = session.And("owner=?", owner)
		} else {
			session = session.And("a.owner=?", owner)
		}
	}
	for _, criteria := range filterCriteria {
		field := criteria.Field
		if strings.Contains(field, ".") {
			parts := strings.Split(field, ".")
			session = session.And(fmt.Sprintf("JSON_EXTRACT(a.%s, '$.%s') = ?", util.SnakeString(parts[0]), parts[1]), criteria.Value)
		} else if util.FilterField(criteria.Field) {
			if offset != -1 {
				field = fmt.Sprintf("a.%s", criteria.Field)
			}
			session = session.And(fmt.Sprintf("%s like ?", util.SnakeString(field)), fmt.Sprintf("%%%s%%", criteria.Value))
		}
	}
	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}

	tableNamePrefix := conf.GetConfigString("tableNamePrefix")
	tableName := tableNamePrefix + "user"
	if offset == -1 {
		if sortOrder == "ascend" {
			session = session.Asc(util.SnakeString(sortField))
		} else {
			session = session.Desc(util.SnakeString(sortField))
		}
	} else {
		if sortOrder == "ascend" {
			session = session.Alias("a").
				Join("INNER", []string{tableName, "b"}, "a.owner = b.owner and a.name = b.name").
				Select("b.*").
				Asc("a." + util.SnakeString(sortField))
		} else {
			session = session.Alias("a").
				Join("INNER", []string{tableName, "b"}, "a.owner = b.owner and a.name = b.name").
				Select("b.*").
				Desc("a." + util.SnakeString(sortField))
		}
	}

	return session
}
