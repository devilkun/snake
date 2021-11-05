package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// AppConfig app config
	type AppConfig struct {
		Name              string
		Version           string
		Mode              string
		PprofPort         string
		URL               string
		JwtSecret         string
		JwtTimeout        int
		SSL               bool
		CtxDefaultTimeout time.Duration
		CSRF              bool
		Debug             bool
	}
	var appConfig AppConfig

	t.Run("using local var", func(t *testing.T) {
		c := New(WithConfigDir("../../config/"))
		if err := c.Scan("app", &appConfig); err != nil {
			t.Fatal(err)
		}
		t.Log(appConfig.Name)
	})

	// test global conf
	t.Run("using global conf", func(t *testing.T) {
		_ = New(WithConfigDir("../../config/"))
		if err := Conf.Scan("app", &appConfig); err != nil {
			t.Fatal(err)
		}

		t.Log(appConfig.Name)
	})
}
