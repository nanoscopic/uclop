package lib

import (
  "fmt"
  "os"
  "strconv"
)

// flags for OPT option creation
const (
    REQ = 1
    FLAG = 2
)

const (
    TYPE_FLAG = iota + 1
)

type Opt struct {
    optType  int
    val      string
    name     string
    descr    string
    required bool
}

type Cmd struct {
    opts      map[string]*Opt
    handler   func( *Cmd )
    argc      int
    argv      []string
    descr     string
    Extrahelp string
}

type Uclop struct {
    cmds map[string]*Cmd
}

type OPTS []*Opt

func (self *Opt) String() string {
    return self.val
}

func (self *Opt) Int() int {
    num, _ := strconv.Atoi( self.val )
    return num 
}

func (self *Opt) Bool() bool {
    if self.val == "true" {
        return true
    }
    return false
}

func OPT( name string, descr string, flags int ) (*Opt) {
    opt := &Opt{
        name:     name,
        descr:    descr,
        required: ( flags & REQ )==REQ,
    }
    if (flags & FLAG)==FLAG {
        opt.optType = TYPE_FLAG
    }
    return opt
}

func (self *Uclop) AddCmd( cmdStr string, descr string, handler func( *Cmd ), opts OPTS ) *Cmd {
    cmd := NewUcmd( descr, handler )
    self.cmds[ cmdStr ] = cmd
    
    if opts != nil {
        for _,opt := range opts {
            oName := opt.name
            cmd.opts[ oName ] = opt
        }
    }
    return cmd
}

func NewUcmd( descr string, handler func(*Cmd) ) ( *Cmd ) {
    return &Cmd{
        opts: make( map[string]*Opt ),
        handler: handler,
        descr: descr,
    }
}

func NewUclop() ( *Uclop ) {
    return &Uclop{
        cmds: make( map[string]*Cmd ),
    }
}

func (self *Cmd) Get(optName string) *Opt {
    return self.opts[ optName ]
}

func (self *Uclop) Run() {
    argv := os.Args[1:]
    argc := len( argv )
    
    // Pull cmd name out of args
    cmdStr := "default"
    
    start := 0
    if argc > 0 && argv[0][0] != '-' {
        start++
        cmdStr = argv[0]
    }
    
    cmd := self.cmds[ cmdStr ]
    
    if cmd == nil {
        self.Usage()
        return
    }
    
    // Store arguments
    i := start
    for {
        if i >= argc {
            break
        }
        arg := argv[i]
        if arg[0] == '-' {
            opt := cmd.opts[ arg ]
            if opt == nil {
                fmt.Printf( "Unknown option %s\n", arg )
                os.Exit(1)
            }
            if opt.optType == TYPE_FLAG {
                opt.val = "true"
            } else {
                i++
                if i >= argc {
                    break
                }
                opt.val = argv[i]
            }
            i++
            if i >= argc {
                break
            }
        } else {
            break
        }
    }
    cmd.argc = argc - i + 1
    cmd.argv = argv[i:]
    
    cmd.ensureRequired()
  
    if cmd.handler == nil {
        self.Usage()
    } else {
        cmd.Run()
    }
}

func (self *Cmd) Run() {
    self.handler( self )
}

func runTest( cmd *Cmd ) {
    opt := cmd.Get("-opt")
    fmt.Printf("Ran test - opt=%s\n", opt.val)
}

func (self *Cmd) UsageInline() {
    for optName, opt := range self.opts {
        if opt.required {
            fmt.Printf(" %s val", optName)
        } else {
            if opt.optType == TYPE_FLAG {
                fmt.Printf(" [%s]", optName)
            } else {
                fmt.Printf(" [%s val]", optName)
            }
        }
    }
    if self.Extrahelp != "" {
        fmt.Printf(" %s", self.Extrahelp )
    }
    fmt.Printf("\n")
}

func (self *Cmd) Usage() {
    for _, opt := range self.opts {
        opt.Usage()
    }
}

func (self *Opt) Usage() {
    reqStr := ""
    if self.required {
        reqStr = "REQUIRED"
    }
    fmt.Printf("  %s %s\n", self.name, reqStr)
    fmt.Printf("    %s\n", self.descr)
}

func (self *Uclop) Usage() {
    fmt.Printf("Usage:\n")
    
    // Get program name, stripping off path
    prog := os.Args[0]
    j := 0
    for i:=0 ; i<len(prog) ; i++ { if prog[i]=='/' { j = i } }
    if j>0 { prog = prog[j+1:] }
    
    for cmdName, cmd := range self.cmds {
        fmt.Printf( prog + " " + cmdName )
        cmd.UsageInline()
        fmt.Printf("  %s\n", cmd.descr)
        cmd.Usage()
    }
}

func (self *Cmd ) ensureRequired() {
    missing := false
    for optName, opt := range self.opts {
        if opt.required && opt.val == "" {
            missing = true
            fmt.Printf("Required option %s missing\n", optName)
            opt.Usage()
        }
    }
    if missing {
        os.Exit(1)
    }
}
