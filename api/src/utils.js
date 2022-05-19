const transactions = require('./transactions');

const makeTheBets = async (holds, users) => {
  // execute all the holds
  for (let i = 0; i < holds.length; i++) {
    await transactions.executeHold(holds[i].holdId);
  }

  // get the random winner
  let winner = Math.floor(Math.random() * 4);

  // get the winner address
  let winnerAddr = users.find((user) => (user.id = holds[winner].holderId))?.address;

  console.log('winnerAddr', winnerAddr);
  // transfer the funds to the winner
  await transactions.transfer(winnerAddr);

  return winner;
};

module.exports = { makeTheBets };
