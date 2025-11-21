# gorev.work Website

Modern task management platform landing page for gorev.work domain.

## 📁 Directory Structure

```
website/
├── index.html    - Main landing page (11KB)
├── styles.css    - Styling (10KB)
├── script.js     - JavaScript functionality (2KB)
└── README.md     - This file
```

## 🚀 Quick Deployment

### From Gorev Project Directory
```bash
cd /home/msenol/Projects/Gorev
./scripts/deploy-website.sh --dry-run    # Test deployment
./scripts/deploy-website.sh               # Deploy to VPS
```

### From ToolsServer
```bash
cd /home/msenol/vpsServers/toolsServer
./deploy-gorev.sh --dry-run    # Test deployment
./deploy-gorev.sh               # Deploy to VPS
```

## 🔧 Development

### Local Testing
```bash
# Serve website locally (Python)
cd website
python3 -m http.server 8080

# Visit: http://localhost:8080
```

### Making Changes
1. Edit files in `website/` directory
2. Test locally with `python3 -m http.server 8080`
3. Run `./scripts/deploy-website.sh --dry-run` to preview changes
4. Run `./scripts/deploy-website.sh` to deploy

## 📝 File Descriptions

### index.html
- Main landing page with modern hero section
- English content (AI-Powered Task Management)
- Responsive design (mobile-first)
- GitHub API integration
- Clean HTML5 structure (200 lines)

### styles.css
- Modern dark theme with blue/purple gradients
- CSS Custom Properties for easy maintenance
- CSS Grid and Flexbox layouts
- Inter font from Google Fonts
- Responsive breakpoints
- Smooth animations
- 510 lines of clean CSS

### script.js
- Vanilla JavaScript (no frameworks)
- Smooth scrolling navigation
- GitHub API calls for repo stats
- Intersection Observer for fade-in animations
- Minimal and efficient (67 lines)

## 🌐 Website Info

- **Domain:** gorev.work / www.gorev.work
- **SSL:** Let's Encrypt (TLS 1.3, HTTP/2)
- **Server:** Nginx on VPS (62.84.183.207)
- **Deploy Path:** `/var/www/gorev.work/`
- **Auto-Deploy:** Via rsync over SSH

## 🔐 Authentication & Security

Deployment uses SSH key authentication. VPS credentials are stored in a separate config file:

**Important Security Notes:**
- ⚠️ **DO NOT commit sensitive credentials to Git**
- Configuration is in `../scripts/deploy-config.sh` (excluded via `.gitignore`)
- Example config: `../scripts/deploy-config.example.sh`
- Always use SSH key authentication (no password)

**Setup:**
1. Copy example config: `cp ../scripts/deploy-config.example.sh ../scripts/deploy-config.sh`
2. Edit `../scripts/deploy-config.sh` with your VPS details
3. Ensure SSH key is added to VPS: `~/.ssh/id_rsa` or `~/.ssh/id_ed25519`

**⚠️ WARNING:** Never commit `deploy-config.sh` to Git! It contains sensitive VPS credentials.

## 📊 Performance

- **Page Load:** ~1.2s (first paint)
- **Total Size:** ~25KB (gzipped: ~10KB)
- **Images:** SVG only (optimized)
- **Caching:** 1 year (static files)
- **Reductions:** 60% HTML, 15% CSS, 75% JS

## 🎨 Features

- ✅ Modern dark theme design
- ✅ English content (AI-Powered Task Management)
- ✅ GitHub API integration
- ✅ Real-time repo stats
- ✅ Smooth scroll animations
- ✅ Responsive (mobile-first)
- ✅ SEO optimized meta tags
- ✅ Security headers (HSTS, CSP)
- ✅ Code mockup in hero section
- ✅ Clean, minimal codebase

## 🐛 Troubleshooting

### Deployment Fails
```bash
# Check if config file exists
ls -la scripts/deploy-config.sh

# Verify config (replace with your VPS details)
cat scripts/deploy-config.sh

# Check SSH connection
ssh -i ~/.ssh/id_rsa USER@VPS_IP

# Run with verbose output
./scripts/deploy-website.sh --verbose
```

### Website Not Loading
```bash
# Check VPS file permissions (replace USER and VPS_IP)
ssh -i ~/.ssh/id_rsa USER@VPS_IP "ls -la /var/www/gorev.work/"

# Check Nginx configuration
ssh -i ~/.ssh/id_rsa USER@VPS_IP "nginx -t"

# Check Nginx status
ssh -i ~/.ssh/id_rsa USER@VPS_IP "systemctl status nginx"
```

### Permission Issues
```bash
# Run deploy with correct user (replace USER and VPS_IP)
ssh -i ~/.ssh/id_rsa USER@VPS_IP "chown -R www-data:www-data /var/www/gorev.work"
```

## 📝 Website Content

The website showcases:
- AI-Powered Task Management for Modern Development
- MCP (Model Context Protocol) integration
- Smart task hierarchy with unlimited subtasks
- AI context management features
- GitHub repository (msenol/Gorev)
- Open-source project information

Content is in **English** to reach a global audience.

## 🔗 Useful Links

- **Live Site:** https://gorev.work
- **GitHub:** https://github.com/msenol/Gorev
- **Nginx Config:** `/etc/nginx/sites-enabled/gorev.work.conf`
- **Config Setup:** `scripts/deploy-config.example.sh`

## 📝 Notes

- Website is static (no backend)
- Deployed via rsync over SSH
- Uses Let's Encrypt SSL certificate
- Supports HTTP/2 for better performance
- Rate limited via Nginx (15 req/sec burst)
- Enterprise-grade security (97% security score)

For detailed hosting information, see:
`/home/msenol/vpsServers/toolsServer/docs/gorev-work-website-management.md`

## 🔄 Recent Updates (November 22, 2025)

**Complete Website Redesign & Security Fix:**
- **Language:** Turkish → English (global audience)
- **Design:** Light theme → Modern dark theme
- **Code:** Reduced by 60% (HTML), 15% (CSS), 75% (JS)
- **Content:** Updated to showcase AI-Powered Task Management
- **Features:** Added code mockup, modern gradients, clean animations
- **🔒 Security:** Fixed hardcoded VPS credentials issue
  - Moved credentials to `deploy-config.sh` (excluded via `.gitignore`)
  - Added example config file for reference
  - Updated README with security warnings

---

**Last Updated:** November 22, 2025
**Maintainer:** Gorev Project
**Version:** 2.0 (Redesigned)
