import { useState, useEffect } from "react";
import Balance from "./components/balance/Balance";
import { MINIMUM_P_CHAIN_BALANCE_AVAX, P_CHAIN_DIVISOR } from "./components/balance/balances";
import { Keys } from "./components/keys/Keys";
import Card from "./lib/Card";
import { useWalletStore } from "./components/balance/walletStore";
import { getLocalStorageWallet } from "./lib/wallet";
import { useAsync } from "./lib/hooks";
import CreateSubnet from "./components/subnet/CreateSubnet";
import { useSubnetStore } from "./components/subnet/subnetStore";
import GenerateGenesis from "./components/genesis/GenerateGenesis";
import { useGenesisStore } from "./components/genesis/genesisStore";
import CreateChain from "./components/chain/CreateChain";
import { useChainStore } from "./components/chain/chainStore";
import ConvertChain from "./components/convert/ConvertChain";

function App() {
  const wallet = useWalletStore(state => state.wallet);
  const pAddress = useWalletStore(state => state.pAddress);
  const cAddress = useWalletStore(state => state.cAddress);
  const pBalance = useWalletStore(state => state.pBalance);
  const subnetId = useSubnetStore(state => state.subnetId);
  const genesis = useGenesisStore(state => state.genesis);
  const chainId = useChainStore(state => state.chainId);

  const setWallet = useWalletStore(state => state.setWallet);
  const generateWalletPromise = useAsync(async () => {
    const wallet = getLocalStorageWallet("https://etna.avax-dev.network");
    await setWallet(wallet);
  });

  const [currentStep, setCurrentStep] = useState(1);


  const STEP_GENERATE_KEYS = 1;
  const STEP_WALLET_BALANCE = 2;
  const STEP_CREATE_SUBNET = 3;
  const STEP_GENERATE_GENESIS = 4;
  const STEP_CREATE_CHAIN = 5;
  const STEP_CONVERT_CHAIN = 6;

  useEffect(() => {
    let step = STEP_GENERATE_KEYS;

    if (cAddress && pAddress) {
      step = STEP_WALLET_BALANCE;
    }

    if (step === STEP_WALLET_BALANCE && pBalance > BigInt(MINIMUM_P_CHAIN_BALANCE_AVAX * P_CHAIN_DIVISOR)) {
      step = STEP_CREATE_SUBNET;
    }

    if (step === STEP_CREATE_SUBNET && subnetId) {
      step = STEP_GENERATE_GENESIS;
    }

    if (step === STEP_GENERATE_GENESIS && genesis.length > 0) {
      step = STEP_CREATE_CHAIN;
    }

    if (step === STEP_CREATE_CHAIN && chainId) {
      step = STEP_CONVERT_CHAIN;
    }

    setCurrentStep(step);
  }, [cAddress, pAddress, pBalance, genesis, chainId]);

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
      {currentStep >= STEP_GENERATE_KEYS && <Card title="🔑 Generate keys" >
        <Keys />
      </Card>}
      {currentStep >= STEP_WALLET_BALANCE && <Card title="💰 Wallet Balance">
        <Balance />
      </Card>}
      {currentStep >= STEP_CREATE_SUBNET && <Card title="🌐 Create subnet">
        <CreateSubnet />
      </Card>}
      {currentStep >= STEP_GENERATE_GENESIS && <Card title="🚀 Generate genesis">
        <GenerateGenesis canRegenerate={currentStep <= STEP_CREATE_CHAIN} />
      </Card>}
      {currentStep >= STEP_CREATE_CHAIN && <Card title="🔗 Create chain">
        <CreateChain />
      </Card>}
      {currentStep >= STEP_CONVERT_CHAIN && <Card title="🔄 Convert chain">
        <ConvertChain />
      </Card>}
    </div >
  )
}

export default App
