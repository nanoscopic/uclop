#include "uclop.h"

void run_default( ucmd *cmd );
void run_test( ucmd *cmd );

int main( int argc, char *argv[] ) {
    uopt *default_options[] = {
        UOPT_REQUIRED("-x","X"),
        UOPT("-y","Y"),
        NULL
    };
    uclop *opts = uclop__new( &run_default, default_options );
    //uclop *opts = uclop__new( NULL, NULL );
    
    uclop__addcmd( opts, "test", "Test Cmd", &run_test, NULL );
    uclop__run( opts, argc, argv );
}

void run_default( ucmd *cmd ) {
    char *x = ucmd__get( cmd, "-x" );
    printf("x=%s\n",x);
}

void run_test( ucmd *cmd ) {
    printf("test\n");
}