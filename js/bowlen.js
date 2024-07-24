document.getElementById('bowlingBookingForm').addEventListener('submit', function(event) {
  event.preventDefault();
  const form = event.target;
  const formData = new FormData(form);
  const data = Object.fromEntries(formData.entries());

  fetch(form.action, {
      method: form.method,
      body: JSON.stringify(data),
      headers: {
          'Content-Type': 'application/json'
      }
  }).then(response => {
      if (!response.ok) {
          return response.text().then(text => { 
              console.error('Error response:', text);
              throw new Error(text); 
          });
      }
      return response.json();
  })
  .then(data => {
      if (data.status === 'success') {
          alert('Reservering succesvol!');
          form.reset();
      } else {
          alert('Fout bij het reserveren: ' + data.message);
      }
  })
  .catch(error => {
      console.error('Fetch error:', error);
      alert('Er is een technische fout opgetreden.');
  });
});

document.addEventListener('DOMContentLoaded', function() {
  const bookingDate = document.getElementById('bookingDate');
  const bookingTime = document.getElementById('bookingTime');
  const numAdults = document.getElementById('numAdults');
  const numChildren = document.getElementById('numChildren');
  const promoCode = document.getElementById('promoCode');
  const totalCostElement = document.getElementById('totalCost');
  const totalCostInput = document.getElementById('totalCostInput');

  let pricePerAdult;
  let pricePerChild;
  let promoCodeDiscount;

  function loadConfiguration() {
      fetch('frontend/config.json')
      .then(response => response.json())
      .then(config => {
          pricePerAdult = config.bookingPricePerAdult;
          pricePerChild = config.bookingPricePerChild;
          promoCodeDiscount = config.promoCodeDiscount;
      })
      .catch(error => console.error('Error loading configuration:', error));
  }

  function updateTotalCost() {
      const adults = parseInt(numAdults.value, 10) || 0;
      const children = parseInt(numChildren.value, 10) || 0;
      let totalCost = (adults * pricePerAdult) + (children * pricePerChild);

      if (promoCode.value.trim().toLowerCase() === 'bowlen10') {
          totalCost *= promoCodeDiscount;
      }

      totalCostElement.textContent = totalCost.toFixed(2);
      totalCostInput.value = totalCost.toFixed(2);
  }

  numAdults.addEventListener('input', updateTotalCost);
  numChildren.addEventListener('input', updateTotalCost);
  promoCode.addEventListener('input', updateTotalCost);

  const today = new Date();
  const minDate = today.toISOString().split('T')[0];
  bookingDate.min = minDate;

  bookingDate.addEventListener('change', function() {
      updateBookingTimes(this.value);
  });

  function updateBookingTimes(date) {
      bookingTime.innerHTML = '';
      const selectedDate = new Date(date);
      const isToday = selectedDate.toISOString().split('T')[0] === today.toISOString().split('T')[0];

      const currentHour = today.getHours();
      const currentMinutes = today.getMinutes();
      const startTime = isToday && currentHour >= 10 && currentHour < 22 ? currentHour : 10;
      const endTime = 22;

      for (let hour = startTime; hour < endTime; hour++) {
          if (!isToday || hour > currentHour || (hour === currentHour && currentMinutes < 30)) {
              addTimeOption(hour, '00');
          }
          if (!isToday || hour > currentHour || (hour === currentHour && currentMinutes < 60)) {
              addTimeOption(hour, '30');
          }
      }

      if (bookingTime.options.length === 0) {
          bookingTime.appendChild(new Option("Geen beschikbare tijden", ""));
      }
  }

  function addTimeOption(hour, minute) {
      const timeString = `${hour.toString().padStart(2, '0')}:${minute}`;
      const option = new Option(timeString, timeString);
      bookingTime.appendChild(option);
  }

  var userIcon = document.getElementById('user-icon');
  var dropdown = document.getElementById('user-dropdown');

  var showDropdown = function() {
      dropdown.style.display = 'block';
  };

  var hideDropdown = function() {
      setTimeout(function() {
          if (!dropdown.matches(':hover') && !userIcon.matches(':hover')) {
              dropdown.style.display = 'none';
          }
      }, 200);
  };

  userIcon.addEventListener('mouseenter', showDropdown);
  userIcon.addEventListener('mouseleave', hideDropdown);

  dropdown.addEventListener('mouseenter', showDropdown);
  dropdown.addEventListener('mouseleave', hideDropdown);
});


// Dit blok code behandelt het inloggen en uitloggen van gebruikers.

// Code om de modal voor het veranderen van het wachtwoord te beheren
var modal = document.getElementById("changePasswordModal"); // Krijg toegang tot de modal-element
var btn = document.querySelector(".user-menu-dropdown ul li a[href='#change-password']"); // Knop die de modal opent
var span = document.getElementsByClassName("close")[0]; // Het element om de modal te sluiten

// Wanneer de gebruiker op de knop klikt, open de modal
btn.onclick = function(event) {
  event.preventDefault(); // Voorkom de standaard link actie
  modal.style.display = "block"; // Toon de modal
}

// Sluit de modal wanneer de gebruiker op het sluit-icoon (x) klikt
span.onclick = function() {
  modal.style.display = "none"; // Verberg de modal
  messageContainer.innerHTML = ""; // Maak de berichtencontainer leeg
}

// Sluit de modal wanneer de gebruiker ergens buiten de modal klikt
window.onclick = function(event) {
  if (event.target == modal) {
      modal.style.display = "none";
      messageContainer.innerHTML = ""; // Maak de berichtencontainer leeg
  }
}

// Code voor het indienen van het formulier om het wachtwoord te wijzigen
document.getElementById("changePasswordForm").onsubmit = function(event) {
  event.preventDefault(); // Voorkom de standaard form submit actie

  var oldPassword = document.getElementById("oldPassword").value;
  var newPassword = document.getElementById("newPassword").value;
  var confirmPassword = document.getElementById("confirmPassword").value;

  // Controleer of de nieuwe wachtwoorden overeenkomen
  if (newPassword !== confirmPassword) {
      messageContainer.innerHTML = "<p style='color: red;'>De nieuwe wachtwoorden komen niet overeen.</p>";
      return; // Stop de functie als de wachtwoorden niet overeenkomen
  }

  // Maak een AJAX verzoek om het wachtwoord te wijzigen
  fetch('/change-password', {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json'
      },
      body: JSON.stringify({
          oldPassword: oldPassword,
          newPassword: newPassword
      })
  })
  .then(response => response.json())
  .then(data => {
      if (data.success) {
          messageContainer.innerHTML = "<p style='color: green;'>Wachtwoord succesvol gewijzigd.</p>";
      } else {
          messageContainer.innerHTML = "<p style='color: red;'>Er is iets misgegaan: " + data.message + "</p>";
      }
  })
  .catch(error => {
      console.error('Error:', error);
      messageContainer.innerHTML = "<p style='color: red;'>Er is een fout opgetreden. Probeer het later opnieuw.</p>";
  });
};

// Code om de zichtbaarheid van wachtwoorden te toggelen
document.querySelectorAll('.toggle-password').forEach(item => {
  item.addEventListener('click', function() {
      var input = document.querySelector(this.getAttribute('toggle'));
      if (input.getAttribute('type') === 'password') {
          input.setAttribute('type', 'text'); // Verander input type naar tekst
          this.src = 'images/eye.png'; // Verander het icoon (indien van toepassing)
      } else {
          input.setAttribute('type', 'password'); // Verander input type terug naar wachtwoord
          this.src = 'images/eye.png'; // Optioneel, verander terug naar het oorspronkelijke icoon
      }
  });
});

// Code om de gebruiker uit te loggen
document.getElementById("logout-link").onclick = function(event) {
  event.preventDefault(); // Voorkom de standaard link actie

  // Maak een AJAX verzoek om uit te loggen
  fetch('/logout', {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json'
      }
  })
  .then(response => response.json())
  .then(data => {
      if (data.success) {
          window.location.href = '/index'; // Stuur de gebruiker naar de startpagina na uitloggen
      } else {
          alert('Failed to log out. Please try again.'); // Toon een foutmelding als het uitloggen mislukt
      }
  })
  .catch(error => {
      console.error('Error:', error);
      alert('An error occurred. Please try again later.'); // Toon een foutmelding bij een serverfout
  });
};
