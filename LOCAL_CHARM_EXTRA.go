

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


###############################################################################################

juju/juju/cmd/juju/application/deploy.go

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

        deploy, err := findDeployerFIFO(        <------
                c.maybeReadLocalBundle,
                func() (deployFn, error) { return c.maybeReadLocalCharm(apiRoot) },      <=========
                c.maybePredeployedLocalCharm,
                c.maybeReadCharmstoreBundleFn(apiRoot),
                c.charmStoreCharm, // This always returns a deployer
        )
        if err != nil {
                return errors.Trace(err)
        }

        return block.ProcessBlockedError(deploy(ctx, apiRoot), block.BlockChange)
}

----------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go:941


type deployFn func(*cmd.Context, DeployAPI) error


func findDeployerFIFO(maybeDeployers ...func() (deployFn, error)) (deployFn, error) {
        for _, d := range maybeDeployers {
                if deploy, err := d(); err != nil {
                        return nil, errors.Trace(err)
                } else if deploy != nil {
                        return deploy, nil
                }
        }
        return nil, errors.NotFoundf("suitable deployer")
}




~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
https://gobyexample.com/variadic-functions

func sum(nums ...int) {

Here’s a function that will take an arbitrary number of ints as arguments.


If you already have multiple args in a slice, apply them to a variadic function using func(slice...) like this.


func main() {

    nums := []int{1, 2, 3, 4}
    sum(nums...)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~




===================================================================================================

juju/juju/cmd/juju/application/deploy.go


        "github.com/juju/juju/cmd/modelcmd"


type DeployCommand struct {
        modelcmd.ModelCommandBase         <===========
        UnitCommandBase

        // CharmOrBundle is either a charm URL, a path where a charm can be found,
        // or a bundle name.
        CharmOrBundle string

        // BundleOverlay refers to config files that specify additional bundle
        // configuration to be merged with the main bundle.
        BundleOverlayFile []string

        // Channel holds the charmstore channel to use when obtaining
        // the charm to be deployed.
        Channel params.Channel

        // Series is the series of the charm to deploy.
        Series string

        // Force is used to allow a charm to be deployed onto a machine
        // running an unsupported series.
        Force bool

        // DryRun is used to specify that the bundle shouldn't actually be
        // deployed but just output the changes.
        DryRun bool

        ApplicationName string
        ConfigOptions   common.ConfigFlag
        ConstraintsStr  string
        Constraints     constraints.Value
        BindToSpaces    string

        // TODO(axw) move this to UnitCommandBase once we support --storage
        // on add-unit too.
        //
        // Storage is a map of storage constraints, keyed on the storage name
        // defined in charm storage metadata.
        Storage map[string]storage.Constraints

        // BundleStorage maps application names to maps of storage constraints keyed on
        // the storage name defined in that application's charm storage metadata.
        BundleStorage map[string]map[string]storage.Constraints

        // Resources is a map of resource name to filename to be uploaded on deploy.
        Resources map[string]string

        Bindings map[string]string
        Steps    []DeployStep

        // UseExisting machines when deploying the bundle.
        UseExisting bool
        // BundleMachines is a mapping for machines in the bundle to machines
        // in the model.
        BundleMachines map[string]string

        // NewAPIRoot stores a function which returns a new API root.
        NewAPIRoot func() (DeployAPI, error)

        machineMap string
        flagSet    *gnuflag.FlagSet
}

----------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go

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

        deploy, err := findDeployerFIFO(        <------
                c.maybeReadLocalBundle,
                func() (deployFn, error) { return c.maybeReadLocalCharm(apiRoot) },      <=========
                c.maybePredeployedLocalCharm,
                c.maybeReadCharmstoreBundleFn(apiRoot),
                c.charmStoreCharm, // This always returns a deployer
        )
        if err != nil {
                return errors.Trace(err)
        }

        return block.ProcessBlockedError(deploy(ctx, apiRoot), block.BlockChange)
}

----------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go:941


type deployFn func(*cmd.Context, DeployAPI) error


func findDeployerFIFO(maybeDeployers ...func() (deployFn, error)) (deployFn, error) {
        for _, d := range maybeDeployers {
                if deploy, err := d(); err != nil {       <==== !!!function are called here!!!
                        return nil, errors.Trace(err)
                } else if deploy != nil {
                        return deploy, nil
                }
        }
        return nil, errors.NotFoundf("suitable deployer")
}

----------------------------------------------------------------------------------------

./cmd/juju/application/deploy.go:997

        "gopkg.in/juju/charmrepo.v2"


func (c *DeployCommand) maybeReadLocalBundle() (deployFn, error) {
        bundleFile := c.CharmOrBundle    // CharmOrBundle is either a charm URL, a path where a charm can be found
        var bundleDir string
        isDir := false
        resolveDir := false

        bundleData, err := charmrepo.ReadBundleFile(bundleFile)
        if err != nil {
                // We may have been given a local bundle archive or exploded directory.
                bundle, _, pathErr := charmrepo.NewBundleAtPath(bundleFile)
                if charmrepo.IsInvalidPathError(pathErr) {
                        return nil, errors.Errorf(""+
                                "The charm or bundle %q is ambiguous.\n"+
                                "To deploy a local charm or bundle, run `juju deploy ./%[1]s`.\n"+
                                "To deploy a charm or bundle from the store, run `juju deploy cs:%[1]s`.",
                                bundleFile,
                        )
                }
                if pathErr != nil {
                        // If the bundle files existed but we couldn't read them,
                        // then return that error rather than trying to interpret
                        // as a charm.
                        if info, statErr := os.Stat(bundleFile); statErr == nil {
                                if info.IsDir() {
                                        if _, ok := pathErr.(*charmrepo.NotFoundError); !ok {
                                                return nil, errors.Annotate(pathErr, "cannot deploy bundle")
                                        }
                                }
                        }

                        logger.Debugf("cannot interpret as local bundle: %v", err)
                        return nil, nil
                }

                bundleData = bundle.Data()
                if info, err := os.Stat(bundleFile); err == nil && info.IsDir() {
                        resolveDir = true
                        isDir = true
                }
        } else {
                resolveDir = true
        }

        if err := c.validateBundleFlags(); err != nil {
                return nil, errors.Trace(err)
        }

        return func(ctx *cmd.Context, apiRoot DeployAPI) error {
                if resolveDir {
                        if isDir {
                                // If we get to here bundleFile is a directory, in which case
                                // we should use the absolute path as the bundFilePath, or it is
                                // an archive, in which case we should pass the empty string.
                                bundleDir = ctx.AbsPath(bundleFile)
                        } else {
                                // If the bundle is defined with just a yaml file, the bundle
                                // path is the directory that holds the file.
                                bundleDir = filepath.Dir(ctx.AbsPath(bundleFile))
                        }
                }
                return errors.Trace(c.deployBundle(
                        ctx,
                        bundleDir,
                        bundleData,
                        c.Channel,
                        apiRoot,
                        c.BundleStorage,
                ))
        }, nil
}


----------------------------------------------------------------------------------------

gopkg.in/juju/charmrepo.v2-unstable/bundlepath.go


        "gopkg.in/juju/charm.v6-unstable"


// ReadBundleFile attempts to read the file at path
// and interpret it as a bundle.
func ReadBundleFile(path string) (*charm.BundleData, error) {
        f, err := os.Open(path)
        if err != nil {
                if isNotExistsError(err) {
                        return nil, BundleNotFound(path)
                }
                return nil, err
        }
        defer f.Close()
        return charm.ReadBundleData(f)
}


----------------------------------------------------------------------------------------

gopkg.in/juju/charm.v6-unstable:242:func ReadBundleData(r io.Reader) (*BundleData, error) {

// ReadBundleData reads bundle data from the given reader.
// The returned data is not verified - call Verify to ensure
// that it is OK.
func ReadBundleData(r io.Reader) (*BundleData, error) {
        bytes, err := ioutil.ReadAll(r)
        if err != nil {
                return nil, err
        }
        var bd BundleData
        if err := yaml.Unmarshal(bytes, &bd); err != nil {
                return nil, fmt.Errorf("cannot unmarshal bundle data: %v", err)
        }
        return &bd, nil
}


----------------------------------------------------------------------------------------
gopkg.in/juju/charm.v6-unstable

// BundleData holds the contents of the bundle.
type BundleData struct {
        // Applications holds one entry for each application
        // that the bundle will create, indexed by
        // the application name.
        Applications map[string]*ApplicationSpec `bson:"applications,omitempty" json:"applications,omitempty" yaml:"applications,omitempty"`

        // Machines holds one entry for each machine referred to
        // by unit placements. These will be mapped onto actual
        // machines at bundle deployment time.
        // It is an error if a machine is specified but
        // not referred to by a unit placement directive.
        Machines map[string]*MachineSpec `bson:",omitempty" json:",omitempty" yaml:",omitempty"`

        // Series holds the default series to use when
        // the bundle chooses charms.
        Series string `bson:",omitempty" json:",omitempty" yaml:",omitempty"`

        // Relations holds a slice of 2-element slices,
        // each specifying a relation between two applications.
        // Each two-element slice holds two endpoints,
        // each specified as either colon-separated
        // (application, relation) pair or just an application name.
        // The relation is made between each. If the relation
        // name is omitted, it will be inferred from the available
        // relations defined in the applications' charms.
        Relations [][]string `bson:",omitempty" json:",omitempty" yaml:",omitempty"`

        // White listed set of tags to categorize bundles as we do charms.
        Tags []string `bson:",omitempty" json:",omitempty" yaml:",omitempty"`

        // Short paragraph explaining what the bundle is useful for.
        Description string `bson:",omitempty" json:",omitempty" yaml:",omitempty"`

        // unmarshaledWithServices holds whether the original marshaled data held a
        // legacy "services" field rather than the "applications" field.
        unmarshaledWithServices bool
}



      
===================================================================================================

juju/juju/cmd/juju/application/deploy.go

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

        deploy, err := findDeployerFIFO(        <------
                c.maybeReadLocalBundle,
                func() (deployFn, error) { return c.maybeReadLocalCharm(apiRoot) },      <=========
                c.maybePredeployedLocalCharm,
                c.maybeReadCharmstoreBundleFn(apiRoot),
                c.charmStoreCharm, // This always returns a deployer
        )
        if err != nil {
                return errors.Trace(err)
        }

        return block.ProcessBlockedError(deploy(ctx, apiRoot), block.BlockChange)
}

----------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go:941


type deployFn func(*cmd.Context, DeployAPI) error


func findDeployerFIFO(maybeDeployers ...func() (deployFn, error)) (deployFn, error) {
        for _, d := range maybeDeployers {
                if deploy, err := d(); err != nil {
                        return nil, errors.Trace(err)
                } else if deploy != nil {
                        return deploy, nil
                }
        }
        return nil, errors.NotFoundf("suitable deployer")
}

----------------------------------------------------------------------------------------

./cmd/juju/application/deploy.go

func (c *DeployCommand) maybeReadLocalCharm(apiRoot DeployAPI) (deployFn, error) {
        // NOTE: Here we select the series using the algorithm defined by
        // `seriesSelector.CharmSeries`. This serves to override the algorithm found in
        // `charmrepo.NewCharmAtPath` which is outdated (but must still be
        // called since the code is coupled with path interpretation logic which
        // cannot easily be factored out).

        // NOTE: Reading the charm here is only meant to aid in inferring the correct
        // series, if this fails we fall back to the argument series. If reading
        // the charm fails here it will also fail below (the charm is read again
        // below) where it is handled properly. This is just an expedient to get
        // the correct series. A proper refactoring of the charmrepo package is
        // needed for a more elegant fix.

        ch, err := charm.ReadCharm(c.CharmOrBundle)
        series := c.Series
        if err == nil {
                modelCfg, err := getModelConfig(apiRoot)
                if err != nil {
                        return nil, errors.Trace(err)
                }

                seriesSelector := seriesSelector{
                        seriesFlag:      series,
                        supportedSeries: ch.Meta().Series,
                        force:           c.Force,
                        conf:            modelCfg,
                        fromBundle:      false,
                }

                series, err = seriesSelector.charmSeries()
                if err != nil {
                        return nil, errors.Trace(err)
                }
        }

        // Charm may have been supplied via a path reference.
        ch, curl, err := charmrepo.NewCharmAtPathForceSeries(c.CharmOrBundle, series, c.Force)
        // We check for several types of known error which indicate
        // that the supplied reference was indeed a path but there was
        // an issue reading the charm located there.
        if charm.IsMissingSeriesError(err) {
                return nil, err
        } else if charm.IsUnsupportedSeriesError(err) {
                return nil, errors.Trace(err)
        } else if errors.Cause(err) == zip.ErrFormat {
                return nil, errors.Errorf("invalid charm or bundle provided at %q", c.CharmOrBundle)
        } else if _, ok := err.(*charmrepo.NotFoundError); ok {
                return nil, errors.Wrap(err, errors.NotFoundf("charm or bundle at %q", c.CharmOrBundle))
        } else if err != nil && err != os.ErrNotExist {
                // If we get a "not exists" error then we attempt to interpret
                // the supplied charm reference as a URL elsewhere, otherwise
                // we return the error.
                return nil, errors.Trace(err)
        } else if err != nil {
                logger.Debugf("cannot interpret as local charm: %v", err)
                return nil, nil
        }

        return func(ctx *cmd.Context, apiRoot DeployAPI) error {
                if err := c.validateCharmFlags(); err != nil {
                        return errors.Trace(err)
                }

                if curl, err = apiRoot.AddLocalCharm(curl, ch); err != nil {
                        return errors.Trace(err)
                }

                id := charmstore.CharmID{
                        URL: curl,
                        // Local charms don't need a channel.
                }

                ctx.Infof("Deploying charm %q.", curl.String())
                return errors.Trace(c.deployCharm(
                        id,
                        (*macaroon.Macaroon)(nil), // local charms don't need one.
                        curl.Series,
                        ctx,
                        apiRoot,
                ))
        }, nil
}

===================================================================================================

juju/juju/cmd/juju/application/deploy.go

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

        deploy, err := findDeployerFIFO(        <------
                c.maybeReadLocalBundle,
                func() (deployFn, error) { return c.maybeReadLocalCharm(apiRoot) },      <=========
                c.maybePredeployedLocalCharm,
                c.maybeReadCharmstoreBundleFn(apiRoot),
                c.charmStoreCharm, // This always returns a deployer
        )
        if err != nil {
                return errors.Trace(err)
        }

        return block.ProcessBlockedError(deploy(ctx, apiRoot), block.BlockChange)
}

----------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go:941


type deployFn func(*cmd.Context, DeployAPI) error


func findDeployerFIFO(maybeDeployers ...func() (deployFn, error)) (deployFn, error) {
        for _, d := range maybeDeployers {
                if deploy, err := d(); err != nil {
                        return nil, errors.Trace(err)
                } else if deploy != nil {
                        return deploy, nil
                }
        }
        return nil, errors.NotFoundf("suitable deployer")
}

----------------------------------------------------------------------------------------

./cmd/juju/application/deploy.go

func (c *DeployCommand) maybePredeployedLocalCharm() (deployFn, error) {
        // If the charm's schema is local, we should definitively attempt
        // to deploy a charm that's already deployed in the
        // environment.
        userCharmURL, err := charm.ParseURL(c.CharmOrBundle)
        if err != nil {
                return nil, errors.Trace(err)
        } else if userCharmURL.Schema != "local" {
                logger.Debugf("cannot interpret as a redeployment of a local charm from the controller")
                return nil, nil
        }

        return func(ctx *cmd.Context, api DeployAPI) error {
                if err := c.validateCharmFlags(); err != nil {
                        return errors.Trace(err)
                }
                formattedCharmURL := userCharmURL.String()
                ctx.Infof("Located charm %q.", formattedCharmURL)
                ctx.Infof("Deploying charm %q.", formattedCharmURL)
                return errors.Trace(c.deployCharm(
                        charmstore.CharmID{URL: userCharmURL},
                        (*macaroon.Macaroon)(nil),
                        userCharmURL.Series,
                        ctx,
                        api,
                ))
        }, nil
}


===================================================================================================


juju/juju/cmd/juju/application/deploy.go

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

        deploy, err := findDeployerFIFO(        <------
                c.maybeReadLocalBundle,
                func() (deployFn, error) { return c.maybeReadLocalCharm(apiRoot) },      <=========
                c.maybePredeployedLocalCharm,
                c.maybeReadCharmstoreBundleFn(apiRoot),
                c.charmStoreCharm, // This always returns a deployer
        )
        if err != nil {
                return errors.Trace(err)
        }

        return block.ProcessBlockedError(deploy(ctx, apiRoot), block.BlockChange)
}

----------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go:941


type deployFn func(*cmd.Context, DeployAPI) error


func findDeployerFIFO(maybeDeployers ...func() (deployFn, error)) (deployFn, error) {
        for _, d := range maybeDeployers {
                if deploy, err := d(); err != nil {
                        return nil, errors.Trace(err)
                } else if deploy != nil {
                        return deploy, nil
                }
        }
        return nil, errors.NotFoundf("suitable deployer")
}

----------------------------------------------------------------------------------------

./cmd/juju/application/deploy.go

func (c *DeployCommand) maybeReadCharmstoreBundleFn(apiRoot DeployAPI) func() (deployFn, error) {
        return func() (deployFn, error) {
                userRequestedURL, err := charm.ParseURL(c.CharmOrBundle)
                if err != nil {
                        return nil, errors.Trace(err)
                }

                modelCfg, err := getModelConfig(apiRoot)
                if err != nil {
                        return nil, errors.Trace(err)
                }

                // Charm or bundle has been supplied as a URL so we resolve and
                // deploy using the store.
                storeCharmOrBundleURL, channel, _, err := apiRoot.Resolve(modelCfg, userRequestedURL)
                if charm.IsUnsupportedSeriesError(err) {
                        return nil, errors.Errorf("%v. Use --force to deploy the charm anyway.", err)
                } else if err != nil {
                        return nil, errors.Trace(err)
                } else if storeCharmOrBundleURL.Series != "bundle" {
                        logger.Debugf(
                                `cannot interpret as charmstore bundle: %v (series) != "bundle"`,
                                storeCharmOrBundleURL.Series,
                        )
                        return nil, nil
                }

                if err := c.validateBundleFlags(); err != nil {
                        return nil, errors.Trace(err)
                }

                return func(ctx *cmd.Context, apiRoot DeployAPI) error {
                        bundle, err := apiRoot.GetBundle(storeCharmOrBundleURL)
                        if err != nil {
                                return errors.Trace(err)
                        }
                        ctx.Infof("Located bundle %q", storeCharmOrBundleURL)
                        data := bundle.Data()

                        return errors.Trace(c.deployBundle(
                                ctx,
                                "", // filepath
                                data,
                                channel,
                                apiRoot,
                                c.BundleStorage,
                        ))
                }, nil
        }
}

===================================================================================================

juju/juju/cmd/juju/application/deploy.go

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

        deploy, err := findDeployerFIFO(        <------
                c.maybeReadLocalBundle,
                func() (deployFn, error) { return c.maybeReadLocalCharm(apiRoot) },      <=========
                c.maybePredeployedLocalCharm,
                c.maybeReadCharmstoreBundleFn(apiRoot),
                c.charmStoreCharm, // This always returns a deployer
        )
        if err != nil {
                return errors.Trace(err)
        }

        return block.ProcessBlockedError(deploy(ctx, apiRoot), block.BlockChange)
}

----------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go:941


type deployFn func(*cmd.Context, DeployAPI) error


func findDeployerFIFO(maybeDeployers ...func() (deployFn, error)) (deployFn, error) {
        for _, d := range maybeDeployers {
                if deploy, err := d(); err != nil {
                        return nil, errors.Trace(err)
                } else if deploy != nil {
                        return deploy, nil
                }
        }
        return nil, errors.NotFoundf("suitable deployer")
}

----------------------------------------------------------------------------------------

./cmd/juju/application/deploy.go

func (c *DeployCommand) charmStoreCharm() (deployFn, error) {
        userRequestedURL, err := charm.ParseURL(c.CharmOrBundle)
        if err != nil {
                return nil, errors.Trace(err)
        }

        return func(ctx *cmd.Context, apiRoot DeployAPI) error {
                // resolver.resolve potentially updates the series of anything
                // passed in. Store this for use in seriesSelector.
                userRequestedSeries := userRequestedURL.Series

                modelCfg, err := getModelConfig(apiRoot)
                if err != nil {
                        return errors.Trace(err)
                }

                // Charm or bundle has been supplied as a URL so we resolve and deploy using the store.
                storeCharmOrBundleURL, channel, supportedSeries, err := apiRoot.Resolve(modelCfg, userRequestedURL)
                if charm.IsUnsupportedSeriesError(err) {
                        return errors.Errorf("%v. Use --force to deploy the charm anyway.", err)
                } else if err != nil {
                        return errors.Trace(err)
                }

                if err := c.validateCharmFlags(); err != nil {
                        return errors.Trace(err)
                }

                selector := seriesSelector{
                        charmURLSeries:  userRequestedSeries,
                        seriesFlag:      c.Series,
                        supportedSeries: supportedSeries,
                        force:           c.Force,
                        conf:            modelCfg,
                        fromBundle:      false,
                }

                // Get the series to use.
                series, err := selector.charmSeries()
                if charm.IsUnsupportedSeriesError(err) {
                        return errors.Errorf("%v. Use --force to deploy the charm anyway.", err)
                }

                // Store the charm in the controller
                curl, csMac, err := addCharmFromURL(apiRoot, storeCharmOrBundleURL, channel)
                if err != nil {
                        if termErr, ok := errors.Cause(err).(*common.TermsRequiredError); ok {
                                return errors.Trace(termErr.UserErr())
                        }
                        return errors.Annotatef(err, "storing charm for URL %q", storeCharmOrBundleURL)
                }

                formattedCharmURL := curl.String()
                ctx.Infof("Located charm %q.", formattedCharmURL)
                ctx.Infof("Deploying charm %q.", formattedCharmURL)
                id := charmstore.CharmID{
                        URL:     curl,
                        Channel: channel,
                }
                return errors.Trace(c.deployCharm(
                        id,
                        csMac,
                        series,
                        ctx,
                        apiRoot,
                ))
        }, nil
}

===================================================================================================
