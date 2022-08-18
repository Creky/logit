// Copyright 2022 FishGoddess. All Rights Reserved.
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

package writer

import (
	"io"
	"os"

	"github.com/go-logit/logit/core"
)

// Flusher is an interface that flushes data to somewhere.
type Flusher interface {
	Flush() (n int, err error)
}

// Writer is an interface which have flush, write and close functions.
type Writer interface {
	io.WriteCloser // WriteCloser is an interface that writes data to somewhere and can be closed.
	Flusher        // Flusher is an interface that flushes data to somewhere.
}

// notStdoutAndStderr returns true if w isn't stdout and stderr.
func notStdoutAndStderr(w io.Writer) bool {
	return w != os.Stdout && w != os.Stderr
}

// Wrap wraps io.writer to Writer.
func Wrap(writer io.Writer) Writer {
	if w, ok := writer.(Writer); ok {
		return w
	}
	return newWrapWriter(writer)
}

// BufferWithSize wraps io.writer with buffer Writer of bufferSize.
func BufferWithSize(writer io.Writer, bufferSize core.ByteSize) Writer {
	if bw, ok := writer.(*bufferWriter); ok {
		return bw
	}
	return newBufferWriter(writer, bufferSize)
}

// Buffer wraps writer to Buffer Writer with default buffer size.
func Buffer(writer io.Writer) Writer {
	return BufferWithSize(writer, core.WriterBufferSize)
}
