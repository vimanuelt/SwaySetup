# SwaySetup

**SwaySetup** is a terminal-based interactive utility for setting up **Sway** and **Wayland** on FreeBSD systems. The application uses the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework for its text-based user interface and simplifies the installation and configuration process for Sway.

---

## Features

SwaySetup provides the following setup actions:

1. **Install Packages**  
   Installs Sway, Wayland, and required dependencies.

2. **Configure Sway**  
   Sets up the default Sway configuration file in the user's home directory.

3. **Setup `seatd`**  
   Enables and starts the `seatd` service for seat management.

4. **Set Environment Variables**  
   Adds necessary environment variables for Wayland compatibility.

---

## Prerequisites

Ensure your FreeBSD system meets the following requirements:

- **FreeBSD Version**: 13 or higher
- **Packages**:
   - `pkg` (Package manager)
   - `sway`, `seatd`, and other Wayland dependencies
- **Root Privileges**:
   - Required for installing packages and configuring system services.

---

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/vimanuelt/SwaySetup.git
cd SwaySetup
```

### 2. Build the Application

Ensure you have **Go 1.18+** installed.

```bash
go build -o SwaySetup
```

---

## Usage

Run the SwaySetup utility with root privileges:

```bash
sudo ./SwaySetup
```

### Navigation

- **↑ / k**: Move Up
- **↓ / j**: Move Down
- **Enter**: Select Action
- **q**: Quit Application

---

## Actions

1. **Install Packages**  
   Installs required Sway and Wayland packages:

   ```bash
   pkg install -y sway swaylock swayidle seatd wayland
   ```

2. **Configure Sway**  
   Copies the default Sway configuration to `~/.config/sway/config` if it doesn't exist.

3. **Setup seatd**  
   Enables and starts the `seatd` service:

   ```bash
   sysrc seatd_enable=YES
   service seatd start
   ```

4. **Set Environment Variables**  
   Appends the following environment variable to `~/.profile`:

   ```bash
   export XDG_RUNTIME_DIR=/tmp/xdg-runtime-$(id -u)
   ```

---

## Screenshots

### Main Menu
```
Sway & Wayland Setup for FreeBSD

Choose an action:
  - Install Packages: Install Sway, Wayland, and dependencies
  - Configure Sway: Set up initial Sway configuration
  - Setup seatd: Enable and start seatd service
  - Set Environment: Set necessary environment variables

Press ↑/↓ to navigate, Enter to select, q to quit.
```

---

## License

This project is licensed under the [BSD 3-Clause License](LICENSE).

---

## Acknowledgements

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework for Go
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling for terminal output

---

## Contributions

Contributions are welcome!  
Feel free to open an issue or submit a pull request.

---

## Author

**Your Name**  
GitHub: [vimanuelt](https://github.com/vimanuelt)

---

## Support

If you encounter any issues or need help, please open an issue on the repository.

