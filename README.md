# 🐶 Watchdog

![Build Status](https://img.shields.io/badge/build-passing-brightgreen) ![Version](https://img.shields.io/badge/version-alpha%200.0.1-blue) ![License](https://img.shields.io/badge/license-MIT-yellowgreen) ![GitHub Stars](https://img.shields.io/github/stars/MarzbanOP/Watchdog?style=social) ![GitHub Forks](https://img.shields.io/github/forks/MarzbanOP/Watchdog?style=social)

Welcome to **Watchdog**, your go-to tool for monitoring and managing proxy usage effectively! Built with Go, Watchdog helps you limit IP connections and block suspicious activity using `iptables` in **Marzban**. With easy JSON configuration and a user-friendly API, managing users and automating controls has never been simpler!

## ✨ Features You'll Love

- **Proxy Management** 🌐: Keep an eye on and manage proxy usage seamlessly.
- **Connection Limits** 🔒: Control the number of simultaneous IP connections to your proxies.
- **Smart Banning** 🚫: Automatically block and unblock IP addresses based on their activity.
- **User Banning** ⛔: Temporarily ban users when necessary, based on specific conditions.
- **Activity Logging** 📊: Track activities and events for better insights.
- **Telegram Notifications** 📲: Get real-time notifications on key events through Telegram.
- **Flexible Storage Options** 💾: Choose your preferred storage method—Redis, SQLite, or JSON.
- **Easy Configuration** ⚙️: Use simple JSON files to set things up without hassle.
- **User Management API** 🛠️: Enjoy a range of API endpoints to manage users easily.

## 🛠️ Requirements

Before diving in, make sure you have the following installed:

- **Docker** and **Docker Compose**

## 🚀 Installation & Usage

Getting started with Watchdog is super easy! Just run the command below to install it:

```bash
sudo bash -c "$(curl -sL https://raw.githubusercontent.com/MarzbanOP/Watchdog/refs/heads/main/watchdog.sh)" @ install
```

Once installed, navigate to the project directory and access the menu:

```bash
cd watchdog
chmod +x watchdog.sh
./watchdog.sh
```

### 📋 Menu Options

When you run the script, you’ll see a menu with the following options:

1. **Install** 🚀: Set up everything by cloning the repository and configuring Docker containers. You’ll choose your storage option (Redis, SQLite, or JSON) here!
2. **Repair** 🔧: This is a future feature to help fix any setup issues.
3. **Monitor** 👀: A placeholder for exciting monitoring features coming soon!
4. **Exit** ❌: Wrap things up and exit the script.

### ⚙️ Configuring Your Environment

During the installation, you’ll be prompted to configure the `.env` file. Here’s what you need to set up:

- **ADDRESS**: The domain or IP address for your application.
- **PORT_ADDRESS**: The port number for your application.
- **SSL**: Do you want SSL? (Answer `true` or `false`)
- **P_USER**: Your chosen username for authentication.
- **P_PASS**: A secure password for authentication.
- **MAX_ALLOW_USERS**: Maximum number of users allowed.
- **BAN_TIME**: Duration (in minutes) for which users will be banned.
- **TG_ENABLE**: Enable Telegram notifications (`true` or `false`).
    - If you choose to enable it, you’ll need:
        - **TG_TOKEN**: Your Telegram bot token.
        - **TG_ADMIN**: Your Telegram admin ID.
- **WHITELIST_ADDRESSES**: A list of IPs or domains that are allowed access, separated by commas.

### 📄 Example `.env` Configuration

Here’s a quick look at what your `.env` file might look like:

```bash
ADDRESS=example.com
PORT_ADDRESS=443
SSL=true
P_USER=admin
P_PASS=admin
MAX_ALLOW_USERS=1
BAN_TIME=5
TG_ENABLE=false
TG_TOKEN=your-telegram-bot-token
TG_ADMIN=your-telegram-admin-id
WHITELIST_ADDRESSES=127.0.0.1,example.com
```

### 🐳 Managing with Docker

The script works with Docker to keep everything running smoothly. It checks if Docker is active and uses Docker Compose for installing or uninstalling the project.

## 🔧 How to Use Watchdog

1. **Install the Project**: Start by selecting option `1` to install and configure everything according to your needs.
2. **Uninstall the Project**: If you ever need to uninstall, just select option `1` again.
3. **Repair**: This option will help you fix any issues in the future.
4. **Monitoring**: Stay tuned for upcoming features related to monitoring!

## 🗑️ Uninstalling Watchdog

When you’re ready to say goodbye to Watchdog, simply select the **Uninstall** option from the main menu, and it will take care of stopping and removing all Docker containers associated with the project.

## 💖 Donate

If you find Watchdog helpful and want to support its development, consider making a donation! Every bit helps keep the project thriving and improving.

## 🤝 Contributing

We’d love to see your contributions! Whether it’s a new feature, bug fix, or documentation improvement, please feel free to create a pull request or open an issue.

## 📄 License

This project is licensed under the MIT License. For more details, check out the [LICENSE](LICENSE) file.

## 📬 Contact Us

Got questions or feedback? Don’t hesitate to reach out at [t.me/MahdiButcher](https://t.me/MahdiButcher) 📬. We’re here to help!