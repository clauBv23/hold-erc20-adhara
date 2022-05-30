const transactions = require('./src/transactions');
const utils = require('./src/utils');
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

// GET: user/holds?id=2
app.get('/user/holds', async (req, res) => {
  // filter the holds of the user
  let holds = currentHolds.filter((hold) => hold.holderId == req.query.id);

  res.send({ holds: holds });
});

// GET: holds
app.get('/holds', (req, res) => {
  res.json(currentHolds);
});

// GET: users
app.get('/users', (req, res) => {
  res.json(registeredUsers);
});

// GET: balance?id=2
app.get('/balance', async (req, res) => {
  // get the user address
  let address = registeredUsers.find((user) => user.id == req.query.id)?.address;
  let balance = 0;

  if (address) {
    balance = await transactions.getBalance(address);
  }

  res.send({ balance: balance });
});

// POST: reg
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

// POST: bet
app.post('/bet', upload.array(), async (req, res) => {
  // get the user address
  let address = registeredUsers.find((user) => user.id == req.body.userId)?.address;
  try {
    if (address) {
      let response = await transactions.holdFrom(address);

      // store the hold
      currentHolds.push({ holderId: req.body.userId, holdId: response.id });

      // check if is the 4th hold
      if (currentHolds.length == 4) {
        // execute the bets
        let winner = await utils.makeTheBets(currentHolds, registeredUsers);

        // empty the current currentHolds
        currentHolds = [];
        res.send({ winner: winner });
      } else {
        res.send({ transaction: response.tx, holdId: response.id });
      }
    } else {
      res.status(500).send('Unregistered user!');
    }
  } catch (err) {
    res.status(500).send('Something went wrong!' + err);
  }
});

app.listen(port, () => {
  console.log(`Listening on port ${port}`);
});
