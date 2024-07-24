// Get the login form
const loginForm = document.getElementById('LoginForm'); 
const errorMessage = document.getElementById("ErrorMessage");

document.addEventListener('DOMContentLoaded', function() {
    submitForm();
});

loginForm.addEventListener('submit', function(event) {
    event.preventDefault();
    event.stopPropagation();
    submitForm();
});

async function submitForm() {
    const response = await fetch ('/admin/', {
        method: 'POST',
    })

    if (response.status === 200) {
        window.location.href = '/admin/dashboard';
    } else {
        errorMessage.innerHTML = "Ongeldige inlogpoging";
    }
};

// Enter key to submit
document.addEventListener('keydown', function(event) {
     if (event.key === 'Enter') {
        submitForm();
    }
});