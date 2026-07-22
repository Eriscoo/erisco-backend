package router

import (
	"github.com/eriscoo/blog-backend/internal/application"
	authsvc "github.com/eriscoo/blog-backend/internal/application/auth"
	catsvc "github.com/eriscoo/blog-backend/internal/application/categories"
	contactsvc "github.com/eriscoo/blog-backend/internal/application/contact"
	postssvc "github.com/eriscoo/blog-backend/internal/application/posts"
	profilesvc "github.com/eriscoo/blog-backend/internal/application/profile"
	tagssvc "github.com/eriscoo/blog-backend/internal/application/tags"
	authhandler "github.com/eriscoo/blog-backend/internal/transport/handler/auth"
	catHandler "github.com/eriscoo/blog-backend/internal/transport/handler/categories"
	contactHandler "github.com/eriscoo/blog-backend/internal/transport/handler/contact"
	oghandler "github.com/eriscoo/blog-backend/internal/transport/handler/og"
	postsHandler "github.com/eriscoo/blog-backend/internal/transport/handler/posts"
	profileHandler "github.com/eriscoo/blog-backend/internal/transport/handler/profile"
	sitemapHandler "github.com/eriscoo/blog-backend/internal/transport/handler/sitemap"
	tagshandler "github.com/eriscoo/blog-backend/internal/transport/handler/tags"
	uploadHandler "github.com/eriscoo/blog-backend/internal/transport/handler/upload"
	"github.com/eriscoo/blog-backend/internal/transport/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, authSvc *authsvc.Service, tagsSvc *tagssvc.Service, catsSvc *catsvc.Service, postsSvc *postssvc.Service, profileSvc *profilesvc.Service, contactSvc *contactsvc.Service, uploadH *uploadHandler.UploadHandler, ogH *oghandler.OGHandler, smH *sitemapHandler.Handler, tokens application.TokenService) {
	authH := authhandler.New(authSvc)
	tagsH := tagshandler.New(tagsSvc)
	catsH := catHandler.New(catsSvc)
	postsH := postsHandler.New(postsSvc, catsSvc, tagsSvc)
	profileH := profileHandler.New(profileSvc)
	contactH := contactHandler.New(contactSvc)

	r.GET("/og/:slug", ogH.HandleOG)
	r.GET("/og/page/:page", ogH.HandleStaticPage)
	r.GET("/sitemap.xml", smH.GetSitemap)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", authH.Register)
		v1.POST("/login", authH.Login)
		v1.GET("/public/posts/all", postsH.GetAllPublishedPosts)
		v1.GET("/public/posts/:slug", postsH.GetPostBySlug)
		v1.GET("/public/posts/categories/:name", postsH.GetPostsByCategory)
		v1.GET("/public/posts/tags/:name", postsH.GetPostsByTag)
		v1.POST("/contact", contactH.Submit)

		authorized := v1.Group("")
		authorized.Use(middleware.AuthRequired(tokens))
		{
			authorized.GET("/tags", tagsH.GetTags)
			authorized.GET("/categories", catsH.GetCategories)
			authorized.GET("/posts", postsH.GetPosts)
			authorized.GET("/posts/:id", postsH.GetPost)
			authorized.GET("/me", authH.GetMe)
			authorized.POST("/tags", tagsH.CreateTag)
			authorized.PUT("/tags/:id", tagsH.UpdateTag)
			authorized.DELETE("/tags/:id", tagsH.DeleteTag)
			authorized.POST("/categories", catsH.CreateCategory)
			authorized.PUT("/categories/:id", catsH.UpdateCategory)
			authorized.DELETE("/categories/:id", catsH.DeleteCategory)
			authorized.POST("/posts", postsH.CreatePost)
			authorized.PUT("/posts/:id", postsH.UpdatePost)
			authorized.DELETE("/posts/:id", postsH.DeletePost)

			authorized.GET("/profile/:user_id", profileH.GetProfile)
			authorized.PUT("/profile/:user_id", profileH.UpdateProfile)
			authorized.POST("/upload", uploadH.Upload)
		}
	}
}
