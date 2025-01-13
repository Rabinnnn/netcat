# TCP Chat

A TCP-based chat server implementation in Go that mimics the basic functionality of NetCat (nc). This project implements a group chat system with multiple clients connecting to a central server.

## Features

- TCP server supporting multiple client connections (1:many relationship)
- Maximum of 10 concurrent client connections
- User name requirement for each client
- Real-time message broadcasting to all connected clients
- Timestamp and username display for each message
- New clients receive chat history upon joining
- Server notifications for client join/leave events
- Empty message filtering
- Default port 8989 if none specified

## Requirements

- Go 1.20 or higher

## Installation

```bash
# Clone the repository
git clone https://learn.zone01kisumu.ke/git/rotieno/net-cat.git
cd net-cat

# Build the project
go build -o TCPChat
```

## Usage

```bash
[USAGE]: ./TCPChat $port
```

### Starting the Server

```bash
# Start with default port (8989)
go run .

# Start with custom port
go run . 2525
```

### Connecting as a Client

Use NetCat to connect to the server:

```bash
nc localhost <port>
```

You will be greeted with:
```
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
```
### Exiting the Chat

```
To exit the chat group use Ctrl+C
```

## Message Format

Messages in the chat are formatted as:
```
[YYYY-MM-DD HH:MM:SS][username]:message
```

## Example Interaction

```
[2024-01-10 16:03:43][Alice]:hello
[2024-01-10 16:03:46][Alice]:How are you?
Bob has joined our chat...
[2024-01-10 16:04:32][Bob]:Hi everyone!
```

## Technical Details

- Written in Go
- Uses goroutines for concurrent client handling
- Implements mutex-based synchronization for shared resources
- Uses channels for communication between goroutines
- Error handling for both server and client sides


## Testing

The project includes unit tests for both server and client functionality. Run tests with:

```bash
go test -v ./utils
```

## Error Handling

- Invalid port numbers
- Connection failures
- Client disconnections
- Maximum connection limit
- Invalid usernames

## License

[MIT License](LICENSE)

## Contributing

Feel free to contribute!

## Authors

[Rabin Otieno](https://learn.zone01kisumu.ke/git/rotieno)

[Stephen Kisengese](https://learn.zone01kisumu.ke/git/skisenge)