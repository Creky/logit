// Copyright 2021 Ye Zi Jie. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/02/29 15:41:09

/*
Package logit provides an easy way to use foundation for your logging operations.

1. basic:

// Create a new logger for use.
	// Default level is debug, so all logs will be logged.
	// Invoke Close() isn't necessary in all situations.
	// If logger's writer has buffer or something like that, it's better to invoke Close() for flushing buffer or something else.
	logger := logit.NewLogger()
	//defer logger.Close()

	// Then, you can log anything you want.
	// Remember, logs will be ignored if their level is smaller than logger's level.
	// End() will do some finishing work, so this invocation is necessary.
	logger.Debug("This is a debug message").End()
	logger.Info("This is a info message").End()
	logger.Warn("This is a warn message").End()
	logger.Error("This is a error message").End()
	logger.Error("This is a %s message, with format", "error").End() // Format with params.

	// As you know, we provide some levels: debug, info, warn, error, off.
	// The lowest is debug and the highest is off.
	// If you want to change the level of your logger, do it at creating.
	logger = logit.NewLogger(logit.Options().WithWarnLevel())
	logger.Debug("This is a debug message, but ignored").End()
	logger.Info("This is a info message, but ignored").End()
	logger.Warn("This is a warn message, not ignored").End()
	logger.Error("This is a error message, not ignored").End()

	// Also, we provide some "old school" log method :)
	// (Don't mistake~ I love old school~)
	logger.Printf("This is a log %s, and it's for compatibility", "printed")
	logger.Print("This is a log printed, and it's for compatibility", 123)
	logger.Println("This is a log printed, and it's for compatibility", 666)

	// If you want to log with some fields, try this:
	logger.Error("This is a structured message").Error("err", io.EOF).Int("trace", 123).End()

	// You may notice logit.Options() which returns an options list.
	// Here is some of them:
	options := logit.Options()
	options.WithCaller()                          // Let logs carry caller information.
	options.WithLevelKey("lvl")                   // Change logger's level key to "lvl".
	options.WithWriter(os.Stderr, true)           // Change logger's writer to os.Stderr with buffer.
	options.WithErrorWriter(os.Stderr, false)     // Change logger's error writer to os.Stderr without buffer.
	options.WithTimeFormat("2006-01-02 15:04:05") // Change the format of time (Only the log's time will apply it).

	// You can bind context with logger and use it as long as you can get the context.
	ctx := logit.NewContext(context.Background(), logger)
	logger = logit.FromContext(ctx)

2. options:

	// We provide some options for you.
	options := logit.Options()
	options.WithDebugLevel()
	options.WithInfoLevel()
	options.WithWarnLevel()
	options.WithErrorLevel()
	options.WithAppender(appender.Text())
	options.WithDebugAppender(appender.Text())
	options.WithInfoAppender(appender.Text())
	options.WithWarnAppender(appender.Text())
	options.WithErrorAppender(appender.Text())
	options.WithWriter(os.Stderr, false)
	options.WithDebugWriter(os.Stderr, false)
	options.WithInfoWriter(os.Stderr, false)
	options.WithWarnWriter(os.Stderr, false)
	options.WithErrorWriter(os.Stderr, false)
	options.WithPid()
	options.WithCaller()
	options.WithMsgKey("msg")
	options.WithTimeKey("time")
	options.WithLevelKey("level")
	options.WithPidKey("pid")
	options.WithFileKey("file")
	options.WithLineKey("line")
	options.WithTimeFormat(appender.UnixTime) // UnixTime means time will be logged as unix time, an int64 number.
	options.WithCallerDepth(3)                // Set caller depth to 3 so the log will get the third depth caller.

	// Remember, these options is only used for creating a logger.
	logger := logit.NewLogger(
		options.WithPid(),
		options.WithWriter(os.Stdout, false),
		options.WithTimeFormat("2006/01/02 15:04:05"),
		options.WithCaller(),
		options.WithCallerDepth(3),
		// ...
	)
	defer logger.Close()
	logger.Info("check options").End()

	// You can use many options at the same time, but some of them is exclusive.
	// So only the last one in order will take effect if you use them at the same time.
	logit.NewLogger(
		options.WithDebugLevel(),
		options.WithInfoLevel(),
		options.WithWarnLevel(),
		options.WithErrorLevel(), // The level of logger is error.
	)

	// You can customize an option for your logger.
	// Actually, Option is just a function like func(logger *Logger).
	// So you can do what you want in creating a logger.
	autoFlushOption := func(logger *logit.Logger) {
		go func() {
			select {
			case <-time.Tick(time.Second):
				logger.Flush()
			}
		}()
	}

	logit.NewLogger(autoFlushOption)

3. appender:

	// We provide some ways to change the form of logs.
	// Actually, appender is an interface with some common methods, see appender.Appender.
	appender.Text()
	appender.Json()

	// Set appender to the one you want to use when creating a logger.
	// Default appender is appender.Text().
	logger := logit.NewLogger()
	logger.Info("appender.Text()").End()

	// You can switch appender to the other one, such appender.Json().
	logger = logit.NewLogger(logit.Options().WithAppender(appender.Json()))
	logger.Info("appender.Json()").End()

	// Every level has its own appender so you can append logs in different level with different appender.
	logger = logit.NewLogger(
		logit.Options().WithDebugAppender(appender.Text()),
		logit.Options().WithInfoAppender(appender.Text()),
		logit.Options().WithWarnAppender(appender.Json()),
		logit.Options().WithErrorAppender(appender.Json()),
	)

	// Appender is an interface so you can implement your own appender.
	// However, we don't recommend you to do that.
	// This interface may change in every version, so you will pay lots of extra attention to it.
	// So you should implement it only if you really need to do.

4. writer:

	// As you know, writer in logit is customized, not io.Writer.
	// The reason why we create a new Writer interface is we want a flushable writer.
	// Then, we notice a flushable writer also need a close method to flush all data in buffer when closing.
	// So, a new Writer is born:
	//
	//     type Writer interface {
	//	       Flusher
	//	       io.WriteCloser
	//     }
	//
	// In package writer, we provide some writers for you.
	writer.Wrapped(os.Stdout)  // Wrap io.Writer to writer.Writer
	writer.Buffered(os.Stderr) // Wrap io.Writer to writer.Writer with buffer, which needs invoking Flush() or Close()

	// Use the writer without buffer.
	logger := logit.NewLogger(logit.Options().WithWriter(os.Stdout, false))
	logger.Info("WriterWithoutBuffer").End()

	// Use the writer with buffer, which is good for io.
	logger = logit.NewLogger(logit.Options().WithWriter(os.Stdout, true))
	defer logger.Close() // Flush data and close writer

	logger.Info("WriterWithBuffer").End()
	logger.Flush() // Remember flushing data or flushing by Close().

	// Every level has its own appender so you can append logs in different level with different appender.
	logger = logit.NewLogger(
		logit.Options().WithDebugWriter(os.Stdout, true),
		logit.Options().WithInfoWriter(os.Stdout, true),
		logit.Options().WithWarnWriter(os.Stdout, false),
		logit.Options().WithErrorWriter(os.Stdout, false),
	)

5. global:

	// There are some global settings for optimizations, and you can set all of them in need.
	//
	//     import "github.com/FishGoddess/logit/core"
	//
	// All global settings are stored in package core.

	// 1. LogMallocSize (The pre-malloc size of a new Log data)
	// If your logs are extremely long, such as 4000 bytes, you can set it to 4096 to avoid re-malloc.
	core.LogMallocSize = 4096 // 4096 Bytes

	// 2. WriterBufferedSize (The default size of buffered writer)
	// If your logs are extremely long, such as 16KB, you can set it to 2048 to avoid re-malloc.
	core.WriterBufferedSize = 32 * writer.KB

	// After setting global settings, just use Logger as normal.
	logger := logit.NewLogger()
	defer logger.Close()

	logger.Info("set global settings").Int("LogMallocSize", core.LogMallocSize).Int("WriterBufferedSize", core.WriterBufferedSize).End()

6. context:

	// By NewContext, you can bind a context with a logger and get it from context again.
	// So you can use this logger from everywhere as long as you can get this context.
	ctx := logit.NewContext(context.Background(), logit.NewLogger())

	// FromContext returns the logger in context.
	logger := logit.FromContext(ctx)
	logger.Info("This is a message logged by logger from context").End()

	// Actually, you also have a chance to specify the key of logger in context.
	// It gives you a way to discriminate different businesses in using logger.
	// For example, you can create two loggers for your two different usages and
	// set them to a context with different key, so you can get each logger from context with each key.
	businessOneKey := "businessOne"
	logger = logit.NewLogger(logit.Options().WithMsgKey("businessOneMsg"))
	ctx = logit.NewContextWithKey(context.Background(), businessOneKey, logger)

	businessTwoKey := "businessTwo"
	logger = logit.NewLogger(logit.Options().WithMsgKey("businessTwoMsg"))
	ctx = logit.NewContextWithKey(ctx, businessTwoKey, logger)

	// Get different logger from the same context with different key.
	logger = logit.FromContextWithKey(ctx, businessOneKey)
	logger.Info("This is a message logged by logger from context with businessOneKey").End()

	logger = logit.FromContextWithKey(ctx, businessTwoKey)
	logger.Info("This is a message logged by logger from context with businessTwoKey").End()

7. maker:

	type testLoggerMaker struct{}

	func (tlm *testLoggerMaker) MakeLogger(ctx context.Context, params ...interface{}) (*logit.Logger, error) {
		if len(params) < 1 {
			return nil, errors.New("testLoggerMaker: len(params) < 1")
		}

		if params[0].(string) == "error" {
			return nil, errors.New("testLoggerMaker: params[0] isn't a string")
		}

		// Customize your creation of logger here.
		return logit.NewLogger(), nil
	}

	makeName := "testLoggerMaker"

	// RegisterLoggerMaker registers maker to logit with given name.
	err := logit.RegisterLoggerMaker(makeName, new(testLoggerMaker))
	if err != nil {
		panic(err)
	}

	// NewLoggerFromMaker creates logger from maker with given params.
	// Panic will be invoked if params is "error" because MakeLogger in testLoggerMaker has this logic.
	logger, err := logit.NewLoggerFromMaker(context.Background(), makeName, "xxx")
	if err != nil {
		panic(err)
	}

	logger.Info("I am made from logger maker!").End()

8. caller:

	// Let's create a logger without caller information.
	logger := logit.NewLogger()
	logger.Info("I am without caller").End()

	// We provide a way to add caller information to log even logger doesn't carry caller.
	logger.Info("Invoke log.WithCaller()").WithCaller().End()
	logger.Close()

	time.Sleep(time.Second)

	// Now, let's create a logger with caller information.
	logger = logit.NewLogger(logit.Options().WithCaller())
	logger.Info("I am with caller").End()

	// We won't carry caller information twice or more if logger carries caller information originally.
	logger.Info("Invoke log.WithCaller() again").WithCaller().End()
	logger.Close()
*/
package logit // import "github.com/FishGoddess/logit"

const (
	// Version is the version string representation of logit.
	Version = "v0.4.15"
)
