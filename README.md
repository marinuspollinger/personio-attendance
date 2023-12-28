# Personio Attendance Tool

This tool is used, to make attendance tracking with Personio easier. It is a command line tool, which can be used to track attendance for the current day.

## Installation

Please Note:\
Executables are built for Amd64, Ubuntu 23.04 and only tested for this OS

If you need them on a different system, you can build them from source.

- Copy personio-server-exectuable & personio-client-executable to your bin folder (e.g. /usr/local/bin)
- Remove -executable suffix
- Make sure, that the files are executable (chmod +x personio-server personio-client)
- set following env vars

  ```sh
  PERSONIO_CLIENT_ID            Secret ID
  PERSONIO_CLIENT_SECRET        Secret Value
  PERSONIO_EMPLOYEE_ID          Your Employee ID
  PERSONIO_HOST                 e.g. https://api.personio.de

  # Optional:
  HTTP_ADDRESS                  Default: 0.0.0.0:33333
  LOG_LEVEL                     Default: "debug"
  CURRENT_TIME_LOOP_INTERVAL    Default: 30s
  ```

- create service for own user in `~/.config/systemd/user/personio-server.service`
  - if needed enable `systemd --user` for own user: `sudo loginctl enable-linger $USER`
- enable service for own user `systemctl --user enable personio-server`

## Personio-Cli Usage

local personio-server must be running for these commands to work

---

`personio-cli` supports the following command-line options:

- `--help`  
  Show this help.

- `--status`  
  Get Current Times.

- `--break-end-time <time>`  
  Set Break End Time (format: 2006-01-02 15:04).

- `--break-start-time <time>`  
  Set Break Start Time (format: 2006-01-02 15:04).

- `--start-break`  
  Start the Break now!

- `--end-break`  
  End the Break now!

- `--send`  
  Send Current Times to Personio API.

- `--yes`  
  Immediately write to personio, without checking times first.
