package handler

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/urlvet/urlvet/internal/service/cache"
)

// cacheEntry is the per-key payload returned by ListCacheHandler.
type cacheEntry struct {
	Key        string          `json:"key"`
	Prefix     string          `json:"prefix"`
	TTLSeconds int64           `json:"ttl_seconds"` // -1 = no expiry, -2 = key gone
	Value      json.RawMessage `json:"value" swaggertype:"object"`
}

// ListCacheHandler godoc
// @Summary      List all cache keys and values
// @Tags         admin
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /admin/cache [get]
func ListCacheHandler(c *gin.Context) {
	cacheInstance, err := cache.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to cache"})
		return
	}
	defer cacheInstance.Close()

	ctx := c.Request.Context()

	keys, err := cacheInstance.Scan(ctx, "*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list keys"})
		return
	}
	sort.Strings(keys)

	entries := make([]cacheEntry, 0, len(keys))
	for _, key := range keys {
		raw, err := cacheInstance.Get(ctx, key)
		if err != nil || raw == "" {
			continue
		}

		ttl, _ := cacheInstance.TTL(ctx, key)
		ttlSec := int64(ttl.Seconds())

		// Try to present value as formatted JSON; fall back to a plain string.
		var valueRaw json.RawMessage
		if json.Valid([]byte(raw)) {
			valueRaw = json.RawMessage(raw)
		} else {
			quoted, _ := json.Marshal(raw)
			valueRaw = json.RawMessage(quoted)
		}

		prefix := key
		for i, ch := range key {
			if ch == ':' {
				prefix = key[:i]
				break
			}
		}

		entries = append(entries, cacheEntry{
			Key:        key,
			Prefix:     prefix,
			TTLSeconds: ttlSec,
			Value:      valueRaw,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(entries),
		"keys":  entries,
	})
}

// DeleteCacheKeyHandler godoc
// @Summary      Delete a single cache key
// @Tags         admin
// @Param        key  path  string  true  "Cache key"
// @Success      200  {object}  map[string]any
// @Router       /admin/cache/{key} [delete]
func DeleteCacheKeyHandler(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}
	// Gin includes the leading slash when using wildcard params; strip it.
	if len(key) > 0 && key[0] == '/' {
		key = key[1:]
	}

	cacheInstance, err := cache.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to cache"})
		return
	}
	defer cacheInstance.Close()

	if err := cacheInstance.Delete(c.Request.Context(), key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "key": key})
}

// FlushCacheHandler godoc
// @Summary      Flush all cache keys
// @Tags         admin
// @Success      200  {object}  map[string]any
// @Router       /admin/cache [delete]
func FlushCacheHandler(c *gin.Context) {
	cacheInstance, err := cache.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  "failed to connect to cache service",
		})
		return
	}
	defer cacheInstance.Close()

	if err := cacheInstance.FlushAll(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
			"error":  "failed to flush cache",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"message": "all cache has been flushed successfully",
	})
}
