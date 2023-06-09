package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/ksw2000/catch_cat_server/cats"
	"github.com/ksw2000/catch_cat_server/config"
	"github.com/ksw2000/catch_cat_server/friends"
	"github.com/ksw2000/catch_cat_server/user"
	"github.com/ksw2000/catch_cat_server/util"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// open database
	util.OpenDB()
	defer util.CloseDB()

	// prepare gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.POST("/register", user.PostRegister)
	r.POST("/login", user.PostLogin)
	r.POST("/logout", user.PostLogout)
	r.POST("/friend/invite", friends.PostFriendInvite)
	r.POST("/friends/inviting_me", friends.PostInvitingMeList)
	r.POST("/friends/list", friends.PostFriendsList)
	r.POST("/friends/position", friends.PostFriendsPosition)
	r.POST("/friends/theme_rank", friends.PostFriendRankAtTheme)
	r.POST("/friend/agree", friends.PostFriendAgree)
	r.POST("/friend/decline", friends.PostFriendDecline)
	r.POST("/friend/delete", friends.PostFriendDelete)
	r.POST("/theme", cats.PostTheme)
	r.POST("/user/update/name", user.PostUpdateName)
	r.POST("/user/update/password", user.PostUpdatePassword)
	r.POST("/user/update/email", user.PostUpdateEmail)
	r.POST("/user/update/gps", user.PostUpdateGPS)
	r.POST("/user/update/share_gps", user.PostUpdateShareGPS)
	r.POST("/user/update/profile", user.PostUpdateProfile)
	r.POST("/user/update/last_login", user.PostUpdateLastLogin)
	r.POST("/user/me", user.PostMe)
	r.POST("/cat/catching", cats.PostCatching)
	r.POST("/cat/my_caught_kind", cats.PostCaughtKind)
	r.POST("/upload/profile", uploadProfile)
	r.GET("/theme_list", getThemeList)
	r.Static("/images", "./images")
	r.Static("/icons", "./web/icons")
	r.Static("/assets", "./web/assets")

	// router for website
	r.LoadHTMLFiles("./web/index.html")
	r.StaticFile("/main.dart.js", "./web/main.dart.js")
	r.StaticFile("/flutter.js", "./web/flutter.js")
	r.StaticFile("/manifest.json", "./web/manifest.json")
	r.StaticFile("/favicon.png", "./web/favicon.png")
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	// start server
	r.Run(":8080")
}

// https://stackoverflow.com/questions/29418478/go-gin-framework-cors
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getThemeList(c *gin.Context) {
	type Theme struct {
		ThemeID     int    `json:"theme_id"`
		Name        string `json:"name"`
		Thumbnail   string `json:"thumbnail"`
		Description string `json:"description"`
	}
	res := struct {
		Error string  `json:"error"`
		List  []Theme `json:"list"`
	}{
		List: []Theme{},
	}

	db := util.OpenDB()

	rows, err := db.Query("SELECT theme_id, name, thumbnail, description FROM theme")
	if err != nil {
		res.Error = fmt.Sprintf("db.Query() error %v", err)
		c.IndentedJSON(http.StatusOK, res)
		return
	}
	defer rows.Close()
	for rows.Next() {
		theme := Theme{}
		rows.Scan(&theme.ThemeID, &theme.Name, &theme.Thumbnail, &theme.Description)
		res.List = append(res.List, theme)
	}

	c.IndentedJSON(http.StatusOK, res)
}

func uploadProfile(c *gin.Context) {
	file, _ := c.FormFile("profile")
	res := struct {
		Error string `json:"error"`
		Path  string `json:"path"`
	}{}

	// generate file name
	res.Path = path.Join(config.UploadRoot, fmt.Sprintf("%s%s", util.RandomString(15), path.Ext(file.Filename)))
	for fileExist(res.Path) {
		res.Path = path.Join(config.UploadRoot, fmt.Sprintf("%s%s", util.RandomString(15), path.Ext(file.Filename)))
	}

	if err := c.SaveUploadedFile(file, res.Path); err != nil {
		res.Error = fmt.Sprintf("upload fail %v", err)
		c.IndentedJSON(http.StatusOK, res)
		return
	}

	res.Path = "/" + res.Path
	c.IndentedJSON(http.StatusCreated, res)
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
