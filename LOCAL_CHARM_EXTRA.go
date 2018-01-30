

juju/juju/juju/osenv/home.go:45

// JujuXDGDataHomePath returns the path to a file in the
// current juju home.
func JujuXDGDataHomePath(names ...string) string {
        all := append([]string{JujuXDGDataHomeDir()}, names...)
        return filepath.Join(all...)
}


-----------------------------------------------------------------------------------------------

juju/juju/juju/osenv/home.go

// JujuXDGDataHomeDir returns the directory where juju should store application-specific files
func JujuXDGDataHomeDir() string {
        JujuXDGDataHomeDir := JujuXDGDataHome()
        if JujuXDGDataHomeDir != "" {
                return JujuXDGDataHomeDir
        }
        JujuXDGDataHomeDir = os.Getenv(JujuXDGDataHomeEnvKey)
        if JujuXDGDataHomeDir == "" {
                if runtime.GOOS == "windows" {
                        JujuXDGDataHomeDir = jujuXDGDataHomeWin()
                } else {
                        JujuXDGDataHomeDir = jujuXDGDataHomeLinux()
                }
        }
        return JujuXDGDataHomeDir
}


-----------------------------------------------------------------------------------------------

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


-----------------------------------------------------------------------------------------------

juju/cmd/supercommand.go:149


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



-----------------------------------------------------------------------------------------------

juju/cmd/supercommand.go

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

-----------------------------------------------------------------------------------------------

juju/cmd/supercommand.go:119

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

-----------------------------------------------------------------------------------------------

juju/cmd.go:86

// CommandBase provides the default implementation for SetFlags, Init, and Help.
type CommandBase struct{}


-----------------------------------------------------------------------------------------------

juju/cmd/supercommand.go:108

type commandReference struct {
        name    string
        command Command
        alias   string
        check   DeprecationCheck
}


-----------------------------------------------------------------------------------------------

github.com/juju/cmd/cmd.go:63

type Command interface {

        IsSuperCommand() bool

        Info() *Info

        SetFlags(f *gnuflag.FlagSet)

        Init(args []string) error

        Run(ctx *Context) error                 <=======

        AllowInterspersedFlags() bool
}


-----------------------------------------------------------------------------------------------

juju/juju/cmd/juju/commands/plugin.go:50


func RunPlugin(ctx *cmd.Context, subcommand string, args []string) error {
        cmdName := JujuPluginPrefix + subcommand
        plugin := modelcmd.Wrap(&PluginCommand{name: cmdName})

        // We process common flags supported by Juju commands.
        // To do this, we extract only those supported flags from the
        // argument list to avoid confusing flags.Parse().
        flags := gnuflag.NewFlagSet(cmdName, gnuflag.ContinueOnError)
        flags.SetOutput(ioutil.Discard)
        plugin.SetFlags(flags)
        jujuArgs := extractJujuArgs(args)
        if err := flags.Parse(false, jujuArgs); err != nil {
                return err
        }
        if err := plugin.Init(args); err != nil {
                return err
        }
        err := plugin.Run(ctx)
        _, execError := err.(*exec.Error)
        // exec.Error results are for when the executable isn't found, in
        // those cases, drop through.
        if !execError {
                return err
        }
        return &cmd.UnrecognizedCommand{Name: subcommand}
}


=============================================================================================

juju/juju/cmd/juju/commands/main.go

        jujucmd "github.com/juju/juju/cmd"
        "github.com/juju/cmd"
        "github.com/juju/juju/juju/osenv"


// NewJujuCommand ...
func NewJujuCommand(ctx *cmd.Context) cmd.Command {
        jcmd := jujucmd.NewSuperCommand(cmd.SuperCommandParams{
                Name:                "juju",
                Doc:                 jujuDoc,
                MissingCallback:     RunPlugin,
                UserAliasesFilename: osenv.JujuXDGDataHomePath("aliases"),
        })
        jcmd.AddHelpTopic("basics", "Basic Help Summary", usageHelp)
        registerCommands(jcmd, ctx)
        return jcmd
}


-----------------------------------------------------------------------------------------------

juju/juju/cmd/juju/commands/main.go:272

        rcmd "github.com/juju/juju/cmd/juju/romulus/commands"


// TODO(ericsnow) Factor out the commands and aliases into a static
// registry that can be passed to the supercommand separately.

// registerCommands registers commands in the specified registry.
func registerCommands(r commandRegistry, ctx *cmd.Context) {


        // Manage and control services
        r.Register(application.NewAddUnitCommand())
        r.Register(application.NewConfigCommand())
        r.Register(application.NewDeployCommand())       <=========


        rcmd.RegisterAll(r)
}


-----------------------------------------------------------------------------------------------


juju/juju/cmd/juju/romulus/commands/commands.go:28:func RegisterAll(r commandRegister) {


// RegisterAll registers all romulus commands with the
// provided command registry.
func RegisterAll(r commandRegister) {
        r.Register(agree.NewAgreeCommand())
        r.Register(listagreements.NewListAgreementsCommand())
        r.Register(budget.NewBudgetCommand())
        r.Register(createwallet.NewCreateWalletCommand())
        r.Register(listplans.NewListPlansCommand())
        r.Register(setwallet.NewSetWalletCommand())
        r.Register(setplan.NewSetPlanCommand())
        r.Register(showwallet.NewShowWalletCommand())
        r.Register(sla.NewSLACommand())
        r.Register(listwallets.NewListWalletsCommand())
}



-----------------------------------------------------------------------------------------------


./cmd/juju/commands/main.go:262

type commandRegistry interface {
        Register(cmd.Command)
        RegisterSuperAlias(name, super, forName string, check cmd.DeprecationCheck)
        RegisterDeprecated(subcmd cmd.Command, check cmd.DeprecationCheck)
}

-----------------------------------------------------------------------------------------------



juju/cmd/supercommand.go:185

// Register makes a subcommand available for use on the command line. The
// command will be available via its own name, and via any supplied aliases.
func (c *SuperCommand) Register(subcmd Command) {
        info := subcmd.Info()
        c.insert(commandReference{name: info.Name, command: subcmd})
        for _, name := range info.Aliases {
                c.insert(commandReference{name: name, command: subcmd, alias: info.Name})
        }
}








////////////////////////////////////////////////////////////////////////////////////////////////



        "github.com/juju/juju/cmd/modelcmd"


juju/juju/cmd/juju/application/deploy.go:225

// NewDeployCommand returns a command to deploy applications.
func NewDeployCommand() modelcmd.ModelCommand {
        steps := []DeployStep{
                &RegisterMeteredCharm{
                        RegisterURL: planURL + "/plan/authorize",
                        QueryURL:    planURL + "/charm",
                },
        }
        deployCmd := &DeployCommand{
                Steps: steps,
        }
        deployCmd.NewAPIRoot = func() (DeployAPI, error) {
                apiRoot, err := deployCmd.ModelCommandBase.NewAPIRoot()
                if err != nil {
                        return nil, errors.Trace(err)
                }
                bakeryClient, err := deployCmd.BakeryClient()
                if err != nil {
                        return nil, errors.Trace(err)
                }
                cstoreClient := newCharmStoreClient(bakeryClient).WithChannel(deployCmd.Channel)

                return &deployAPIAdapter{
                        Connection:        apiRoot,
                        apiClient:         &apiClient{Client: apiRoot.Client()},
                        charmsClient:      &charmsClient{Client: apicharms.NewClient(apiRoot)},
                        applicationClient: &applicationClient{Client: application.NewClient(apiRoot)},
                        modelConfigClient: &modelConfigClient{Client: modelconfig.NewClient(apiRoot)},
                        charmstoreClient:  &charmstoreClient{Client: cstoreClient},
                        annotationsClient: &annotationsClient{Client: annotations.NewClient(apiRoot)},
                        charmRepoClient:   &charmRepoClient{CharmStore: charmrepo.NewCharmStoreFromClient(cstoreClient)},
                }, nil
        }

        return modelcmd.Wrap(deployCmd)
}

-----------------------------------------------------------------------------------------------


juju/juju/cmd/modelcmd/modelcommand.go:316

// Wrap wraps the specified ModelCommand, returning a ModelCommand
// that proxies to each of the ModelCommand methods.
// Any provided options are applied to the wrapped command
// before it is returned.
func Wrap(c ModelCommand, options ...WrapOption) ModelCommand {
        wrapper := &modelCommandWrapper{
                ModelCommand:    c,
                skipModelFlags:  false,
                useDefaultModel: true,
        }
        for _, option := range options {
                option(wrapper)
        }
        // Define a new type so that we can embed the ModelCommand
        // interface one level deeper than cmd.Command, so that
        // we'll get the Command methods from WrapBase
        // and all the ModelCommand methods not in cmd.Command
        // from modelCommandWrapper.
        type embed struct {
                *modelCommandWrapper
        }
        return struct {
                embed
                cmd.Command
        }{
                Command: WrapBase(wrapper),
                embed:   embed{wrapper},
        }
}


-----------------------------------------------------------------------------------------------

juju/juju/cmd/modelcmd/modelcommand.go


// ModelCommand extends cmd.Command with a SetModelName method.
type ModelCommand interface {
        Command

        // SetClientStore is called prior to the wrapped command's Init method
        // with the default controller store. It may also be called to override the
        // default controller store for testing.
        SetClientStore(jujuclient.ClientStore)

        // ClientStore returns the controller store that the command is
        // associated with.
        ClientStore() jujuclient.ClientStore

        // SetModelName sets the model name for this command. Setting the model
        // name will also set the related controller name. The model name can
        // be qualified with a controller name (controller:model), or
        // unqualified, in which case it will be assumed to be within the
        // current controller.
        //
        // Passing an empty model name will choose the default
        // model, or return an error if there isn't one.
        //
        // SetModelName is called prior to the wrapped command's Init method
        // with the active model name. The model name is guaranteed
        // to be non-empty at entry of Init.
        SetModelName(modelName string, allowDefault bool) error

        // ModelName returns the name of the model.
        ModelName() (string, error)

        // ControllerName returns the name of the controller that contains
        // the model returned by ModelName().
        ControllerName() (string, error)

        // initModel initializes the model name, resolving empty
        // model or controller parts to the current model or controller if
        // needed. It fails a model cannot be determined.
        initModel() error
}

-----------------------------------------------------------------------------------------------

juju/juju/cmd/modelcmd/base.go:36:type Command interface {

        "github.com/juju/cmd"


// Command extends cmd.Command with a closeContext method.
// It is implicitly implemented by any type that embeds CommandBase.
type Command interface {
        cmd.Command

        // SetAPIOpen sets the function used for opening an API connection.
        SetAPIOpen(opener api.OpenFunc)

        // SetModelAPI sets the api used to access model information.
        SetModelAPI(api ModelAPI)

        // closeAPIContexts closes any API contexts that have been opened.
        closeAPIContexts()
        initContexts(*cmd.Context)
        setRunStarted()
}


-----------------------------------------------------------------------------------------------

juju/cmd.go:63

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

-----------------------------------------------------------------------------------------------

./cmd/juju/application/deploy.go:915

func (c *DeployCommand) Run(ctx *cmd.Context) error {
        var err error
        c.Constraints, err = common.ParseConstraints(ctx, c.ConstraintsStr)
        if err != nil {
                return err
        }
        apiRoot, err := c.NewAPIRoot()
        if err != nil {
                return errors.Trace(err)
        }
        defer apiRoot.Close()

        deploy, err := findDeployerFIFO(
                c.maybeReadLocalBundle,
                func() (deployFn, error) { return c.maybeReadLocalCharm(apiRoot) },
                c.maybePredeployedLocalCharm,
                c.maybeReadCharmstoreBundleFn(apiRoot),
                c.charmStoreCharm, // This always returns a deployer
        )
        if err != nil {
                return errors.Trace(err)
        }

        return block.ProcessBlockedError(deploy(ctx, apiRoot), block.BlockChange)
}

-----------------------------------------------------------------------------------------------

./cmd/juju/application/deploy.go:563

func (c *DeployCommand) Init(args []string) error {
        if c.Force && c.Series == "" && c.PlacementSpec == "" {
                return errors.New("--force is only used with --series")
        }
        switch len(args) {
        case 2:
                if !names.IsValidApplication(args[1]) {
                        return errors.Errorf("invalid application name %q", args[1])
                }
                c.ApplicationName = args[1]
                fallthrough
        case 1:
                c.CharmOrBundle = args[0]
        case 0:
                return errors.New("no charm or bundle specified")
        default:
                return cmd.CheckEmpty(args[2:])
        }

        if err := c.parseBind(); err != nil {
                return err
        }

        useExisting, mapping, err := parseMachineMap(c.machineMap)
        if err != nil {
                return errors.Annotate(err, "error in --map-machines")
        }
        c.UseExisting = useExisting
        c.BundleMachines = mapping

        return c.UnitCommandBase.Init(args)
}

-----------------------------------------------------------------------------------------------

./cmd/juju/application/deploy.go:538

func (c *DeployCommand) SetFlags(f *gnuflag.FlagSet) {
        c.ConfigOptions.SetPreserveStringValue(true)
        // Keep above charmOnlyFlags and bundleOnlyFlags lists updated when adding
        // new flags.
        c.UnitCommandBase.SetFlags(f)
        c.ModelCommandBase.SetFlags(f)
        f.IntVar(&c.NumUnits, "n", 1, "Number of application units to deploy for principal charms")
        f.StringVar((*string)(&c.Channel), "channel", "", "Channel to use when getting the charm or bundle from the charm store")
        f.Var(&c.ConfigOptions, "config", "Either a path to yaml-formatted application config file or a key=value pair ")
        f.Var(cmd.NewAppendStringsValue(&c.BundleOverlayFile), "overlay", "Bundles to overlay on the primary bundle, applied in order")
        f.StringVar(&c.ConstraintsStr, "constraints", "", "Set application constraints")
        f.StringVar(&c.Series, "series", "", "The series on which to deploy")
        f.BoolVar(&c.DryRun, "dry-run", false, "Just show what the bundle deploy would do")
        f.BoolVar(&c.Force, "force", false, "Allow a charm to be deployed to a machine running an unsupported series")
        f.Var(storageFlag{&c.Storage, &c.BundleStorage}, "storage", "Charm storage constraints")
        f.Var(stringMap{&c.Resources}, "resource", "Resource to be uploaded to the controller")
        f.StringVar(&c.BindToSpaces, "bind", "", "Configure application endpoint bindings to spaces")
        f.StringVar(&c.machineMap, "map-machines", "", "Specify the existing machines to use for bundle deployments")

        for _, step := range c.Steps {
                step.SetFlags(f)
        }
        c.flagSet = f
}


-----------------------------------------------------------------------------------------------



./cmd/juju/application/deploy.go:516

func (c *DeployCommand) Info() *cmd.Info {
        return &cmd.Info{
                Name:    "deploy",
                Args:    "<charm or bundle> [<application name>]",
                Purpose: "Deploy a new application or bundle.",
                Doc:     deployDoc,
        }
}




////////////////////////////////////////////////////////////////////////////////////////////////








-----------------------------------------------------------------------------------------------

juju/cmd/supercommand.go

func (c *SuperCommand) insert(value commandReference) {
        if _, found := c.subcmds[value.name]; found {
                panic(fmt.Sprintf("command already registered: %q", value.name))
        }
        c.subcmds[value.name] = value
}


-----------------------------------------------------------------------------------------------

juju/cmd/supercommand.go:108

type commandReference struct {
        name    string
        command Command
        alias   string
        check   DeprecationCheck
}


-----------------------------------------------------------------------------------------------

github.com/juju/cmd/supercommand.go:119:type SuperCommand struct {

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

=============================================================================================

juju/cmd/supercommand.go:429


// Run executes the subcommand that was selected in Init.
func (c *SuperCommand) Run(ctx *Context) error {
        if c.showDescription {
                if c.Purpose != "" {
                        fmt.Fprintf(ctx.Stdout, "%s\n", c.Purpose)
                } else {
                        fmt.Fprintf(ctx.Stdout, "%s: no description available\n", c.Info().Name)
                }
                return nil
        }
        if c.action.command == nil {
                panic("Run: missing subcommand; Init failed or not called")
        }
        if c.Log != nil {
                if err := c.Log.Start(ctx); err != nil {
                        return err
                }
        }
        if c.notifyRun != nil {
                name := c.Name
                if c.usagePrefix != "" && c.usagePrefix != name {
                        name = c.usagePrefix + " " + name
                }
                c.notifyRun(name)
        }
        if deprecated, replacement := c.action.Deprecated(); deprecated {
                ctx.Infof("WARNING: %q is deprecated, please use %q", c.action.name, replacement)
        }
        err := c.action.command.Run(ctx)
        if err != nil && !IsErrSilent(err) {
                WriteError(ctx.Stderr, err)
                logger.Debugf("error stack: \n%v", errors.ErrorStack(err))
                // Now that this has been logged, don't log again in cmd.Main.
                if !IsRcPassthroughError(err) {
                        err = ErrSilent
                }
        } else {
                logger.Infof("command finished")
        }
        return err
}

-----------------------------------------------------------------------------------------------

juju/cmd/supercommand.go:119

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

-----------------------------------------------------------------------------------------------

juju/cmd/supercommand.go:108

type commandReference struct {
        name    string
        command Command
        alias   string
        check   DeprecationCheck
}

