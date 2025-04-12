package middlewares

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sharing-vision-id/pkg"
	"sync"
	"time"
)

type rateLimitEntry struct {
	count     int
	timestamp time.Time
}

var (
	rateLimitMap = make(map[string]*rateLimitEntry)
	mu           sync.Mutex
)

func RateLimitMiddleware(limit int, period time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := c.IP()
		mu.Lock()
		defer mu.Unlock()

		entry, exists := rateLimitMap[clientIP]
		if !exists {
			rateLimitMap[clientIP] = &rateLimitEntry{
				count:     1,
				timestamp: time.Now(),
			}
		} else {
			if time.Since(entry.timestamp) > period {
				entry.count = 1
				entry.timestamp = time.Now()
			} else {
				entry.count++
			}
		}

		if rateLimitMap[clientIP].count > limit {
			//return c.Status(fiber.StatusTooManyRequests).JSON(map[string]interface{}{"message": "Too many requests"})
			return c.Status(fiber.StatusTooManyRequests).SendFile("./public/views/429.html")
		}

		return c.Next()
	}
}

func SignatureMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		payloadString := fmt.Sprintf("%s:%s", pkg.USERNAME, pkg.PASSWORD)
		h := hmac.New(sha512.New, []byte(pkg.APIKEY))
		h.Write([]byte(payloadString))
		expectedHash := h.Sum(nil)
		expectedSignature := base64.StdEncoding.EncodeToString(expectedHash)
		fmt.Println("expectedSignature", expectedSignature)
		if c.Get("X-SIGNATURE") != expectedSignature {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "SIGNATURE INVALID",
			})
		}
		return c.Next()
	}
}

const maxRequests = 10

func RateLimitRequestMiddleware(c *fiber.Ctx) error {
	ip := c.IP()
	value, _ := pkg.RateLimitMap.LoadOrStore(ip, &struct {
		count     int
		lastReset time.Time
	}{0, time.Now()})

	data := value.(*struct {
		count     int
		lastReset time.Time
	})

	if time.Since(data.lastReset) > time.Minute {
		data.count = 0
		data.lastReset = time.Now()
	}

	if data.count >= maxRequests {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Request limit exceeded, Upgrade your account",
		})
	}

	data.count++
	pkg.RateLimitMap.Store(ip, data)
	return c.Next()
}
