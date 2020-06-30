#include "headers/heclib7.h"
#include "headers/hecdss7.h"

// Example from:  https://www.hec.usace.army.mil/confluence/dsscprogrammer/dss-progammers-guide-for-c/accessing-hec-dss-files-open-and-close
// To Compile: gcc hello_dss.c -L. -ljavaheclib  -o hello_dss_C

int ExampleOpen()
{
    long long ifltab[250];
    int status;

    status = zopen(ifltab, "data/G14.dss");
    if (status != STATUS_OKAY)
        return status;

    zclose(ifltab);

    return 0;
}

int main()
{

    ExampleOpen();

    return 0;
}