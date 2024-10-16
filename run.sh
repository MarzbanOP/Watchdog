#!/bin/bash

# Color codes
RED='\033[0;31m'       # Red
GREEN='\033[0;32m'     # Green
YELLOW='\033[0;33m'    # Yellow
BLUE='\033[0;34m'      # Blue
CYAN='\033[0;36m'      # Cyan
NC='\033[0m'           # No Color

# Function to display the menu
show_menu() {
    clear
    echo -e "${BLUE}--- Simple Menu ---${NC}"
    echo -e "${YELLOW}Use arrow keys to navigate and press Enter to select an option.${NC}"

    # Display options with highlighting
    options=("Option 1" "Option 2" "Option 3" "Option 4" "Exit")
    for i in "${!options[@]}"; do
        if [[ $i -eq $selected_option ]]; then
            echo -e "${CYAN}> ${options[i]} ${NC}"  # Highlight selected option
        else
            echo -e "  ${options[i]}"
        fi
    done
}

# Function to display a big message
display_big_message() {
    clear
    echo -e "${GREEN}BIGITY!${NC}"  # Display big message in green
    echo -e "${YELLOW}You selected Option 1!${NC}"  # Additional message
    echo
    echo -e "${GREEN}Press any key to return to the menu...${NC}"
    read -n 1 -s  # Wait for user input
}

# Initialize selected option
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
        case ${options[selected_option]} in
            "Option 1")
                display_big_message  # Show big message for Option 1
                ;;
            "Option 2")
                echo -e "${GREEN}You selected Option 2!${NC}"
                sleep 1
                ;;
            "Option 3")
                echo -e "${GREEN}You selected Option 3!${NC}"
                sleep 1
                ;;
            "Option 4")
                echo -e "${GREEN}You selected Option 4!${NC}"
                sleep 1
                ;;
            "Exit")
                echo -e "${GREEN}Exiting...${NC}"
                exit 0
                ;;
        esac
    fi
done
