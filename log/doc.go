// Package log provides an open-source simple logging api providing two log loglevels
// When using, first Initialize the log if you wish to change the device to use with the
// logger has its set to use os.Stdout as default and also pass in a function which returns
// the desired initial log level i.e USER or DEV.
//
//   Init(os.Stderr,func() int {
//        return USER
//   })
//
// When logging in the alternate log Level ensure to switch the log levels to
// that specific level else all logs of that level would be ignored, to do so
// simple call the SwitchLevel function giving the specified LogLevel constant
//
//    SwitchLevel(USER)
//  or
//   SwitchLevel(DEV)
//
// Logging is done using the provided logging functions User and Dev, each loggimg
// to their respective levels when active
//
//  User("32","BuildTower","Initializing game tower build process for %d",43)
//
//  Dev("32","BuildTower","Intializing game tower build process failed due to %s","No Cash")
//
package log
