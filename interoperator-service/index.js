const express = require('express');
const app = express();

app.get('/', (req, res) => res.send('Interoperator Service: ready'));
app.listen(3000, '0.0.0.0');
