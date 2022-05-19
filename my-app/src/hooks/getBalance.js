import { useCallback, useEffect, useState } from 'react';

import {getCurrentWalletConnected, loadContract} from '../util/interact.js';
import {utils} from 'ethers';
// Other imports...

export function useGetBalance() {
  const [balance, setBalance] = useState(0);

  const fetchBalance = useCallback(async () => {
    const contract = await loadContract();
    const address = (await getCurrentWalletConnected()).address;

    console.log(address);

    const rawBalance = await contract.methods.balanceOf(address).call();
    const value = utils.formatEther(rawBalance);
    setBalance(value);
  }, []);

  useEffect(() => {
    fetchBalance();
    
  }, [fetchBalance]);

  return balance;
}