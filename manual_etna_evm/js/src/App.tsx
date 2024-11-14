import { useState, useEffect } from "react";
import Balance from "./components/balance/Balance";
import { MINIMUM_P_CHAIN_BALANCE_AVAX } from "./components/balance/balances";
import { Keys } from "./components/keys/Keys";
import Card from "./lib/Card";
import { useWalletStore } from "./lib/store";
import { getLocalStorageWallet } from "./lib/wallet";
import { useAsync } from "./lib/hooks";
import CreateSubnet from "./components/subnet/CreateSubnet";

function App() {
  const wallet = useWalletStore(state => state.wallet);
  const pAddress = useWalletStore(state => state.pAddress);
  const cAddress = useWalletStore(state => state.cAddress);
  const pBalance = useWalletStore(state => state.pBalance);

  const setWallet = useWalletStore(state => state.setWallet);
  const generateWalletPromise = useAsync(async () => {
    const wallet = getLocalStorageWallet("https://etna.avax-dev.network");
    await setWallet(wallet);
  });

  const [currentStep, setCurrentStep] = useState(1);


  const STEP_GENERATE_KEYS = 1;
  const STEP_WALLET_BALANCE = 2;
  const STEP_CREATE_SUBNET = 3;

  useEffect(() => {
    let step = STEP_GENERATE_KEYS;

    if (cAddress && pAddress) {
      step = STEP_WALLET_BALANCE;
    }

    if (step === STEP_WALLET_BALANCE && pBalance > MINIMUM_P_CHAIN_BALANCE_AVAX) {
      step = STEP_CREATE_SUBNET;
    }

    setCurrentStep(step);
  }, [cAddress, pAddress, pBalance]);

  useEffect(() => {
    if (!wallet) {
      generateWalletPromise.execute();
    }
  }, [wallet]);

  if (generateWalletPromise.error) {
    return <div>Error generating wallet: {generateWalletPromise.error}</div>;
  }

  if (!wallet || generateWalletPromise.loading) {
    return <div>Generating wallet...</div>;
  }

  return (
    <div className="h-full container mx-auto">
      <Card title="🔑 Generate keys" >
        <Keys />
      </Card>
      {currentStep >= STEP_WALLET_BALANCE && <Card title="💰 Wallet Balance">
        <Balance />
      </Card>}
      {currentStep >= STEP_CREATE_SUBNET && <Card title="Create subnet">
        <CreateSubnet />
      </Card>}
    </div >
  )
}

export default App
