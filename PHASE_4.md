# Phase 4: CLI Migration to Cobra and Viper

## Overview
Migrate from `flag` package to `cobra` for better CLI structure, auto-generated `--help`, and improved configuration management with `viper`.

## Motivation
- Auto-generated `--help` with proper formatting and documentation
- Better subcommand support for future extensibility
- Unified configuration management (CLI flags + config file + environment variables)
- Standard Go CLI patterns that developers are familiar with
- Better error handling and exit codes

## Phase 4.1: Setup and Dependencies

### Files to Create:
- `cmd/david/cli/root.go` - Root command setup
- `cmd/david/cli/server.go` - Server subcommand
- `cmd/dcrypt/cli/root.go` - BCPT root command
- `cmd/dcrypt/cli/passwd.go` - BCPT passwd subcommand

### Files to Modify:
- `cmd/david/main.go` - Replace with simple CLI invocation
- `cmd/dcrypt/main.go` - Replace with simple CLI invocation
- `go.mod` - Add cobra and viper dependencies

### Dependencies to Add:
```go
github.com/spf13/cobra
github.com/spf13/viper
```

## Phase 4.2: Root Command Structure

### Root Command (`cmd/david/cli/root.go`):
```go
var rootCmd = &cobra.Command{
    Use:   "david",
    Short: "A simple WebDAV server... extended.",
    Long: `david is a simple WebDAV server that provides:
- Single binary that runs under Windows, Linux and OSX.
- Authentication via HTTP-Basic.
- CRUD operation permissions
- TLS support
- A simple user management which allows user-directory-jails
- Live config reload
- A CLI tool to generate BCrypt password hashes (dcrypt)`,
    Run: func(cmd *cobra.Command, args []string) {
        // Default behavior: run server
        server.Run()
    },
}
```

### Server Command (`cmd/david/cli/server.go`):
```go
var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Start the WebDAV server",
    Long:  `Start the WebDAV server with optional configuration overrides.`,
    Run: func(cmd *cobra.Command, args []string) {
        configPath, _ := cmd.Flags().GetString("config")
        host, _ := cmd.Flags().GetString("host")
        port, _ := cmd.Flags().GetString("port")
        
        // Use viper to merge config sources
        // CLI flags > config file > defaults
        
        // Run server logic
        app.RunServer(configPath, host, port)
    },
}

func init() {
    rootCmd.AddCommand(serverCmd)
    
    serverCmd.Flags().StringP("config", "c", "", "Path to configuration file")
    serverCmd.Flags().StringP("host", "H", "", "Override host address")
    serverCmd.Flags().StringP("port", "p", "", "Override port")
    serverCmd.Flags().BoolP("debug", "d", false, "Enable debug logging")
    serverCmd.Flags().BoolP("production", "", false, "Enable production (JSON) logging")
    
    viper.BindPFlags(serverCmd.Flags())
}
```

## Phase 4.3: Viper Integration

### Configuration Loading (`app/config.go`):
```go
func ParseConfig(configPath string) *Config {
    // Use viper for config management
    viper.SetConfigFile(configPath)
    viper.SetConfigType("yaml")
    
    // Try to read config
    if err := viper.ReadInConfig(); err != nil {
        log.Warnf("No config file found, using defaults: %v", err)
    }
    
    // Unmarshal into config struct
    cfg := &Config{}
    if err := viper.Unmarshal(cfg); err != nil {
        log.Fatalf("Error parsing config: %v", err)
    }
    
    return cfg
}
```

### Viper Configuration Sources:
1. **CLI flags** (highest priority)
2. **Config file** (`config.yaml`)
3. **Environment variables** (optional, e.g., `DAVID_HOST`, `DAVID_PORT`)
4. **Defaults** (lowest priority)

## Phase 4.4: Help Documentation

### Auto-generated Help:
Cobra automatically generates help based on command structure:
- `david --help` - Shows all available commands and flags
- `david server --help` - Shows server-specific flags and documentation
- `david server -h` - Short help

### Documentation Best Practices:
- Each flag needs `Short` and `Long` descriptions
- Use `Example` field for usage examples
- Add `Aliases` for common abbreviations
- Add `Annotations` for help categorization

### Example Help Output:
```
$ david server --help
Start the WebDAV server with optional configuration overrides.

Usage:
  david server [flags]

Flags:
  -c, --config string   Path to configuration file
  -d, --debug           Enable debug logging
  -H, --host string     Override host address (default "0.0.0.0")
  -p, --port string     Override port (default "8000")
  -h, --help            Help for server
      --production      Enable production (JSON) logging

Global Flags:
      --version   Print version information
```

## Phase 4.5: BCPT CLI Migration

### BCPT Command Structure:
```go
var bcptCmd = &cobra.Command{
    Use:   "dcrypt",
    Short: "BCrypt password hash generator",
    Long:  `dcrypt is a CLI tool to generate BCrypt password hashes for david configuration.`,
}

var passwdCmd = &cobra.Command{
    Use:   "passwd",
    Short: "Generate a BCrypt password hash",
    Long:  `Generate a BCrypt password hash for use in david configuration.`,
    Run: func(cmd *cobra.Command, args []string) {
        password, _ := cmd.Flags().GetString("password")
        cost, _ := cmd.Flags().GetInt("cost")
        
        hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
        if err != nil {
            log.Fatalf("Error generating hash: %v", err)
        }
        
        fmt.Printf("$2a$%02d$%s\n", cost, string(hash))
    },
}

func init() {
    bcptCmd.AddCommand(passwdCmd)
    
    passwdCmd.Flags().StringP("password", "p", "", "Password to hash (required)")
    passwdCmd.Flags().IntP("cost", "c", 10, "BCrypt cost factor")
    passwdCmd.MarkFlagRequired("password")
}
```

## Phase 4.6: Migration Steps

### Step 1: Add Dependencies
```bash
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
```

### Step 2: Create CLI Structure
- Create `cmd/david/cli/` directory
- Create `cmd/dcrypt/cli/` directory
- Add root commands
- Add subcommands

### Step 3: Migrate Main Functions
- Replace `main.go` with simple CLI entry point
- Delegate to CLI command structure

### Step 4: Update Config Parsing
- Integrate viper into `app/config.go`
- Add support for environment variables
- Update config file scanning logic

### Step 5: Test CLI
- Test all flag combinations
- Test help output
- Test config file loading
- Test environment variable overrides

### Step 6: Update Documentation
- Update README with CLI flags
- Add examples for common use cases
- Document environment variables

## Phase 4.7: Testing Strategy

### Unit Tests:
- Test CLI command structure
- Test flag parsing
- Test config loading with viper

### Integration Tests:
- Test full CLI workflow
- Test config file overrides
- Test environment variable overrides

## Phase 4.8: Migration Checklist

- [ ] Add cobra and viper dependencies
- [ ] Create root command structure
- [ ] Implement server subcommand
- [ ] Implement passwd subcommand
- [ ] Integrate viper for config management
- [ ] Add environment variable support
- [ ] Test all flag combinations
- [ ] Generate and verify help output
- [ ] Update README documentation
- [ ] Remove old flag package usage
- [ ] Update CI/CD if needed

## Phase 4.9: Timeline Estimate
- Phase 4.1-4.2: 2-3 hours (setup and structure)
- Phase 4.3-4.4: 2-3 hours (viper integration and help)
- Phase 4.5-4.6: 2-3 hours (BCPT migration and main refactoring)
- Phase 4.7-4.8: 1-2 hours (testing and documentation)
- **Total: 7-11 hours**

## Phase 4.10: Benefits
1. **Better UX**: Auto-generated help with proper formatting
2. **Flexibility**: Easy to add new subcommands (e.g., `david validate-config`)
3. **Unified config**: CLI flags, config file, and environment variables
4. **Standard patterns**: Developers familiar with cobra/viper can contribute
5. **Future-proof**: Foundation for additional CLI features