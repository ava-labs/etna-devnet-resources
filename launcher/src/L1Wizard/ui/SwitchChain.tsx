import { useState, useEffect, ReactNode } from 'react';

interface Props {
    children: ReactNode;
}

type Status = 'not_started' | 'wrong_chain' | 'success';

export default function SwitchChain({ children }: Props) {
    const [chainStatus, setChainStatus] = useState<Status>('not_started');

    // Check if user is on the right chain
    const checkChain = async () => {
        if (!window.ethereum) {
            setChainStatus('wrong_chain');
            return;
        }

        try {
            const chainId = await window.ethereum.request({ method: 'eth_chainId' });
            if (chainId === '0xa869') { // 43113 in hex (Fuji)
                setChainStatus('success');
            } else {
                setChainStatus('wrong_chain');
            }
        } catch (error) {
            setChainStatus('wrong_chain');
        }
    };

    useEffect(() => {
        checkChain();
        // Listen for chain changes
        if (window.ethereum) {
            window.ethereum.on('chainChanged', checkChain);
            return () => {
                window.ethereum.removeListener('chainChanged', checkChain);
            };
        }
    }, []);

    const switchToFuji = async () => {
        if (!window.ethereum) return;

        try {
            await window.ethereum.request({
                method: 'wallet_switchEthereumChain',
                params: [{ chainId: '0xa869' }], // 43113 in hex
            });
        } catch (error: any) {
            // If the chain hasn't been added to MetaMask, add it
            if (error.code === 4902) {
                try {
                    await window.ethereum.request({
                        method: 'wallet_addEthereumChain',
                        params: [{
                            chainId: '0xa869',
                            chainName: 'Avalanche Fuji Testnet',
                            nativeCurrency: {
                                name: 'Avalanche',
                                symbol: 'AVAX',
                                decimals: 18
                            },
                            rpcUrls: ['https://api.avax-test.network/ext/bc/C/rpc'],
                            blockExplorerUrls: ['https://testnet.snowtrace.io/']
                        }]
                    });
                } catch (addError) {
                    console.error('Failed to add network:', addError);
                }
            }
        }
    };

    if (chainStatus === 'success') {
        return <>{children}</>;
    }

    return (
        <div className="p-4 border rounded-lg">
            <h3 className="font-medium mb-4">Network Check</h3>
            <button
                onClick={switchToFuji}
                className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
            >
                Switch to Avalanche Fuji
            </button>
        </div>
    );
}
