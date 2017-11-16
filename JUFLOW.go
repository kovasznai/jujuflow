


cat >> ~/maas.yaml << EOF
clouds:/
  maas:
    type: maas
    auth-types: [oauth1]
    endpoint: http://192.168.100.3/MAAS/
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
        Run(ctx *Context) error                                                <------------

        // AllowInterspersedFlags returns whether the command allows flag
        // arguments to be interspersed with non-flag arguments.
        AllowInterspersedFlags() bool
}


====================================================================================

vi ./juju/juju/cmd/juju/cloud/add.go

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



juju/juju/cmd/juju/cloud/add.go

// Run executes the add cloud command, adding a cloud based on a passed-in yaml
// file or interactive queries.
func (c *AddCloudCommand) Run(ctxt *cmd.Context) error {
        if c.CloudFile == "" {
                return c.runInteractive(ctxt)
        }

        specifiedClouds, err := c.cloudMetadataStore.ParseCloudMetadataFile(c.CloudFile)          <========== 1
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

====================================================================================
__1__

        specifiedClouds, err := c.cloudMetadataStore.ParseCloudMetadataFile(c.CloudFile)


vi ./juju/juju/cmd/juju/cloud/add.go


type CloudMetadataStore interface {
        ParseCloudMetadataFile(path string) (map[string]cloud.Cloud, error)
        ParseOneCloud(data []byte) (cloud.Cloud, error)
        PublicCloudMetadata(searchPaths ...string) (result map[string]cloud.Cloud, fallbackUsed bool, _ error)
        PersonalCloudMetadata() (map[string]cloud.Cloud, error)
        WritePersonalCloudMetadata(cloudsMap map[string]cloud.Cloud) error
}



vi ./juju/juju/cloud/personalclouds.go

import (
        "io/ioutil"

        "github.com/juju/juju/juju/osenv"


// ParseCloudMetadataFile loads any cloud metadata defined
// in the specified file.
func ParseCloudMetadataFile(file string) (map[string]Cloud, error) {
        data, err := ioutil.ReadFile(file)
        if err != nil {
                return nil, err
        }
        clouds, err := ParseCloudMetadata(data)         <------
        if err != nil {
                return nil, err
        }
        return clouds, err
}

--------------------------------------------------------------

cat >> ~/maas.yaml << EOF
clouds:/
  maas:
    type: maas
    auth-types: [oauth1]
    endpoint: http://192.168.100.3/MAAS/
EOF

--------------------------------------------------------------


./cloud/clouds.go

        "gopkg.in/yaml.v2"



var defaultCloudDescription = map[string]string{
        "aws":         "Amazon Web Services",
        "aws-china":   "Amazon China",
        "aws-gov":     "Amazon (USA Government)",
        "google":      "Google Cloud Platform",
        "azure":       "Microsoft Azure",
        "azure-china": "Microsoft Azure China",
        "rackspace":   "Rackspace Cloud",
        "joyent":      "Joyent Cloud",
        "cloudsigma":  "CloudSigma Cloud",
        "lxd":         "LXD Container Hypervisor",
        "maas":        "Metal As A Service",
        "openstack":   "Openstack Cloud",
        "oracle":      "Oracle Compute Cloud Service",
}


// ParseCloudMetadata parses the given yaml bytes into Clouds metadata.
func ParseCloudMetadata(data []byte) (map[string]Cloud, error) {
        var metadata cloudSet
        if err := yaml.Unmarshal(data, &metadata); err != nil {
                return nil, errors.Annotate(err, "cannot unmarshal yaml cloud metadata")
        }

        // Translate to the exported type. For each cloud, we store
        // the first region for the cloud as its default region.
        clouds := make(map[string]Cloud)
        for name, cloud := range metadata.Clouds {
                details := cloudFromInternal(cloud)
                details.Name = name
                if details.Description == "" {
                        var ok bool
                        if details.Description, ok = defaultCloudDescription[name]; !ok {
                                details.Description = defaultCloudDescription[cloud.Type]
                        }
                }
                clouds[name] = details
        }
        return clouds, nil
}




func cloudFromInternal(in *cloud) Cloud {
        var regions []Region
        if len(in.Regions.Map) > 0 {
                for _, item := range in.Regions.Slice {
                        name := fmt.Sprint(item.Key)
                        r := in.Regions.Map[name]
                        if r == nil {
                                // r will be nil if none of the fields in
                                // the YAML are set.
                                regions = append(regions, Region{Name: name})
                        } else {
                                regions = append(regions, Region{
                                        name,
                                        r.Endpoint,
                                        r.IdentityEndpoint,
                                        r.StorageEndpoint,
                                })
                        }
                }
        }
        meta := Cloud{
                Name:             in.Name,
                Type:             in.Type,
                AuthTypes:        in.AuthTypes,
                Endpoint:         in.Endpoint,
                IdentityEndpoint: in.IdentityEndpoint,
                StorageEndpoint:  in.StorageEndpoint,
                Regions:          regions,
                Config:           in.Config,
                RegionConfig:     in.RegionConfig,
                Description:      in.Description,
        }
        meta.denormaliseMetadata()
        return meta
}




meta

{
  "Name": "",
  "Type": "openstack",
  "Description": "",
  "AuthTypes": [
    "access-key",
    "userpass"
  ],
  "Endpoint": "",
  "IdentityEndpoint": "",
  "StorageEndpoint": "",
  "Regions": [
    {
      "Name": "reg1",
      "Endpoint": "https://openstack.example.com:35574/v3.0/",
      "IdentityEndpoint": "https://graph.windows.net",
      "StorageEndpoint": "https://core.windows.net"
    },
    {
      "Name": "reg2",
      "Endpoint": "https://openstack.example.com:35574/v3.0/",
      "IdentityEndpoint": "https://graph.windows.net",
      "StorageEndpoint": "https://core.windows.net"
    }
  ],
  "Config": null,
  "RegionConfig": null
}





import (
        "reflect"
        "strings"
)





func init() {
        RegisterStructTags(Cloud{}, Region{})
}

// RegisterStructTags ensures the yaml tags for the given structs are able to be used
// when parsing cloud metadata.
func RegisterStructTags(vals ...interface{}) {
        tags := mkTags(vals...)
        for k, v := range tags {
                tagsForType[k] = v
        }
}

func mkTags(vals ...interface{}) map[reflect.Type]map[string]int {
        typeMap := make(map[reflect.Type]map[string]int)
        for _, v := range vals {
                t := reflect.TypeOf(v)
                typeMap[t] = yamlTags(t)
        }
        return typeMap
}

// yamlTags returns a map from yaml tag to the field index for the string fields in the given type.
func yamlTags(t reflect.Type) map[string]int {
        if t.Kind() != reflect.Struct {
                panic(errors.Errorf("cannot get yaml tags on type %s", t))
        }
        tags := make(map[string]int)
        for i := 0; i < t.NumField(); i++ {
                f := t.Field(i)
                if f.Type != reflect.TypeOf("") {
                        continue
                }
                if tag := f.Tag.Get("yaml"); tag != "" {
                        if i := strings.Index(tag, ","); i >= 0 {
                                tag = tag[0:i]
                        }
                        if tag == "-" {
                                continue
                        }
                        if tag != "" {
                                f.Name = tag
                        }
                }
                tags[f.Name] = i
        }
        return tags
}










func (cloud Cloud) denormaliseMetadata() {
        for name, region := range cloud.Regions {
                r := region
                inherit(&r, &cloud)
                cloud.Regions[name] = r
        }
}



// inherit sets any blank fields in dst to their equivalent values in fields in src that have matching json tags.
// The dst parameter must be a pointer to a struct.
func inherit(dst, src interface{}) {
        for tag := range tags(dst) {
                setFieldByTag(dst, tag, fieldByTag(src, tag), false)
        }
}



type structTags map[reflect.Type]map[string]int

var tagsForType structTags = make(structTags)



// tags returns the field offsets for the JSON tags defined by the given value, which must be
// a struct or a pointer to a struct.
func tags(x interface{}) map[string]int {
        t := reflect.TypeOf(x)
        if t.Kind() == reflect.Ptr {
                t = t.Elem()
        }
        if t.Kind() != reflect.Struct {
                panic(errors.Errorf("expected struct, not %s", t))
        }

        if tagm := tagsForType[t]; tagm != nil {
                return tagm
        }
        panic(errors.Errorf("%s not found in type table", t))
}



// fieldByTag returns the value for the field in x with the given JSON tag, or "" if there is no such field.
func fieldByTag(x interface{}, tag string) string {
        tagm := tags(x)
        v := reflect.ValueOf(x)
        if v.Kind() == reflect.Ptr {
                v = v.Elem()
        }
        if i, ok := tagm[tag]; ok {
                return v.Field(i).Interface().(string)
        }
        return ""
}



// setFieldByTag sets the value for the field in x with the given JSON tag to val.
// The override parameter specifies whether the value will be set even if the original value is non-empty.
func setFieldByTag(x interface{}, tag, val string, override bool) {
        i, ok := tags(x)[tag]
        if !ok {
                return
        }
        v := reflect.ValueOf(x).Elem()
        f := v.Field(i)
        if override || f.Interface().(string) == "" {
                f.Set(reflect.ValueOf(val))
        }
}




// inherit sets any blank fields in dst to their equivalent values in fields in src that have matching json tags.
// The dst parameter must be a pointer to a struct.
func inherit(dst, src interface{}) {
        for tag := range tags(dst) {
                setFieldByTag(dst, tag, fieldByTag(src, tag), false)
        }
}



func (cloud Cloud) denormaliseMetadata() {
        for name, region := range cloud.Regions {
                r := region
                inherit(&r, &cloud)
                cloud.Regions[name] = r
        }
}


--------------------------------------------------------------

../../../gopkg.in/yaml.v2/yaml.go

func Unmarshal(in []byte, out interface{}) (err error) {
        return unmarshal(in, out, false)
}


--------------------------------------------------------------

./cloud/clouds.go

// cloudSet contains cloud definitions, used for marshalling and
// unmarshalling.
type cloudSet struct {
        // Clouds is a map of cloud definitions, keyed on cloud name.
        Clouds map[string]*cloud `yaml:"clouds"`
}


--------------------------------------------------------------

./cloud/clouds.go

// RegionConfig holds a map of regions and the attributes that serve as the
// region specific configuration options. This allows model inheritance to
// function, providing a place to store configuration for a specific region
// which is  passed down to other models under the same controller.
type RegionConfig map[string]Attrs


// Attrs serves as a map to hold regions specific configuration attributes.
// This serves to reduce confusion over having a nested map, i.e.
// map[string]map[string]interface{}
type Attrs map[string]interface{}



// AuthType is the type of authentication used by the cloud.
type AuthType string

// AuthTypes is defined to allow sorting AuthType slices.
type AuthTypes []AuthType


--------------------------------------------------------------

./cloud/clouds.go

// cloud is equivalent to Cloud, for marshalling and unmarshalling.
type cloud struct {
        Name             string                 `yaml:"name,omitempty"`
        Type             string                 `yaml:"type"`
        Description      string                 `yaml:"description,omitempty"`
        AuthTypes        []AuthType             `yaml:"auth-types,omitempty,flow"`
        Endpoint         string                 `yaml:"endpoint,omitempty"`
        IdentityEndpoint string                 `yaml:"identity-endpoint,omitempty"`
        StorageEndpoint  string                 `yaml:"storage-endpoint,omitempty"`
        Regions          regions                `yaml:"regions,omitempty"`
        Config           map[string]interface{} `yaml:"config,omitempty"`
        RegionConfig     RegionConfig           `yaml:"region-config,omitempty"`
}


// regions is a collection of regions, either as a map and/or
// as a yaml.MapSlice.
//
// When marshalling, we populate the Slice field only. This is
// necessary for us to control the order of map items.
//
// When unmarshalling, we populate both Map and Slice. Map is
// populated to simplify conversion to Region objects. Slice
// is populated so we can identify the first map item, which
// becomes the default region for the cloud.
type regions struct {
        Map   map[string]*region
        Slice yaml.MapSlice
}

// region is equivalent to Region, for marshalling and unmarshalling.
type region struct {
        Endpoint         string `yaml:"endpoint,omitempty"`
        IdentityEndpoint string `yaml:"identity-endpoint,omitempty"`
        StorageEndpoint  string `yaml:"storage-endpoint,omitempty"`
}


--------------------------------------------------------------

./cloud/clouds.go

// Cloud is a cloud definition.
type Cloud struct {
        // Name of the cloud.
        Name string

        // Type is the type of cloud, eg ec2, openstack etc.
        // This is one of the provider names registered with
        // environs.RegisterProvider.
        Type string

        // Description describes the type of cloud.
        Description string

        // AuthTypes are the authentication modes supported by the cloud.
        AuthTypes AuthTypes

        // Endpoint is the default endpoint for the cloud regions, may be
        // overridden by a region.
        Endpoint string

        // IdentityEndpoint is the default identity endpoint for the cloud
        // regions, may be overridden by a region.
        IdentityEndpoint string

        // StorageEndpoint is the default storage endpoint for the cloud
        // regions, may be overridden by a region.
        StorageEndpoint string

        // Regions are the regions available in the cloud.
        //
        // Regions is a slice, and not a map, because order is important.
        // The first region in the slice is the default region for the
        // cloud.
        Regions []Region

        // Config contains optional cloud-specific configuration to use
        // when bootstrapping Juju in this cloud. The cloud configuration
        // will be combined with Juju-generated, and user-supplied values;
        // user-supplied values taking precedence.
        Config map[string]interface{}

        // RegionConfig contains optional region specific configuration.
        // Like Config above, this will be combined with Juju-generated and user
        // supplied values; with user supplied values taking precedence.
        RegionConfig RegionConfig
}

--------------------------------------------------------------

./cloud/clouds.go

// Region is a cloud region.
type Region struct {
        // Name is the name of the region.
        Name string

        // Endpoint is the region's primary endpoint URL.
        Endpoint string

        // IdentityEndpoint is the region's identity endpoint URL.
        // If the cloud/region does not have an identity-specific
        // endpoint URL, this will be empty.
        IdentityEndpoint string

        // StorageEndpoint is the region's storage endpoint URL.
        // If the cloud/region does not have a storage-specific
        // endpoint URL, this will be empty.
        StorageEndpoint string
}


--------------------------------------------------------------


clouds:
  <cloud_name>:
    type: <type_of_cloud>
    auth-types: <[access-key, oauth, userpass]>
    regions:
      <region-name>:
        endpoint: <https://xxx.yyy.zzz:35574/v3.0/>


--------------------------------------------------------------


#####################################################################################################
__2__

juju/juju/cloud/clouds.go

// Cloud is a cloud definition.
type Cloud struct {
        // Name of the cloud.
        Name string

        // Type is the type of cloud, eg ec2, openstack etc.
        // This is one of the provider names registered with
        // environs.RegisterProvider.
        Type string

        // Description describes the type of cloud.
        Description string

        // AuthTypes are the authentication modes supported by the cloud.
        AuthTypes AuthTypes

        // Endpoint is the default endpoint for the cloud regions, may be
        // overridden by a region.
        Endpoint string

        // IdentityEndpoint is the default identity endpoint for the cloud
        // regions, may be overridden by a region.
        IdentityEndpoint string

        // StorageEndpoint is the default storage endpoint for the cloud
        // regions, may be overridden by a region.
        StorageEndpoint string

        // Regions are the regions available in the cloud.
        //
        // Regions is a slice, and not a map, because order is important.
        // The first region in the slice is the default region for the
        // cloud.
        Regions []Region

        // Config contains optional cloud-specific configuration to use
        // when bootstrapping Juju in this cloud. The cloud configuration
        // will be combined with Juju-generated, and user-supplied values;
        // user-supplied values taking precedence.
        Config map[string]interface{}

        // RegionConfig contains optional region specific configuration.
        // Like Config above, this will be combined with Juju-generated and user
        // supplied values; with user supplied values taking precedence.
        RegionConfig RegionConfig
}





juju/juju/cmd/juju/cloud/add.go


type CloudMetadataStore interface {
        ParseCloudMetadataFile(path string) (map[string]cloud.Cloud, error)
        ParseOneCloud(data []byte) (cloud.Cloud, error)
        PublicCloudMetadata(searchPaths ...string) (result map[string]cloud.Cloud, fallbackUsed bool, _ error)
        PersonalCloudMetadata() (map[string]cloud.Cloud, error)
        WritePersonalCloudMetadata(cloudsMap map[string]cloud.Cloud) error
}


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



// Run executes the add cloud command, adding a cloud based on a passed-in yaml
// file or interactive queries.
func (c *AddCloudCommand) Run(ctxt *cmd.Context) error {
        if c.CloudFile == "" {
                return c.runInteractive(ctxt)
        }

        specifiedClouds, err := c.cloudMetadataStore.ParseCloudMetadataFile(c.CloudFile)          <========== 1
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

        return addCloud(c.cloudMetadataStore, newCloud)                                            <========== 2
}



func addCloud(cloudMetadataStore CloudMetadataStore, newCloud cloud.Cloud) error {
        personalClouds, err := cloudMetadataStore.PersonalCloudMetadata()
        if err != nil {
                return err
        }
        if personalClouds == nil {
                personalClouds = make(map[string]cloud.Cloud)
        }
        personalClouds[newCloud.Name] = newCloud
        return cloudMetadataStore.WritePersonalCloudMetadata(personalClouds)
}


=======================================================================================

src/github.com/juju/juju/cloud/personalclouds.go


import (
        "github.com/juju/juju/juju/osenv"
)



// PersonalCloudMetadata loads any personal cloud metadata defined
// in the Juju Home directory. If not cloud metadata is found,
// that is not an error; nil is returned.
func PersonalCloudMetadata() (map[string]Cloud, error) {
        clouds, err := ParseCloudMetadataFile(JujuPersonalCloudsPath())
        if err != nil && os.IsNotExist(err) {
                return nil, nil
        }
        return clouds, err
}



// ParseCloudMetadataFile loads any cloud metadata defined
// in the specified file.
func ParseCloudMetadataFile(file string) (map[string]Cloud, error) {
        data, err := ioutil.ReadFile(file)
        if err != nil {
                return nil, err
        }
        clouds, err := ParseCloudMetadata(data)
        if err != nil {
                return nil, err
        }
        return clouds, err
}




// JujuPersonalCloudsPath is the location where personal cloud information is
// expected to be found. Requires JUJU_HOME to be set.
func JujuPersonalCloudsPath() string {
        return osenv.JujuXDGDataHomePath("clouds.yaml")
}





./juju/osenv/home.go


import (
        "os"
        "path/filepath"
        "runtime"
        "sync"

        "github.com/juju/utils"
)

var (
        jujuXDGDataHomeMu sync.Mutex
        jujuXDGDataHome   string
)


// JujuXDGDataHomePath returns the path to a file in the
// current juju home.
func JujuXDGDataHomePath(names ...string) string {
        all := append([]string{JujuXDGDataHomeDir()}, names...)
        return filepath.Join(all...)
}

s
// JujuXDGDataHome returns the current juju home.
func JujuXDGDataHome() string {
        jujuXDGDataHomeMu.Lock()
        defer jujuXDGDataHomeMu.Unlock()
        return jujuXDGDataHome
}



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
\// jujuXDGDataHomeLinux returns the directory where juju should store application-specific files on Linux.
func jujuXDGDataHomeLinux() string {
        xdgConfig := os.Getenv(XDGDataHome)
        if xdgConfig != "" {
                return filepath.Join(xdgConfig, "juju")
        }
        // If xdg config home is not defined, the standard indicates that its default value
        // is $HOME/.local/share
        home := utils.Home()
        return filepath.Join(home, ".local", "share", "juju")
}




------------------------------------------------------------------



cloud/personalclouds.go

// WritePersonalCloudMetadata marshals to YAMl and writes the cloud metadata
// to the personal cloud file.
func WritePersonalCloudMetadata(cloudsMap map[string]Cloud) error {
        data, err := marshalCloudMetadata(cloudsMap)
        if err != nil {
                return errors.Trace(err)
        }
        return ioutil.WriteFile(JujuPersonalCloudsPath(), data, os.FileMode(0600))
}

------------------------------------------------------------------


cloud/clouds.go

// marshalCloudMetadata marshals the given clouds to YAML.
func marshalCloudMetadata(cloudsMap map[string]Cloud) ([]byte, error) {
        clouds := cloudSet{make(map[string]*cloud)}
        for name, metadata := range cloudsMap {
                clouds.Clouds[name] = cloudToInternal(metadata, false)
        }
        data, err := yaml.Marshal(clouds)
        if err != nil {
                return nil, errors.Annotate(err, "cannot marshal cloud metadata")
        }
        return data, nil
}


// cloudSet contains cloud definitions, used for marshalling and
// unmarshalling.
type cloudSet struct {
        // Clouds is a map of cloud definitions, keyed on cloud name.
        Clouds map[string]*cloud `yaml:"clouds"`
}

// cloud is equivalent to Cloud, for marshalling and unmarshalling.
type cloud struct {
        Name             string                 `yaml:"name,omitempty"`
        Type             string                 `yaml:"type"`
        Description      string                 `yaml:"description,omitempty"`
        AuthTypes        []AuthType             `yaml:"auth-types,omitempty,flow"`
        Endpoint         string                 `yaml:"endpoint,omitempty"`
        IdentityEndpoint string                 `yaml:"identity-endpoint,omitempty"`
        StorageEndpoint  string                 `yaml:"storage-endpoint,omitempty"`
        Regions          regions                `yaml:"regions,omitempty"`
        Config           map[string]interface{} `yaml:"config,omitempty"`
        RegionConfig     RegionConfig           `yaml:"region-config,omitempty"`
}

func cloudToInternal(in Cloud, withName bool) *cloud {
        var regions regions
        for _, r := range in.Regions {
                regions.Slice = append(regions.Slice, yaml.MapItem{
                        r.Name, region{
                                r.Endpoint,
                                r.IdentityEndpoint,
                                r.StorageEndpoint,
                        },
                })
        }
        name := in.Name
        if !withName {
                name = ""
        }
        return &cloud{
                Name:             name,
                Type:             in.Type,
                AuthTypes:        in.AuthTypes,
                Endpoint:         in.Endpoint,
                IdentityEndpoint: in.IdentityEndpoint,
                StorageEndpoint:  in.StorageEndpoint,
                Regions:          regions,
                Config:           in.Config,
                RegionConfig:     in.RegionConfig,
        }
}



