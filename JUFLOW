

cat​ ​ >>​ ​ ~/maas.yaml​ ​ <<​ ​ EOF
clouds:/
​ ​ maas:
​ ​ ​ ​ type:​ ​ maas
​ ​ ​ ​ auth-types:​ ​ [oauth1]
​ ​ ​ ​ endpoint:​ ​ http://192.168.100.3/MAAS/
EOF

juju​ ​ add-cloud​ ​ maas​ ​ ~/maas.yaml



/usr/bin/juju

#!/bin/sh
export PATH=/usr/lib/juju-2.2/bin:"$PATH"
exec juju "$@"


ls /usr/lib/juju-2.2/bin
juju  jujud  juju-metadata





juju/cmd/README.md

## func Main
``` go
func Main(c Command, ctx *Context, args []string) int
```
Main runs the given Command in the supplied Context with the given
arguments, which should not include the command name. It returns a code
suitable for passing to os.Exit.






juju/cmd/cmd.go


// Main runs the given Command in the supplied Context with the given
// arguments, which should not include the command name. It returns a code
// suitable for passing to os.Exit.
func Main(c Command, ctx *Context, args []string) int {
        f := gnuflag.NewFlagSet(c.Info().Name, gnuflag.ContinueOnError)
        f.SetOutput(ioutil.Discard)
        c.SetFlags(f)
        if rc, done := handleCommandError(c, ctx, f.Parse(c.AllowInterspersedFlags(), args), f); done {
                return rc
        }
        // Since SuperCommands can also return gnuflag.ErrHelp errors, we need to
        // handle both those types of errors as well as "real" errors.
        if rc, done := handleCommandError(c, ctx, c.Init(f.Args()), f); done {
                return rc
        }
        if err := c.Run(ctx); err != nil {
                if IsRcPassthroughError(err) {
                        return err.(*RcPassthroughError).Code
                }
                if err != ErrSilent {
                        WriteError(ctx.Stderr, err)
                }
                return 1
        }
        return 0
}



// Command is implemented by types that interpret command-line arguments.
type Command interface {
        // IsSuperCommand returns true if the command is a super command.
        IsSuperCommand() bool

        // Info returns information about the Command.
        Info() *Info

        // SetFlags adds command specific flags to the flag set.
        SetFlags(f *gnuflag.FlagSet)

        // Init initializes the Command before running.
        Init(args []string) error

        // Run will execute the Command as directed by the options and positional
        // arguments passed to Init.
        Run(ctx *Context) error

        // AllowInterspersedFlags returns whether the command allows flag
        // arguments to be interspersed with non-flag arguments.
        AllowInterspersedFlags() bool
}



====================================================================================




./cmd/juju/commands/main.go


        "github.com/juju/cmd"


// Run is the main entry point for the juju client.
func (m main) Run(args []string) int {
        ctx, err := cmd.DefaultContext()
        if err != nil {
                cmd.WriteError(os.Stderr, err)
                return 2
        }

        // note that this has to come before we init the juju home directory,
        // since it relies on detecting the lack of said directory.
        newInstall := m.maybeWarnJuju1x()

        if err = juju.InitJujuXDGDataHome(); err != nil {
                cmd.WriteError(ctx.Stderr, err)
                return 2
        }

        if err := installProxy(); err != nil {
                cmd.WriteError(ctx.Stderr, err)
                return 2
        }

        if newInstall {
                fmt.Fprintf(ctx.Stderr, "Since Juju %v is being run for the first time, downloading latest cloud information.\n", jujuversion.Current.Major)
                updateCmd := cloud.NewUpdateCloudsCommand()
                if err := updateCmd.Run(ctx); err != nil {
                        cmd.WriteError(ctx.Stderr, err)
                }
        }

        for i := range x {
                x[i] ^= 255
        }
        if len(args) == 2 && args[1] == string(x[0:2]) {
                os.Stdout.Write(x[2:])
                return 0
        }
 
        jcmd := NewJujuCommand(ctx)                <============
        return cmd.Main(jcmd, ctx, args[1:])
}


        jujucmd "github.com/juju/juju/cmd"


// NewJujuCommand ...
func NewJujuCommand(ctx *cmd.Context) cmd.Command {
        jcmd := jujucmd.NewSuperCommand(cmd.SuperCommandParams{
                Name:                "juju",
                Doc:                 jujuDoc,
                MissingCallback:     RunPlugin,
                UserAliasesFilename: osenv.JujuXDGDataHomePath("aliases"),
        })
        jcmd.AddHelpTopic("basics", "Basic Help Summary", usageHelp)
        registerCommands(jcmd, ctx)                                         <============
        return jcmd
}





// TODO(ericsnow) Factor out the commands and aliases into a static
// registry that can be passed to the supercommand separately.

// registerCommands registers commands in the specified registry.
func registerCommands(r commandRegistry, ctx *cmd.Context) {
        // Creation commands.
        r.Register(newBootstrapCommand())
        r.Register(application.NewAddRelationCommand())

...
        // Manage clouds and credentials
        r.Register(cloud.NewUpdateCloudsCommand())
        r.Register(cloud.NewListCloudsCommand())
        r.Register(cloud.NewListRegionsCommand())
        r.Register(cloud.NewShowCloudCommand())
        r.Register(cloud.NewAddCloudCommand(&cloudToCommandAdapter{}))      =====> juju/cmd/supercommand.go 
                                                                                   func (c *SuperCommand) Register(subcmd Command)
        r.Register(cloud.NewRemoveCloudCommand())
        r.Register(cloud.NewListCredentialsCommand())
        r.Register(cloud.NewDetectCredentialsCommand())
        r.Register(cloud.NewSetDefaultRegionCommand())
        r.Register(cloud.NewSetDefaultCredentialCommand())
        r.Register(cloud.NewAddCredentialCommand())
        r.Register(cloud.NewRemoveCredentialCommand())
        r.Register(cloud.NewUpdateCredentialCommand())


type cloudToCommandAdapter struct{}



./cmd/juju/commands/main.go


type commandRegistry interface {
        Register(cmd.Command)
        RegisterSuperAlias(name, super, forName string, check cmd.DeprecationCheck)
        RegisterDeprecated(subcmd cmd.Command, check cmd.DeprecationCheck)
}



------------------------------------------------------------------------------

juju/juju/cmd/supercommand.go

        "github.com/juju/cmd"


// NewSuperCommand is like cmd.NewSuperCommand but
// it adds juju-specific functionality:
// - The default logging configuration is taken from the environment;
// - The version is configured to the current juju version;
// - The command emits a log message when a command runs.
func NewSuperCommand(p cmd.SuperCommandParams) *cmd.SuperCommand {
        p.Log = &cmd.Log{
                DefaultConfig: os.Getenv(osenv.JujuLoggingConfigEnvKey),
        }
        current := version.Binary{
                Number: jujuversion.Current,
                Arch:   arch.HostArch(),
                Series: series.MustHostSeries(),
        }

        // p.Version should be a version.Binary, but juju/cmd does not
        // import juju/juju/version so this cannot happen. We have
        // tests to assert that this string value is correct.
        p.Version = current.String()
        p.NotifyRun = runNotifier
        return cmd.NewSuperCommand(p)
}

-----------------------------------------

juju/cmd/supercommand.go

// NewSuperCommand creates and initializes a new `SuperCommand`, and returns
// the fully initialized structure.
func NewSuperCommand(params SuperCommandParams) *SuperCommand {
        command := &SuperCommand{
                Name:                params.Name,
                Purpose:             params.Purpose,
                Doc:                 params.Doc,
                Log:                 params.Log,
                usagePrefix:         params.UsagePrefix,
                missingCallback:     params.MissingCallback,
                Aliases:             params.Aliases,
                version:             params.Version,
                notifyRun:           params.NotifyRun,
                notifyHelp:          params.NotifyHelp,
                userAliasesFilename: params.UserAliasesFilename,
        }
        command.init()
        return command
}


func (c *SuperCommand) init() {
        if c.subcmds != nil {
                return
        }
        c.help = &helpCommand{
                super: c,
        }
        c.help.init()
        c.subcmds = map[string]commandReference{
                "help": commandReference{command: c.help},
        }
        if c.version != "" {
                c.subcmds["version"] = commandReference{
                        command: newVersionCommand(c.version),
                }
        }

        c.userAliases = ParseAliasFile(c.userAliasesFilename)
}



// SuperCommand is a Command that selects a subcommand and assumes its
// properties; any command line arguments that were not used in selecting
// the subcommand are passed down to it, and to Run a SuperCommand is to run
// its selected subcommand.
type SuperCommand struct {
        CommandBase
        Name                string
        Purpose             string
        Doc                 string
        Log                 *Log
        Aliases             []string
        version             string
        usagePrefix         string
        userAliasesFilename string
        userAliases         map[string][]string
        subcmds             map[string]commandReference
        help                *helpCommand
        commonflags         *gnuflag.FlagSet
        flags               *gnuflag.FlagSet
        action              commandReference
        showHelp            bool
        showDescription     bool
        showVersion         bool
        noAlias             bool
        missingCallback     MissingCallback
        notifyRun           func(string)
        notifyHelp          func([]string)
}


// SuperCommandParams provides a way to have default parameter to the
// `NewSuperCommand` call.
type SuperCommandParams struct {
        // UsagePrefix should be set when the SuperCommand is
        // actually a subcommand of some other SuperCommand;
        // if NotifyRun is called, it name will be prefixed accordingly,
        // unless UsagePrefix is identical to Name.
        UsagePrefix string

        // Notify, if not nil, is called when the SuperCommand
        // is about to run a sub-command.
        NotifyRun func(cmdName string)

        // NotifyHelp is called just before help is printed, with the
        // arguments received by the help command. This can be
        // used, for example, to load command information for external
        // "plugin" commands, so that their documentation will show up
        // in the help output.
        NotifyHelp func([]string)

        Name            string
        Purpose         string
        Doc             string
        Log             *Log
        MissingCallback MissingCallback
        Aliases         []string
        Version         string

        // UserAliasesFilename refers to the location of a file that contains
        //   name = cmd [args...]
        // values, that is used to change default behaviour of commands in order
        // to add flags, or provide short cuts to longer commands.
        UserAliasesFilename string
}



vi juju/cmd/cmd.go

// CommandBase provides the default implementation for SetFlags, Init, and Help.
type CommandBase struct{}




------------------------------------------------------------------------------


vi ./juju/juju/cmd/juju/cloud/add.go

        "github.com/juju/cmd"

// AddCloudCommand is the command that allows you to add a cloud configuration
// for use with juju bootstrap.
type AddCloudCommand struct {
        cmd.CommandBase

        // Replace, if true, existing cloud information is overwritten.
        Replace bool

        // Cloud is the name fo the cloud to add.
        Cloud string

        // CloudFile is the name of the cloud YAML file.
        CloudFile string

        // Ping contains the logic for pinging a cloud endpoint to know whether or
        // not it really has a valid cloud of the same type as the provider.  By
        // default it just calls the correct provider's Ping method.
        Ping func(p environs.EnvironProvider, endpoint string) error

        cloudMetadataStore CloudMetadataStore
}



// Info returns help information about the command.
func (c *AddCloudCommand) Info() *cmd.Info {
        return &cmd.Info{
                Name:    "add-cloud",
                Args:    "<cloud name> <cloud definition file>",
                Purpose: usageAddCloudSummary,
                Doc:     usageAddCloudDetails,
        }
}



vi juju/cmd/cmd.go

// CommandBase provides the default implementation for SetFlags, Init, and Help.
type CommandBase struct{}



./cmd/juju/cloud/add.go

// NewAddCloudCommand returns a command to add cloud information.
func NewAddCloudCommand(cloudMetadataStore CloudMetadataStore) *AddCloudCommand {
        // Ping is provider.Ping except in tests where we don't actually want to
        // require a valid cloud.
        return &AddCloudCommand{
                cloudMetadataStore: cloudMetadataStore,
                Ping: func(p environs.EnvironProvider, endpoint string) error {
                        return p.Ping(endpoint)
                },
        }
}




./cmd/juju/cloud/add.go

type CloudMetadataStore interface {
        ParseCloudMetadataFile(path string) (map[string]cloud.Cloud, error)
        ParseOneCloud(data []byte) (cloud.Cloud, error)
        PublicCloudMetadata(searchPaths ...string) (result map[string]cloud.Cloud, fallbackUsed bool, _ error)
        PersonalCloudMetadata() (map[string]cloud.Cloud, error)
        WritePersonalCloudMetadata(cloudsMap map[string]cloud.Cloud) error
}




==================================================================================

juju/cmd Register <===== juju/juju AddCloudCommand


juju/juju/cmd/juju/commands/main.go

func registerCommands(r commandRegistry, ctx *cmd.Context) {

...
        // Manage clouds and credentials
        r.Register(cloud.NewAddCloudCommand(&cloudToCommandAdapter{}))



juju/juju/cmd/juju/cloud/add.go

// NewAddCloudCommand returns a command to add cloud information.
func NewAddCloudCommand(cloudMetadataStore CloudMetadataStore) *AddCloudCommand {
        // Ping is provider.Ping except in tests where we don't actually want to
        // require a valid cloud.
        return &AddCloudCommand{
                cloudMetadataStore: cloudMetadataStore,
                Ping: func(p environs.EnvironProvider, endpoint string) error {
                        return p.Ping(endpoint)
                },
        }
}



vi ./juju/juju/cmd/juju/cloud/add.go

        "github.com/juju/cmd"

// AddCloudCommand is the command that allows you to add a cloud configuration
// for use with juju bootstrap.
type AddCloudCommand struct {
        cmd.CommandBase                      <========== juju/cmd/cmd.go  CommandBase struct

        // Replace, if true, existing cloud information is overwritten.
        Replace bool

        // Cloud is the name fo the cloud to add.
        Cloud string

        // CloudFile is the name of the cloud YAML file.
        CloudFile string

        // Ping contains the logic for pinging a cloud endpoint to know whether or
        // not it really has a valid cloud of the same type as the provider.  By
        // default it just calls the correct provider's Ping method.
        Ping func(p environs.EnvironProvider, endpoint string) error

        cloudMetadataStore CloudMetadataStore
}


// Info returns help information about the command.
func (c *AddCloudCommand) Info() *cmd.Info {
        return &cmd.Info{
                Name:    "add-cloud",
                Args:    "<cloud name> <cloud definition file>",
                Purpose: usageAddCloudSummary,
                Doc:     usageAddCloudDetails,
        }
}



====================================================================================


juju/cmd/cmd.go

// CommandBase provides the default implementation for SetFlags, Init, and Help.
type CommandBase struct{}





juju/cmd/supercommand.go

// Register makes a subcommand available for use on the command line. The
// command will be available via its own name, and via any supplied aliases.
func (c *SuperCommand) Register(subcmd Command) {
        info := subcmd.Info()
        c.insert(commandReference{name: info.Name, command: subcmd})
        for _, name := range info.Aliases {
                c.insert(commandReference{name: name, command: subcmd, alias: info.Name})
        }
}




juju/cmd/cmd.go

// Command is implemented by types that interpret command-line arguments.
type Command interface {
        // IsSuperCommand returns true if the command is a super command.
        IsSuperCommand() bool

        // Info returns information about the Command.
        Info() *Info

        // SetFlags adds command specific flags to the flag set.
        SetFlags(f *gnuflag.FlagSet)

        // Init initializes the Command before running.
        Init(args []string) error

        // Run will execute the Command as directed by the options and positional
        // arguments passed to Init.
        Run(ctx *Context) error

        // AllowInterspersedFlags returns whether the command allows flag
        // arguments to be interspersed with non-flag arguments.
        AllowInterspersedFlags() bool
}


====================================================================================


juju/juju/cmd/juju/cloud/add.go

// Run executes the add cloud command, adding a cloud based on a passed-in yaml
// file or interactive queries.
func (c *AddCloudCommand) Run(ctxt *cmd.Context) error {
        if c.CloudFile == "" {
                return c.runInteractive(ctxt)
        }

        specifiedClouds, err := c.cloudMetadataStore.ParseCloudMetadataFile(c.CloudFile)
        if err != nil {
                return err
        }
        if specifiedClouds == nil {
                return errors.New("no personal clouds are defined")
        }
        newCloud, ok := specifiedClouds[c.Cloud]
        if !ok {
                return errors.Errorf("cloud %q not found in file %q", c.Cloud, c.CloudFile)
        }

        // first validate cloud input
        data, err := ioutil.ReadFile(c.CloudFile)
        if err != nil {
                return errors.Trace(err)
        }
        if err = cloud.ValidateCloudSet([]byte(data)); err != nil {
                ctxt.Warningf(err.Error())
        }

        // validate cloud data
        provider, err := environs.Provider(newCloud.Type)
        if err != nil {
                return errors.Trace(err)
        }
        schemas := provider.CredentialSchemas()
        for _, authType := range newCloud.AuthTypes {
                if _, defined := schemas[authType]; !defined {
                        return errors.NotSupportedf("auth type %q", authType)
                }
        }
        if err := c.verifyName(c.Cloud); err != nil {
                return errors.Trace(err)
        }

        return addCloud(c.cloudMetadataStore, newCloud)
}






