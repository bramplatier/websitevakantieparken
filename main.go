package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type Config struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type BookingData struct {
	FirstName     string  `json:"firstName"`
	LastName      string  `json:"lastName"`
	PhoneNumber   string  `json:"phoneNumber"`
	BookingDate   string  `json:"bookingDate"`
	BookingLength int     `json:"bookingLength"`
	BookingTime   string  `json:"bookingTime"`
	NumAdults     int     `json:"numAdults"`
	NumChildren   int     `json:"numChildren"`
	PromoCode     string  `json:"promoCode"`
	TotalCost     float64 `json:"totalCost"`
}

// Define a struct to hold the data you want to send to the HTML template
type FormData struct {
	FirstName    string
	LastName     string
	Email        string
	LicensePlate string
	Password     string
}

type Data struct {
	ErrorMessage string
}

// Session struct to hold Session information
type Session struct {
	username string
	expiry   time.Time
}

// Global variable to store sessions
var sessions = map[string]Session{}

// Function to check if session is valid
func (s Session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

type Accommodation struct {
	ID       int
	Name     string
	Des      string
	Price    float64
	ImageURL string
	Location string
}

type Reservation struct {
	ID        int
	Name      string
	Location  string
	Email     string
	StartDate string
	EndDate   string
}
type AccommodationPageData struct {
	FirstName      string
	Accommodations []Accommodation
}

type UsageData struct {
	ID               int     `json:"id"`
	Datum            string  `json:"datum"`
	ParkName         string  `json:"park_name"`
	ElectricityUsage int     `json:"electricity_usage"`
	GasUsage         int     `json:"gas_usage"`
	WaterUsage       int     `json:"water_usage"`
	Temperature      float64 `json:"temperature"`
}

func init() {
	// Build connection string

	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file: %v", err)
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error decoding config JSON: %v", err)
	}

	// Build connection string using config values
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		config.Server, config.User, config.Password, config.Port, config.Database)
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

}
func main() {
	fmt.Println("Starting server...")
	http.HandleFunc("/", mainForm)
	http.HandleFunc("/parken", parkenForm)
	http.HandleFunc("/contact", contactForm)
	http.HandleFunc("/login", loginForm)
	http.HandleFunc("/register-submit", submitForm)
	http.HandleFunc("/register", serveRegisterForm)
	http.HandleFunc("/admin/", adminForm)
	http.HandleFunc("/bowlen/", bowlenForm)
	http.HandleFunc("/bowlenadmin/", reservationDashboard)
	http.HandleFunc("/remove-reservation", removeReservationHandler)
	http.HandleFunc("/reserve", reserveHandler)
	http.HandleFunc("/carbonfootprint", carbonfootprintform)
	http.Handle("/accommodation", authenticate(http.HandlerFunc(accomendatieForm)))
	http.Handle("/admin/dashboard", authenticate(http.HandlerFunc(dashboardForm)))
	http.Handle("/admin/bowlingadmin", authenticate(http.HandlerFunc(reservationDashboard)))
	http.Handle("/admin/carbonfootprint", authenticate(http.HandlerFunc(carbonfootprintform)))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/plugins/", http.StripPrefix("/plugins/", http.FileServer(http.Dir("plugins"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("reservation/styles"))))
	http.ListenAndServe(":80", nil)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// hoofdpagina mainForm index.html
func mainForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for main page")

	// Check if the headers exist in the request
	email := r.Header.Get("upn")
	password := r.Header.Get("connectionstring")

	// If both headers exist, redirect to the admin
	if email != "" && password != "" {
		http.Redirect(w, r, "/admin/", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "index.html")
}

// handle the accommodation form and query the database for all availabaa accommodatis (1)
func accomendatieForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for accommodation page")

	// Haal de sessie op en controleer of deze geldig is daarna kijk hij of hij cookies heeft
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	sessionToken := cookie.Value
	session, exists := sessions[sessionToken]
	if !exists || session.isExpired() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//querries the database for all available accommodations (1)
	query := "SELECT ID, naam, beschrijving, prijs, imgurl, locatie FROM accommodations WHERE beschikbaar = 1"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create local Accommodation structs to hold the data
	var accommodations []Accommodation

	for rows.Next() {
		var acc Accommodation
		if err := rows.Scan(&acc.ID, &acc.Name, &acc.Des, &acc.Price, &acc.ImageURL, &acc.Location); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		accommodations = append(accommodations, acc)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows", http.StatusInternalServerError)
		return
	}

	// stuct aan maken voor de data die naar de template moet
	data := struct {
		FirstName      string
		SuccessMessage bool
		Accommodations []Accommodation
	}{
		FirstName:      session.username,
		SuccessMessage: r.URL.Query().Get("reserved") == "1",
		Accommodations: accommodations,
	}

	tmpl := template.Must(template.ParseFiles("reservation/accommodation.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func getBookings() ([]BookingData, error) {
	rows, err := db.Query("SELECT firstName, lastName, phoneNumber, bookingDate, bookingLength, bookingTime, numAdults, numChildren, promoCode, totalCost FROM bowling")
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var bookings []BookingData
	for rows.Next() {
		var booking BookingData
		if err := rows.Scan(&booking.FirstName, &booking.LastName, &booking.PhoneNumber, &booking.BookingDate, &booking.BookingLength, &booking.BookingTime, &booking.NumAdults, &booking.NumChildren, &booking.PromoCode, &booking.TotalCost); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return bookings, nil
}

func toJsonFunc(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(v); err != nil {
		return "", fmt.Errorf("error encoding to JSON: %w", err)
	}
	return buf.String(), nil
}

func reservationDashboard(w http.ResponseWriter, r *http.Request) {
	bookings, err := getBookings()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching bookings from database: %v", err), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"toJson": func(v interface{}) string {
			jsonStr, err := toJsonFunc(v)
			if err != nil {
				log.Printf("Error encoding JSON: %v", err)
				return ""
			}
			return jsonStr
		},
	}

	tmpl, err := template.New("bowlenadmin.html").Funcs(funcMap).ParseFiles("admin/bowlenadmin.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Bookings []BookingData
	}{
		Bookings: bookings,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

func bowlenForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for bowlen page")
	http.ServeFile(w, r, "reservation/bowlen.html")
}

func submitBowling(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var booking BookingData
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received booking data: %+v\n", booking)

	// Insert the data into the database
	query := `
		INSERT INTO bowling (FirstName, LastName, PhoneNumber, BookingDate, BookingLength, BookingTime, NumAdults, NumChildren, PromoCode, TotalCost)
		VALUES (@FirstName, @LastName, @PhoneNumber, @BookingDate, @BookingLength, @BookingTime, @NumAdults, @NumChildren, @PromoCode, @TotalCost)
	`
	_, err = db.Exec(query,
		sql.Named("FirstName", booking.FirstName),
		sql.Named("LastName", booking.LastName),
		sql.Named("PhoneNumber", booking.PhoneNumber),
		sql.Named("BookingDate", booking.BookingDate),
		sql.Named("BookingLength", booking.BookingLength),
		sql.Named("BookingTime", booking.BookingTime),
		sql.Named("NumAdults", booking.NumAdults),
		sql.Named("NumChildren", booking.NumChildren),
		sql.Named("PromoCode", booking.PromoCode),
		sql.Named("TotalCost", booking.TotalCost),
	)
	if err != nil {
		http.Error(w, "Error inserting data into the database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Booking successfully submitted"})
}

// loginForm to handle the login form login.html and check if the user is authenticated
func loginForm(w http.ResponseWriter, r *http.Request) {
	var data Data

	if r.Method != "POST" {
		fmt.Println("Ontvangen verzoek voor inlogpagina")
		tmpl := template.Must(template.ParseFiles("reservation/login.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Fout bij het parseren van het formulier", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	var hashedPassword string
	var firstname string
	var ssoUser bool
	query := "SELECT wachtwoord, voornaam, sso_user FROM persoonsgegevens WHERE email = @Email"
	err = db.QueryRow(query, sql.Named("Email", email)).Scan(&hashedPassword, &firstname, &ssoUser)
	if err != nil {
		if err == sql.ErrNoRows {
			data.ErrorMessage = "Ongeldige inloggegevens"
			fmt.Println("Ongeldige inloggegevens")
		} else {
			log.Println("Fout bij het raadplegen van de database:", err)
			http.Error(w, "Fout bij het raadplegen van de database", http.StatusInternalServerError)
			return
		}
	} else if ssoUser {
		// Check if the user is an SSO user and block login
		data.ErrorMessage = "Login without SSO is not allowed for this user."
		fmt.Println("Login without SSO is not allowed for this user.")

	} else {
		// Vergelijk hashwachtwoord met ingevoerd wachtwoord
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			fmt.Println("Ongeldige hashinloggegevens")
			data.ErrorMessage = "Ongeldige inloggegevens"
		} else {
			// Als alles goed is, log dat de gebruiker is ingelogd
			fmt.Println("Gebruiker", firstname, "is ingelogd")

			//https://www.sohamkamani.com/golang/session-cookie-authentication/
			// Create a new random session token
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(120 * time.Minute)

			// Set the token in the session map, along with the session information
			sessions[sessionToken] = Session{
				username: firstname,
				expiry:   expiresAt,
			}

			// Finally, we set the client cookie for "session_token" as the session token we just generated
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expiresAt,
			})

			// Stuur de gebruiker door naar "accommodation.html"
			http.Redirect(w, r, "/accommodation", http.StatusSeeOther)
			return
		}
	}

	// Stel de foutmelding in op de template
	tmpl := template.Must(template.ParseFiles("reservation/login.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// parkenForm to handle the parken form parken.html
func parkenForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for parken page")
	http.ServeFile(w, r, "parken.html")
}

// contactForm to handle the contact form + query the database for all accommodations
func contactForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for contact page")

	query := "SELECT naam, beschrijving, prijs, imgurl, locatie FROM accommodations"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var accommodations []Accommodation

	for rows.Next() {
		var acc Accommodation
		if err := rows.Scan(&acc.Name, &acc.Des, &acc.Price, &acc.ImageURL, &acc.Location); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		accommodations = append(accommodations, acc)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating rows", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("contact.html"))
	err = tmpl.Execute(w, accommodations)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// authenticate to check if a user is authenticated
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// Als er geen cookie is, redirect naar login
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			// Voor andere fouten, een interne serverfout retourneren
			http.Error(w, "Fout bij het verkrijgen van cookie", http.StatusInternalServerError)
			return
		}

		sessionToken := cookie.Value
		userSession, exists := sessions[sessionToken]
		if !exists || userSession.isExpired() {
			// Als de sessie niet bestaat of is verlopen, redirect naar login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Als de sessie geldig is, roep de volgende handler aan
		next.ServeHTTP(w, r)
	})
}

// serveRegisterForm
func serveRegisterForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for regestratie page")
	http.ServeFile(w, r, "reservation/register.html")
}

// submitForm to submit the form
func submitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Haal de gegevens uit het formulier op
	firstname := r.Form.Get("firstname")
	lastname := r.Form.Get("lastname")
	email := r.Form.Get("email")
	licenseplate := r.Form.Get("license_plate")
	password := r.Form.Get("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password contact the administrator for more information", http.StatusInternalServerError)
		return
	}
	// Voer de query uit om de gebruiker op te slaan in de database
	query := "INSERT INTO persoonsgegevens (voornaam, achternaam, email, wachtwoord, kenteken) VALUES (@FirstName, @LastName, @Email, @Password, @LicensePlate)"
	// Execute the query naar de database
	_, err = db.Exec(query,
		sql.Named("FirstName", firstname),
		sql.Named("LastName", lastname),
		sql.Named("Email", email),
		sql.Named("Password", hashedPassword),
		sql.Named("LicensePlate", licenseplate))
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	// Retrieve the saved user to verify and display confirmation
	var storedFirstName string
	var storedHashedPassword string
	query = "SELECT wachtwoord, voornaam FROM persoonsgegevens WHERE email = @Email"
	row := db.QueryRow(query, sql.Named("Email", email))
	err = row.Scan(&storedHashedPassword, &storedFirstName)
	if err != nil {
		// If there's an error or no row found for the given email, return unauthorized status
		http.Error(w, "Invalid email or password contact the administrator for more information", http.StatusUnauthorized)
		return
	}
	formData := FormData{
		FirstName:    firstname,
		LastName:     lastname,
		Email:        email,
		LicensePlate: licenseplate,
		Password:     password,
	}
	tmpl := template.Must(template.ParseFiles("reservation/submit.html"))
	err = tmpl.Execute(w, formData)
	if err != nil {
		http.Error(w, "Error rendering template contact the administrator for more information", http.StatusInternalServerError)
		return
	}
}

// adminForm to handle the admin form
func adminForm(w http.ResponseWriter, r *http.Request) {
	var data Data

	// Show the login page to the user
	if r.Method == "GET" && r.URL.Path == "/admin/" {
		fmt.Println("Ontvangen verzoek voor inlogpagina")
		tmpl := template.Must(template.ParseFiles("admin/admin.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	// Workaround to serve the login.js file
	if r.Method == "GET" && r.URL.Path == "/admin/login.js" {
		http.ServeFile(w, r, "admin/login.js")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Fout bij het parseren van het formulier", http.StatusBadRequest)
		return
	}

	email := r.Header.Get("upn")
	password := r.Header.Get("connectionstring")

	var hashedPassword string
	var firstname string
	var ssoUser bool
	query := "SELECT wachtwoord, voornaam, sso_user FROM persoonsgegevens WHERE email = @Email"
	err = db.QueryRow(query, sql.Named("Email", email)).Scan(&hashedPassword, &firstname, &ssoUser)
	if err != nil {
		if err == sql.ErrNoRows {
			// Als er geen gebruiker is gevonden, stel de foutmelding in
			data.ErrorMessage = "Ongeldige inloggegevens"
			fmt.Println("Ongeldige inloggegevens")
			http.Error(w, "Ongeldige inloggegevens", http.StatusUnauthorized)
		} else {
			log.Println("Fout bij het raadplegen van de database:", err)
			http.Error(w, "Fout bij het raadplegen van de database", http.StatusInternalServerError)
			return
		}
	} else if !ssoUser { // Check if ssoUser is false
		data.ErrorMessage = "Toegang geweigerd. Alleen SSO-gebruikers mogen inloggen."
		fmt.Println("Toegang geweigerd. Alleen SSO-gebruikers mogen inloggen.")

	} else {
		// Vergelijk het ingevoerde wachtwoord met het gehashte wachtwoord in de database
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			// Als het wachtwoord niet overeenkomt, stel de foutmelding in
			fmt.Println("Ongeldige hashinloggegevens")
			data.ErrorMessage = "Ongeldige inloggegevens"
		} else {
			// Als alles goed is, log dat de gebruiker is ingelogd
			fmt.Println("Gebruiker", firstname, "is ingelogd")

			// Create a new random session token
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(120 * time.Minute)

			// Set the token in the session map, along with the session information
			sessions[sessionToken] = Session{
				username: firstname,
				expiry:   expiresAt,
			}

			// Finally, we set the client cookie for "session_token" as the session token we just generated
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expiresAt,
			})
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}
	}

	// Stel de foutmelding in op de template
	tmpl := template.Must(template.ParseFiles("admin/admin.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// dashboardForm to handle the dashboard form
func dashboardForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for dashboard page")
	query := "SELECT ID, naam, locatie, email, indatum, uitdatum FROM reservering"

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		log.Println("Error querying database:", err)
		return
	}
	defer rows.Close()

	var reservations []Reservation

	for rows.Next() {
		var res Reservation
		err := rows.Scan(&res.ID, &res.Name, &res.Location, &res.Email, &res.StartDate, &res.EndDate)
		if err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			log.Println("Error scanning data:", err)
			return
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		log.Println("Error iterating over rows:", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("admin/dashboard.html"))
	err = tmpl.Execute(w, reservations)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Error rendering template:", err)
	}
}

// removeReservationHandler to handle the remove reservation
func removeReservationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "Missing reservation ID", http.StatusBadRequest)
		return
	}

	// Haal de locatie (locatie) van de accommodatie op basis van de reserverings-ID
	var location string
	query := "SELECT locatie FROM reservering WHERE id = @id"
	err := db.QueryRow(query, sql.Named("id", id)).Scan(&location)
	if err != nil {
		http.Error(w, "Error retrieving accommodation location", http.StatusInternalServerError)
		log.Println("Error retrieving accommodation location:", err)
		return
	}

	// Voer de query uit om de beschikbaarheid van de accommodatie in te stellen op 1
	updateQuery := "UPDATE accommodations SET beschikbaar = 1 WHERE locatie = @locatie"
	_, err = db.Exec(updateQuery, sql.Named("locatie", location))
	if err != nil {
		http.Error(w, "Error updating accommodation availability", http.StatusInternalServerError)
		log.Println("Error updating accommodation availability:", err)
		return
	}

	// Voer de SQL DELETE-query uit om de reservering te verwijderen
	deleteQuery := "DELETE FROM reservering WHERE id = @id"
	_, err = db.Exec(deleteQuery, sql.Named("id", id))
	if err != nil {
		http.Error(w, "Error deleting reservation", http.StatusInternalServerError)
		log.Println("Error deleting reservation:", err)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

// reserveHandler to handle de reserveren
func reserveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Haal de sessie op
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	sessionToken := cookie.Value
	session, exists := sessions[sessionToken]
	if !exists || session.isExpired() {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Haal de accommodatiegegevens op uit het formulier
	accommodationID := r.FormValue("accommodation_id")
	accommodationName := r.FormValue("accommodation_name")
	accommodationLocation := r.FormValue("accommodation_location")
	if accommodationID == "" || accommodationName == "" || accommodationLocation == "" {
		http.Error(w, "Missing accommodation details", http.StatusBadRequest)
		return
	}

	// Haal de gebruikersgegevens (naam en e-mail) op basis van de sessie
	var userEmail, userName string
	query := "SELECT email, voornaam FROM persoonsgegevens WHERE voornaam = @FirstName"
	err = db.QueryRow(query, sql.Named("FirstName", session.username)).Scan(&userEmail, &userName)
	if err != nil {
		http.Error(w, "Error retrieving user details", http.StatusInternalServerError)
		log.Println("Error retrieving user details:", err)
		return
	}

	// Voeg de reservering toe aan de database
	insertQuery := "INSERT INTO reservering (locatie, email, indatum, uitdatum, naam) VALUES (@locatie, @Email, @indatum, @uitdatum, @naam)"
	_, err = db.Exec(insertQuery,
		sql.Named("locatie", accommodationLocation),
		sql.Named("Email", userEmail),
		sql.Named("indatum", time.Now()),
		sql.Named("uitdatum", time.Now().AddDate(0, 0, 7)),
		sql.Named("naam", accommodationName))
	if err != nil {
		http.Error(w, "Error creating reservation", http.StatusInternalServerError)
		log.Println("Error creating reservation:", err)
		return
	}

	// Zet de beschikbaarheid van de accommodatie op 0
	updateQuery := "UPDATE accommodations SET beschikbaar = 0 WHERE ID = @ID"
	_, err = db.Exec(updateQuery, sql.Named("ID", accommodationID))
	if err != nil {
		http.Error(w, "Error updating accommodation availability", http.StatusInternalServerError)
		log.Println("Error updating accommodation availability:", err)
		// Als er een fout optreedt, zou het goed zijn om de reservering ongedaan te maken bekijk fuction undoReservation
		undoReservation(userEmail, accommodationName)
		return
	}

	// Redirect naar de accommodatiepagina met een bevestigingsbericht
	http.Redirect(w, r, "/accommodation?reserved=1", http.StatusSeeOther)
}

// Functie om reservering ongedaan te maken als er een fout optreedt
func undoReservation(email, accommodationName string) {
	undoQuery := "DELETE FROM reservering WHERE email = @Email AND naam = @naam"
	_, err := db.Exec(undoQuery, sql.Named("Email", email), sql.Named("naam", accommodationName))
	if err != nil {
		log.Println("Error undoing reservation:", err)
	}
}

// carbonfootprint to handle the carbonfootprint form
func carbonfootprintform(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for carbonfootprint page")

	// Query om alle data op te halen
	query := `
SELECT
    p.id AS park_id,
    convert(varchar, e.datum, 101) AS datum,
    p.naam AS park_name,
    SUM(e.verbruik) AS electricity_usage,
    SUM(g.verbruik) AS gas_usage,
    SUM(w.verbruik) AS water_usage,
    FORMAT(AVG(t.temperatuur), 'N2') AS temperature
FROM
    dbo.electricity_usage e
LEFT JOIN
    dbo.park p ON e.parkid = p.id
LEFT JOIN
    dbo.gas_usage g ON e.datum = g.datum AND e.parkid = g.parkid
LEFT JOIN
    dbo.water_usage w ON e.datum = w.datum AND e.parkid = w.parkid
LEFT JOIN
    dbo.temperature t ON e.datum = t.datum AND e.parkid = t.parkid
GROUP BY
    p.id, e.datum, p.naam
ORDER BY
    p.naam, e.datum;
		`

	// Verwerken van de resultaten
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		log.Println("Error querying database:", err)
		return
	}
	defer rows.Close()

	var usagedata []UsageData
	for rows.Next() {
		var data UsageData
		if err := rows.Scan(
			&data.ID, &data.Datum, &data.ParkName, &data.ElectricityUsage, &data.GasUsage, &data.WaterUsage, &data.Temperature,
		); err != nil {
			log.Fatal(err)
		}
		usagedata = append(usagedata, data)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseFiles("admin/carbonfootprint.html"))

	// Convert usagedata to JSON
	usagedataJSON, err := json.Marshal(usagedata)
	if err != nil {
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		log.Println("Error processing data:", err)
		return
	}

	// Pass the JSON string to the template
	data := struct {
		UsageDataJSON string
	}{
		UsageDataJSON: string(usagedataJSON),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Error rendering template:", err)
	}
}
