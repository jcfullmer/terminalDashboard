# terminalDashboard

Terminal dashboard is a program designed to run as you boot up your desired terminal or CLI program.
When running the program for the first time it will look for a config file and if none is found take you through the initial setup. You can run this setup at any time by running "terminalDashboard --setup"

## Getting Started

- [Install Go](https://go.dev/doc/install)

- Build the binaries:

```
go build .
```

- *(Optional)* Run the program for the first time to go through setup:

```
./terminalDashboard
```

## Configuring Program to run on terminal start

### Linux (GNOME, KDE, XFCE, etc.)

      After creating the binaries it is recommended to move it to /bin 
      Open Home folder in your terminal.
      Open .bashrc in your text editor
      Add:

```
terminalDashboard
```

      to the end of the .bashrc
      Save and exit

### macOS: **Not Tested**

        Open System Preferences > Users & Groups > Login Items.
        Click the + button and add your script or application. You can also create a script and use the Automator to wrap it as an application to add it as a login item.

### Windows: **Not Tested**

        Press the Windows Key + R, then type shell:startup and press Enter to open the current user's Startup folder.
        Place a shortcut to your program or a batch file (.bat or .cmd) in this folder.
        For a program to run with administrator privileges, use the Task Scheduler instead and set a trigger for user login. 
