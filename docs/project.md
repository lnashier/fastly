# Project Details

## Scenario

Mission is to create an HTTP API that can store and retrieve files from Memcached. To complete this mission, we will
build a library that will accept a file (somewhere between 0 and 50 MB) as an input and store it in Memcached. The
library must also be able to retrieve the file from Memcached and return it. HTTP API will utilize the library to store
and retrieve the objects.

#### Notes

- Using the default slab size, Memcached can only store up to 1 MB per key. That means we'll have to implement some
  means of chunking the file to store it in Memcached.
    - In memcached 1.4.2 and higher, the maximum supported object size can be configured by using the -I command-line
      option. For example, to increase the maximum object size to 5 MB: `$ memcached -I 5m`
- Additionally, Memcached can evict keys when it runs out of memory. A complete solution should detect these cases and
  handle them appropriately.

#### Assumptions / Trade-Offs

- App will use memcached with default settings.
- App will not implement authentication and authorization.
- App will do simply loggin to stdout.
- Library will not type-cast objects.
- Library will remain agnostic to data-type of payload.
    - Library takes slice/array of bytes and returns slice/array of bytes.
- Library will not set expiration time means the stored items have no expiration time.
- Library will only perform basic sanity checks as specified in the requirements:
    - Payload size limit (0 - 50 MB]

## Deliverables

There are two deliverables for this project:

- [x] A library to store and retrieve files in Memcached.
- [x] An HTTP API that utilizes the library to store and retrieve files.

## Specs

### Library

- [x] Library should be small and self-contained.
- [x] Library should utilize a Memcached client, as well as any other libraries required.
- [ ] Library must accept any file size from 0 to 50 MB. It must reject files larger than 50 MB.
- [ ] Library must accept a file, chunk it, and store as bytes in Memcached with a minimum amount of overhead.
- [ ] Library must retrieve a file's chunks from Memcached and return a single stream of bytes.
- [ ] Library should chunk the file in any way appropriate.
- [ ] Library should key the chunks in any way appropriate.
- [ ] Library must check for file consistency to ensure the data retrieved is the same as the original data stored.
- [ ] Library must handle edge cases appropriately by raising an Exception or similar. Some examples of edge cases may
  include storing a file that already exists, trying to retrieve a file that does not exist, or when a file retrieved is
  inconsistent/corrupt.
- [x] Library must have at least one test.

### API

- [ ] Application should be a REST API.
- [ ] Application must accept a POST request with file contents in the payload and store it using the library. It may be
  convenient to return an identifier used for retrieval at a later time.
- [ ] Application must accept a GET request with a file name/identifier and retrieve it using the library. The file
  contents must be returned in the response.
- [ ] Application should appropriately handle edge cases (return an error response) when a file does not exist or is not
  consistent.
- [ ] Application must have at least one test.
