package gttools

import (
	"fmt"
	"os/exec"

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
	// The JDBC URL to use to connect to the database.
	Url string
	// The user to use to connect to the database.
	User string
	// The password to use to connect to the database.
	Password string
	// The fully qualified class name of the JDBC driver to use to connect to the database.
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
	// Allows you to override Flyway's logging autodetection and specify one or multiple loggers to use.
	Loggers []FlywayLogger
	// Whether to allow mixing transactional and non-transactional statements within the same migration.
	Mixed *bool
	// Allows migrations to be run "out of order".
	OutOfOrder *bool
	// Filename for the report file.
	ReportFilename string
	// Whether default built-in callbacks (SQL) should be skipped.
	SkipDefaultCallbacks *bool
	// Whether default built-in resolvers (SQL and JDBC) should be skipped.
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
	// Whether Flyway should attempt to create the schemas specified in the schemas' property.
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
	// The file name separator for SQL migrations.
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
	// Whether transactional advisory locks should be used with PostgreSQL.
	PostgresqlTransactionalLock *bool
}

func (tool *FlywayTool) Baseline(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "baseline")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

func (tool *FlywayTool) Clean(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "clean")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

func (tool *FlywayTool) Info(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "info")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

func (tool *FlywayTool) Migrate(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "migrate")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

func (tool *FlywayTool) Repair(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "repair")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

func (tool *FlywayTool) Validate(settings *FlywaySettings) error {
	args := tool.buildArguments(settings)
	args = append(args, settings.CustomArguments...)
	args = append(args, "validate")

	cmd := exec.Command(settings.ToolPath, args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

func (tool *FlywayTool) buildArguments(settings *FlywaySettings) []string {
	args := []string{}

	// Connection
	args = addString(args, settings.Url, addSettings{prefix: "-url="})
	args = addString(args, settings.User, addSettings{prefix: "-user="})
	args = addString(args, settings.Password, addSettings{prefix: "-password="})
	args = addString(args, settings.Driver, addSettings{prefix: "-driver="})
	args = addInt(args, settings.ConnectRetries, addSettings{prefix: "-connectRetries="})
	args = addInt(args, settings.ConnectRetriesInterval, addSettings{prefix: "-connectRetriesInterval="})
	args = addString(args, settings.InitSql, addSettings{prefix: "-initSql="})

	// General
	args = addStringList(args, settings.Callbacks, addSettings{prefix: "-callbacks=", listSeparator: ","})
	args = addStringList(args, settings.ConfigFiles, addSettings{prefix: "-configFiles=", listSeparator: ","})
	args = goext.AppendIf(args, len(settings.Encoding) > 0, fmt.Sprintf("-encoding=%s", settings.Encoding))
	args = addString(args, settings.Environment, addSettings{prefix: "-environment="})
	args = addBoolean(args, settings.ExecuteInTransaction, addSettings{prefix: "-executeInTransaction="})
	args = addBoolean(args, settings.Group, addSettings{prefix: "-group="})
	args = addString(args, settings.InstalledBy, addSettings{prefix: "-installedBy="})
	args = addStringList(args, settings.JarDirs, addSettings{prefix: "-jarDirs=", listSeparator: ","})
	args = addStringList(args, settings.Locations, addSettings{prefix: "-locations=", listSeparator: ","})
	args = addBoolean(args, settings.FailOnMissingLocations, addSettings{prefix: "-failOnMissingLocations="})
	args = addInt(args, settings.LockRetryCount, addSettings{prefix: "-lockRetryCount="})
	args = goext.AppendIf(args, len(settings.Loggers) > 0, fmt.Sprintf("-loggers=%s", goext.StringsJoinAny(settings.Loggers, ",")))
	args = addBoolean(args, settings.Mixed, addSettings{prefix: "-mixed="})
	args = addBoolean(args, settings.OutOfOrder, addSettings{prefix: "-outOfOrder="})
	args = addString(args, settings.ReportFilename, addSettings{prefix: "-reportFilename="})
	args = addBoolean(args, settings.SkipDefaultCallbacks, addSettings{prefix: "-skipDefaultCallbacks="})
	args = addBoolean(args, settings.SkipDefaultResolvers, addSettings{prefix: "-skipDefaultResolvers="})
	args = addString(args, settings.Table, addSettings{prefix: "-table="})
	args = addBoolean(args, settings.ValidateMigrationNaming, addSettings{prefix: "-validateMigrationNaming="})
	args = addBoolean(args, settings.ValidateOnMigrate, addSettings{prefix: "-validateOnMigrate="})
	args = addString(args, settings.WorkingDirectory, addSettings{prefix: "-workingDirectory="})

	// Schema
	args = addBoolean(args, settings.CreateSchemas, addSettings{prefix: "-createSchemas="})
	args = addString(args, settings.DefaultSchema, addSettings{prefix: "-defaultSchema="})
	args = addStringList(args, settings.Schemas, addSettings{prefix: "-schemas=", listSeparator: ","})

	// Baseline
	args = addString(args, settings.BaselineDescription, addSettings{prefix: "-baselineDescription="})
	args = addBoolean(args, settings.BaselineOnMigrate, addSettings{prefix: "-baselineOnMigrate="})
	args = addString(args, settings.BaselineVersion, addSettings{prefix: "-baselineVersion="})

	// Clean
	args = addBoolean(args, settings.CleanDisabled, addSettings{prefix: "-cleanDisabled="})
	args = addBoolean(args, settings.CleanOnValidationError, addSettings{prefix: "-cleanOnValidationError="})

	// Validate
	args = addStringList(args, settings.IgnoreMigrationPatterns, addSettings{prefix: "-ignoreMigrationPatterns=", listSeparator: ","})

	// Migrations
	args = addString(args, settings.RepeatableSqlMigrationPrefix, addSettings{prefix: "-repeatableSqlMigrationPrefix="})
	args = addStringList(args, settings.Resolvers, addSettings{prefix: "-resolvers=", listSeparator: ","})
	args = addString(args, settings.SqlMigrationPrefix, addSettings{prefix: "-sqlMigrationPrefix="})
	args = addString(args, settings.SqlMigrationSeparator, addSettings{prefix: "-sqlMigrationSeparator="})
	args = addStringList(args, settings.SqlMigrationSuffixes, addSettings{prefix: "-sqlMigrationSuffixes=", listSeparator: ","})

	// Placeholders
	args = addString(args, settings.PlaceholderPrefix, addSettings{prefix: "-placeholderPrefix="})
	args = addBoolean(args, settings.PlaceholderReplacement, addSettings{prefix: "-placeholderReplacement="})
	args = addStringList(args, settings.Placeholders, addSettings{prefix: "-placeholders=", listSeparator: ","})
	args = addString(args, settings.PlaceholderSeparator, addSettings{prefix: "-placeholderSeparator="})
	args = addString(args, settings.PlaceholderSuffix, addSettings{prefix: "-placeholderSuffix="})

	// Command Line
	args = addString(args, string(settings.Color), addSettings{prefix: "-color="})
	args = addString(args, string(settings.OutputType), addSettings{prefix: "-outputType="})

	// PostgreSQL
	args = addBoolean(args, settings.PostgresqlTransactionalLock, addSettings{prefix: "-postgresqlTransactionalLock="})

	return args
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
