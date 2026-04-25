# cashew
My very own Redis inspired caching service

This is my take of [Build Your Own Redis Server](https://codingchallenges.fyi/challenges/challenge-redis)

## Roadmap

| Goal                                                                      | Status                |
| ------------------------------------------------------------------------- | --------------------- |
| Serialise RESP (Simple Strings, Errors, Integers, Bulk Strings, Arrays)   | :white_check_mark:    |
| De-serialise RESP                                                         | :white_check_mark:    |
| A server that listens (`PING`, `ECHO`)                                    | :white_check_mark:    |
| Set and Get keys                                                          | :white_check_mark:    |
| Concurrent access                                                         | :white_check_mark:    |
| Implement expiry options                                                  | :white_check_mark:    |
| Implement `EXISTS`, `DEL`, `INCR`, `DECR`, `LPUSH`, `RPUSH`               | :white_check_mark:    |
