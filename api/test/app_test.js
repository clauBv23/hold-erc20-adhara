let chai = require('chai');
let chaiHttp = require('chai-http');
const expect = require('chai').expect;
var assert = require('assert');

chai.use(chaiHttp);
const url = 'http://localhost:3000';

let newUsersIds = [0, 0, 0];
let userBetIds = [0, 0, 0];

let usersAdress = [
  '0x37D284F6E4EF4F174cE326d0F9450133c89Cf071',
  '0x65210F41E2450c1374cF3f668ed5B3d2Cfd07216',
  '0xeCA850e66C2557B9651F6979fA01b0C6a4a7933c',
];

let winner;

// check a complete bet winning process
describe('the bet-winning process', () => {
  it('should register user 1', (done) => {
    chai
      .request(url)
      .post('/reg')
      .send({ addr: usersAdress[0] })
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // store the user id
        newUsersIds[0] = res.body.registerId;
        // check the user id is correct
        assert.notEqual(newUsersIds[0], 0);
        done();
      });
  });

  it('should register user 2', (done) => {
    chai
      .request(url)
      .post('/reg')
      .send({ addr: usersAdress[1] })
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // store the user id
        newUsersIds[1] = res.body.registerId;
        // check the user id is correct
        assert.notEqual(newUsersIds[1], 0);
        done();
      });
  });

  it('should register user 3', (done) => {
    chai
      .request(url)
      .post('/reg')
      .send({ addr: usersAdress[2] })
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // store the user id
        newUsersIds[2] = res.body.registerId;
        // check the user id is correct
        assert.notEqual(newUsersIds[2], 0);
        done();
      });
  });

  it('check there are 3 registered users', (done) => {
    chai
      .request(url)
      .get('/users')
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // 3 registered users
        assert.equal(res.body.length, 3);
        done();
      });
  });

  it("check the first user's balance is 100", (done) => {
    chai
      .request(url)
      .get('/balance?id=' + newUsersIds[0])
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // check the balances
        assert.equal(res.body.balance, 100);
        done();
      });
  });

  it('should make a bet to user 1', (done) => {
    chai
      .request(url)
      .post('/bet')
      .send({ userId: newUsersIds[0] })
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // store the bets ids
        userBetIds[0] = res.body.holdId;
        //check the bets id is correct
        assert.notEqual(userBetIds[0], 0);
        done();
      });
  });

  it('should make a bet to user 2', (done) => {
    chai
      .request(url)
      .post('/bet')
      .send({ userId: newUsersIds[1] })
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // store the bets ids
        userBetIds[1] = res.body.holdId;
        //check the bets id is correct
        assert.notEqual(userBetIds[1], 0);
        done();
      });
  });

  it('should make a bet to user 3', (done) => {
    chai
      .request(url)
      .post('/bet')
      .send({ userId: newUsersIds[2] })
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // store the bets ids
        userBetIds[2] = res.body.holdId;
        //check the bets id is correct
        assert.notEqual(userBetIds[2], 0);
        done();
      });
  });

  it('check the holds stored', (done) => {
    chai
      .request(url)
      .get('/holds')
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // there most be 3 holds one for each user
        assert.equal(res.body.length, 3);
        done();
      });
  });

  it('check the user 1 hold', (done) => {
    chai
      .request(url)
      .get('/user/holds?id=' + newUsersIds[0])
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // check there are one hold
        assert.equal(res.body.holds.length, 1);
        done();
      });
  });

  it('check the users balances were decreased', (done) => {
    chai
      .request(url)
      .get('/balance?id=' + newUsersIds[0])
      .end(function (err, res) {
        expect(res).to.have.status(200);
        // check the balance decreased in 5 tokens
        assert.equal(res.body.balance, 95);
        done();
      });
  });

  // execute the last bet to run the process
  describe('Make the 4th bet', () => {
    it('call the 4th bet', (done) => {
      chai
        .request(url)
        .post('/bet')
        .send({ userId: newUsersIds[0] })
        .end(function (err, res) {
          winner = res.body.winner;

          expect(res).to.have.status(200);
          done();
        });
    });

    it('check the the tokens were transferred to the winner', (done) => {
      chai
        .request(url)
        .get('/balance?id=' + newUsersIds[winner])
        .end(function (err, res) {
          // check the hold is created
          assert.equal(res.body.balance, 115);
          expect(res).to.have.status(200);
          done();
        });
    });
  });
});
