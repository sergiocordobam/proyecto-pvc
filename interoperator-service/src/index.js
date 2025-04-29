const express = require('express');
const operatorRoutes = require('./Routes/Operator_Routes'); // Import the routes file

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(express.json());

// Routes
app.use('/comunication', operatorRoutes); // Mount the routes under '/comunication'

app.get('/', (req, res) => {
    res.send('Welcome to the Node.js App!');
});

// Start the server
app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});