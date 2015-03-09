# cs733

This is a key-value store written in GOlang. It listens on port 9000.

# Features
* Keys can be strings and values can be any binary data.
* The key value pair can be assigned an expiry time after which it may not be available.
* It uses a priority queue internally to keep track of key value pairs to be deleted after expiry.

# Instructions: How to run
1.Obtain a copy of the project using:	go get github.com/swapniel99/keyvalstore
2.Change directory to keyvalstore : cd $GOPATH/src/github.com/swapniel99/keyvalstore
3.Run the test script using following commands: go test

# Protocol specification

This specification describes the on-the-wire protocol, the exchange of commands and their results on a TCP connection.
The client opens up a TCP connection
Each command is a line in this format.The options are explained in the next section. The actual value data is on the next line.

1.  Set: create the key-value pair, or update the value if it already exists.

    set <key> <exptime> <numbytes> [noreply]\r\n
    <value bytes>\r\n

    The server responds with:

    OK <version>\r\n  

    where version is a unique 64-bit number (in decimal format) assosciated with the key.

2.  Get: Given a key, retrieve the corresponding key-value pair

    get <key>\r\n

    The server responds with the following format (or one of the errors described later)

    VALUE <numbytes>\r\n
    <value bytes>\r\n

3.  Get Meta: Retrieve value, version number and expiry time left

     getm <key>\r\n

    The server responds with the following format (or one of the errors described below)

    VALUE <version> <exptime> <numbytes>\r\n
    <value bytes>\r\n

4.  Compare and swap. This replaces the old value (corresponding to key) with the new value only if the version is still the same.

    cas <key> <exptime> <version> <numbytes> [noreply]\r\n
    <value bytes>\r\n

    The server responds with the new version if successful (or one of the errors described late)

      OK <version>\r\n

5.  Delete key-value pair

     delete <key>\r\n

    Server response (if successful)

      DELETED\r\n

Options:

    key : an ascii text string (max 250 bytes) without spaces
    numbytes: size of the value block, not including the trailing \r\n. It is in an ascii text format.
    version: A 64-bit number generated by the server, in ascii text format.
    exptime: An offset in seconds after which the value may not be available. 0 indicates no expiry at all.

Errors that can be returned.

    “ERR_VERSION \r\n” (the value was not changed because of a version mismatch)
    “ERRNOTFOUND\r\n” (the key doesn’t exist)
    “ERRCMDERR\r\n” (the command line is not formatted correctly)
    “ERR_INTERNAL\r\n

Thank you!

