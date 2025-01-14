const express = require('express');
const app = express();
const port = 3000;

app.use(express.static('public'));  // Serves index.html
app.use('/assets', express.static('assets'));  // Serves the font

app.listen(port, () => {
    console.log(`Server running at http://localhost:${port}`);
});
