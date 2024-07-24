const express = require('express');
const bodyParser = require('body-parser');
const sql = require('mssql');

// SQL Server configuration
const sqlConfig = {
    user: 'Bramplatier',
    password: 'Brapla_51',
    database: 'sqlres',
    server: 'sql-srv-res.database.windows.net',
    pool: {
        max: 10,
        min: 0,
        idleTimeoutMillis: 30000
    },
    options: {
        encrypt: true,
        trustServerCertificate: false
    }
};

const app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

app.post('/submit-bowling', async (req, res) => {
    try {
        const { firstName, lastName, phoneNumber, bookingDate, bookingLength, bookingTime, numAdults, numChildren, promoCode, totalCost } = req.body;

        let pool = await sql.connect(sqlConfig);

        await pool.request()
            .input('firstName', sql.NVarChar, firstName)
            .input('lastName', sql.NVarChar, lastName)
            .input('phoneNumber', sql.NVarChar, phoneNumber)
            .input('bookingDate', sql.Date, bookingDate)
            .input('bookingLength', sql.Int, bookingLength)
            .input('bookingTime', sql.Time, bookingTime)
            .input('numAdults', sql.Int, numAdults)
            .input('numChildren', sql.Int, numChildren)
            .input('promoCode', sql.NVarChar, promoCode)
            .input('totalCost', sql.Decimal(10, 2), totalCost)
            .query(`INSERT INTO Bookings 
                    (FirstName, LastName, PhoneNumber, BookingDate, BookingLength, BookingTime, NumAdults, NumChildren, PromoCode, TotalCost)
                    VALUES 
                    (@firstName, @lastName, @phoneNumber, @bookingDate, @bookingLength, @bookingTime, @numAdults, @numChildren, @promoCode, @totalCost)`);

        res.status(200).json({ status: 'success', message: 'Booking successfully submitted' });
    } catch (err) {
        console.error('SQL error', err);
        res.status(500).json({ status: 'error', message: 'Error submitting booking' });
    }
});

const PORT = process.env.PORT || 1433;
app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
});
