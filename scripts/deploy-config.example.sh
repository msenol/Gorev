################################################################################
# Gorev Website Deployment Configuration
# Copy this file to deploy-config.sh and customize your settings
#
# SECURITY WARNING:
# - This file contains sensitive information (VPS credentials)
# - deploy-config.sh is in .gitignore and will NOT be committed
# - NEVER share your VPS credentials or add them to Git
################################################################################

# VPS Configuration
# Replace with your actual VPS details
VPS_HOST="your-vps-ip-here"
VPS_USER="your-username"
VPS_PORT="22"

# Deployment Configuration
# Path to your website directory on the VPS
DEPLOY_DIR="/var/www/gorev.work"

################################################################################
# Instructions:
#
# 1. Copy this file:
#    cp deploy-config.example.sh deploy-config.sh
#
# 2. Edit deploy-config.sh and replace:
#    - your-vps-ip-here with your actual VPS IP
#    - your-username with your actual username
#
# 3. Make sure you have SSH key access to your VPS
#
# 4. Run deployment:
#    ./deploy-website.sh --dry-run    # Test
#    ./deploy-website.sh               # Deploy
################################################################################
