const express = require('express');
const app = express();
const port = 3000;

var registeredUsers = [];
var counter = 0;

var currentHolds = [];

app.get('/user/holds', async (req, res) => {
  let holdId = currentHolds.find((hold) => hold.holderId == req.query.id)?.holdId;
  console.log(holdId);

  res.send({ holdId: holdId });
});

app.get('/holds', (req, res) => {
  res.json(currentHolds);
});

app.get('/users', (req, res) => {
  res.json(registeredUsers);
});


app.listen(port, () => {
  console.log(`Example app listening on port ${port}`);
});
