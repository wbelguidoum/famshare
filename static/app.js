document.addEventListener('DOMContentLoaded', () => {
    // UI Elements
    const loginView = document.getElementById('login-view');
    const appView = document.getElementById('app-view');
    const loginForm = document.getElementById('login-form');
    const logoutBtn = document.getElementById('logout-btn');
    const uploadForm = document.getElementById('uploadForm');
    const contentDiv = document.getElementById('content');
    const welcomeUserSpan = document.getElementById('welcome-user');

    // UI State Management
    function showLoginView() {
        loginView.style.display = 'block';
        appView.style.display = 'none';
    }

    function showAppView(user) {
        welcomeUserSpan.textContent = user.Username;
        loginView.style.display = 'none';
        appView.style.display = 'block';
        refreshPhotos();
    }

    // API Functions
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
            
            let html = `<h2>Family Gallery</h2>`;
            if (photos && photos.length > 0) {
                 html += `<div class="gallery">`;
                photos.forEach(photo => {
                    const photoDate = new Date(photo.date).toLocaleDateString();
                    html += `
                        <div class="photo-card">
                            <button class="delete-btn" data-id="${photo.ID}" title="Delete Photo">X</button>
                            <img src="/uploads/${photo.filename}" alt="${photo.title}">
                            <h3>${photo.title}</h3>
                            <p>By: ${photo.owner}</p>
                            <p><em>${photoDate}</em></p>
                        </div>
                    `;
                });
                html += `</div>`;
            } else {
                html += `<p>The gallery is empty. Upload a photo to get started!</p>`;
            }
            contentDiv.innerHTML = html;
        } catch (error) {
            contentDiv.innerHTML = `<p style="color:red;">Error loading photos. You may need to log in.</p>`;
        }
    }

    // --- Event Listeners ---
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('username').value;
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username }),
        });

        if (response.ok) {
            checkSession(); // On successful login, verify session and switch view
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
        
        const response = await fetch('/api/photos/upload', {
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
        if (e.target.classList.contains('delete-btn')) {
            const photoId = e.target.getAttribute('data-id');
            if (confirm('Are you sure you want to permanently delete this photo?')) {
                const response = await fetch(`/api/photos/delete/${photoId}`, {
                    method: 'DELETE',
                });

                if (response.ok) {
                    refreshPhotos();
                } else {
                    alert('Delete failed! Your session may have expired.');
                }
            }
        }
    });

    // Initial Application Load
    checkSession();
});