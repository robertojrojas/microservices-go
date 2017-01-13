package api

import (
	"testing"
	"os"
)


func Test_getConfig_Default(t *testing.T) {

	config := getConfig()

	if config.mySQLDBURI != defaultMYSQLURI {
		t.Error("default MYSQL URI does not match")
	}

}


func Test_getConfig_ModifiedViaEnvs(t *testing.T) {

	expectedMYSQLDBURI := "root:testlocation@tcp(:3306)/cats_db"
	os.Setenv(mysqlDBURI_key, expectedMYSQLDBURI)
	defer func() {
		os.Unsetenv(mysqlDBURI_key)
	}()

	config := getConfig()

	if config.mySQLDBURI != expectedMYSQLDBURI {
		t.Errorf("Expected mySQLDBURI: %s, but instead got: %s", expectedMYSQLDBURI, config.mySQLDBURI)
	}
}