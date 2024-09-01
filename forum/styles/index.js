document.addEventListener('DOMContentLoaded', function() {
    const profileButton = document.getElementById('profile-button');
    const profileDropdown = document.getElementById('profile-dropdown');
    const darkModeToggle = document.getElementById('dark-mode-toggle');
    const postGrid = document.querySelector('.post-grid');
    const categoryButtons = document.querySelectorAll('.category-btn');

    // Profile dropdown
    profileButton.addEventListener('click', function(e) {
        e.stopPropagation();
        profileDropdown.style.display = profileDropdown.style.display === 'block' ? 'none' : 'block';
    });

    // Close dropdown when clicking outside
    document.addEventListener('click', function(event) {
        if (!profileButton.contains(event.target) && !profileDropdown.contains(event.target)) {
            profileDropdown.style.display = 'none';
        }
    });

    // Dark mode toggle
    darkModeToggle.addEventListener('click', function() {
        document.body.dataset.theme = document.body.dataset.theme === 'dark' ? 'light' : 'dark';
        localStorage.setItem('theme', document.body.dataset.theme);
        const icon = this.querySelector('i');
        icon.classList.toggle('fa-moon');
        icon.classList.toggle('fa-sun');
    });

    // Check for saved theme preference
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) {
        document.body.dataset.theme = savedTheme;
        if (savedTheme === 'dark') {
            darkModeToggle.querySelector('i').classList.replace('fa-moon', 'fa-sun');
        }
    }
    
    // Like and dislike functionality
    postGrid.addEventListener('click', function(e) {
        if (e.target.classList.contains('like-btn') || e.target.classList.contains('dislike-btn')) {
            const button = e.target.closest('.action-btn');
            const postId = button.closest('.post-card').dataset.postId;
            const action = button.dataset.action;

            fetch(`/updateVote?id=${postId}&action=${action}`, { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        button.querySelector('.like-count, .dislike-count').textContent = data.count;
                        button.classList.add('voted');
                        setTimeout(() => button.classList.remove('voted'), 1000);
                    }
                })
                .catch(error => console.error('Error:', error));
        }
    });

    // Initial load of all posts
    fetchPosts('all');
});function toggleCategory(categoryId) {
    const categoryContent = document.getElementById(categoryId);
    if (categoryContent.style.display === 'block') {
        categoryContent.style.display = 'none';
    } else {
        document.querySelectorAll('.category-content').forEach(content => {
            content.style.display = 'none';
        });
        categoryContent.style.display = 'block';
    }
}