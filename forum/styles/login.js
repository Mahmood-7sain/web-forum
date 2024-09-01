const container = document.getElementById('container');
const registerBtn = document.getElementById('register');
const loginBtn = document.getElementById('login');
const darkModeToggle = document.getElementById('dark-mode-toggle');

registerBtn.addEventListener('click', () => {
    container.classList.add("active");
});

loginBtn.addEventListener('click', () => {
    container.classList.remove("active");
});

darkModeToggle.addEventListener('click', function() {
    document.body.dataset.theme = document.body.dataset.theme === 'dark' ? 'light' : 'dark';
    localStorage.setItem('theme', document.body.dataset.theme);
    updateThemeIcon();
});

function updateThemeIcon() {
    const icon = darkModeToggle.querySelector('i');
    icon.classList.toggle('fa-moon', document.body.dataset.theme !== 'dark');
    icon.classList.toggle('fa-sun', document.body.dataset.theme === 'dark');

    const create = document.getElementById("createacc")
    const signin = document.getElementById("signin")
    
    if(document.body.dataset.theme !== 'dark'){
        create.style.color = "black"
        signin.style.color = "black"
    } else {
        create.style.color = "white"
        signin.style.color = "white"
    }
}

const savedTheme = localStorage.getItem('theme');
if (savedTheme) {
    document.body.dataset.theme = savedTheme;
    updateThemeIcon();
}