import Balance from "./components/balance/Balance";
import { MINIMUM_P_CHAIN_BALANCE } from "./components/balance/balances";
import { Keys } from "./components/keys/Keys";
import Card from "./lib/Card";
import { useWalletStore } from "./lib/store";
import { getLocalStorageWallet } from "./lib/wallet";

function App() {
  const wallet = getLocalStorageWallet("https://etna.avax-dev.network");
  const pAddress = useWalletStore(state => state.pAddress);
  const cAddress = useWalletStore(state => state.cAddress);
  const pBalance = useWalletStore(state => state.pBalance);


  return (
    <div className="h-full container mx-auto">
      <Card title="🔑 Generate keys" >
        <Keys wallet={wallet} />
      </Card>
      {cAddress && pAddress && <Card title="💰 Wallet Balance">
        <Balance wallet={wallet} />
      </Card>}
      {pBalance > MINIMUM_P_CHAIN_BALANCE && <Card title="Create subnet">
        ...
      </Card>}
    </div >
  )
}

export default App
