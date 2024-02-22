package gttools

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

type FlywayTool struct {
}

func CreateFlywayTool() *FlywayTool {
	return &FlywayTool{}
}

type FlywaySettings struct {
	ToolSettingsBase

	// The path to the Flyway executable.
	ToolPath string

	//// Connection
	// The jdbc url to use to connect to the database.
	Url string
	// The user to use to connect to the database.
	User string
	// The password to use to connect to the database.
	Password string
	// The fully qualified class name of the jdbc driver to use to connect to the database.
	Driver string
	// The maximum number of retries when attempting to connect to the database.
	ConnectRetries *int
	// The maximum time between retries when attempting to connect to the database in seconds.
	ConnectRetriesInterval *int
	// SQL statements to be run immediately after a database connection has been established.
	InitSql string

	//// General
	// List of fully qualified class names of Callback implementations to use to hook into the Flyway lifecycle, or packages to scan for these classes.
	Callbacks []string
	// The Flyway configuration files to load.
	ConfigFiles []string
	// The encoding of SQL migrations.
	Encoding FlywayEncoding
	// The name of the environment you wish Flyway to load configuration from.
	Environment string
	// Whether Flyway should execute SQL within a transaction.
	ExecuteInTransaction *bool
	// Whether to group all pending migrations together in the same transaction when applying them.
	Group *bool
	// The username that will be recorded in the schema history table as having applied the migration.
	InstalledBy string
	// List of directories containing JDBC drivers and Java-based migrations.
	JarDirs []string
	// List of locations to scan recursively for migrations.
	Locations []string
	// Whether to fail if a location specified in the locations option doesn't exist.
	FailOnMissingLocations *bool
	// At the start of a migration, Flyway will attempt to take a lock to prevent competing instances executing in parallel.
	LockRetryCount *int
	// Allows you to override Flyway's logging auto-detection and specify one or multiple loggers to use.
	Loggers []FlywayLogger
	// Whether to allow mixing transactional and non-transactional statements within the same migration.
	Mixed *bool
	// Allows migrations to be run "out of order".
	OutOfOrder *bool
	// Filename for the report file.
	ReportFilename string
	// Whether default built-in callbacks (sql) should be skipped.
	SkipDefaultCallbacks *bool
	// Whether default built-in resolvers (sql and jdbc) should be skipped.
	SkipDefaultResolvers *bool
	// The name of Flyway's schema history table.
	Table string
	// Whether to ignore migration files whose names do not match the naming conventions.
	ValidateMigrationNaming *bool
	// Whether to automatically call validate or not when running migrate.
	ValidateOnMigrate *bool
	// The working directory to consider when dealing with relative paths for both config files and locations.
	WorkingDirectory string

	//// Schema
	// Whether Flyway should attempt to create the schemas specified in the schemas property.
	CreateSchemas *bool
	// The default schema managed by Flyway.
	DefaultSchema string
	// List of schemas managed by Flyway.
	Schemas []string

	//// Baseline
	// The Description to tag an existing schema with when executing baseline.
	BaselineDescription string
	// Whether to automatically call baseline when migrate is executed against a non-empty schema with no schema history table.
	BaselineOnMigrate *bool
	// The version to tag an existing schema with when executing baseline.
	BaselineVersion string

	//// Clean
	// Whether to disable clean.
	CleanDisabled *bool
	// Whether to automatically call clean or not when a validation error occurs.
	CleanOnValidationError *bool

	//// Validate
	// Ignore migrations during validate and repair according to a given list of patterns.
	IgnoreMigrationPatterns []string

	//// Migrations
	// The file name prefix for repeatable SQL migrations.
	RepeatableSqlMigrationPrefix string
	// List of fully qualified class names of custom MigrationResolver implementations to be used in addition to the built-in ones for resolving Migrations to apply.
	Resolvers []string
	// The file name prefix for versioned SQL migrations.
	SqlMigrationPrefix string
	// The file name separator for Sql migrations.
	SqlMigrationSeparator string
	// List of file name suffixes for SQL migrations.
	SqlMigrationSuffixes []string

	//// Placeholders
	// The prefix of every placeholder.
	PlaceholderPrefix string
	// Whether placeholders should be replaced.
	PlaceholderReplacement *bool
	// Placeholders to replace in SQL migrations. In the form of "key=value"
	Placeholders []string
	// The separator of default placeholders.
	PlaceholderSeparator string
	// The suffix of every placeholder.
	PlaceholderSuffix string

	//// Command Line
	// Whether to colorize output.
	Color FlywayColor
	// Change the console output formatting.
	OutputType FlywayOutputType

	//// PostgreSQL
	// Whether or not transactional advisory locks should be used with PostgreSQL.
	PostgresqlTransactionalLock *bool
}

func (tool *FlywayTool) Baseline(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "baseline")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

func (tool *FlywayTool) Clean(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "clean")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

func (tool *FlywayTool) Info(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "info")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

func (tool *FlywayTool) Migrate(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "migrate")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

func (tool *FlywayTool) Repair(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "repair")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

func (tool *FlywayTool) Validate(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "validate")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

func (tool *FlywayTool) buildArguments(settings *FlywaySettings) []string {
	arguments := []string{}

	// Connection
	arguments = tool.addString(arguments, settings.Url, "-url=%s")
	arguments = tool.addString(arguments, settings.User, "-user=%s")
	arguments = tool.addString(arguments, settings.Password, "-password=%s")
	arguments = tool.addString(arguments, settings.Driver, "-driver=%s")
	arguments = tool.addInt(arguments, settings.ConnectRetries, "-connectRetries=%d")
	arguments = tool.addInt(arguments, settings.ConnectRetriesInterval, "-connectRetriesInterval=%d")
	arguments = tool.addString(arguments, settings.InitSql, "-initSql=%s")

	// General
	arguments = tool.addStringList(arguments, settings.Callbacks, "-callbacks=%s")
	arguments = tool.addStringList(arguments, settings.ConfigFiles, "-configFiles=%s")
	arguments = goext.AppendIf(arguments, len(settings.Encoding) > 0, fmt.Sprintf("-encoding=%s", settings.Encoding))
	arguments = tool.addString(arguments, settings.Environment, "-environment=%s")
	arguments = tool.addBoolean(arguments, settings.ExecuteInTransaction, "-executeInTransaction=%t")
	arguments = tool.addBoolean(arguments, settings.Group, "-group=%t")
	arguments = tool.addString(arguments, settings.InstalledBy, "-installedBy=%s")
	arguments = tool.addStringList(arguments, settings.JarDirs, "-jarDirs=%s")
	arguments = tool.addStringList(arguments, settings.Locations, "-locations=%s")
	arguments = tool.addBoolean(arguments, settings.FailOnMissingLocations, "-failOnMissingLocations=%t")
	arguments = tool.addInt(arguments, settings.LockRetryCount, "-lockRetryCount=%d")
	arguments = goext.AppendIf(arguments, len(settings.Loggers) > 0, fmt.Sprintf("-loggers=%s", goext.StringsJoinAny(settings.Loggers, ",")))
	arguments = tool.addBoolean(arguments, settings.Mixed, "-mixed=%t")
	arguments = tool.addBoolean(arguments, settings.OutOfOrder, "-outOfOrder=%t")
	arguments = tool.addString(arguments, settings.ReportFilename, "-reportFilename=%s")
	arguments = tool.addBoolean(arguments, settings.SkipDefaultCallbacks, "-skipDefaultCallbacks=%t")
	arguments = tool.addBoolean(arguments, settings.SkipDefaultResolvers, "-skipDefaultResolvers=%t")
	arguments = tool.addString(arguments, settings.Table, "-table=%s")
	arguments = tool.addBoolean(arguments, settings.ValidateMigrationNaming, "-validateMigrationNaming=%t")
	arguments = tool.addBoolean(arguments, settings.ValidateOnMigrate, "-validateOnMigrate=%t")
	arguments = tool.addString(arguments, settings.WorkingDirectory, "-workingDirectory=%s")

	// Schema
	arguments = tool.addBoolean(arguments, settings.CreateSchemas, "-createSchemas=%t")
	arguments = tool.addString(arguments, settings.DefaultSchema, "-defaultSchema=%s")
	arguments = tool.addStringList(arguments, settings.Schemas, "-schemas=%s")

	// Baseline
	arguments = tool.addString(arguments, settings.BaselineDescription, "-baselineDescription=%s")
	arguments = tool.addBoolean(arguments, settings.BaselineOnMigrate, "-baselineOnMigrate=%t")
	arguments = tool.addString(arguments, settings.BaselineVersion, "-baselineVersion=%s")

	// Clean
	arguments = tool.addBoolean(arguments, settings.CleanDisabled, "-cleanDisabled=%t")
	arguments = tool.addBoolean(arguments, settings.CleanOnValidationError, "-cleanOnValidationError=%t")

	// Validate
	arguments = tool.addStringList(arguments, settings.IgnoreMigrationPatterns, "-ignoreMigrationPatterns=%s")

	// Migrations
	arguments = tool.addString(arguments, settings.RepeatableSqlMigrationPrefix, "-repeatableSqlMigrationPrefix=%s")
	arguments = tool.addStringList(arguments, settings.Resolvers, "-resolvers=%s")
	arguments = tool.addString(arguments, settings.SqlMigrationPrefix, "-sqlMigrationPrefix=%s")
	arguments = tool.addString(arguments, settings.SqlMigrationSeparator, "-sqlMigrationSeparator=%s")
	arguments = tool.addStringList(arguments, settings.SqlMigrationSuffixes, "-sqlMigrationSuffixes=%s")

	// Placeholders
	arguments = tool.addString(arguments, settings.PlaceholderPrefix, "-placeholderPrefix=%s")
	arguments = tool.addBoolean(arguments, settings.PlaceholderReplacement, "-placeholderReplacement=%t")
	arguments = tool.addStringList(arguments, settings.Placeholders, "-placeholders=%s")
	arguments = tool.addString(arguments, settings.PlaceholderSeparator, "-placeholderSeparator=%s")
	arguments = tool.addString(arguments, settings.PlaceholderSuffix, "-placeholderSuffix=%s")

	// Command Line
	arguments = goext.AppendIf(arguments, len(settings.Color) > 0, fmt.Sprintf("-color=%s", settings.Color))
	arguments = goext.AppendIf(arguments, len(settings.OutputType) > 0, fmt.Sprintf("-outputType=%s", settings.OutputType))

	// PostgreSQL
	arguments = tool.addBoolean(arguments, settings.PostgresqlTransactionalLock, "-postgresqlTransactionalLock=%t")

	return arguments
}

////////// Helpers

func (tool *FlywayTool) addString(slice []string, value string, format string) []string {
	if len(value) > 0 {
		slice = append(slice, fmt.Sprintf(format, value))
	}
	return slice
}

func (tool *FlywayTool) addStringList(slice []string, values []string, format string) []string {
	if len(values) > 0 {
		slice = append(slice, fmt.Sprintf(format, strings.Join(values, ",")))
	}
	return slice
}

func (tool *FlywayTool) addBoolean(slice []string, value *bool, format string) []string {
	if value != nil {
		slice = append(slice, fmt.Sprintf(format, *value))
	}
	return slice
}

func (tool *FlywayTool) addInt(slice []string, value *int, format string) []string {
	if value != nil {
		slice = append(slice, fmt.Sprintf(format, *value))
	}
	return slice
}

////////// Types

type FlywayEncoding string

const (
	FLYWAY_ENCODING_DEFAULT    FlywayEncoding = ""
	FLYWAY_ENCODING_UTF_8      FlywayEncoding = "UTF-8"
	FLYWAY_ENCODING_UTF_16BE   FlywayEncoding = "UTF-16BE"
	FLYWAY_ENCODING_UTF_16LE   FlywayEncoding = "UTF-16LE"
	FLYWAY_ENCODING_UTF_16     FlywayEncoding = "UTF-16"
	FLYWAY_ENCODING_US_ASCII   FlywayEncoding = "US-ASCII"
	FLYWAY_ENCODING_ISO_8859_1 FlywayEncoding = "ISO-8859-1"
)

type FlywayLogger string

const (
	FLYWAY_LOGGER_DEFAULT        FlywayLogger = ""
	FLYWAY_LOGGER_AUTO           FlywayLogger = "auto"
	FLYWAY_LOGGER_CONSOLE        FlywayLogger = "console"
	FLYWAY_LOGGER_SLF4J          FlywayLogger = "slf4j"
	FLYWAY_LOGGER_LOG4J2         FlywayLogger = "log4j2"
	FLYWAY_LOGGER_APACHE_COMMONS FlywayLogger = "apache-commons"
)

type FlywayColor string

const (
	FLYWAY_COLOR_DEFAULT FlywayColor = ""
	FLYWAY_COLOR_AUTO    FlywayColor = "auto"
	FLYWAY_COLOR_ALWAYS  FlywayColor = "always"
	FLYWAY_COLOR_NEVER   FlywayColor = "never"
)

type FlywayOutputType string

const (
	FLYWAY_OUTPUT_TYPE_DEFAULT FlywayOutputType = ""
	FLYWAY_OUTPUT_TYPE_JSON    FlywayOutputType = "json"
)
