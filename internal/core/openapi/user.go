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
	"net/http"
	"strconv"
	"time"

	"github.com/acmestack/envcd/internal/core/storage/dao"
	"github.com/acmestack/envcd/internal/pkg/entity"
	"github.com/acmestack/envcd/pkg/entity/result"
	"github.com/acmestack/godkits/gox/stringsx"
	"github.com/acmestack/godkits/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// loginParam Login
type loginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// userParam Create User Param
type userParam struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Identity int    `json:"identity"`
	State    bool   `json:"state"`
}

const (
	// hmacSecret secret
	hmacSecret = "9C035514A15F78"
	userIdKey  = "userId"
	tokenKey   = "token"
)

// claims claims
type claims struct {
	*jwt.RegisteredClaims
	userId   int
	userName string
}

// newJWTToken secret
func newJWTToken(authClaims claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		// todo
		return ""
	}
	return tokenString
}

func (openapi *Openapi) login(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		param := loginParam{}
		if err := ginCtx.ShouldBindJSON(&param); err != nil {
			log.Error("Bind error, %v", err)
			return result.InternalServerErrorFailure("Illegal params !")
		}

		users, err := dao.New(openapi.storage).SelectUser(entity.User{
			Name: param.Username,
		})
		if err != nil {
			log.Error("Query User error: %v", err)
			// todo error code : result.Failure0(code, message, httpStatusCode)
			return result.Failure("System Error!", http.StatusBadRequest)
		}

		if len(users) == 0 {
			// todo error code : result.Failure0(code, message, httpStatusCode)
			log.Error("User does not exist : %v", param)
			return result.Failure("User does not exist!", http.StatusOK)
		}
		user := users[0]
		if saltPassword(param.Password, user.Salt) != user.Password {
			// todo error code : result.Failure0(code, message, httpStatusCode)
			return result.Failure("password error!", http.StatusOK)
		}
		token := newJWTToken(claims{
			RegisteredClaims: &jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			},
			userId:   user.Id,
			userName: user.Name,
		})
		return result.Success(map[string]interface{}{
			userIdKey: user.Id,
			tokenKey:  token,
		})
	})
}

func (openapi *Openapi) logout(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		// UserDao.save(),
		// LogDao.save()
		return nil
	})
}

func (openapi *Openapi) createUser(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		param := userParam{}
		if er := ginCtx.ShouldBindJSON(&param); er != nil {
			log.Error("Bind error, %v", er)
			return result.InternalServerErrorFailure("Illegal params !")
		}
		daoApi := dao.New(openapi.storage)
		// check if the user already exists in the database
		users, er := daoApi.SelectUser(entity.User{
			Name: param.Name,
		})
		if er != nil {
			log.Error("Query User error: %v", er)
			return result.InternalServerErrorFailure("System Error!")
		}
		if len(users) > 0 {
			log.Error("User Has exists: %v", users)
			return result.InternalServerErrorFailure("User Has Exists!")
		}
		// generate database password by salt
		salt := randomSalt()
		password := saltPassword(param.Password, salt)
		state := 1
		if !param.State {
			state = 2
		}
		user := entity.User{
			Name:      param.Name,
			Password:  password,
			Salt:      salt,
			Identity:  param.Identity,
			State:     state,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// save user
		if _, _, err := daoApi.InsertUser(user); err != nil {
			log.Error("insert error=%v", err)
			return result.InternalServerErrorFailure("Save User Error!")
		}
		// fixme update success message or response token and id ?
		return result.Success("ok")
	})
}

func (openapi *Openapi) updateUser(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		return nil
	})
}

func (openapi *Openapi) user(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		id := stringsx.ToInt(ginCtx.Param("userId"))
		param := entity.User{Id: id}
		// query user by param
		users, err := dao.New(openapi.storage).SelectUser(param)
		if err != nil {
			log.Error("select user error = %v", err)
			return result.InternalServerErrorFailure("Get User Error!")
		}
		if len(users) == 0 {
			log.Error("User does not exist : %v", param)
			return result.Failure("User does not exist!", http.StatusOK)
		}
		return result.Success(userVO{
			Id:        users[0].Id,
			Name:      users[0].Name,
			Identity:  users[0].Identity,
			State:     users[0].State,
			CreatedAt: users[0].CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: users[0].UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	})
}

type pageUserVO struct {
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	List     []userVO `json:"list"`
}

type userVO struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Identity  int    `json:"identity"`
	State     int    `json:"state"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (openapi *Openapi) removeUser(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		return nil
	})
}

func (openapi *Openapi) users(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		// receive params from request
		page := stringsx.ToInt(ginCtx.Query("page"))
		pageSize := stringsx.ToInt(ginCtx.Query("pageSize"))
		if page == 0 {
			page = 1
		}
		if pageSize == 0 {
			pageSize = 10
		}
		nameParam := ginCtx.Query("name")

		// construct sql and params
		builder := stringsx.Builder{}
		position := 0
		builder.JoinString("SELECT id, `name`, identity, state, created_at, updated_at FROM `user`")
		params := []interface{}{}
		if !stringsx.Empty(nameParam) {
			builder.JoinString(" where name like '%${", strconv.Itoa(position), "}%'")
			params = append(params, nameParam)
			position++
		}
		builder.JoinString(" limit ${", strconv.Itoa(position))
		position++
		builder.JoinString("},${", strconv.Itoa(position), "}")
		position++
		params = append(params, strconv.Itoa((page-1)*pageSize), strconv.Itoa(pageSize))

		// query users by param
		ret := []entity.User{}
		err := openapi.storage.NewSession().Select(builder.String()).Param(params...).Result(&ret)
		if err != nil {
			log.Error("select users error = %v", err)
			return result.InternalServerErrorFailure("Get Users Error!")
		}
		return result.Success(pageUserVO{
			page, pageSize, userTransfer(ret),
		})
	})
}

func userTransfer(users []entity.User) []userVO {
	back := []userVO{}
	if users == nil || len(users) == 0 {
		return back
	}
	for _, user := range users {
		back = append(back, userVO{
			Id:        user.Id,
			Name:      user.Name,
			Identity:  user.Identity,
			State:     user.State,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return back
}

func (openapi *Openapi) userScopeSpaces(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		return nil
	})
}

func (openapi *Openapi) userDictionaries(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		return nil
	})
}

func (openapi *Openapi) userDictionariesUnderScopeSpace(ginCtx *gin.Context) {
	openapi.response(ginCtx, nil, func() *result.EnvcdResult {
		fmt.Println("hello world")
		return nil
	})
}
