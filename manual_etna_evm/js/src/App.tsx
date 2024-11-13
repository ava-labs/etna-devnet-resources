import Balance from "./components/balance/Balance";
import { getLocalStorageWallet } from "./lib/wallet";

function App() {
  const wallet = getLocalStorageWallet("https://etna.avax-dev.network");
  return (
    <>
      <Balance wallet={wallet} />
    </>
  )
}

export default App
