# IPWatchdog

**IPWatchdog** is a Go Fiber-based proxy management system that limits IP connections and blocks suspicious activity via `iptables`. It uses JSON for configuration and includes an API for managing users and automating controls.

## Features

- Limit the number of concurrent IP connections to proxies.
- Automatically block and unblock IP addresses based on activity.
- Lightweight configuration using JSON files.
- Flexible API for user management.

## Requirements

- Go (1.17+)
- Node.js
- iptables
- gawk
- csvtool

## Installation

1. RUN COMMAND 
```bash
sudo bash -c "$(curl -sL https://raw.githubusercontent.com/MarzbanOP/Watchdog/refs/heads/main/run.sh)" @ install
```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Configure your environment in the `.env` file.

4. Set permissions for scripts:
   ```bash
   chmod +x ./ipban.sh
   chmod +x ./ipunban.sh
   ```

## Usage

Start the application:
```bash
npm start
```

To stop the application:
```bash
pm2 kill
```

## API Reference

- **GET /api/token**: Get access token.
- **POST /api/add**: Add a new user.
- **POST /api/update**: Update user limits.
- **GET /api/delete/<email>**: Delete a user.
- **GET /api/clear**: Clear the database.

## FAQ

**Q: How do I reset the application?**  
A: Navigate to the project directory and run:
```bash
pm2 kill
npm start
```

## Contributing

Contributions are welcome! Please create a pull request or open an issue.

## License

This project is licensed under the MIT License.

## Contact

For questions or feedback, please contact [t.me/MahdiButcher].
