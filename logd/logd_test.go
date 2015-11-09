package logd

import "testing"

func TestDevLog(t *testing.T) {

	// => Level: app.Dev Type: INFO Time: 2015-02-01 12:32:40:4032 Func: TestDevLog  Line:30:322 Message: ...
	Dev.Log(struct{}{}, InfoLevel, "TestDevLog", "Write to bridget")
	Dev.Logf(struct{}{}, InfoLevel, "TestDevLog", "Write to %s", "john")

	// => Level: app.Dev Type: Debug Time: 2015-02-01 12:32:40:4032 Func: TestDevLog  Line:30:322  Message: ...
	Dev.Debug(struct{}{}, "TestDevLog", "Write to sumer")
	Dev.Debugf(struct{}{}, "TestDevLog", "Write to %s", "sally")

	// => Level: app.Dev Type: Error Time: 2015-02-01 12:32:40:4032 Func: TestDevLog  Line:30:322 Message: ...
	Dev.Error(struct{}{}, "TestDevLog", "Error at procer", 100)
	Dev.Errorf(struct{}{}, "TestDevLog", "Error at %d", 100)

	// => Level: app.Dev Type: Info Time: 2015-02-01 12:32:40:4032 Func: TestDevLog  Line:30:322 Message: ...
	Dev.Info(struct{}{}, "TestDevLog", "Write to work")
	Dev.Infof(struct{}{}, "TestDevLog", "Write to %s", "slow")

	//dumps the given data into json formatted output
	// => Level: app.Dev Type: DataDump Time: 2015-02-01 12:32:40:4032 Func: TestDevLog Line:30:322 Message: JSON Request Body...
	/*
		map:
		 name: alex
		 sid: 32
	*/
	Dev.Dump(struct{}{}, "TestDevLog", map[string]interface{}{"name": "alex", "sid": 32}, "JSON Requests Body")
	Dev.Dumpf(struct{}{}, "TestDevLog", map[string]interface{}{"name": "alex", "sid": 32}, "JSON Requests Body %s", "url")
}

func TestUserLog(t *testing.T) {
	// => Level: app.User Type: Time: 2015-02-01 12:32:40:4032 Message: ...
	User.Log(struct{}{}, InfoLevel, "TestUserLog", "Write to bridget")
	User.Logf(struct{}{}, InfoLevel, "TestUserLog", "Write to %s", "john")

	// => Level: app.User Type: Debug Time: 2015-02-01 12:32:40:4032 Message: ...
	User.Debug(struct{}{}, "TestUserLog", "Write to sumer")
	User.Debugf(struct{}{}, "TestUserLog", "Write to %s", "sally")

	// => Level: app.User Type: Error Time: 2015-02-01 12:32:40:4032 Message: ...
	User.Error(struct{}{}, "TestUserLog", "Error at procer", 100)
	User.Errorf(struct{}{}, "TestUserLog", "Error at %d", 100)

	// => Level: app.User Type: Info Time: 2015-02-01 12:32:40:4032 Message: ...
	User.Info(struct{}{}, "TestUserLog", "Write to work")
	User.Infof(struct{}{}, "TestUserLog", "Write to %s", "slow")

	//dumps the given data into json formatted output
	// => Level: app.User Type: DataDump Time: 2015-02-01 12:32:40:4032 Message:
	/*
		map:
			name: alex
			sid: 32
	*/
	User.Dump(struct{}{}, "TestDevLog", map[string]interface{}{"name": "alex", "sid": 32}, "JSON Requests Body")
	User.Dumpf(struct{}{}, "TestDevLog", map[string]interface{}{"name": "alex", "sid": 32}, "JSON Requests Body %s", "url")
}
