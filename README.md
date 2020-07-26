# MicroModels in Spiff

## Introduction

This Dockerfile implements a base container for MicroModels used in the Spiff backend.
The container does the following things:

- It runs a grpc server on port 8181.
- It includes a a ready-to-use health check to allow for zero-downtime updates.
- It dynamically loads the Go plugin /app/service.so, which must be provided by the container that you build.
- It runs a [grpcui](https://github.com/fullstorydev/grpcui) on port 8080.


## How to get started.

I recommend that you clone the following repo as a starting point:

https://github.com/barebaric/spiff-mm-demo.git

It includes a ready to use Docker file, and everything you need to build
a MicroModel.

You can then define your MicroModel options in proto/options.proto, and
implement your service in service/service.go.

## General MicroModel information

### What is a Spiff MicroModel?

Spiff MicroModels are Docker containers with a standardized API that model
a financial instrument and generate a time series.
The standard API allows them to be used as add-ins in Spiff.

For example, a "Fixed Term Deposit" model would take parameters like the interest
rate, the start and end of the deposit, whether interest accumulates or is payed
out, etc. From this data, the model generates a series of payments that could,
for example, be put into a graph to be visualized.

### Assets vs. income

It is important to distinguish between asset growth and generated income.
For example, a fixed term deposit may not generate any income if interest isn't
paid out - it just accumulates in value (assets).

Therefore, each MicroModel needs different APIs to request either the value of
the asset, or the generated income.

### Performance considerations

Since we are passing a large amount of times series data between MicroModels,
it make sense to use a more efficient serialization mechanism than HTTP/JSON.
We therefore settled to using grpc+protobuf.


## The MicroModel API

For the exact API, please check [proto/micromodel.proto](proto/micromodel.proto)
and [proto/options.proto](proto/options.proto).

### GetUI (v2 only, due to native hard coded frontend in v1)

Returns the user interface (form) for the app in HTML/Javascript
The Javascript must contain a serializeForm() function that can be used by
the frontend to get the form data.

It is frontend responsibility to
- render the HTML
- provide a .css file and style the form elements
- fetch the input data from the WebView by calling the serializeForm() function included in the response of the backend.
- store the result in a storage backend.
