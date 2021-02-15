package mod

import (
  "fmt"
  "os"
}

type Opt struct {
    optType int
    val string
    name string
    descr string
    required bool
}

func OPT( name string, descr string ) (*Opt) {
    return &Opt{
        name: name,
        descr: descr,
    }
}

func OPT_REQ( name string, descr string ) (*Opt) {
    opt := OPT( name, descr )
    opt.required = true
    return opt
}

type Cmd struct {
    opts map[string]*Opt
    handler func( *Cmd )
    argc int
    argv []string
    descr string
}

type Uclop struct {
    cmds map[string]*Cmd
}

type OPTS []*Opt

func (self *Uclop) AddCmd( cmdStr string, descr string, handler func( *Cmd ), opts OPTS ) {
    cmd := NewUcmd( descr, handler )
    self.cmds[ cmdStr ] = cmd
    
    if opts != nil {
        for _,opt := range opts {
            oName := opt.name
            cmd.opts[ oName ] = opt
        }
    }
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
    } else {
        //fmt.Printf("Found cmd %s\n", cmdStr )
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
                fmt.Fprintf( os.Stderr, "Unknown option %s\n", arg )
                os.Exit(1)
            }
            if opt.optType == 2 {
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
            fmt.Printf(" %s val", optName )
        } else {
            if opt.optType == 2 {
                fmt.Printf(" [%s]", optName )
            } else {
                fmt.Printf(" [%s val]", optName )
            }
        }
    }
    //if self.extrahelp {
    //    fmt.Printf(" %s", self.extrahelp )
    //}
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
    fmt.Printf("  %s %s\n", self.name, reqStr )
    fmt.Printf("    %s\n", self.descr );
}

func (self *Uclop) Usage() {
    fmt.Printf("Usage:\n");
    
    prog := "test"//os.Args[0]
    for cmdName, cmd := range self.cmds {
        fmt.Printf("%s %s", prog, cmdName )
        cmd.UsageInline()
        fmt.Printf("  %s\n", cmd.descr )
        cmd.Usage()
    }
}

func (self *Cmd ) ensureRequired() {
    missing := false
    for optName, opt := range self.opts {
        if opt.required && opt.val == "" {
            missing = true
            fmt.Fprintf(os.Stderr,"Required option %s missing\n", optName )
            opt.Usage()
        }
    }
    if missing {
        os.Exit(1)
    }
}