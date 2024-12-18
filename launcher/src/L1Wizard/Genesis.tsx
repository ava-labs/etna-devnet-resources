import { useState } from 'react';
import { getWalletAddress } from './wallet';
import { useWizardStore } from './store';

export default function Genesis() {
    const [address, setAddress] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    const { advanceFrom } = useWizardStore();

    const handleConnectWallet = async () => {
        setIsLoading(true);
        setError('');
        try {
            const walletAddress = await getWalletAddress();
            setAddress(walletAddress);
        } catch (err: any) {
            setError(err.message || 'Failed to connect wallet');
        } finally {
            setIsLoading(false);
        }
    };

    const handleContinue = () => {
        advanceFrom('genesis')
    }

    return (
        <div className="max-w-2xl">
            <h1 className="text-2xl font-medium mb-6">Genesis Settings</h1>

            {error && (
                <div className="mb-4 p-3 text-sm text-red-500 bg-red-50 rounded-md">
                    {error}
                </div>
            )}

            <div className="mb-6">
                <div className="flex gap-3 items-start">
                    <div className="flex-grow">
                        <input
                            type="text"
                            value={address}
                            onChange={(e) => setAddress(e.target.value)}
                            placeholder="Wallet Address"
                            className="w-full p-2 border border-gray-200 rounded-md"
                        />
                    </div>
                    <button
                        onClick={handleConnectWallet}
                        disabled={isLoading}
                        className={`px-4 py-2 rounded-md ${isLoading
                            ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                            : 'bg-gray-100 hover:bg-gray-200'
                            }`}
                    >
                        {isLoading ? 'Loading...' : 'Fill from Wallet'}
                    </button>
                </div>
                <p className="mt-2 text-sm text-gray-500">
                    This address will receive all tokens and control in case of Proof of Authority chain.
                </p>
            </div>

            <button
                onClick={handleContinue}
                disabled={!isValidEthereumAddress(address)}
                className={`px-4 py-2 rounded-md ${!isValidEthereumAddress(address)
                    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                    : 'bg-blue-500 text-white hover:bg-blue-600'
                    }`}
            >
                Continue
            </button>
        </div>
    );
}

function isValidEthereumAddress(address: string) {
    return /^0x[a-fA-F0-9]{40}$/.test(address);
}
