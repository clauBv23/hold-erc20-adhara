const transactions = require('./src/transactions');
const express = require('express');
const app = express();
const port = 3000;

const bodyParser = require('body-parser');
const multer = require('multer');
const upload = multer(); // for parsing multipart/form-data

app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({ extended: true })); // for parsing application/x-www-form-urlencoded

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

app.post('/reg', upload.array(), async (req, res) => {
  // check if the user has been registered
  let userId = registeredUsers.find((user) => user.address == req.body.addr)?.id;

  if (userId) {
    // user already registered
    res.send({ registerId: userId });
  } else {
    //register the user
    counter++;
    registeredUsers.push({ id: counter, address: req.body.addr });

    // mint the tokens
    let mintTx = await transactions.mint(req.body.addr);

    // send the transction and the id of the user registered
    res.send({ registerId: counter, mintTx: mintTx });
  }
});

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`);
});
