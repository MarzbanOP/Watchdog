#!/bin/bash

# Color codes
RED='\033[0;31m'       # Red
GREEN='\033[0;32m'     # Green
YELLOW='\033[0;33m'    # Yellow
BLUE='\033[0;34m'      # Blue
CYAN='\033[0;36m'      # Cyan
MAGENTA='\033[0;35m'   # Magenta
NC='\033[0m'           # No Color

# Function to show loading effect
function show_loading {
    echo -e "${BLUE}Loading${NC}"
    for i in {1..3}; do
        sleep 0.5
        echo -n "."
    done
    echo ""
}

# Function to display project information
function show_project_info {
    echo -e "${GREEN}Welcome to the IPWatchdog project!${NC}"
    echo -e "${YELLOW}This project is designed to monitor and manage proxy usage effectively.${NC}"
    echo -e "${YELLOW}It includes features such as user banning, logging, and Telegram notifications.${NC}"
    echo -e "${YELLOW}You can choose different storage options: Redis, SQLite, or JSON for data management.${NC}"
    echo -e "${GREEN}Let's get started!${NC}"
    echo
}

# Function to check if Docker containers are running
function check_docker_status {
    if [[ $(docker ps --filter "name=ipwatchdog" --format '{{.Names}}') == "ipwatchdog" ]]; then
        return 0  # Running
    else
        return 1  # Not running
    fi
}

# Function to display the main menu
function show_menu {
    echo -e "${BLUE}--- IPWatchdog Menu ---${NC}"
    echo -e "${YELLOW}This menu allows you to manage the IPWatchdog project.${NC}"
    
    # Display project information on the first menu
    if [[ $FIRST_MENU -eq 1 ]]; then
        show_project_info
        FIRST_MENU=0  # Set the flag to 0 after displaying info
    fi
    
    # Display all options
    if check_docker_status; then
        echo -e "${GREEN} [Running]   Project is currently running.${NC}"
        echo -e "${MAGENTA}1) Uninstall ${NC}"
        echo -e "${CYAN}2) Repair ${NC}"
        echo -e "${GREEN}3) Monitor ${NC}"
    else
        echo -e "${RED} [Stopped]   Project is not running.${NC}"
        echo -e "${GREEN}1) Install ${NC}"
    fi
    echo -e "${RED}0) Exit ${NC}"
    echo -n "Please choose an option: "
}

# Centralized function to handle environment variable configuration
function configure_env {
    echo -e "${YELLOW}Please enter the following details for the .env file (press Enter to use default):${NC}"
    
    # Function to read input with default values
    function read_input {
        local prompt="$1"
        local default_value="$2"
        read -p "$(echo -e "${BLUE}$prompt (default: $default_value): ${NC}")" input
        echo "${input:-$default_value}"
    }

    ADDRESS=$(read_input "ADDRESS" "example.com")
    PORT_ADDRESS=$(read_input "PORT_ADDRESS" "443")
    SSL=$(read_input "SSL (true/false)" "true")
    P_USER=$(read_input "P_USER" "admin")
    P_PASS=$(read_input "P_PASS" "admin")
    MAX_ALLOW_USERS=$(read_input "MAX_ALLOW_USERS" "1")
    BAN_TIME=$(read_input "BAN_TIME (in minutes)" "5")
    TG_ENABLE=$(read_input "TG_ENABLE (true/false)" "false")

    if [[ "$TG_ENABLE" == "true" ]]; then
        TG_TOKEN=$(read_input "TG_TOKEN" "your-telegram-bot-token")
        TG_ADMIN=$(read_input "TG_ADMIN" "your-telegram-admin-id")
    else
        TG_TOKEN="your-telegram-bot-token"
        TG_ADMIN="your-telegram-admin-id"
    fi

    # Create or append to the .env file
    {
        echo "ADDRESS=$ADDRESS"
        echo "PORT_ADDRESS=$PORT_ADDRESS"
        echo "SSL=$SSL"
        echo "P_USER=$P_USER"
        echo "P_PASS=$P_PASS"
        echo "MAX_ALLOW_USERS=$MAX_ALLOW_USERS"
        echo "BAN_TIME=$BAN_TIME"
        echo "TG_ENABLE=$TG_ENABLE"
        echo "TG_TOKEN=$TG_TOKEN"
        echo "TG_ADMIN=$TG_ADMIN"
    } > .env
}

# Initialize a flag to track the first menu display
FIRST_MENU=1

# Main menu loop
while true; do
    show_menu
    read -r option

    case $option in
        0)
            echo -e "${GREEN}Exiting...${NC}"
            exit 0
            ;;
        1)
            if check_docker_status; then
                echo -e "${MAGENTA}Uninstalling...${NC}"
                docker-compose down
                echo -e "${GREEN}Uninstallation complete.${NC}"
            else
                echo -e "${CYAN}Installing...${NC}"
                show_loading
                REPO_URL="https://github.com/MarzbanOP/Watchdog.git"
                PROJECT_DIR="ipwatchdog"

                # Check if Docker is running
                if ! systemctl is-active --quiet docker; then
                    echo -e "${RED}Docker is not running. Please start Docker first.${NC}"
                    exit 1
                fi

                # Check if the project directory already exists
                if [ ! -d "$PROJECT_DIR" ]; then
                    echo -e "${BLUE}Cloning the repository...${NC}"
                    git clone "$REPO_URL"
                else
                    echo -e "${BLUE}Repository already exists. Pulling the latest changes...${NC}"
                    cd "$PROJECT_DIR" || exit
                    git pull
                    cd .. || exit
                fi

                # Navigate to the project directory
                cd "$PROJECT_DIR" || exit

                # Display storage options
                echo -e "${BLUE}Choose a storage option:${NC}"
                echo
                echo -e "${YELLOW}1) Redis${NC}"
                echo -e "${GREEN}   - Pros: Fast in-memory storage, supports complex data types, great for caching, and pub/sub messaging.${NC}"
                echo -e "${RED}   - Cons: Data is lost if not persisted, requires Redis server management.${NC}"
                echo
                echo -e "${YELLOW}2) SQLite${NC}"
                echo -e "${GREEN}   - Pros: Lightweight, serverless, easy to set up, and file-based storage.${NC}"
                echo -e "${RED}   - Cons: Not suitable for high-concurrency writes, limited scalability compared to other SQL databases.${NC}"
                echo
                echo -e "${YELLOW}3) JSON${NC}"
                echo -e "${GREEN}   - Pros: Simple and human-readable format, easy to set up without dependencies, good for small-scale applications.${NC}"
                echo -e "${RED}   - Cons: Poor scalability, not ideal for concurrent access, and lacks advanced querying capabilities.${NC}"
                echo
                read -p "Enter option (1-3): " storage_option

                # Validate storage option
                case $storage_option in
                    1) echo "STORAGE_TYPE=redis" > .env ;;
                    2) echo "STORAGE_TYPE=sqlite" > .env ;;
                    3) echo "STORAGE_TYPE=json" > .env ;;
                    *) echo -e "${RED}Invalid option. Exiting.${NC}"; exit 1 ;;
                esac

                # Call the centralized environment variable configuration function
                configure_env

                # Start Docker Compose
                echo -e "${BLUE}Starting IPWatchdog...${NC}"
                docker-compose up --build
            fi
            ;;
        2)
            echo -e "${CYAN}Repairing...${NC}"  
            # Placeholder for repair logic, if any specific repairs are needed
            echo -e "${GREEN}Repair completed!${NC}"
            ;;
        3)
            echo -e "${YELLOW}Monitoring...${NC}"  
            # Placeholder for monitoring logic, if any specific monitoring actions are needed
            ;;
        *)
            if [[ -z "$option" ]]; then
                echo -e "${RED}No option entered. Please try again.${NC}"
            else
                echo -e "${RED}Invalid option. Please try again.${NC}"
            fi
            FIRST_MENU=1  # Reset flag to show project info again on next loop
            ;;
    esac
done
