package my_errors

import "errors"

var ErrParseJSON = errors.New("error decoding json")
var ErrNoSlink = errors.New("error searching for short link in database: short link not found")
var ErrNoLlink = errors.New("error searching for long link in database: long link not found")
var ErrEqualJSON = errors.New("error read information in JSON format: empty")
var ErrWriteJSONerr = errors.New("error writing response to JSON")
var ErrWriteJSON = errors.New("error writing error to JSON")
var ErrMethodNotAllowed = errors.New("this method is not provided")
var ErrParseDuration = errors.New("error when parsing links lifetime")
var ErrNoRedirect = errors.New("error in trying to get the redirect")
var ErrNoDataDeath = errors.New("error in trying to get the data_death")
