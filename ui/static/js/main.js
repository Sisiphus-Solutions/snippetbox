let navLinks = document.querySelectorAll("nav a");
for (const element of navLinks) {
    let link = element
    if (link.getAttribute('href') === window.location.pathname) {
        link.classList.add("live");
        break;
    }
}