# Project Details

## Scenario

Mission is to create an HTTP API that can store and retrieve files from Memcached. To complete this mission, we will
build a library that will accept a file (somewhere between 0 and 50 MB) as an input and store it in Memcached. The
library must also be able to retrieve the file from Memcached and return it. HTTP API will utilize the library to store
and retrieve the objects.

#### Assumptions / Limitations / Future work

- App will not allow user custom keys.
- App will support the following content types:
    - plain/text
    - application/octet-stream
    - multipart/form-data
        - App will allow one file per request.
- App will support the following functions:
    - Store the object
    - Retrieve the object by key
    - Delete the object by key (optional)
- App will use single memcached instance with default settings.
    - In memcached 1.4.2 and higher, the maximum supported object size can be configured by using the -I command-line
      option. For example, to increase the maximum object size to 5 MB: `$ memcached -I 5m`
- App will not implement authentication and authorization.
- App logs will be sent to stdout.
- Library will reject 0 size content.
- Library will not type-cast objects.
- Library will remain agnostic to data-type of payload.
- Library will not set expiration time that means the stored items have no expiration time.
- Library will not compress content.
- Library will keep most recent stored content.
- Library will not retry on failure.
- Library allows resetting in order to naturally fixing of partially evicted content.

## Deliverables

There are two deliverables for this project:

- [x] A library to store and retrieve files in Memcached.
- [x] An HTTP API that utilizes the library to store and retrieve files.

## Specs

### Library

- [x] Library should be small and self-contained.
- [x] Library should utilize a Memcached client, as well as any other libraries required.
- [x] Library must accept any file size from 0 to 50 MB. It must reject files larger than 50 MB.
- [x] Using the default slab size, Memcached can only store up to 1 MB per key. Library must accept a file, chunk it,
  and store as bytes in Memcached with a minimum amount of overhead.
    - Library should chunk the file in any way appropriate.
    - Library should key the chunks in any way appropriate.
- [x] Library must retrieve a file's chunks from Memcached and return a single stream of bytes.
- [x] Library must check for file consistency to ensure the data retrieved is the same as the original data stored.
- [x] Library must handle edge cases appropriately by raising an exception or similar when:
    - [x] Trying to retrieve a file that does not exist.
    - [x] A file retrieved is inconsistent/corrupt.
    - [ ] Storing a file that already exists.
- [x] Memcached can evict keys when it runs out of memory. Library should detect these cases and handle them
  appropriately.
- [x] Library must have at least one test.

### API

- [x] Application should be a REST API.
- [x] Application must accept a POST request with file contents in the payload and store it using the library. It may be
  convenient to return an identifier used for retrieval at a later time.
- [x] Application must accept a GET request with a file name/identifier and retrieve it using the library. The file
  contents must be returned to the caller in the response.
- [x] Application should appropriately handle edge cases (return an error response) when:
    - [x] A key does not exist.
    - [x] Content is not consistent.
- [x] Application must have at least one test.
