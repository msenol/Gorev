// Minimal JavaScript for Gorev Website

// Smooth scrolling for navigation links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Intersection Observer for fade-in animations
const observerOptions = {
    threshold: 0.1,
    rootMargin: '0px 0px -50px 0px'
};

const observer = new IntersectionObserver(function(entries) {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.opacity = '1';
            entry.target.style.transform = 'translateY(0)';
        }
    });
}, observerOptions);

// Animate feature cards on scroll
document.querySelectorAll('.feature-card').forEach(card => {
    card.style.opacity = '0';
    card.style.transform = 'translateY(30px)';
    card.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
    observer.observe(card);
});

// GitHub API integration
async function fetchGitHubStats() {
    const repo = 'msenol/Gorev';
    const apiUrl = `https://api.github.com/repos/${repo}`;

    try {
        const response = await fetch(apiUrl);
        if (response.ok) {
            const data = await response.json();

            const starsElement = document.getElementById('stars');
            const forksElement = document.getElementById('forks');
            const issuesElement = document.getElementById('issues');

            if (starsElement) starsElement.textContent = `⭐ ${data.stargazers_count}`;
            if (forksElement) forksElement.textContent = `🍴 ${data.forks_count}`;
            if (issuesElement) issuesElement.textContent = `🐛 ${data.open_issues_count}`;
        }
    } catch (error) {
        console.log('GitHub API unavailable');
    }
}

// Initialize on load
document.addEventListener('DOMContentLoaded', () => {
    fetchGitHubStats();
});
