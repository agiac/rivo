// Package rivo is a library for stream processing. It provides a simple and flexible way to create and compose streams of data.
//
// There are three main types in this library: Item, Stream, and Pipeline.
//
// Item is a struct which contains a value and an optional error. Just like errors are returned next to the result
// of a function in synchronous code, they should be passed along into asynchronous code and handled where more appropriate.
//
// Stream is a read only channel of items. As the name suggests, it represents a stream of data.
//
// Pipeline is a function that takes a context.Context and a Stream of one type and returns a Stream of the same or a different type.
// Pipeables can be composed together using the one of the Pipe functions.
// Pipeables are divided in three categories: generators, sinks and transformers.
//   - Generator is a pipeable that does not read from its input stream. It starts a new stream from scratch.
//   - Sync is a pipeable function that does not emit any items. It is used at the end of a pipeline.
//   - Pipeline is a pipeable that reads from its input stream and emits items to its output stream.
package rivo
