
vi openstack-swift-storage.yaml

swift-storage:
  block-device: /srv/swift.img|5G
  overwrite: "true"


juju deploy --to=0 --config ~/openstack-swift-storage.yaml "/home/ubuntu/CHARMS/charm-swift-storage"
juju add-unit --to=1 swift-storage
juju add-unit --to=2 swift-storage


~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

vi ~/neutron-gateway.yaml

neutron-gateway:
  ext-port: 'eth1'
neutron-api:
  openstack-origin: distro
  neutron-security-groups: True


juju deploy --to=0 --config ~/neutron-gateway.yaml "/home/ubuntu/CHARMS/charm-neutron-gateway"

juju add-relation nova-cloud-controller neutron-gateway
juju add-relation neutron-api neutron-gateway
juju add-relation neutron-gateway:amqp rabbitmq-server:amqp

=============================================================================================

juju/juju/cmd/juju/commands/main.go

        "os/exec"
        "github.com/juju/cmd"


// Main registers subcommands for the juju executable, and hands over control
// to the cmd package. This function is not redundant with main, because it
// provides an entry point for testing with arbitrary command line arguments.
// This function returns the exit code, for main to pass to os.Exit.
func Main(args []string) int {
        return main{
                execCommand: exec.Command,
        }.Run(args)
}


-----------------------------------------------------------------------------------------------
juju/juju/cmd/juju/commands/main.go


// main is a type that captures dependencies for running the main function.
type main struct {
        // execCommand abstracts away exec.Command.
        execCommand func(command string, args ...string) *exec.Cmd
}

-----------------------------------------------------------------------------------------------
juju/juju/cmd/juju/commands/main.go

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

        jcmd := NewJujuCommand(ctx)             <=============   |---> invokes: registerCommands(jcmd, ctx) | returns: cmd.Command interface and SuperCommand struct
        return cmd.Main(jcmd, ctx, args[1:])    <-------------
}

-----------------------------------------------------------------------------------------------
juju/cmd/cmd.go

// DefaultContext returns a Context suitable for use in non-hosted situations.
func DefaultContext() (*Context, error) {
        dir, err := os.Getwd()
        if err != nil {
                return nil, err
        }
        abs, err := filepath.Abs(dir)
        if err != nil {
                return nil, err
        }
        return &Context{
                Dir:    abs,
                Stdin:  os.Stdin,
                Stdout: os.Stdout,
                Stderr: os.Stderr,
        }, nil
}


-----------------------------------------------------------------------------------------------
juju/cmdo/cmd.go:110

// Context represents the run context of a Command. Command implementations
// should interpret file names relative to Dir (see AbsPath below), and print
// output and errors to Stdout and Stderr respectively.
type Context struct {
        Dir     string
        Env     map[string]string
        Stdin   io.Reader
        Stdout  io.Writer
        Stderr  io.Writer
        quiet   bool
        verbose bool
}

=============================================================================================

juju/cmd.go:305


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
        if err := c.Run(ctx); err != nil {              <============
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

juju/juju/cmd/juju/commands/main.go:389

        rcmd "github.com/juju/juju/cmd/juju/romulus/commands"


// registerCommands registers commands in the specified registry.
func registerCommands(r commandRegistry, ctx *cmd.Context) {

        // Manage and control services

        r.Register(application.NewDeployCommand())          <===========

...

        rcmd.RegisterAll(r)
}


-----------------------------------------------------------------------------------------------

juju/juju/cmd/supercommand.go:37

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

-----------------------------------------------------------------------------------------------

juju/cmd/supercommand.go

func (c *SuperCommand) insert(value commandReference) {
        if _, found := c.subcmds[value.name]; found {
                panic(fmt.Sprintf("command already registered: %q", value.name))
        }
        c.subcmds[value.name] = value
}


-----------------------------------------------------------------------------------------------

github.com/juju/cmd/cmd.go:63

type Command interface {

        IsSuperCommand() bool

        Info() *Info

        SetFlags(f *gnuflag.FlagSet)

        Init(args []string) error

        Run(ctx *Context) error

        AllowInterspersedFlags() bool
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
juju/cmd/supercommand.go:108

type commandReference struct {
        name    string
        command Command
        alias   string
        check   DeprecationCheck
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

juju/juju/cmd/juju/commands/main.go

        "github.com/juju/cmd"

type commandRegistry interface {
        Register(cmd.Command)           
        RegisterSuperAlias(name, super, forName string, check cmd.DeprecationCheck)
        RegisterDeprecated(subcmd cmd.Command, check cmd.DeprecationCheck)
}

-----------------------------------------------------------------------------------------------

./cmd/juju/romulus/commands/commands.go:28

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

github.com/juju/cmd/cmd.go:63:type Command interface {


type Command interface {

        IsSuperCommand() bool

        Info() *Info

        SetFlags(f *gnuflag.FlagSet)

        Init(args []string) error

        Run(ctx *Context) error

        AllowInterspersedFlags() bool
}

-----------------------------------------------------------------------------------------------

github.com/juju/cmd/supercommand.go:185:func (c *SuperCommand) Register(subcmd Command) {

func (c *SuperCommand) Register(subcmd Command) {
        info := subcmd.Info()
        c.insert(commandReference{name: info.Name, command: subcmd})
        for _, name := range info.Aliases {
                c.insert(commandReference{name: name, command: subcmd, alias: info.Name})
        }
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

-----------------------------------------------------------------------------------------------

github.com/juju/cmd/supercommand.go

func (c *SuperCommand) insert(value commandReference) {
        if _, found := c.subcmds[value.name]; found {
                panic(fmt.Sprintf("command already registered: %q", value.name))
        }
        c.subcmds[value.name] = value
}

-----------------------------------------------------------------------------------------------

github.com/juju/cmd/supercommand.go:429:func (c *SuperCommand) Run(ctx *Context) error {

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
                c.notifyRun(name)          <=============
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

-----------------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go:813:func (c *DeployCommand) Run(ctx *cmd.Context) error {


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



===============================================================================================

juju/juju/cmd/juju/application/deploy.go

        "github.com/juju/errors"

        "gopkg.in/juju/charmrepo.v2/csclient"

        "github.com/juju/juju/api/application"
        apicharms "github.com/juju/juju/api/charms"
        "github.com/juju/juju/api/modelconfig"
        "github.com/juju/juju/api/annotations"
        "gopkg.in/juju/charmrepo.v2"

        "github.com/juju/juju/cmd/modelcmd"


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
                apiRoot, err := deployCmd.ModelCommandBase.NewAPIRoot()   =>  type Connection interface
                if err != nil {
                        return nil, errors.Trace(err)
                }
                bakeryClient, err := deployCmd.BakeryClient()
                if err != nil {
                        return nil, errors.Trace(err)
                }
                cstoreClient := newCharmStoreClient(bakeryClient).WithChannel(deployCmd.Channel)   =>  *csclient.Client  WithChannel(

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

        return modelcmd.Wrap(deployCmd)         <=========
}


-----------------------------------------------------------------------------------------------

./cmd/modelcmd/modelcommand.go:337

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
                Command: WrapBase(wrapper),          <======== 
                embed:   embed{wrapper},
        }
}

-----------------------------------------------------------------------------------------------

./cmd/modelcmd/base.go:370

// WrapBase wraps the specified Command. This should be
// used by any command that embeds CommandBase.
func WrapBase(c Command) Command {
        return &baseCommandWrapper{
                Command: c,
        }
}


-----------------------------------------------------------------------------------------------


juju/juju/cmd/modelcmd/base.go


// Run implements Command.Run.
func (w *baseCommandWrapper) Run(ctx *cmd.Context) error {
        defer w.closeAPIContexts()
        w.initContexts(ctx)
        w.setRunStarted()
        return w.Command.Run(ctx)
}


-----------------------------------------------------------------------------------------------


juju/juju/cmd/modelcmd/base.go


type baseCommandWrapper struct {
        Command
}



-----------------------------------------------------------------------------------------------


juju/juju/cmd/modelcmd/base.go

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


juju/cmd/cmd.go:63:type Command interface {

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


==================================================================================================================


juju/juju/cmd/modelcmd/base.go

        "github.com/juju/cmd"

// CommandBase is a convenience type for embedding that need
// an API connection.
type CommandBase struct {
        cmd.CommandBase
        cmdContext    *cmd.Context
        apiContexts   map[string]*apiContext
        modelAPI_     ModelAPI
        apiOpenFunc   api.OpenFunc
        authOpts      AuthOpts
        runStarted    bool
        refreshModels func(jujuclient.ClientStore, string) error
}

-----------------------------------------------------------------------------------------------

juju/juju/cmd/modelcmd/modelcommand.go

// ModelCommandBase is a convenience type for embedding in commands
// that wish to implement ModelCommand.
type ModelCommandBase struct {
        CommandBase

        // store is the client controller store that contains information
        // about controllers, models, etc.
        store jujuclient.ClientStore

        // _modelName and _controllerName hold the current
        // model and controller names. They are only valid
        // after initModel is called, and should in general
        // not be accessed directly, but through ModelName and
        // ControllerName respectively.
        _modelName      string
        _controllerName string

        allowDefaultModel bool

        // doneInitModel holds whether initModel has been called.
        doneInitModel bool

        // initModelError holds the result of the initModel call.
        initModelError error
}

-----------------------------------------------------------------------------------------------

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

=============================================================================================

./cmd/juju/application/deploy.go

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

=============================================================================================


juju/juju/cmd/juju/application/deploy.go

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
                return errors.Trace(c.deployCharm(                <===============
                        id,
                        (*macaroon.Macaroon)(nil), // local charms don't need one.
                        curl.Series,
                        ctx,
                        apiRoot,
                ))
        }, nil
}


------------------------------------------------------------------------------------------------


juju/juju/cmd/juju/application/deploy.go

func (c *DeployCommand) deployCharm(
        id charmstore.CharmID,
        csMac *macaroon.Macaroon,
        series string,
        ctx *cmd.Context,
        apiRoot DeployAPI,
) (rErr error) {
        charmInfo, err := apiRoot.CharmInfo(id.URL.String())
        if err != nil {
                return err
        }

        if len(c.AttachStorage) > 0 && apiRoot.BestFacadeVersion("Application") < 5 {
                // DeployArgs.AttachStorage is only supported from
                // Application API version 5 and onwards.
                return errors.New("this juju controller does not support --attach-storage")
        }

        numUnits := c.NumUnits
        if charmInfo.Meta.Subordinate {
                if !constraints.IsEmpty(&c.Constraints) {
                        return errors.New("cannot use --constraints with subordinate application")
                }
                if numUnits == 1 && c.PlacementSpec == "" {
                        numUnits = 0
                } else {
                        return errors.New("cannot use --num-units or --to with subordinate application")
                }
        }
        applicationName := c.ApplicationName
        if applicationName == "" {
                applicationName = charmInfo.Meta.Name
        }

        // Process the --config args.
        // We may have a single file arg specified, in which case
        // it points to a YAML file keyed on the charm name and
        // containing values for any charm settings.
        // We may also have key/value pairs representing
        // charm settings which overrides anything in the YAML file.
        // If more than one file is specified, that is an error.
        var configYAML []byte
        files, err := c.ConfigOptions.AbsoluteFileNames(ctx)
        if err != nil {
                return errors.Trace(err)
        }
        if len(files) > 1 {
                return errors.Errorf("only a single config YAML file can be specified, got %d", len(files))
        }
        if len(files) == 1 {
                configYAML, err = ioutil.ReadFile(files[0])
                if err != nil {
                        return errors.Trace(err)
                }
        }

        attr, err := c.ConfigOptions.ReadConfigPairs(ctx)
        if err != nil {
                return errors.Trace(err)
        }
        appConfig := make(map[string]string)
        for k, v := range attr {
                appConfig[k] = v.(string)
        }

        // Application facade V5 expects charm config to either all be in YAML
        // or config map. If config map is specified, that overrides YAML.
        // So we need to combine the two here to have only one.
        if apiRoot.BestFacadeVersion("Application") < 6 && len(appConfig) > 0 {
                var configFromFile map[string]map[string]string
                err := yaml.Unmarshal(configYAML, &configFromFile)
                if err != nil {
                        return errors.Annotate(err, "badly formatted YAML config file")
                }
                if configFromFile == nil {
                        configFromFile = make(map[string]map[string]string)
                }
                charmSettings, ok := configFromFile[applicationName]
                if !ok {
                        charmSettings = make(map[string]string)
                }
                for k, v := range appConfig {
                        charmSettings[k] = v
                }
                appConfig = nil
                configFromFile[applicationName] = charmSettings
                configYAML, err = yaml.Marshal(configFromFile)
                if err != nil {
                        return errors.Trace(err)
                }
        }

        bakeryClient, err := c.BakeryClient()
        if err != nil {
                return errors.Trace(err)
        }

        uuid, ok := apiRoot.ModelUUID()
        if !ok {
                return errors.New("API connection is controller-only (should never happen)")
        }

        deployInfo := DeploymentInfo{
                CharmID:         id,
                ApplicationName: applicationName,
                ModelUUID:       uuid,
                CharmInfo:       charmInfo,
        }

        for _, step := range c.Steps {
                err = step.RunPre(apiRoot, bakeryClient, ctx, deployInfo)
                if err != nil {
                        return errors.Trace(err)
                }
        }

        defer func() {
                for _, step := range c.Steps {
                        err = errors.Trace(step.RunPost(apiRoot, bakeryClient, ctx, deployInfo, rErr))
                        if err != nil {
                                rErr = err
                        }
                }
        }()

        if id.URL != nil && id.URL.Schema != "local" && len(charmInfo.Meta.Terms) > 0 {
                ctx.Infof("Deployment under prior agreement to terms: %s",
                        strings.Join(charmInfo.Meta.Terms, " "))
        }

        ids, err := resourceadapters.DeployResources(
                applicationName,
                id,
                csMac,
                c.Resources,
                charmInfo.Meta.Resources,
                apiRoot,
        )
        if err != nil {
                return errors.Trace(err)
        }

        if len(appConfig) == 0 {
                appConfig = nil
        }
        args := application.DeployArgs{
                CharmID:          id,
                Cons:             c.Constraints,
                ApplicationName:  applicationName,
                Series:           series,
                NumUnits:         numUnits,
                ConfigYAML:       string(configYAML),
                Config:           appConfig,
                Placement:        c.Placement,
                Storage:          c.Storage,
                AttachStorage:    c.AttachStorage,
                Resources:        ids,
                EndpointBindings: c.Bindings,
        }
        if len(appConfig) > 0 {
                args.Config = appConfig
        }
        return errors.Trace(apiRoot.Deploy(args))
}




