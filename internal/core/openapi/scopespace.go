/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openapi

import (
	"fmt"

	"github.com/acmestack/envcd/pkg/entity/result"
	"github.com/gin-gonic/gin"
)

func (openapi *Openapi) scopeSpace(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		// create config
		// ConfigDao.save();
		// go LogDao.save()
		// openapi.exchange.Put("key", "value")
		return nil
	})
}

func (openapi *Openapi) createScopeSpace(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		// create config
		// ConfigDao.save();
		// go LogDao.save()
		// openapi.exchange.Put("key", "value")
		return nil
	})
}

func (openapi *Openapi) updateScopeSpace(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		// create config
		// ConfigDao.save();
		// go LogDao.save()
		// openapi.exchange.Put("key", "value")
		return nil
	})
}

func (openapi *Openapi) removeScopeSpace(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		// create config
		// ConfigDao.save();
		// go LogDao.save()
		// openapi.exchange.Put("key", "value")
		return nil
	})
}

func (openapi *Openapi) scopespaces(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		// create config
		// ConfigDao.save();
		// go LogDao.save()
		// openapi.exchange.Put("key", "value")
		return nil
	})
}
