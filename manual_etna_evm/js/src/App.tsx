import { useEffect } from "react";
import Balance from "./components/balance/Balance";
import { MINIMUM_P_CHAIN_BALANCE_AVAX } from "./components/balance/balances";
import { Keys } from "./components/keys/Keys";
import Card from "./lib/Card";
import { useWalletStore } from "./lib/store";
import { getLocalStorageWallet } from "./lib/wallet";
import { useAsync } from "./lib/hooks";

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

  useEffect(() => {
    if (!wallet) {
      generateWalletPromise.execute();
    }
  }, []);

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
      {cAddress && pAddress && <Card title="💰 Wallet Balance">
        <Balance />
      </Card>}
      {pBalance > MINIMUM_P_CHAIN_BALANCE_AVAX && <Card title="Create subnet">
        ...
      </Card>}
    </div >
  )
}

export default App
