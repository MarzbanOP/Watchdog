#!/bin/bash

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'  # No Color

# Array of menu options
options=()

# Function to display project information
show_project_info() {
    echo -e "${GREEN}Welcome to the watchdog project!${NC}"
    echo -e "${YELLOW}This project is designed to monitor and manage proxy usage effectively.${NC}"
    echo -e "${YELLOW}It includes features such as user banning, logging, and Telegram notifications.${NC}"
    echo -e "${YELLOW}You can choose different storage options: Redis, SQLite, or JSON for data management.${NC}"
    echo -e "${GREEN}Let's get started!${NC}"
    echo
}

# Function to check if Docker containers are running
check_docker_status() {
    if [[ $(docker ps --filter "name=watchdog" --format '{{.Names}}') == "watchdog" ]]; then
        return 0  # Running
    else
        return 1  # Not running
    fi
}

# Function to display the menu
show_menu() {
    clear
    echo -e "${BLUE}--- Watchdog Menu ---${NC}"
<<<<<<< HEAD
    echo -e "${YELLOW}This menu allows you to manage the Watchdog project.${NC}"
    echo -e "${YELLOW}Use arrow keys to navigate and press Enter to select an option.${NC}"

    # Display project information on the first menu
    if [[ $FIRST_MENU -eq 1 ]]; then
        show_project_info
        FIRST_MENU=0  # Set the flag to 0 after displaying info
    fi

    # Display all options with highlighting
=======
    
>>>>>>> 64312b3e6b03fcac7e483281e7ae156749bc4339
    if check_docker_status; then
        echo -e "${GREEN} [Running]   Project is currently running.${NC}"
        options=("Uninstall" "Repair" "Monitor" "Exit")
    else
        echo -e "${RED} [Stopped]   Project is not running.${NC}"
        options=("Install" "Exit")
    fi

    for i in "${!options[@]}"; do
        if [[ $i -eq $selected ]]; then
            echo -e "${CYAN}> ${options[$i]} ${NC}"
        else
            echo -e "  ${options[$i]}"
        fi
    done
}

# Function to handle menu actions
handle_action() {
    case ${options[$selected]} in
        "Uninstall")
            echo -e "${YELLOW}Uninstalling...${NC}"
            docker-compose down
            echo -e "${GREEN}Uninstallation complete.${NC}"
            ;;
        "Repair")
            echo -e "${YELLOW}Repairing...${NC}"
            echo -e "${GREEN}Repair completed!${NC}"
            ;;
        "Monitor")
            echo -e "${YELLOW}Monitoring...${NC}"
            # Add monitoring logic here
            ;;
        "Install")
            echo -e "${YELLOW}Installing...${NC}"
            # Add installation logic here
            ;;
        "Exit")
            echo -e "${GREEN}Exiting...${NC}"
            exit 0
            ;;
    esac
    read -n 1 -s -r -p "Press any key to continue..."
}

# Initialize selection
selected=0

# Show project info at the start
show_project_info

# Main loop
while true; do
    show_menu

    # Read a single character without requiring Enter
    read -s -n 1 key

    # Handle the arrow keys
    if [[ $key == $'\e' ]]; then
        read -s -n 2 key
        if [[ $key == '[A' ]]; then  # Up arrow
            ((selected--))
            [[ $selected -lt 0 ]] && selected=$((${#options[@]} - 1))
        elif [[ $key == '[B' ]]; then  # Down arrow
            ((selected++))
            [[ $selected -ge ${#options[@]} ]] && selected=0
        fi
    elif [[ $key == '' ]]; then  # Enter key
        handle_action
    fi
done

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
    WHITELIST_ADDRESSES=$(read_input "WHITELIST_ADDRESSES (comma-separated)" "127.0.0.1")
    USER_DELETE_DELAY=$(read_input "USER_DELETE_DELAY (in seconds)" "30")

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
        echo "USER_DELETE_DELAY=$USER_DELETE_DELAY"
        echo "TG_ENABLE=$TG_ENABLE"
        echo "TG_TOKEN=$TG_TOKEN"
        echo "TG_ADMIN=$TG_ADMIN"
        echo "WHITELIST_ADDRESSES=$WHITELIST_ADDRESSES"
    } > .env
}

# Initialize a flag to track the first menu display
FIRST_MENU=1
selected_option=0  # Start at the first option

# Main menu loop
while true; do
    show_menu

    # Read single character input (including arrow keys)
    read -rsn1 input
    if [[ $input == $'\e' ]]; then
        read -rsn2 input # read the two characters after escape
        case "$input" in
            '[A')  # Up arrow
                ((selected_option--))
                if [[ $selected_option -lt 0 ]]; then
                    selected_option=0
                fi
                ;;
            '[B')  # Down arrow
                ((selected_option++))
                if [[ $selected_option -ge ${#options[@]} ]]; then
                    selected_option=$((${#options[@]} - 1))
                fi
                ;;
        esac
    elif [[ $input == $'\n' ]]; then
        # Check if Enter is pressed without a valid selection
        if [[ ${options[selected_option]} ]]; then
            # Check the selected option and execute the corresponding action
            case ${options[selected_option]} in
                "Uninstall")
                    echo -e "${MAGENTA}Uninstalling...${NC}"
                    docker-compose down
                    echo -e "${GREEN}Uninstallation complete.${NC}"
                    ;;
                "Repair")
                    echo -e "${CYAN}Repairing...${NC}"  
                    # Placeholder for repair logic, if any specific repairs are needed
                    echo -e "${GREEN}Repair completed!${NC}"
                    ;;
                "Monitor")
                    echo -e "${YELLOW}Monitoring...${NC}"  
                    # Placeholder for monitoring logic, if any specific monitoring actions are needed
                    ;;
                "Install")
                    echo -e "${CYAN}Installing...${NC}"
                    show_loading
                    REPO_URL="https://github.com/MarzbanOP/Watchdog.git"
                    PROJECT_DIR="Watchdog"

                    # Check if Docker is running
                    if ! systemctl is-active --quiet docker; then
                        echo -e "${RED}Docker is not running. Please start Docker first.${NC}"
                        exit 1
                    fi

                    # Check if the project directory already exists
                    if [ -d "$PROJECT_DIR" ]; then
                        echo -e "${RED}Directory '$PROJECT_DIR' already exists. Deleting it...${NC}"
                        rm -rf "$PROJECT_DIR"
                    fi

                    # Clone the repository
                    echo -e "${BLUE}Cloning the repository...${NC}"
                    git clone "$REPO_URL"

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
                    echo -e "${BLUE}Starting watchdog...${NC}"
                    docker-compose up --build || { echo -e "${RED}Failed to start Docker compose.${NC}"; exit 1; }
                    ;;
                "Exit")
                    echo -e "${GREEN}Exiting...${NC}"
                    exit 0
                    ;;
            esac
        fi
    fi
done
