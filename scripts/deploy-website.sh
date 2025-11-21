#!/bin/bash

################################################################################
# Gorev Website Deployment Script
# Deploys gorev.work website to VPS server
#
# Usage:
#   ./deploy-website.sh [--dry-run] [--verbose]
#
# Configuration:
#   Copy deploy-config.example.sh to deploy-config.sh and customize
################################################################################

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Load configuration (if exists)
CONFIG_FILE="$(dirname "$0")/deploy-config.sh"
if [ -f "$CONFIG_FILE" ]; then
    source "$CONFIG_FILE"
else
    echo -e "${YELLOW}⚠ Warning: Configuration file not found at $CONFIG_FILE${NC}"
    echo -e "${YELLOW}  Please copy deploy-config.example.sh to deploy-config.sh and customize it${NC}"
    echo ""
    echo "Example configuration:"
    echo "  VPS_HOST=\"your-vps-ip\""
    echo "  VPS_USER=\"your-username\""
    echo "  DEPLOY_DIR=\"/var/www/gorev.work\""
    exit 1
fi

# Set defaults if not provided in config
VPS_HOST="${VPS_HOST:-}"
VPS_USER="${VPS_USER:-}"
VPS_PORT="${VPS_PORT:-22}"
DEPLOY_DIR="${DEPLOY_DIR:-/var/www/gorev.work}"
LOCAL_DIR="$(dirname "$0")/../website"
LOG_FILE="/tmp/gorev-deploy-$(date +%Y%m%d-%H%M%S).log"

# Parse arguments
DRY_RUN=false
VERBOSE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [--dry-run] [--verbose]"
            echo ""
            echo "Options:"
            echo "  --dry-run    Test deployment without actually uploading"
            echo "  --verbose    Show detailed rsync output"
            echo "  -h, --help   Show this help message"
            echo ""
            echo "Configuration:"
            echo "  Create deploy-config.sh from deploy-config.example.sh"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Validate configuration
if [ -z "$VPS_HOST" ] || [ -z "$VPS_USER" ]; then
    echo -e "${RED}✗ Error: VPS_HOST and VPS_USER must be configured${NC}"
    echo -e "${YELLOW}  Please set these in deploy-config.sh${NC}"
    exit 1
fi

# Check if local directory exists
if [ ! -d "$LOCAL_DIR" ]; then
    echo -e "${RED}✗ Error: Local directory not found: $LOCAL_DIR${NC}"
    exit 1
fi

# Function to log messages
log() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "[$timestamp] ${2:-$NC}ℹ $1${NC}" | tee -a "$LOG_FILE"
}

log_error() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "[$timestamp] ${RED}✗ $1${NC}" | tee -a "$LOG_FILE"
}

log_success() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "[$timestamp] ${GREEN}✓ $1${NC}" | tee -a "$LOG_FILE"
}

log_warning() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "[$timestamp] ${YELLOW}⚠ $1${NC}" | tee -a "$LOG_FILE"
}

# Print header
echo -e "${BLUE}ℹ ==========================================${NC}"
echo -e "${BLUE}ℹ Gorev Website Deployment${NC}"
echo -e "${BLUE}ℹ ==========================================${NC}"
log "Started at: $(date)" "$BLUE"
echo -e "${BLUE}ℹ ==========================================${NC}"
echo ""

# Check requirements
log "Checking requirements..." "$BLUE"

# Check if rsync is installed
if ! command -v rsync &> /dev/null; then
    log_error "rsync is not installed. Please install it first."
    exit 1
fi

# Check if SSH key exists
SSH_KEY="$HOME/.ssh/id_rsa"
if [ ! -f "$SSH_KEY" ]; then
    SSH_KEY="$HOME/.ssh/id_ed25519"
    if [ ! -f "$SSH_KEY" ]; then
        log_error "SSH key not found at $SSH_KEY or $HOME/.ssh/id_rsa"
        exit 1
    fi
fi

log_success "All requirements met"
echo ""

# Test SSH connection
log "Testing SSH connection to $VPS_USER@$VPS_HOST..." "$BLUE"

if ssh -o ConnectTimeout=10 -o BatchMode=yes -i "$SSH_KEY" -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" "echo 'SSH connection successful'" &> /dev/null; then
    log_success "SSH connection successful"
else
    log_error "SSH connection failed. Please check:"
    echo "  1. VPS IP and username are correct"
    echo "  2. SSH key is added to VPS"
    echo "  3. Network connectivity"
    exit 1
fi
echo ""

# Create backup on VPS
if [ "$DRY_RUN" = false ]; then
    log "Creating backup on VPS..." "$BLUE"
    BACKUP_DIR="/var/backups/gorev-work-$(date +%Y%m%d-%H%M%S)"
    ssh -i "$SSH_KEY" -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" \
        "mkdir -p $BACKUP_DIR && cp -r $DEPLOY_DIR/* $BACKUP_DIR/ 2>/dev/null || true" || true
    log_success "Backup created at $BACKUP_DIR" "$GREEN"
    echo ""
fi

# Deploy files
log "Deploying website files..." "$BLUE"
log "From: $LOCAL_DIR" "$BLUE"
log "To: $VPS_USER@$VPS_HOST:$DEPLOY_DIR/" "$BLUE"
echo ""

if [ "$DRY_RUN" = true ]; then
    log "DRY RUN: Would create backup of $DEPLOY_DIR"
    log "DRY RUN: Would deploy website files"
    RSYNC_CMD="rsync -azq --progress --delete $LOCAL_DIR/ $VPS_USER@$VPS_HOST:$DEPLOY_DIR/"
    log "Command: $RSYNC_CMD"
    echo ""
    log "DRY RUN: Files that would be deployed:" "$BLUE"
    ls -lah "$LOCAL_DIR"
    echo ""
    log "DRY RUN: Would set permissions (www-data:www-data 644 for files, 755 for dirs)"
    echo ""
    log "DRY RUN: Skipping verification" "$BLUE"
else
    # Perform actual deployment
    if [ "$VERBOSE" = true ]; then
        rsync -azv --progress --delete \
            -e "ssh -i $SSH_KEY -p $VPS_PORT" \
            "$LOCAL_DIR/" "$VPS_USER@$VPS_HOST:$DEPLOY_DIR/"
    else
        rsync -azq --delete \
            -e "ssh -i $SSH_KEY -p $VPS_PORT" \
            "$LOCAL_DIR/" "$VPS_USER@$VPS_HOST:$DEPLOY_DIR/" >> "$LOG_FILE" 2>&1
    fi

    if [ $? -eq 0 ]; then
        log_success "Website files deployed successfully"
    else
        log_error "Deployment failed"
        exit 1
    fi
    echo ""

    # Set permissions
    log "Setting file permissions..." "$BLUE"
    ssh -i "$SSH_KEY" -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" \
        "find $DEPLOY_DIR -type f -exec chmod 644 {} \; && \
         find $DEPLOY_DIR -type d -exec chmod 755 {} \; && \
         chown -R www-data:www-data $DEPLOY_DIR 2>/dev/null || true" || true
    log_success "Permissions set successfully"
    echo ""

    # Verify deployment
    log "Verifying deployment..." "$BLUE"
    ssh -i "$SSH_KEY" -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" "test -f $DEPLOY_DIR/index.html && echo 'index.html exists'" >> "$LOG_FILE" 2>&1

    LOCAL_FILE_COUNT=$(find "$LOCAL_DIR" -type f | wc -l)
    VPS_FILE_COUNT=$(ssh -i "$SSH_KEY" -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" "find $DEPLOY_DIR -type f | wc -l")

    if [ "$LOCAL_FILE_COUNT" -eq "$VPS_FILE_COUNT" ]; then
        log_success "File count matches: $VPS_FILE_COUNT files" "$GREEN"
    else
        log_warning "File count mismatch: Local=$LOCAL_FILE_COUNT, VPS=$VPS_FILE_COUNT"
    fi
    echo ""
fi

# Print summary
echo -e "${BLUE}ℹ ==========================================${NC}"
echo -e "${BLUE}ℹ Deployment Summary${NC}"
echo -e "${BLUE}ℹ ==========================================${NC}"
log "Mode: $([ "$DRY_RUN" = true ] && echo 'DRY RUN' || echo 'LIVE')" "$BLUE"
log "Target: $VPS_USER@$VPS_HOST:$DEPLOY_DIR" "$BLUE"
log "Local Directory: $LOCAL_DIR" "$BLUE"
log "Log File: $LOG_FILE" "$BLUE"
echo -e "${BLUE}ℹ ==========================================${NC}"

if [ "$DRY_RUN" = false ]; then
    echo -e "${BLUE}ℹ Website should be available at: https://gorev.work${NC}"
    echo -e "${BLUE}ℹ Test the deployment by visiting the URL above${NC}"
    echo ""
fi

if [ "$DRY_RUN" = true ]; then
    log_success "Deployment completed successfully!" "$GREEN"
else
    log_success "Deployment completed successfully!" "$GREEN"
fi
