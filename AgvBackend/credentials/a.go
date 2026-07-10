package credentials

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var APP_NAME = "INTL"

var APP_NAME_SMALL = "intl"

func LicenceDisabled() bool {
	return true
}
func LicenseSource() string {
	return ""
}
func GetLicense() string {
	return ""
}

// EmptyObject is used for handling errors in argument parsing
type EmptyObject struct{}

// ParseArguments parses command-line arguments
func ParseArguments() map[string]string {
	// Define flags (command-line arguments)
	dbName := flag.String("db-name", "-x-", "Database name")
	dbConnStr := flag.String("db-conn-str", "-x-", "Database connection string")
	ipAddress := flag.String("ip-address", "-x-", "IP address")
	port := flag.String("port", "-x-", "Port number")

	// Parse flags
	flag.Parse()

	// Return a map of parsed arguments
	return map[string]string{
		"db_name":     *dbName,
		"db_conn_str": *dbConnStr,
		"ip_address":  *ipAddress,
		"port":        *port,
	}
}

var arguments = ParseArguments()

// GetValue retrieves a value from the arguments map, returns default if not found or empty
func GetValue(args map[string]string, key string) string {
	if val, ok := args[key]; ok && val != "-x-" {
		return val
	}
	return ""
}

func GetEnvOrDefault(key string, defaultValue string) string {
	var value = os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

// GetParamFromEnvOrArg retrieves a value from environment variables or arguments
func GetParamFromEnvOrArg(envName, defaultValue, argValue string) string {
	if argValue != "" {
		return argValue
	}
	if envValue, ok := os.LookupEnv(envName); ok {
		return envValue
	}
	return defaultValue
}

// Database-related configuration
const (
	TWIN_DB_NAME             = "TURTLE_DB_NAME"
	TWIN_DB_NAME_VAL_DEFAULT = "turtle"
	TWIN_DB_NAME_VAL_TEST    = "turtle_test"
)

func GetAppName() string {
	return TWIN_DB_NAME_VAL_DEFAULT
}

// GetDBName returns the database name
func GetDBName() string {
	return GetParamFromEnvOrArg(TWIN_DB_NAME, TWIN_DB_NAME_VAL_DEFAULT, GetValue(arguments, "db_name"))
} // GetDBName returns the database name

// GetDBConnStr returns the database connection string
func GetDBConnStr() string {
	return GetParamFromEnvOrArg("TURTLE_DB_CONN_STRING", "mongodb://localhost:27017/", GetValue(arguments, "db_conn_str"))
}

// GetStorageBackend returns which db.IStorageController implementation
// db.SC should use. Defaults to "fs" (local filesystem). Set to "mongo" to
// store files in MongoDB via db.MongoStorageClient.
func GetStorageBackend() string {
	return GetEnvOrDefault("TURTLE_STORAGE_BACKEND", "fs") //mongo
}

// Network-related configuration

// GetIpAddress returns the IP address
func GetIpAddress() string {
	return GetParamFromEnvOrArg("TURTLE_IP", "0.0.0.0", GetValue(arguments, "ip_address"))
}

// GetPort returns the port number as an integer
func GetPort() int {
	portStr := GetParamFromEnvOrArg("TURTLE_PORT", "5000", GetValue(arguments, "port"))
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 5000
	}
	return port
}

// Authentication-related configuration
const (
	AUTH_PROVIDER_NONE     = ""
	AUTH_PROVIDER_INFINITY = "infinity"
	AUTH_PROVIDER_AUTH0    = "auth0"
	AUTH_PROVIDER_ALL      = "all"
)

// AuthProvider returns the authentication provider
func AuthProvider() string {
	return os.Getenv("TURTLE_AUTH_PROVIDER")
}

// AuthDefaultUser returns the default user for authentication
func AuthDefaultUser() string {
	return GetEnvOrDefault("TURTLE_DEFAULT_USER", "turtle@turtle.sk")
}

// AuthDefaultUser returns the default user for authentication
func AuthDefaultPassword() string {
	return GetEnvOrDefault("TURTLE_DEFAULT_PASSWORD", "turtle@turtle.sk")
}

// Auth0-related configuration

// Auth0Enabled checks if Auth0 is enabled
func Auth0Enabled() bool {
	return AuthProvider() == AUTH_PROVIDER_AUTH0 || AuthProvider() == AUTH_PROVIDER_ALL
}

// Auth0ClientID returns the Auth0 client ID
func Auth0ClientID() string {
	return os.Getenv("TURTLE_AUTH0_CLIENT_ID")
}

// Auth0Secret returns the Auth0 secret
func Auth0Secret() string {
	return os.Getenv("TURTLE_AUTH0_SECRET")
}

// Auth0Scope returns the Auth0 scope
func Auth0Scope() map[string]string {
	return map[string]string{
		"scope": "openid profile email",
	}
}

// Auth0Domain returns the Auth0 domains
func Auth0Domain() string {
	return os.Getenv("TURTLE_AUTH0_DOMAIN")
}

// Auth0Metadata returns the Auth0 metadata URL
func Auth0Metadata() string {
	return fmt.Sprintf("https://%s/.well-known/openid-configuration", Auth0Domain())
}

// Auth0AccessToken returns the Auth0 access token
func Auth0AccessToken() string {
	return "access_token"
}

// Infinity Authentication-related configuration

// AuthInfinityEnabled checks if Infinity authentication is enabled
func AuthInfinityEnabled() bool {
	return AuthProvider() == AUTH_PROVIDER_INFINITY || AuthProvider() == AUTH_PROVIDER_ALL
}

// AuthInfinityJwtSecret returns the JWT secret for Infinity authentication
func AuthInfinityJwtSecret() string {
	return GetEnvOrDefault("TURTLE_AUTH_JWT_SECRET", "infinitysecret")
}

func AuthInfinityJwtKey() string {
	return GetEnvOrDefault("TURTLE_AUTH_JWT_SECRET", "infinity")
}

func AuthInfinityTokenExpire() string {
	//In seconds

	return GetEnvOrDefault("TURTLE_AUTH_TOKEN_EXPIRE", "3600")
}

// AuthInfinityChangePasswordURL returns the change password URL for Infinity authentication
func AuthInfinityChangePasswordURL() string {
	return os.Getenv("TURTLE_AUTH_CHANGE_PASSWORD_URL")
}

// AuthInfinityChangePasswordSession returns the session duration for changing passwords in minutes
func AuthInfinityChangePasswordSession() int {
	timeStr := os.Getenv("TURTLE_AUTH_CHANGE_PASSWORD_MINUTE_SESSION")
	time, err := strconv.Atoi(timeStr)
	if err != nil {
		return 1
	}
	return time
}

// Azure-related configuration

// RunHttps checks if the service should run with HTTPS
func RunHttps() bool {
	return os.Getenv("TURTLE_HTTPS") == "1"
}

// AzureMailCredentials returns the Azure mail credentials
func AzureMailCredentials() string {
	return os.Getenv("TURTLE_AZURE_MAIL_KEY")
}

// GetAzureStorageCredentials returns the Azure storage credentials
func GetAzureStorageCredentials() string {
	return os.Getenv("TURTLE_AZURE_STORAGE_KEY")
}

// FillTestCredentials sets tool_calling credentials for the environment
func FillTestCredentials() {
	os.Setenv(TWIN_DB_NAME, TWIN_DB_NAME_VAL_TEST)
}

// SamlEnabled checks if SAML is enabled
func SamlEnabled() bool {
	return os.Getenv("TURTLE_SAML") == "1"
}

// SamlEnabled checks if SAML is enabled
func LinuxWorkspace() string {
	return GetEnvOrDefault("TURTLE_LINUX_WORKSPACE", "../turtle_storage") + "/" + GetAppName()
}
