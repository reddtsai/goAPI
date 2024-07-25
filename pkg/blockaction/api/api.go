package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/reddtsai/goAPI/pkg/blockaction/storage"
)

type IBlockActionApi interface {
	http.Handler

	Health(c *gin.Context)
	Signup(c *gin.Context)
	Signin(c *gin.Context)
	GetPersonalInfo(c *gin.Context)
}

const (
	SECRET            = "8w+uSC4vD136hOaT1m1fWeuuULid9LiLSJ52kAPHC33bY9lK5/3ZAS+nm+HaJz93qCWDXcrJLl/9mcXzwUL3DQ=="
	TOKEN_EXPIRE_TIME = 900
	CTX_REQUEST_ID    = "c_request_id"
	CTX_USER_ID       = "c_user_id"
	CTX_USER_ACCOUNT  = "c_user_account"
)

var (
	_ IBlockActionApi = (*BlockActionApi)(nil)

	_snowNode      *snowflake.Node
	_routerMetrics *RouterMetrics
)

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	_snowNode = node
	_routerMetrics = &RouterMetrics{
		RequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_router_requests_total",
				Help: "Requests total count",
			},
			[]string{"http_code", "http_method", "http_url_path", "http_service"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "api_router_request_duration_seconds",
				Help: "Requests duration seconds histogram",
			},
			[]string{"http_code", "http_method", "http_url_path", "http_service"},
		),
	}
}

func NewBlockActionApi(opts ...BlockActionApiOption) (*BlockActionApi, error) {
	api := new(BlockActionApi)
	api.opts = DefaultOptions()
	for _, opt := range opts {
		opt(&api.opts)
	}
	if api.opts.storage == nil {
		return nil, fmt.Errorf("storage is nil")
	}

	api.Engine = gin.New()
	api.Engine.Use(gin.Logger())
	api.Engine.Use(gin.Recovery())
	api.Engine.Use(cors.New(cors.Config{
		AllowOrigins:     api.opts.allowOrigins,
		AllowMethods:     api.opts.allowMethods,
		AllowHeaders:     api.opts.allowHeaders,
		ExposeHeaders:    api.opts.exposeHeaders,
		AllowCredentials: api.opts.allowCredentials,
		MaxAge:           api.opts.maxAge,
	}))
	api.Engine.GET("/health", api.Health)
	privateGroup := api.Engine.Group("/_")
	{
		prometheus.Register(_routerMetrics)
		privateGroup.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}
	// TODO : swagger
	v1Group := api.Engine.Group("/v1")
	{
		v1Group.Use(middleware)
		v1Group.POST("/signup", api.Signup)
		v1Group.POST("/signin", api.Signin)

		userGroup := v1Group.Group("/user")
		userGroup.Use(authMiddleware)
		userGroup.GET("/personal-info", api.GetPersonalInfo)
	}

	return api, nil
}

type BlockActionApiOptions struct {
	allowOrigins     []string
	allowMethods     []string
	allowHeaders     []string
	exposeHeaders    []string
	allowCredentials bool
	maxAge           time.Duration
	storage          storage.IStorage
}

type BlockActionApiOption func(*BlockActionApiOptions)

func DefaultOptions() BlockActionApiOptions {
	return BlockActionApiOptions{
		allowOrigins:     []string{"*"},
		allowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		allowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authorization", "X-Request-ID", "X-API-Key"},
		exposeHeaders:    []string{"*"},
		allowCredentials: true,
		maxAge:           24 * time.Hour,
	}
}

func SetStorage(storage storage.IStorage) BlockActionApiOption {
	return func(o *BlockActionApiOptions) {
		o.storage = storage
	}
}

func signaturePwd(pwd string) (mac string, err error) {
	mac, err = hmacSignature([]byte(pwd), []byte(SECRET))
	return
}

func validatePwd(pwd, secret string) (bool, error) {
	mac, err := hmacSignature([]byte(pwd), []byte(SECRET))
	if err != nil {
		return false, err
	}

	return (mac == secret), nil
}

func hmacSignature(msg []byte, key []byte) (string, error) {
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(msg); err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func genToken(claims jwt.Claims) (token string, err error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(SECRET))

	return
}

func middleware(c *gin.Context) {
	requestID := _snowNode.Generate().String()
	c.Set(CTX_REQUEST_ID, requestID)
	c.Header("X-Request-ID", requestID)

	startTime := time.Now()
	defer func() {
		statusCode := strconv.Itoa(c.Writer.Status())
		serviceName := "im2"
		duration := time.Since(startTime).Seconds()
		_routerMetrics.RequestDuration.WithLabelValues(statusCode, c.Request.Method, c.Request.URL.Path, serviceName).Observe(duration)
		_routerMetrics.RequestTotal.WithLabelValues(statusCode, c.Request.Method, c.Request.URL.Path, serviceName).Inc()
	}()

	c.Next()
}

type RouterMetrics struct {
	RequestTotal    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func (m *RouterMetrics) Collect(ch chan<- prometheus.Metric) {
	m.RequestTotal.Collect(ch)
	m.RequestDuration.Collect(ch)
}

func (m *RouterMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.RequestTotal.Describe(ch)
	m.RequestDuration.Describe(ch)
}

func authMiddleware(c *gin.Context) {
	token := getBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization"})
		return
	}
	claims := &UserClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization"})
		return
	}

	c.Set(CTX_USER_ID, claims.ID)
	c.Set(CTX_USER_ACCOUNT, claims.Account)

	c.Next()
}

func getBearerToken(auth string) (token string) {
	bearer := strings.Split(auth, "Bearer ")
	if len(bearer) == 2 {
		token = bearer[1]
	}

	return
}
