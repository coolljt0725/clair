// Copyright 2015 clair authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysql

import (
	"time"

	"github.com/coreos/clair/database"
	cerrors "github.com/coreos/clair/utils/errors"
)

func (mySQL *mySQL) insertNamespaceWithTransaction(queryer Queryer, namespace database.Namespace) (int, error) {
	if namespace.Name == "" {
		return 0, cerrors.NewBadRequestError("could not find/insert invalid Namespace")
	}

	if mySQL.cache != nil {
		database.PromCacheQueriesTotal.WithLabelValues("namespace").Inc()
		if id, found := mySQL.cache.Get("namespace:" + namespace.Name); found {
			database.PromCacheHitsTotal.WithLabelValues("namespace").Inc()
			return id.(int), nil
		}
	}

	// We do `defer database.ObserveQueryTime` here because we don't want to observe cached namespaces.
	defer database.ObserveQueryTime("insertNamespace", "all", time.Now())

	var id int
	res, err := queryer.Exec(insertNamespace, namespace.Name, namespace.Name)
	if err != nil {
		return 0, handleError("insertNamespace", err)
	}
	tmpid, err := res.LastInsertId()
	if err != nil {
		return 0, handleError("insertNamespace", err)
	}
	id = int(tmpid)
	if id == 0 {
		err = queryer.QueryRow(soiNamespace, namespace.Name).Scan(&id)
		if err != nil {
			return 0, handleError("soiNamespace", err)
		}
	}
	if mySQL.cache != nil {
		mySQL.cache.Add("namespace:"+namespace.Name, id)
	}

	return id, nil

}

func (mySQL *mySQL) insertNamespace(namespace database.Namespace) (int, error) {
	if namespace.Name == "" {
		return 0, cerrors.NewBadRequestError("could not find/insert invalid Namespace")
	}

	if mySQL.cache != nil {
		database.PromCacheQueriesTotal.WithLabelValues("namespace").Inc()
		if id, found := mySQL.cache.Get("namespace:" + namespace.Name); found {
			database.PromCacheHitsTotal.WithLabelValues("namespace").Inc()
			return id.(int), nil
		}
	}

	// We do `defer database.ObserveQueryTime` here because we don't want to observe cached namespaces.
	defer database.ObserveQueryTime("insertNamespace", "all", time.Now())

	var id int
	res, err := mySQL.Exec(insertNamespace, namespace.Name, namespace.Name)
	if err != nil {
		return 0, handleError("insertNamespace", err)
	}
	tmpid, err := res.LastInsertId()
	if err != nil {
		return 0, handleError("insertNamespace", err)
	}
	id = int(tmpid)
	if id == 0 {
		err = mySQL.QueryRow(soiNamespace, namespace.Name).Scan(&id)
		if err != nil {
			return 0, handleError("soiNamespace", err)
		}
	}
	if mySQL.cache != nil {
		mySQL.cache.Add("namespace:"+namespace.Name, id)
	}

	return id, nil
}

func (mySQL *mySQL) ListNamespaces() (namespaces []database.Namespace, err error) {
	rows, err := mySQL.Query(listNamespace)
	if err != nil {
		return namespaces, handleError("listNamespace", err)
	}
	defer rows.Close()

	for rows.Next() {
		var namespace database.Namespace

		err = rows.Scan(&namespace.ID, &namespace.Name)
		if err != nil {
			return namespaces, handleError("listNamespace.Scan()", err)
		}

		namespaces = append(namespaces, namespace)
	}
	if err = rows.Err(); err != nil {
		return namespaces, handleError("listNamespace.Rows()", err)
	}

	return namespaces, err
}