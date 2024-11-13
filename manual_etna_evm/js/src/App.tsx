import Balance from "./components/balance/Balance";
import { Keys } from "./components/keys/Keys";
import Card from "./lib/Card";
import { getLocalStorageWallet } from "./lib/wallet";

function App() {
  const wallet = getLocalStorageWallet("https://etna.avax-dev.network");
  return (
    <div className="h-full container mx-auto">
      <Card title="🔑 Generate keys" >
        <Keys wallet={wallet} />
      </Card>
      <Card title="💰 Wallet Balance">
        <Balance wallet={wallet} />
      </Card>
    </div >
  )
}

export default App
