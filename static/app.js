document.addEventListener('DOMContentLoaded', () => {
    const loginView = document.getElementById('login-view');
    const appView = document.getElementById('app-view');
    const loginForm = document.getElementById('login-form');
    const logoutBtn = document.getElementById('logout-btn');
    const uploadForm = document.getElementById('uploadForm');
    const contentDiv = document.getElementById('content');
    const welcomeUserSpan = document.getElementById('welcome-user');
    const galleryToolbar = document.getElementById('gallery-toolbar');
    const filterDropdown = document.getElementById('filter-dropdown'); 

    Object.defineProperty(String.prototype, 'capitalize', {
        value: function() {
            return this.charAt(0).toUpperCase() + this.slice(1);
        },
        enumerable: false
    });

    let photoCache = [];
    let currentFilter = 'all';

    function showLoginView() {
        loginView.style.display = 'block';
        appView.style.display = 'none';
    }

    function showAppView(user) {
        welcomeUserSpan.textContent = user.Username.capitalize();
        loginView.style.display = 'none';
        appView.style.display = 'block';
        refreshPhotos();
    }

    async function checkSession() {
        try {
            const response = await fetch('/api/session/check');
            if (response.ok) {
                const user = await response.json();
                showAppView(user);
            } else {
                showLoginView();
            }
        } catch (error) {
            console.error("Session check failed:", error);
            showLoginView();
        }
    }
    async function refreshPhotos() {
        try {
            const response = await fetch('/api/photos');
            const photos = await response.json();
            photoCache = photos || [];

            const owners = [...new Set(photoCache.map(p => p.owner))].sort();
            let optionsHTML = `<option value="all">All Photos</option>`;
            owners.forEach(owner => {
                optionsHTML += `<option value="${owner}" ${currentFilter === owner ? 'selected' : ''}>${owner.capitalize()}'s Photos</option>`;
            });
            filterDropdown.innerHTML = optionsHTML;

            const filteredPhotos = currentFilter === 'all'
                ? photoCache
                : photoCache.filter(p => p.owner === currentFilter);

            let html = `<br/>`;
            if (filteredPhotos.length > 0) {
                html += `<div class="gallery">`;
                filteredPhotos.forEach(photo => {
                    const photoDate = new Date(photo.date).toLocaleDateString("fr-FR");
                    const toggleVisibilityText = photo.is_public ? 'Make Private' : 'Make Public';
                    const privacyClass = photo.is_public ? 'public' : 'private';
                    const statusHTML = photo.is_public ? `<span>🌎 Public</span>` : `<b>🔒 Private</b>`;
                    html += `
                        <div class="photo-card ${privacyClass}">
                            <div class="menu-container">
                                <button class="menu-btn">⋮</button>
                                <div class="menu-dropdown">
                                    <a href="#" class="toggle-visibility-link" data-id="${photo.id}">${toggleVisibilityText}</a>
                                    <a href="#" class="delete-link" data-id="${photo.id}">Delete</a>
                                </div>
                            </div>
                            <img src="/api/photos/${photo.id}" alt="${photo.title}" class="zoomable">
                            <h3>${photo.title}</h3>
                            <p>By: ${photo.owner.capitalize()}</p>
                            <p><em>${photoDate}</em></p>
                            <p>Status: ${statusHTML}</p>
                        </div>
                    `;
                });
                html += `</div>`;
            } else {
                html += `<p>No photos match the current filter.</p>`;
            }
            contentDiv.innerHTML = html;

        } catch (error) {
            console.log(error)
            contentDiv.innerHTML = `<p style="color:red;">Error loading photos.</p>`;
        }
    }


    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('username').value;
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username }),
        });

        if (response.ok) {
            checkSession();
        } else {
            alert('Login failed: User not found.');
        }
    });

    logoutBtn.addEventListener('click', async () => {
        await fetch('/api/logout', { method: 'POST' });
        showLoginView();
    });

    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(uploadForm);
        const response = await fetch('/api/photos', {
            method: 'POST',
            body: formData,
        });

        if (response.ok) {
            uploadForm.reset();
            refreshPhotos();
        } else {
            alert('Upload failed! Your session may have expired.');
        }
    });

    contentDiv.addEventListener('click', async (e) => {
        if (e.target.classList.contains('menu-btn')) {
            e.stopPropagation();
            const currentMenu = e.target.nextElementSibling;
            document.querySelectorAll('.menu-dropdown.show').forEach(openMenu => {
                if (openMenu !== currentMenu) {
                    openMenu.classList.remove('show');
                }
            });
            currentMenu.classList.toggle('show');
        }

        if (e.target.classList.contains('toggle-visibility-link')) {
            e.preventDefault();
            const photoId = parseInt(e.target.getAttribute('data-id'), 10);
            const photoToUpdate = photoCache.find(p => p.id === photoId);

            if (!photoToUpdate) {
                alert('Could not find photo data to update.');
                return;
            }

            const payload = { ...photoToUpdate, is_public: !photoToUpdate.is_public };

            const response = await fetch(`/api/photos/${photoId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (response.ok) {
                refreshPhotos();
            } else {
                alert('Failed to update photo visibility. You may not be the owner.');
            }
        }

        if (e.target.classList.contains('delete-link')) {
            e.preventDefault();
            const photoId = e.target.getAttribute('data-id');
            if (confirm('Are you sure you want to permanently delete this photo?')) {
                const response = await fetch(`/api/photos/${photoId}`, {
                    method: 'DELETE',
                });
                if (response.ok) {
                    refreshPhotos();
                } else {
                    alert('Delete failed! You may not have the right to delete this photo.');
                }
            }
        }
    });

    window.addEventListener('click', (e) => {
        if (!e.target.matches('.menu-btn')) {
            document.querySelectorAll('.menu-dropdown.show').forEach(openMenu => {
                openMenu.classList.remove('show');
            });
        }
    });

    filterDropdown.addEventListener('change', (e) => {
        currentFilter = e.target.value;
        refreshPhotos();
    });


    galleryToolbar.addEventListener('click', (e) => {
        if (e.target.classList.contains('filter-btn')) {
            currentFilter = e.target.dataset.filter;
            refreshPhotos(); 
        }
    });


    checkSession();
});